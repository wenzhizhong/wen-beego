package upload

import (
	"WenBeego/apps/common/dto_vo/upload_dto"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models_ar"
	"os"
	"path/filepath"
	"strconv"
)

func GetTempUploadDir(userId string, fileMd5 string, originFilepath string) (string, error) {
	originFileDir := filepath.Dir(originFilepath)
	if len(originFileDir) <= 1 {
		originFileDir = ""
	}

	path := filepath.Join(os.TempDir(), "upload", userId, fileMd5, originFileDir)
	err := helper.MkdirAll(path)
	return path, err
}

func ClearTempUploadData(userId string, fileMd5 string, originFilepath string) error {
	tempUploadDir, err := GetTempUploadDir(userId, fileMd5, originFilepath)
	if err != nil {
		return err
	}
	os.RemoveAll(tempUploadDir)
	return (&models_ar.FileSliceAr{}).Delete(fileMd5)
}

func SaveInfoToDB(requestDto upload_dto.UploadFileReqDto, uploadResult *upload_dto.UploadFileRespDto, userId string, unitId string, moduleName string) error {
	tmpUniqueName := helper.Md5(uploadResult.FilePath)
	exists, err := (&models_ar.FileAr{}).GetByName(tmpUniqueName)
	if err == nil && exists.ID != "" {
		uploadResult.FileId = exists.ID
		return nil
	}

	return (&models_ar.FileAr{}).Insert(&models.File{
		ID:         uploadResult.FileId,
		BatchId:    requestDto.BatchId,
		FileMd5:    requestDto.FileMd5,
		Name:       tmpUniqueName,
		Path:       uploadResult.FilePath,
		RealName:   requestDto.FileName,
		Size:       strconv.FormatInt(requestDto.TotalSize, 10),
		Suffix:     requestDto.FileExt,
		Type:       requestDto.FileMime,
		UnitId:     unitId,
		UnitType:   moduleName,
		CreateBy:   userId,
		CreateTime: helper.GetTime(),
		UpdateTime: helper.GetTime(),
	})
}

func SaveSliceToDB(requestDto upload_dto.UploadFileReqDto, uploadResult upload_dto.UploadFileRespDto, userId string) error {
	uuid, err := helper.GetUuid()
	if err != nil {
		return err
	}

	modelObj := &models_ar.FileSliceAr{}
	exists, err := modelObj.GetListByFileMd5(requestDto.FileMd5, userId, requestDto.SliceIndex)
	if err != nil {
		return err
	}
	if len(exists) > 0 {
		return nil
	}

	return modelObj.Insert(&models.FileSlice{
		ID:         uuid,
		FileName:   strconv.FormatInt(requestDto.SliceIndex, 10),
		FileMd5:    requestDto.FileMd5,
		FilePath:   uploadResult.TempSlicePath,
		CreateBy:   userId,
		CreateTime: helper.GetTime(),
		UpdateTime: helper.GetTime(),
		SliceIndex: requestDto.SliceIndex,
		SliceTotal: requestDto.SliceTotal,
		SliceSize:  requestDto.FileSize,
		TotalSize:  requestDto.TotalSize,
	})
}
