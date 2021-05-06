package judger

import (
	"FrontEndOJudger/models"
	"github.com/chromedp/chromedp"
	"strconv"
	"strings"
	"time"
)

type NormalJudger struct {
}


func(n *NormalJudger) RunTests(url string, labTestcase *models.LabTestcase, output *interface{}) chromedp.Action {
	task := chromedp.Tasks{
		chromedp.Navigate(url),
	}
	if labTestcase.WaitBefore != 0 {
		var temp *interface{}
		// 在sleep之前执行一下，需要注意两次执行代码一样，但结果不同，为了保持核心代码和数据表整洁
		task = append(task, chromedp.EvaluateAsDevTools(strings.ReplaceAll(labTestcase.TestcaseCode, "\n", ""), &temp))
		task = append(task, chromedp.Sleep(time.Duration(labTestcase.WaitBefore)*time.Millisecond))
	}
	task = append(task, chromedp.EvaluateAsDevTools(strings.ReplaceAll(labTestcase.TestcaseCode, "\n", ""), &output))
	return task
}

func (n *NormalJudger) PrepareOutput(output *interface{}, testcase models.LabTestcase, testResult *TestResult) {
	if o, ok := (*output).(string); ok {
		testResult.SubmitOutput = o
	}

	if o, ok := (*output).(float64); ok {
		testResult.SubmitOutput = strconv.FormatFloat(o,'f',-1,64)
	}

	if o, ok := (*output).(int64); ok {
		testResult.SubmitOutput = strconv.FormatInt(o, 10)
	}

	if o, ok := (*output).(bool); ok {
		testResult.SubmitOutput = strconv.FormatBool(o)
	}

	if o, ok := (*output).(uint64); ok {
		testResult.SubmitOutput = strconv.FormatUint(o, 10)
	}
}

func (n *NormalJudger) JudgeOutput(testcase models.LabTestcase, testResult *TestResult) {
	testResult.TestcaseOutput = testcase.Output

	if testResult.SubmitOutput == testcase.Output {
		testResult.Status = models.LABSUBMITSTATUS_ACCEPTED
		//log.Printf("#### ACC OUTPUT[%v] id[%d] TESTCASEOUTPUT[%v]", testResult.SubmitOutput, testcase.ID, testcase.Output)
	} else {
		testResult.Status = models.LABSUBMITSTATUS_WRONG_ANSWER
		//log.Printf("#### WA OUTPUT[%v] id[%v] TESTCASEOUTPUT[%v]", testResult.SubmitOutput, testcase.ID, testcase.Output)
	}
}
