package upload

import (
	"WenBeego/apps/common/dto/upload_dto"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models_ar"
	"WenBeego/apps/common/services/upload/itf"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type UploadToLocal struct {
}

var _ itf.UploadItf = (*UploadToLocal)(nil)

func (s *UploadToLocal) Upload(requestDto upload_dto.UploadFileReqDto, userId string, unitId string) (result upload_dto.UploadFileRespDto, err error) {
	var uploadPathStr string
	uploadPathStr, err = GetUploadpath(userId, requestDto.FilePath)
	if err != nil {
		return
	}

	err = os.WriteFile(uploadPathStr, *requestDto.FileData, 0755)
	if err != nil {
		return
	}
	result.HttpCode = 200
	result.FilePath = uploadPathStr
	return
}
func (s *UploadToLocal) SliceUpload(requestDto upload_dto.UploadFileReqDto, userId string, unitId string) (result upload_dto.UploadFileRespDto, err error) {

	_, err = GetUploadConfigPath()
	if err != nil {
		return
	}

	tempUploadDir, err1 := GetTempUploadDir(userId, requestDto.FileMd5, requestDto.FilePath)
	if err1 != nil {
		err = err1
		return
	}
	uploadSlicePathStr := filepath.Join(tempUploadDir, strconv.FormatInt(requestDto.SliceIndex, 10))
	err = helper.MkdirAll(tempUploadDir)
	if err != nil {
		return
	}

	err = os.WriteFile(uploadSlicePathStr, *requestDto.FileData, 0755)
	if err != nil {
		return
	}
	result.TempSlicePath = tempUploadDir
	return
}

func (s *UploadToLocal) SaveInfoToDB(requestDto upload_dto.UploadFileReqDto, uploadResult *upload_dto.UploadFileRespDto, userId string, unitId string, moduleName string) error {
	return SaveInfoToDB(requestDto, uploadResult, userId, unitId, moduleName)
}

func (s *UploadToLocal) SaveSliceToDB(requestDto upload_dto.UploadFileReqDto, uploadResult upload_dto.UploadFileRespDto, userId string) error {
	return SaveSliceToDB(requestDto, uploadResult, userId)
}

func (s *UploadToLocal) MergeSliceFile(requestDto upload_dto.UploadFileReqDto, uploadResult *upload_dto.UploadFileRespDto, userId string) (bool, error) {
	if !requestDto.IsSlice {
		return false, fmt.Errorf("当前文件%s不需要合并", requestDto.FileName)
	}
	if uploadResult.TempSlicePath == "" {
		return false, fmt.Errorf("当前文件%s没有返回临时文件目录", requestDto.FileName)
	}

	modelObj := &models_ar.FileSliceAr{}
	exists, err := modelObj.GetListByFileMd5(requestDto.FileMd5, userId)
	if err != nil {
		return false, err
	}
	if len(exists) < int(requestDto.SliceTotal) {
		return false, nil
	} else if len(exists) != int(requestDto.SliceTotal) {
		err = ClearTempUploadData(userId, requestDto.FileMd5, requestDto.FilePath)
		if err != nil {
			return false, err
		}
		return false, fmt.Errorf("当前文件%s的分片数量有误, 请重新上传", requestDto.FileName)
	}

	var uploadPathStr string
	uploadPathStr, err = GetUploadpath(userId, requestDto.FilePath)
	if err != nil {
		return false, err
	}

	// 读取目录下的所有分片文件
	data := make([]byte, 0)
	for _, v := range exists {
		slicePath := filepath.Join(uploadResult.TempSlicePath, strconv.FormatInt(v.SliceIndex, 10))
		tmpData, err1 := os.ReadFile(slicePath)
		if err1 != nil {
			return false, err1
		}
		data = append(data, tmpData...)
	}

	// 创建文件
	err = os.WriteFile(uploadPathStr, data, 0755)
	if err != nil {
		return false, err
	}

	err = ClearTempUploadData(userId, requestDto.FileMd5, requestDto.FilePath)
	if err != nil {
		return false, err
	}

	uploadResult.FilePath = uploadPathStr
	uploadResult.HttpCode = 200
	uploadResult.TempSlicePath = ""
	uploadResult.Uploaded = make([]int64, 0)

	return true, nil
}
