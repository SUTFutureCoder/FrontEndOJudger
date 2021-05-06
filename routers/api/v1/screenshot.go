package v1

import (
	"FrontEndOJudger/caroline/judger"
	"FrontEndOJudger/models"
	"FrontEndOJudger/pkg/file"
	"FrontEndOJudger/pkg/net"
	"FrontEndOJudger/pkg/utils"
	"context"
	"encoding/json"
	"github.com/chromedp/chromedp"
	"net/http"
)

type httpScreenShotJudgerReq struct {
	LabId uint64 `json:"lab_id"`
	LabTestcase models.LabTestcase `json:"lab_testcase"`
}

func ScreenShot(w http.ResponseWriter, req *http.Request) {

	screenReq := &httpScreenShotJudgerReq{}
	testResult := &judger.TestResult{}
	decoder := json.NewDecoder(req.Body)
	decoder.Decode(screenReq)

	testcaseInput := &models.LabTestcaseInputImitate{}
	err := json.Unmarshal([]byte(screenReq.LabTestcase.Input), testcaseInput)
	if err != nil {
		net.Writer(w, net.ERROR, err.Error(), nil)
		return
	}

	// Start Chrome
	// Remove the 2nd param if you don't need debug information logged
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Run Tasks
	// List of actions to run in sequence (which also fills our image buffer)
	var imageBuf interface{}
	if err := chromedp.Run(ctx, utils.ScreenshotTasks(testcaseInput.TestcaseUrl, testcaseInput.TestcaseUrlCookies, &imageBuf, testcaseInput.TestcaseWidth, testcaseInput.TestcaseHeight, screenReq.LabTestcase.WaitBefore)); err != nil {
		net.Writer(w, net.ERROR, err.Error(), nil)
		return
	}

	// get path
	var imageBufByte []byte
	if o, ok := imageBuf.([]byte); ok {
		imageBufByte = o
	}
	filepath, reqpath := file.GetPathByBytes(screenReq.LabTestcase.CreatorId, imageBufByte)
	err = file.PutLocal(imageBufByte, filepath + ".png")
	if err != nil {
		net.Writer(w, net.ERROR, err.Error(), nil)
		return
	}

	testResult.SubmitOutput = reqpath  + ".png"
	net.Writer(w, net.SUCCESS, "", testResult)
}

