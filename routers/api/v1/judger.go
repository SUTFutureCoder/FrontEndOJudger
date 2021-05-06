package v1

import (
	"FrontEndOJudger/caroline"
	"FrontEndOJudger/caroline/judger"
	"FrontEndOJudger/models"
	"FrontEndOJudger/pkg/file"
	"FrontEndOJudger/pkg/net"
	"FrontEndOJudger/pkg/setting"
	"context"
	"encoding/json"
	"fmt"
	"github.com/chromedp/chromedp"
	"net/http"
	"time"
)

type httpJudgerReq struct {
	LabId uint64 `json:"lab_id"`
	LabTestcase models.LabTestcase `json:"lab_testcase"`
}

func HttpJudger(w http.ResponseWriter, req *http.Request) {
	judgerReq := &httpJudgerReq{}
	testResult := &judger.TestResult{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(judgerReq)

	// check fields which must be filled
	if err != nil || !checkReq(judgerReq) {
		net.Writer(w, net.INVALID_PARAMS, "please check your params", nil)
		return
	}

	// get lab data
	lab := &models.Lab{}
	err = lab.GetFullInfo(judgerReq.LabId)
	if err != nil {
		net.Writer(w, net.ERROR, err.Error(), nil)
		return
	}
	if lab.Status != models.STATUS_CONSTRUCTING {
		net.Writer(w, net.INVALID_PARAMS, "lab status is not STATUS_CONSTRUCTING", nil)
		return
	}

	// submit data
	labSubmit, err := submit(judgerReq.LabId, lab.LabSample)
	if err != nil {
		net.Writer(w, net.ERROR, err.Error(), nil)
		return
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	_, destFile, err := file.GetDest(labSubmit.ID)
	if err != nil {
		return
	}
	err = file.PutLocal([]byte(labSubmit.SubmitData), destFile)
	if err != nil {
		net.Writer(w, net.ERROR, err.Error(), nil)
		return
	}

	testChamber := fmt.Sprintf("%s:%s/%d", setting.JudgerSetting.TestChamberAddr, setting.JudgerSetting.TestChamberPort, labSubmit.ID)
	caroline.ExecTestCase(testChamber, judgerReq.LabTestcase, testResult, lab, &ctx)

	// write back result
	_, err = labSubmit.UpdateStatusResult(labSubmit.Status, labSubmit.Status, testResult.SubmitOutput)
	if err != nil {
		net.Writer(w, net.ERROR, err.Error(), nil)
		return
	}
	net.Writer(w, net.SUCCESS, "", testResult)
}

func checkReq(req *httpJudgerReq) bool {
	if req.LabId == 0 || req.LabTestcase.TestcaseCode == "" {
		return false
	}
	return true
}

func submit(labId uint64, labTemplate string) (models.LabSubmit, error) {
	var labSubmit models.LabSubmit
	labSubmit.LabID = labId
	labSubmit.SubmitData = labTemplate
	labSubmit.Status = models.LABSUBMITSTATUS_TEST
	labSubmit.CreateTime = time.Now().UnixNano() / 1e6
	id, err := labSubmit.Insert()
	labSubmit.ID = uint64(id)
	return labSubmit, err
}