package models

import "time"

// File 文件上传信息
type File struct {
	ID         string    `json:"id" gorm:"column:id;type:char(36);primaryKey"`
	BatchId    string    `json:"batchId" gorm:"column:batch_id;type:varchar(24)"`
	RealName   string    `json:"realName" gorm:"column:real_name;type:varchar(255)"`
	Name       string    `json:"name" gorm:"column:name;type:varchar(255)"`
	Suffix     string    `json:"suffix" gorm:"column:suffix;type:varchar(255)"`
	Path       string    `json:"path" gorm:"column:path;type:varchar(255)"`
	Type       string    `json:"type" gorm:"column:type;type:varchar(255)"`
	Size       string    `json:"size" gorm:"column:size;type:varchar(100)"`
	FileMd5    string    `json:"fileMd5" gorm:"column:file_md5;type:varchar(255)"`
	UnitId     string    `json:"unitId" gorm:"column:unit_id;type:varchar(255)"`
	UnitType   string    `json:"unitType" gorm:"column:unit_type;type:varchar(255)"`
	CreateBy   string    `json:"createBy" gorm:"column:create_by;type:varchar(255)"`
	UpdateBy   string    `json:"updateBy" gorm:"column:update_by;type:varchar(255)"`
	Createtime time.Time `json:"createtime" gorm:"column:createtime;type:timestamp"`
	Updatetime time.Time `json:"updatetime" gorm:"column:updatetime;type:timestamp"`
}

func (File) TableName() string {
	return "file"
}
