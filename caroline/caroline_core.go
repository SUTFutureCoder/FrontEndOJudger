package caroline

import (
	"FrontEndOJudger/caroline/judger"
	"FrontEndOJudger/models"
	"context"
	"errors"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"log"
	"strings"
	"time"
)

//var ctx context.Context
func ExecCaroline(testChamber string, testcases []models.LabTestcase, id uint64, labId uint64) []*judger.TestResult {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var testResults []*judger.TestResult

	// get lab info
	lab := &models.Lab{}
	err := lab.GetFullInfo(labId)
	if err != nil {
		testResults = append(testResults, &judger.TestResult {
			Id: id,
			TestCaseId: 0,
			Status: models.LABSUBMITSTATUS_SYSTEM_ERROR,
		})
		log.Printf("get lab info error:%#v", err)
		return testResults
	}

	// loop testcases
	for testcaseId, testcase := range testcases {
		testResult := &judger.TestResult {
			Id: id,
			TestCaseId: testcaseId,
		}
		ExecTestCase(testChamber, testcase, testResult, lab, &ctx)
		testResults = append(testResults, testResult)
	}

	return testResults
}

func ExecTestCase(testChamber string, testcase models.LabTestcase, testResult *judger.TestResult, lab *models.Lab, ctx *context.Context) {
	testResult.TestCaseInput = testcase.Input
	var output interface{}

	// 实际执行
	run := func () error {
		testCaseCtx := *ctx
		if testcase.TimeLimit != 0 {
			var testCaseCancel context.CancelFunc
			testCaseCtx, testCaseCancel = context.WithTimeout(*ctx, time.Duration(testcase.TimeLimit)*time.Millisecond)
			defer testCaseCancel()
		}

		// 控制台及异常监听
		var exceptions *runtime.ExceptionDetails
		chromedp.ListenTarget(testCaseCtx, func(ev interface{}) {
			switch ev := ev.(type) {
			case *runtime.EventExceptionThrown:
				exceptions = ev.ExceptionDetails
			}
		})

		if err := chromedp.Run(testCaseCtx, judger.GetJudger(lab).RunTests(testChamber, &testcase, &output)); err != nil {
			testResult.Err = err.Error()
			testResult.Status = models.LABSUBMITSTATUS_ERROR
			if testResult.Err == "context deadline exceeded" {
				testResult.Status = models.LABSUBMITSTATUS_TIME_LIMIT_EXCEEDED
			}
			if strings.Contains(testResult.Err, "encountered exception 'Uncaught'") && exceptions != nil{
				testResult.Status = models.LABSUBMITSTATUS_RUNTIME_ERROR
				byteException, errException := exceptions.MarshalJSON()
				if errException != nil {
					testResult.Status = models.LABSUBMITSTATUS_SYSTEM_ERROR
					testResult.Err = errException.Error()
					return errors.New(testResult.Err)
				}
				testResult.Err = string(byteException)
			}
			//log.Printf("#### ERR err[%v] id[%d] testcase[%v] testResult[%v]", err, testcase.ID, testcase.TimeLimit, testResult)
			return err
		}
		return nil
	}

	err := run()
	if err != nil {
		log.Printf("run error[%#v]", err)
		testResult.Status = models.LABSUBMITSTATUS_WRONG_ANSWER
		return
	}

	// prepare result output
	judger.GetJudger(lab).PrepareOutput(&output, testcase, testResult)

	// judge result output
	judger.GetJudger(lab).JudgeOutput(testcase, testResult)

	return
}

func runWithTimeOut(ctx *context.Context, timeout time.Duration, tasks chromedp.Tasks) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		timeoutContext, cancel := context.WithTimeout(ctx, timeout * time.Second)
		defer cancel()
		return tasks.Do(timeoutContext)
	}
}
