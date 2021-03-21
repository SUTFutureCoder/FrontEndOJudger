package models

import (
	"strconv"
)

// LabTestcaseMap 实验室、测试用例关联表
type LabTestcaseMap struct {
	Model
	// LabID 实验室id
	LabID uint64 `json:"lab_id"`
	// TestcaseID 测试用例id
	TestcaseID uint64 `json:"testcase_id"`
}

func (l *LabTestcaseMap) GetByLabId() ([]interface{}, error) {
	var testcaseIds []interface{}
	stmt, err := DB.Prepare("SELECT testcase_id FROM lab_testcase_map WHERE lab_id = ? AND status = 1")
	rows, err := stmt.Query(
		&l.LabID,
	)
	defer stmt.Close()
	for rows.Next() {
		var testcaseId int
		rows.Scan(&testcaseId)
		testcaseIds = append(testcaseIds, strconv.Itoa(testcaseId))
	}
	return testcaseIds, err
}
