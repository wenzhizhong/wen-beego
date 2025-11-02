package upload

import (
	"WenBeego/apps/common/dto/upload_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models_ar"
	"WenBeego/apps/common/services/upload/itf"
	"errors"
	"io"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
)

type Upload struct {
}

func (s *Upload) Upload(userId string, unitId string, file *multipart.File, fileHeader *multipart.FileHeader, postData url.Values, moduleName string) (upload_dto.UploadFileRespDto, error) {
	requestDto, err := s.getUploadReqDto(file, fileHeader, postData)
	if err != nil {
		return upload_dto.UploadFileRespDto{}, err
	}

	return s.doUpload(requestDto, userId, unitId, moduleName)
}

func (s *Upload) VueSliceUpload(userId string, unitId string, file *multipart.File, fileHeader *multipart.FileHeader, postData url.Values, moduleName string) (upload_dto.UploadFileRespDto, error) {
	requestDto, err := s.getUploadReqDto(file, fileHeader, postData)
	if err != nil {
		return upload_dto.UploadFileRespDto{}, err
	}
	return s.doUpload(requestDto, userId, unitId, moduleName)
}
func (s *Upload) VueSliceUploadCheck(userId string, unitId string, fileMd5 string, sliceIndex string, sliceTotal string) (checkDto upload_dto.UploadFileRespDto, err error) {
	checkDto.HttpCode = 200
	sliceTotalInt, err1 := strconv.ParseInt(sliceTotal, 10, 64)
	uploaded, err2 := s.getSliceUploaded(fileMd5, userId)
	if err1 != nil {
		err = err1
		checkDto.HttpCode = 210
		return
	}
	if err2 != nil {
		err = err2
		checkDto.HttpCode = 210
		return
	}
	checkDto.Uploaded = uploaded
	checkDto.SkipUpload = sliceTotalInt == int64(len(uploaded))

	return
}

func (s *Upload) getUploadReqDto(file *multipart.File, fileHeader *multipart.FileHeader, postData url.Values) (upload_dto.UploadFileReqDto, error) {

	data := upload_dto.UploadFileReqDto{}
	data.FileName = fileHeader.Filename
	data.FileSize = fileHeader.Size
	data.FileMime = fileHeader.Header.Get("Content-Type")
	data.FileExt = filepath.Ext(fileHeader.Filename)

	data.FileMd5 = postData.Get("identifier")
	data.FilePath = postData.Get("relativePath")

	sliceTotal, err := strconv.ParseInt(postData.Get("totalChunks"), 10, 64)
	sliceIndex, err2 := strconv.ParseInt(postData.Get("chunkNumber"), 10, 64)
	totalSize, err3 := strconv.ParseInt(postData.Get("totalSize"), 10, 64)
	currentSize, err4 := strconv.ParseInt(postData.Get("currentChunkSize"), 10, 64)
	chunkSize, err5 := strconv.ParseInt(postData.Get("chunkSize"), 10, 64)
	if err != nil {
		return data, err2
	}
	if err2 != nil {
		return data, err2
	}
	if err3 != nil {
		return data, err3
	}
	if err4 != nil {
		return data, err4
	}
	if err5 != nil {
		return data, err5
	}

	data.IsSlice = sliceTotal > 1
	data.BatchId = postData.Get("batchId")
	data.BucketACL = postData.Get("bucketACL")
	data.SliceTotal = sliceTotal
	data.SliceIndex = sliceIndex
	data.CurrentSize = currentSize
	data.SliceSize = chunkSize
	data.TotalSize = totalSize

	fileData := make([]byte, fileHeader.Size)
	_, err = io.ReadFull(*file, fileData) // 使用 io.ReadFull 将文件内容读入缓冲区
	if err != nil {
		return data, err // 处理错误
	}
	data.FileData = &fileData
	return data, nil
}

func (s *Upload) doUpload(requestDto upload_dto.UploadFileReqDto, userId string, unitId string, moduleName string) (result upload_dto.UploadFileRespDto, err error) {

	uploadType, err := global.GetConfigDiy("upload.type")
	uploadTypes, err2 := global.GetConfigDiy("upload.types")

	if err != nil || err2 != nil || uploadType == nil || uploadTypes == nil {
		return result, errors.New("上传配置错误")
	}
	if res, err := helper.InArray(uploadType.(string), uploadTypes); !res || err != nil {
		return result, errors.New("不支持上传类型：" + uploadType.(string))
	}

	var uploadObj itf.UploadItf
	uploadTypeStr := uploadType.(string)
	switch uploadTypeStr {
	case "local":
		uploadObj = &UploadToLocal{}
	case "minio":
		// TODO: minio
		err = errors.New("minio 上传暂不支持")
	case "aliyun":
		// TODO: aliyun
		err = errors.New("aliyun 上传暂不支持")
	}
	if err != nil {
		return result, err
	}

	// upload file
	mergeRes := false
	if requestDto.IsSlice {
		result, err = uploadObj.SliceUpload(requestDto, userId, unitId)
		if err != nil {
			return result, err
		}

		err = uploadObj.SaveSliceToDB(requestDto, result, userId)
		if err != nil {
			return result, err
		}

		mergeRes, err = uploadObj.MergeSliceFile(requestDto, &result, userId)
	} else {
		mergeRes = true
		result, err = uploadObj.Upload(requestDto, userId, unitId)
	}
	if err != nil {
		return result, err
	}

	if mergeRes {
		// save info to db
		result.FileId, err = helper.GetUuid()
		if err != nil {
			return result, err
		}
		result.FilePath = strings.ReplaceAll(result.FilePath, global.RootPath, "")
		result.FilePath = strings.ReplaceAll(result.FilePath, "\\", "/")
		err = uploadObj.SaveInfoToDB(requestDto, &result, userId, unitId, moduleName)
	}
	if err == nil {
		result.HttpCode = 200
	}
	return result, err
}

// get slice uploaded
func (s *Upload) getSliceUploaded(fileMd5 string, userId string) ([]int64, error) {
	data, err := (&models_ar.FileSliceAr{}).GetListByFileMd5(fileMd5, userId)

	uploaded := make([]int64, 0)
	if err != nil {
		return uploaded, err
	}
	dataLen := len(data)
	if dataLen > 0 {
		for _, v := range data {
			uploaded = append(uploaded, v.SliceIndex)
		}
	}
	return uploaded, nil
}
