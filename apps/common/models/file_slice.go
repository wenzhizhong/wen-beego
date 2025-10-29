package models

import "time"

// FileSlice 文件上传分片信息
type FileSlice struct {
	ID          string    `json:"id" gorm:"column:id;type:char(36);primaryKey"`
	FileName    string    `json:"fileName" gorm:"column:file_name;type:varchar(255)"`
	CurrentSize int64     `json:"currentSize" gorm:"column:current_size;type:bigint"`
	SliceSize   int64     `json:"sliceSize" gorm:"column:slice_size;type:bigint"`
	TotalSize   int64     `json:"totalSize" gorm:"column:total_size;type:bigint"`
	SliceIndex  int64     `json:"sliceIndex" gorm:"column:slice_index;type:bigint"`
	SliceTotal  int64     `json:"sliceTotal" gorm:"column:slice_total;type:bigint"`
	FileMd5     string    `json:"fileMd5" gorm:"column:file_md5;type:varchar(45)"`
	FilePath    string    `json:"filePath" gorm:"column:file_path;type:varchar(255)"`
	CreateBy    string    `json:"createBy" gorm:"column:create_by;type:varchar(255)"`
	Createtime  time.Time `json:"createtime" gorm:"column:createtime;type:timestamp"`
	Updatetime  time.Time `json:"updatetime" gorm:"column:updatetime;type:timestamp"`
}

func (FileSlice) TableName() string {
	return "file_slice"
}
