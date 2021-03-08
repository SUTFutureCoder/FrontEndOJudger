package models

import "log"

// Lab 实验室表
type Lab struct {
	Model
	// LabName 实验室名称
	LabName string `json:"lab_name"`
	// LabDesc 实验室描述
	LabDesc string `json:"lab_desc"`
	// LabType 实验室类型
	LabType int `json:"lab_type"`
	// LabSample 实验室样例或地址
	LabSample string `json:"lab_sample"`
	// LabTemplate 实验室模板代码
	LabTemplate string `json:"lab_template"`
}

const (
	LABTYPE_NORMAL = iota
	LABTYPE_IMITATE
)

func GetLabFullInfo(id uint64) (Lab, error) {
	var lab Lab
	stmt, err := DB.Prepare("SELECT id, lab_name, lab_desc, lab_type, lab_sample, lab_template, status, creator_id, creator, create_time, update_time FROM lab WHERE id=?")
	if err != nil {
		return lab, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(&id)
	err = row.Scan(
		&lab.ID, &lab.LabName, &lab.LabDesc, &lab.LabType, &lab.LabSample, &lab.LabTemplate, &lab.Status, &lab.CreatorId, &lab.Creator, &lab.CreateTime, &lab.UpdateTime,
	)
	if err != nil {
		log.Printf("get lab info error [%v]\n", err)
		return lab, err
	}
	return lab, err
}