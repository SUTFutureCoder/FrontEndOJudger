package caroline

import (
	"FrontEndOJudger/models"
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
	var judgerReq httpJudgerReq
	var testResult TestResult
	decoder := json.NewDecoder(req.Body)
	decoder.Decode(&judgerReq)

	// check fields which must be filled
	if !checkReq(judgerReq) {
		net.Writer(w, net.INVALID_PARAMS, "please check your params", nil)
		return
	}

	// get lab data
	lab, err := models.GetLabFullInfo(judgerReq.LabId)
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

	_, testChamberUrlName, err := WriteSubmitToFile(&labSubmit)
	if err != nil {
		net.Writer(w, net.ERROR, err.Error(), nil)
		return
	}

	testChamber := fmt.Sprintf("%s:%s/%s", setting.JudgerSetting.TestChamberAddr, setting.JudgerSetting.TestChamberPort, testChamberUrlName)
	ExecTestCase(testChamber, judgerReq.LabTestcase, &testResult, &ctx)

	// write back result
	_, err = models.UpdateSubmitStatusResult(labSubmit.ID, labSubmit.Status, labSubmit.Status, testResult.SubmitOutput)
	if err != nil {
		net.Writer(w, net.ERROR, err.Error(), nil)
		return
	}
	net.Writer(w, net.SUCCESS, "", testResult)
}

func checkReq(req httpJudgerReq) bool {
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