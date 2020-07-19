package caroline

import (
	"FrontEndOJudger/models"
	"context"
	"fmt"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"log"
	"strings"
	"time"
)

// 测试结果构造
type TestResult struct {
	Id uint64
	TestCaseId int
	SubmitOutput string
	TestcaseOutput string
	Status	int
	Err 	string
}

//var ctx context.Context
func ExecCaroline(testChamber string, testcases []models.LabTestcase, id uint64) []TestResult {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var testResults []TestResult

	for testcaseId, testcase := range testcases {

		testResult := TestResult{
			Id:         id,
			TestCaseId: testcaseId,
		}

		var output interface{}

		testCaseCtx := ctx
		if testcase.TimeLimit != 0 {
			var testCaseCancel context.CancelFunc
			testCaseCtx, testCaseCancel = context.WithTimeout(ctx, time.Duration(testcase.TimeLimit) * time.Millisecond)
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

		if err := chromedp.Run(testCaseCtx, runTests(testChamber, &testcase, &output)); err != nil {
			testResult.Err = err.Error()
			testResult.Status = models.LABSUBMITSTATUS_ERROR
			if testResult.Err == "context deadline exceeded" {
				testResult.Status = models.LABSUBMITSTATUS_TIME_LIMIT_EXCEEDED
			}
			if strings.Contains(testResult.Err, "encountered exception 'Uncaught'") {
				testResult.Status = models.LABSUBMITSTATUS_RUNTIME_ERROR
				byteException, errException :=  exceptions.MarshalJSON()
				if errException != nil {
					testResult.Status = models.LABSUBMITSTATUS_SYSTEM_ERROR
					testResult.Err = errException.Error()
					continue
				}
				testResult.Err = string(byteException)
			}
			log.Printf("#### ERR err[%v] id[%d] testcase[%v] testResult[%v]", err, testcase.ID, testcase.TimeLimit, testResult)
			continue
		}


		testResult.SubmitOutput = output.(string)
		testResult.TestcaseOutput = testcase.Output

		if output == testcase.Output {
			testResult.Status = models.LABSUBMITSTATUS_ACCEPTED
			log.Printf("#### ACC OUTPUT[%v] id[%d] TESTCASEOUTPUT[%v]", output, testcase.ID, testcase.Output)
		} else {
			testResult.Status = models.LABSUBMITSTATUS_WRONG_ANSWER
			log.Printf("#### WA OUTPUT[%v] id[%v] TESTCASEOUTPUT[%v]", output, testcase.ID, testcase.Output)
		}

	}

	return testResults
}

func runTests(url string, labTestcase *models.LabTestcase, output *interface{}) chromedp.Action {
	task := chromedp.Tasks{
		chromedp.Navigate(url),
	}
	if labTestcase.WaitBefore != 0 {
		var temp *interface{}
		// 在sleep之前执行一下，需要注意两次执行代码一样，但结果不同，为了保持核心代码和数据表整洁
		task = append(task, chromedp.EvaluateAsDevTools(strings.ReplaceAll(labTestcase.TestcaseCode, "\n", ""), &temp))
		task = append(task, chromedp.Sleep(time.Duration(labTestcase.WaitBefore)*time.Millisecond))
		fmt.Println(labTestcase.WaitBefore)
	}
	task = append(task, chromedp.EvaluateAsDevTools(strings.ReplaceAll(labTestcase.TestcaseCode, "\n", ""), &output))
	return task
}