package upload_dto

type UploadFileReqDto struct {
	FileName string
	FileSize int64
	FilePath string
	FileData *[]byte
	FileMd5  string
	FileExt  string
	FileMime string

	IsSlice     bool
	BatchId     string
	SliceTotal  int64
	SliceIndex  int64
	CurrentSize int64
	SliceSize   int64
	TotalSize   int64
	BucketACL   string // public | private
}
type UploadFileRespDto struct {
	FileId        string `json:"fileId"`
	FilePath      string `json:"filePath"`
	TempSlicePath string `json:"tempSlicePath"`
	UploadFileCheckDto
}

type UploadFileCheckDto struct {
	HttpCode   int     `json:"httpCode"`
	SkipUpload bool    `json:"skipUpload"`
	Uploaded   []int64 `json:"uploaded"`
}
