package models

import (
	"log"
)

// Lab 实验室表
type Lab struct {
	Model
	// LabName 实验室名称
	LabName string `json:"lab_name"`
	// LabDesc 实验室描述
	LabDesc string `json:"lab_desc"`
	// LabType 实验室类型
	LabType int8 `json:"lab_type"`
	// LabSample 实验室样例或地址
	LabSample string `json:"lab_sample"`
}

const (
	LABTYPE_HTML = iota
	LABTYPE_CSS
	LABTYPE_JS
	LABTYPE_VUE
	LABTYPE_COMPLEX
	LABTYPE_PRD
	LABTYPE_IMITATE
	LABTYPE_OTHER
)

func (lab *Lab) Insert() error {
	stmt, err := DB.Prepare("INSERT INTO lab (lab_name, lab_desc, lab_type, lab_sample, creator, create_time) VALUES(?,?,?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		log.Printf("[ERROR] database exec error input[%v] err[%v]", lab, err)
		return err
	}
	_, err = stmt.Exec(
		lab.LabName,
		lab.LabDesc,
		lab.LabType,
		lab.LabSample,
		lab.Creator,
		lab.CreateTime,
	)
	return nil
}
