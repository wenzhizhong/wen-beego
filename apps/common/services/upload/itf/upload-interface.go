package itf

import (
	"WenBeego/apps/common/dto_vo/upload_dto"
)

type UploadItf interface {
	Upload(requestDto upload_dto.UploadFileReqDto, userId string, unitId string) (upload_dto.UploadFileRespDto, error)
	SliceUpload(requestDto upload_dto.UploadFileReqDto, userId string, unitId string) (upload_dto.UploadFileRespDto, error)
	MergeSliceFile(requestDto upload_dto.UploadFileReqDto, uploadResult *upload_dto.UploadFileRespDto, userId string) (bool, error)
	SaveInfoToDB(requestDto upload_dto.UploadFileReqDto, uploadResult *upload_dto.UploadFileRespDto, userId string, unitId string, moduleName string) error
	SaveSliceToDB(requestDto upload_dto.UploadFileReqDto, uploadResult upload_dto.UploadFileRespDto, userId string) error
}
