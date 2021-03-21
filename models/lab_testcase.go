package models

import (
	"strings"
)

// LabTestcase 测试用例表
type LabTestcase struct {
	Model
	// TestcaseDesc 测试用例描述
	TestcaseDesc string `json:"testcase_desc"`
	// TestcaseCode 测试用例代码
	TestcaseCode string `json:"testcase_code"`
	// Input 测试用例输入
	Input string `json:"input"`
	// Output 测试用例输出
	Output string `json:"output"`
	// TimeLimit 测试用例时间限制
	TimeLimit int `json:"time_limit"`
	// MemLimit 测试用例内存限制
	MemLimit int `json:"mem_limit"`
	// WaitBefore 测试用例执行前等待
	WaitBefore int `json:"wait_before"`
}

func (labTestcase *LabTestcase) GetByIds(testcaseIds []interface{}) ([]LabTestcase, error) {
	rows, err := DB.Query("SELECT id, testcase_code, testcase_desc, input, output, time_limit, mem_limit, wait_before, status, creator, create_time, update_time FROM lab_testcase WHERE id IN (?"+strings.Repeat(",?", len(testcaseIds)-1)+") AND status = 1", testcaseIds...)
	defer rows.Close()
	var testcases []LabTestcase
	for rows.Next() {
		var testcase LabTestcase
		err = rows.Scan(&testcase.ID, &testcase.TestcaseCode, &testcase.TestcaseDesc, &testcase.Input, &testcase.Output, &testcase.TimeLimit, &testcase.MemLimit, &testcase.WaitBefore, &testcase.Status, &testcase.Creator, &testcase.CreateTime, &testcase.UpdateTime)
		testcases = append(testcases, testcase)
	}
	rows.Close()
	return testcases, err

}
