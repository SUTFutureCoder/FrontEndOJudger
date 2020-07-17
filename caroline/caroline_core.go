package caroline

import (
	"FrontEndOJudger/models"
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
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
		if err := chromedp.Run(ctx, runTests(testChamber, &testcase, &output)); err != nil {
			testResult.Err = err.Error()
		}

		testResult.SubmitOutput = output.(string)
		testResult.TestcaseOutput = testcase.Output

		if output == testcase.Output {
			testResult.Status = models.LABSUBMITSTATUS_ACCEPTED
		} else {
			testResult.Status = models.LABSUBMITSTATUS_WRONG_ANSWER
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
