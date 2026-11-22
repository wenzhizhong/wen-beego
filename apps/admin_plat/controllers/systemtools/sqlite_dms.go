package systemtools

import (
	"WenBeego/apps/admin_plat/services/systemtools"
	commonControllers "WenBeego/apps/common/controller"
	"WenBeego/apps/common/dto_vo/page_dto"
	"WenBeego/apps/common/dto_vo/sqlite_dms_dto"
	"WenBeego/apps/common/global/constant"
	"WenBeego/apps/common/helper"
)

type SqliteDmsController struct {
	commonControllers.AdminBaseController
	SqliteDmsService systemtools.SqliteDmsService
}

// @Summary 获取SQLite数据源列表
// @Description 获取SQLite数据源列表
// @Tags 系统工具-SQLite DMS
// @Accept application/json
// @Produce application/json
// @Param keyword query string false "名称/路径关键字"
// @Success 200 {object} dto.RespDataListDto
// @Router /admin_plat/system-tools/sqlite-dms/data-source/list [get]
func (c *SqliteDmsController) GetDataSourceList() {
	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	reqDataListDto, err1 := helper.GetReqDataListDto(&c.Controller)
	if err != nil || err1 != nil {
		err = helper.Ternary(err != nil, err, err1)
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	reqDto := page_dto.SqliteDataSourceListReqDto{
		BaseParamDto:   baseParamDto,
		ReqDataListDto: reqDataListDto,
		Keyword:        c.GetString("keyword"),
	}
	data, err := c.SqliteDmsService.GetDataSourceList(reqDto)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
	} else {
		c.Data["json"] = helper.Response(200, "success", data)
	}
	c.ServeJSON()
}

// @Summary 获取全部SQLite数据源
// @Description 获取全部SQLite数据源（下拉选择）
// @Tags 系统工具-SQLite DMS
// @Accept application/json
// @Produce application/json
// @Success 200 {object} dto.Response
// @Router /admin_plat/system-tools/sqlite-dms/data-source/all [get]
func (c *SqliteDmsController) GetAllDataSources() {
	data, err := c.SqliteDmsService.GetAllDataSources()
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
	} else {
		c.Data["json"] = helper.Response(200, "success", data)
	}
	c.ServeJSON()
}

// @Summary 测试SQLite连接
// @Description 测试SQLite数据源连接
// @Tags 系统工具-SQLite DMS
// @Accept application/json
// @Produce application/json
// @Param data body sqlite_dms_dto.TestConnectionDto true "连接参数"
// @Success 200 {object} dto.Response
// @Router /admin_plat/system-tools/sqlite-dms/data-source/test [post]
func (c *SqliteDmsController) TestConnection() {
	reqDto, err := helper.GetReqBody[sqlite_dms_dto.TestConnectionDto](c.Ctx)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	data, err := c.SqliteDmsService.TestConnection(reqDto)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
	} else {
		c.Data["json"] = helper.Response(200, "连接成功", data)
	}
	c.ServeJSON()
}

// @Summary 保存SQLite数据源
// @Description 保存SQLite数据源
// @Tags 系统工具-SQLite DMS
// @Accept application/json
// @Produce application/json
// @Param data body sqlite_dms_dto.SaveDataSourceDto true "数据源参数"
// @Success 200 {object} dto.Response
// @Router /admin_plat/system-tools/sqlite-dms/data-source/save [post]
func (c *SqliteDmsController) SaveDataSource() {
	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	reqDto, err1 := helper.GetReqBody[sqlite_dms_dto.SaveDataSourceDto](c.Ctx)
	if err != nil || err1 != nil {
		err = helper.Ternary(err != nil, err, err1)
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	err = c.SqliteDmsService.SaveDataSource(baseParamDto, reqDto)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
	} else {
		c.Data["json"] = helper.Response(200, "success", nil)
	}
	c.ServeJSON()
}

// @Summary 删除SQLite数据源
// @Description 删除SQLite数据源
// @Tags 系统工具-SQLite DMS
// @Accept application/json
// @Produce application/json
// @Param data body sqlite_dms_dto.DeleteDataSourceDto true "数据源ID集合"
// @Success 200 {object} dto.Response
// @Router /admin_plat/system-tools/sqlite-dms/data-source/del [post]
func (c *SqliteDmsController) DeleteDataSource() {
	reqDto, err := helper.GetReqBody[sqlite_dms_dto.DeleteDataSourceDto](c.Ctx)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	err = c.SqliteDmsService.DeleteDataSources(reqDto.Ids)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
	} else {
		c.Data["json"] = helper.Response(200, "success", nil)
	}
	c.ServeJSON()
}

// @Summary 获取SQLite数据表
// @Description 获取SQLite数据源下的数据表
// @Tags 系统工具-SQLite DMS
// @Accept application/json
// @Produce application/json
// @Param data body sqlite_dms_dto.DataSourceIdDto true "数据源ID"
// @Success 200 {object} dto.Response
// @Router /admin_plat/system-tools/sqlite-dms/table/list [post]
func (c *SqliteDmsController) GetTables() {
	reqDto, err := helper.GetReqBody[sqlite_dms_dto.DataSourceIdDto](c.Ctx)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	data, err := c.SqliteDmsService.GetTables(reqDto.DataSourceId)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
	} else {
		c.Data["json"] = helper.Response(200, "success", data)
	}
	c.ServeJSON()
}

// @Summary 获取SQLite表结构
// @Description 获取SQLite表结构
// @Tags 系统工具-SQLite DMS
// @Accept application/json
// @Produce application/json
// @Param data body sqlite_dms_dto.TableNameDto true "数据源ID与表名"
// @Success 200 {object} dto.Response
// @Router /admin_plat/system-tools/sqlite-dms/table/schema [post]
func (c *SqliteDmsController) GetTableSchema() {
	reqDto, err := helper.GetReqBody[sqlite_dms_dto.TableNameDto](c.Ctx)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	data, err := c.SqliteDmsService.GetTableSchema(reqDto.DataSourceId, reqDto.TableName)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
	} else {
		c.Data["json"] = helper.Response(200, "success", data)
	}
	c.ServeJSON()
}

// @Summary 获取SQLite表数据
// @Description 分页获取SQLite表数据
// @Tags 系统工具-SQLite DMS
// @Accept application/json
// @Produce application/json
// @Param data body sqlite_dms_dto.TableDataDto true "数据源ID、表名与分页参数"
// @Success 200 {object} dto.Response
// @Router /admin_plat/system-tools/sqlite-dms/table/data [post]
func (c *SqliteDmsController) GetTableData() {
	reqDto, err := helper.GetReqBody[sqlite_dms_dto.TableDataDto](c.Ctx)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	data, err := c.SqliteDmsService.GetTableData(reqDto)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
	} else {
		c.Data["json"] = helper.Response(200, "success", data)
	}
	c.ServeJSON()
}

// @Summary 执行SQLite SQL
// @Description 执行DML/DDL，查询语句返回结果集
// @Tags 系统工具-SQLite DMS
// @Accept application/json
// @Produce application/json
// @Param data body sqlite_dms_dto.ExecuteSqlDto true "数据源ID与SQL"
// @Success 200 {object} dto.Response
// @Router /admin_plat/system-tools/sqlite-dms/execute [post]
func (c *SqliteDmsController) ExecuteSql() {
	reqDto, err := helper.GetReqBody[sqlite_dms_dto.ExecuteSqlDto](c.Ctx)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	data, err := c.SqliteDmsService.ExecuteSql(reqDto)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
	} else {
		c.Data["json"] = helper.Response(200, "success", data)
	}
	c.ServeJSON()
}
