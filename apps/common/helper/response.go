package helper

import dto "WenBeego/apps/common/dto_vo"

func Response(code int, message string, data interface{}) dto.Response {
	if data == nil {
		data = new(interface{})
	}
	return dto.Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// 设置并返回页面列表数据
func GetRespDataListDto(pageSize int, currentPage int, total int64, list interface{}) (dto.RespDataListDto, error) {
	data := dto.RespDataListDto{}
	data.Total = total
	data.List = list
	data.PageSize = pageSize
	data.CurrentPage = currentPage
	return data, nil
}
