package judger

import (
	"FrontEndOJudger/models"
	"github.com/chromedp/chromedp"
	"log"
)

type JudgerI interface {
	// prepare action run tests
	RunTests(url string, labTestcase *models.LabTestcase, output *interface{}) chromedp.Action

	// prepare result output
	PrepareOutput(output *interface{}, testcase models.LabTestcase, testResult *TestResult)

	// judge result output
	JudgeOutput(testcase models.LabTestcase, testResult *TestResult)
}

var JudgerMap map[int]JudgerI

func init() {
	JudgerMap = make(map[int]JudgerI)
	JudgerMap[models.LABTYPE_NORMAL] = &NormalJudger{}
	JudgerMap[models.LABTYPE_IMITATE] = &DeeplearningJudger{}
}

func GetJudger(lab *models.Lab) JudgerI {
	if _, ok := JudgerMap[lab.LabType]; !ok {
		log.Println("lab type not exist when get judger")
		return &NormalJudger{}
	}
	return JudgerMap[lab.LabType]
}
