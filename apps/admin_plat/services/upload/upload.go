package upload

import (
	"WenBeego/apps/common/dto_vo/upload_dto"
	commonUpload "WenBeego/apps/common/services/upload"
	"mime/multipart"
	"net/url"
)

type UploadService struct {
	commonUploadService commonUpload.Upload
}

func (s *UploadService) Upload(userId string, unitId string, file *multipart.File, fileHeader *multipart.FileHeader, postData url.Values, moduleName string) (upload_dto.UploadFileRespDto, error) {
	return s.commonUploadService.Upload(userId, unitId, file, fileHeader, postData, moduleName)
}

func (s *UploadService) VueSliceUpload(userId string, unitId string, file *multipart.File, fileHeader *multipart.FileHeader, postData url.Values, moduleName string) (upload_dto.UploadFileRespDto, error) {
	return s.commonUploadService.VueSliceUpload(userId, unitId, file, fileHeader, postData, moduleName)
}

func (s *UploadService) VueSliceUploadCheck(userId string, unitId string, fileMd5 string, sliceIndex string, sliceTotal string) (upload_dto.UploadFileRespDto, error) {
	return s.commonUploadService.VueSliceUploadCheck(userId, unitId, fileMd5, sliceIndex, sliceTotal)
}

func (s *UploadService) LinkSign(host, urls string) (interface{}, error) {
	return s.commonUploadService.LinkSign(host, urls)
}
func (s *UploadService) GetLinkById(host, ids string) (interface{}, error) {
	return s.commonUploadService.GetLinkById(host, ids)
}
