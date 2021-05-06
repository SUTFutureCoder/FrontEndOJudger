package judger

import (
	"FrontEndOJudger/models"
	"FrontEndOJudger/pkg/file"
	"FrontEndOJudger/pkg/net"
	"FrontEndOJudger/pkg/setting"
	"FrontEndOJudger/pkg/utils"
	"encoding/json"
	"github.com/chromedp/chromedp"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type DeeplearningJudger struct {
}


func (d *DeeplearningJudger) RunTests(url string, labTestcase *models.LabTestcase, output *interface{}) chromedp.Action {
	testcaseInput := &models.LabTestcaseInputImitate{}
	err := json.Unmarshal([]byte(labTestcase.Input), testcaseInput)
	if err != nil {
		return nil
	}
	return utils.ScreenshotTasks(url, "", output, testcaseInput.TestcaseWidth, testcaseInput.TestcaseHeight, labTestcase.WaitBefore)
}

func (d *DeeplearningJudger) PrepareOutput(output *interface{}, testcase models.LabTestcase, testResult *TestResult) {
	// get path
	var imageBuf []byte
	if o, ok := (*output).([]byte); ok {
		imageBuf = o
	}

	dirPath, _, err := file.GetDest(testResult.Id)
	screenshotFileName := dirPath + "screenshot.png"
	labTestFileName := dirPath + "test.png"
	err = file.PutLocal(imageBuf, screenshotFileName)
	if err != nil {
		log.Printf("file put local error[%#v]", err)
		testResult.Status = models.LABSUBMITSTATUS_SYSTEM_ERROR
		return
	}

	err = net.GetAndSaveFile(testcase.TestcaseCode, labTestFileName)
	if err != nil {
		log.Printf("get and save file error[%#v]", err)
		testResult.Status = models.LABSUBMITSTATUS_SYSTEM_ERROR
		return
	}

	// 执行python逻辑
	cmd := exec.Command(setting.JudgerSetting.DeepLearningPython, setting.JudgerSetting.DeepLearningJudger, labTestFileName, screenshotFileName)
	log.Println(cmd.String())
	outputByte, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("deeplearning exec combineoutput err[%#v]", err)
		return
	}
	testResult.SubmitOutput = strings.Trim(string(outputByte), "\n")
}

func (d *DeeplearningJudger) JudgeOutput(testcase models.LabTestcase, testResult *TestResult) {
	stdoutFloat, err := strconv.ParseFloat(testResult.SubmitOutput, 64)
	if err != nil {
		log.Printf("deeplearning judge output parse float error[%#v]", err)
		testResult.Status = models.LABSUBMITSTATUS_SYSTEM_ERROR
		return
	}
	testcaseOutput := &models.LabTestcaseOutputImitate{}
	err = json.Unmarshal([]byte(testcase.Output), testcaseOutput)
	if err != nil {
		log.Printf("deeplearning judger output unmarshal error[%#v]", err)
		testResult.Status = models.LABSUBMITSTATUS_SYSTEM_ERROR
		return
	}
	// 判断rate
	if stdoutFloat * 100 >= testcaseOutput.AcSimilarityRate {
		testResult.Status = models.LABSUBMITSTATUS_ACCEPTED
	} else {
		testResult.Status = models.LABSUBMITSTATUS_WRONG_ANSWER
	}
}