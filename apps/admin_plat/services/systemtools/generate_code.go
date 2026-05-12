package systemtools

import (
	"WenBeego/apps/admin_plat/models_ar"
	"WenBeego/apps/common/dto"
	"WenBeego/apps/common/dto/generate_code_dto"
	"WenBeego/apps/common/dto/page_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
	"archive/zip"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"gorm.io/gorm"
)

type GenerateCodeService struct {
	GenerateCodeModel models.GenerateCode
	GenerateCodeAr    models_ar.GenerateCodeAr
}

var CODE_TYPE_MODEL = "model"
var CODE_TYPE_AR = "ar"
var CODE_TYPE_SERVICE = "service"
var CODE_TYPE_CONTROLLER = "controller"
var CODE_TYPE_VIEW = "view"
var CODE_TYPE_ALL = "all"

var VIEW_TYPE_ELEMENT_PLUS = "element-plus"

var CODE_TYPES_MAP = map[string]string{
	CODE_TYPE_MODEL:      "模型",
	CODE_TYPE_AR:         "数据访问层",
	CODE_TYPE_SERVICE:    "服务层",
	CODE_TYPE_CONTROLLER: "控制器",
	CODE_TYPE_VIEW:       "视图",
	CODE_TYPE_ALL:        "全部",
}

var VIEW_TYPES_MAP = map[string]string{
	VIEW_TYPE_ELEMENT_PLUS: "Element-Plus",
}

func (s *GenerateCodeService) GetDbTableList() (interface{}, error) {
	data, err := s.GenerateCodeAr.GetAllDbTables()
	if err != nil {
		return nil, err
	}
	return data, nil
}

type TableDetailData struct {
	FullName string                  `json:"fullName"`
	Comment  string                  `json:"comment"`
	Columns  []models_ar.DbColumnInfo `json:"columns"`
}

func (s *GenerateCodeService) GetDbTableDetail(tableName string) (interface{}, error) {
	if tableName == "" {
		return nil, fmt.Errorf("表名不能为空")
	}
	columns, err := s.GenerateCodeAr.GetDbTableColumns(tableName)
	if err != nil {
		return nil, err
	}
	comment, err := s.GenerateCodeAr.GetTableComment(tableName)
	if err != nil {
		comment = ""
	}
	data := TableDetailData{
		FullName: tableName,
		Comment:  comment,
		Columns:  columns,
	}
	return data, nil
}

func (s *GenerateCodeService) SaveFormDetail(baseParamDto dto.BaseParamDto, data generate_code_dto.SaveFormDetailDto) error {
	if data.TableName == "" {
		return fmt.Errorf("表名不能为空")
	}
	genDto := generate_code_dto.GenerateCodeDto{}
	genDto.TableName = data.TableName

	if data.Data != "" {
		genDto.Data = data.Data
	}

	err := global.GetWriteDb().Transaction(func(tx *gorm.DB) error {
		if data.Id != "" {
			genDto.Id = data.Id
			return s.GenerateCodeAr.Update(tx, genDto)
		}
		id, err := helper.GetUuid()
		if err != nil {
			return err
		}
		genDto.Id = id
		return s.GenerateCodeAr.Insert(tx, genDto)
	})
	return err
}

func (s *GenerateCodeService) DelFormDetail(ids []string) error {
	if len(ids) == 0 {
		return fmt.Errorf("ids不能为空")
	}
	err := global.GetWriteDb().Transaction(func(tx *gorm.DB) error {
		return s.GenerateCodeAr.Delete(tx, ids)
	})
	return err
}

func (s *GenerateCodeService) GetGenerateCodeList(reqDto page_dto.GenerateCodeListReqDto) (*dto.RespDataListDto, error) {
	data, count, err := s.GenerateCodeAr.GetList(reqDto)
	if err != nil {
		return nil, err
	}
	res := &dto.RespDataListDto{}
	res.List = data
	res.Total = count
	res.PageSize = reqDto.PageSize
	res.CurrentPage = reqDto.CurrentPage
	return res, nil
}

func (s *GenerateCodeService) GetGenerateCodeParam() (interface{}, error) {
	data := generate_code_dto.GenCodeParamDto{
		ViewTypes:         VIEW_TYPES_MAP,
		GenerateCodeTypes: CODE_TYPES_MAP,
		MenuName:          "",
	}
	return data, nil
}

func (s *GenerateCodeService) DownloadCode(zipPath string) (string, error) {
	if zipPath == "" {
		return "", fmt.Errorf("zipPath不能为空")
	}
	fullPath := filepath.Join(global.UploadsDir, "public", "code", zipPath)
	if !helper.PathIsExist(fullPath) {
		return "", fmt.Errorf("文件不存在: %s", zipPath)
	}
	return zipPath, nil
}

type ColumnConfig struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Required   bool   `json:"required"`
	DefVal     interface{} `json:"defVal"`
	FormType   string `json:"formType"`
	FormParam  string `json:"formParam"`
	Comment    string `json:"comment"`
	GoFieldName string `json:"-"`
	GoType      string `json:"-"`
	PgType      string `json:"-"`
	JsonName    string `json:"-"`
}

type TemplateData struct {
	ModelName      string
	ModelNameLower string
	TableName      string
	AppModule      string
	MenuModule     string
	BizModule      string
	MenuName       string
	ApiUrlPrefix   string
	ApiNamePrefix  string
	Columns        []ColumnConfig
	HasDeleted     bool
	HasCreateTime  bool
	IsMultiApp     bool
	IsMultiTenant  bool
	PlatModelName  string
	MchntModelName string
	PlatTableName  string
	MchntTableName string
}

func (s *GenerateCodeService) GenerateCode(reqDto generate_code_dto.GenCodeRunDto) (map[string]string, error) {
	if reqDto.TableGenerateCodeId == "" {
		return nil, fmt.Errorf("tableGenerateCodeId不能为空")
	}
	appModules := reqDto.AppModules
	if len(appModules) == 0 {
		appModules = []string{"admin_plat"}
	}
	if reqDto.MenuModule == "" {
		reqDto.MenuModule = "system"
	}

	genCodeDetail, err := s.GenerateCodeAr.GetById(reqDto.TableGenerateCodeId)
	if err != nil {
		return nil, fmt.Errorf("获取生成代码配置失败: %v", err)
	}

	var columnConfigs []ColumnConfig
	if genCodeDetail.Data != "" {
		if err := json.Unmarshal([]byte(genCodeDetail.Data), &columnConfigs); err != nil {
			return nil, fmt.Errorf("解析字段配置失败: %v", err)
		}
	}

	for i := range columnConfigs {
		columnConfigs[i].GoFieldName = snakeToPascal(columnConfigs[i].Name)
		columnConfigs[i].JsonName = columnConfigs[i].Name
		goType, pgType := getGoPgType(columnConfigs[i].Type)
		columnConfigs[i].GoType = goType
		columnConfigs[i].PgType = pgType
	}

	codeTypes := reqDto.CodeType
	if len(codeTypes) == 0 || (len(codeTypes) == 1 && codeTypes[0] == "") {
		return nil, fmt.Errorf("codeType不能为空")
	}

	if contains(codeTypes, CODE_TYPE_ALL) {
		codeTypes = []string{CODE_TYPE_MODEL, CODE_TYPE_AR, CODE_TYPE_SERVICE, CODE_TYPE_CONTROLLER, CODE_TYPE_VIEW}
	}

	viewType := reqDto.ViewType
	if viewType == "" {
		viewType = VIEW_TYPE_ELEMENT_PLUS
	}

	bizModule := reqDto.BizModule
	if bizModule == "" {
		bizModule = snakeToPascal(genCodeDetail.TableName)
	}
	tableName := genCodeDetail.TableName
	menuModule := reqDto.MenuModule

	isMultiApp := len(appModules) > 1
	isMultiTenant := strings.HasPrefix(tableName, "plat_") || strings.HasPrefix(tableName, "mchnt_")

	modelName := bizModule
	modelNameLower := snakeToCamel(modelName)
	platModelName := "Plat" + modelName
	mchntModelName := "Mchnt" + modelName
	platTableName := tableName
	mchntTableName := tableName
	if isMultiTenant && strings.HasPrefix(tableName, "plat_") {
		mchntTableName = "mchnt_" + strings.TrimPrefix(tableName, "plat_")
	} else if isMultiTenant && strings.HasPrefix(tableName, "mchnt_") {
		platTableName = "plat_" + strings.TrimPrefix(tableName, "mchnt_")
	}

	timestamp := time.Now().Format("20060102150405")
	tempDir := filepath.Join(global.TempDir, "code", modelNameLower)
	zipDir := filepath.Join(global.UploadsDir, "public", "code")
	zipName := modelNameLower + "-" + timestamp + ".zip"
	zipPath := filepath.Join(zipDir, zipName)

	if err := os.RemoveAll(tempDir); err != nil {
		return nil, fmt.Errorf("清理临时目录失败: %v", err)
	}
	if err := helper.MkdirAll(tempDir); err != nil {
		return nil, fmt.Errorf("创建临时目录失败: %v", err)
	}
	if err := helper.MkdirAll(zipDir); err != nil {
		return nil, fmt.Errorf("创建zip目录失败: %v", err)
	}

	tplDir := filepath.Join(global.AppDir, "common", "codeTpl")

	td := TemplateData{
		ModelName:      modelName,
		ModelNameLower: modelNameLower,
		TableName:      tableName,
		MenuModule:     menuModule,
		BizModule:      bizModule,
		MenuName:       reqDto.MenuName,
		Columns:        columnConfigs,
		HasDeleted:     hasColumn(columnConfigs, "deleted"),
		HasCreateTime:  hasColumn(columnConfigs, "create_time") || hasColumn(columnConfigs, "created_at"),
		IsMultiApp:     isMultiApp,
		IsMultiTenant:  isMultiTenant,
		PlatModelName:  platModelName,
		MchntModelName: mchntModelName,
		PlatTableName:  platTableName,
		MchntTableName: mchntTableName,
	}

	// shared code: base_model, model, dto (always generated once)
	if contains(codeTypes, CODE_TYPE_MODEL) {
		if err := s.renderTemplate(tplDir, "admin_plat", "base_model.tpl", tempDir, "apps/common/models/base_model/"+modelNameLower+".go", td); err != nil {
			return nil, err
		}
		if err := s.renderTemplate(tplDir, "admin_plat", "model.tpl", tempDir, "apps/common/models/"+modelNameLower+".go", td); err != nil {
			return nil, err
		}
		if err := s.renderTemplate(tplDir, "admin_plat", "dto.tpl", tempDir, "apps/common/dto/"+menuModule+"_dto/"+modelNameLower+".go", td); err != nil {
			return nil, err
		}
		// multi-tenant: additionally generate Plat/Mchnt entity models
		if isMultiTenant {
			if err := s.renderTemplate(tplDir, "admin_plat", "plat_model.tpl", tempDir, "apps/common/models/plat_"+modelNameLower+".go", td); err != nil {
				return nil, err
			}
			if err := s.renderTemplate(tplDir, "admin_plat", "mchnt_model.tpl", tempDir, "apps/common/models/mchnt_"+modelNameLower+".go", td); err != nil {
				return nil, err
			}
		}
	}

	if isMultiApp {
		// multi-app: generate common shared service + AR, then per-app wrappers
		if contains(codeTypes, CODE_TYPE_AR) {
			// always generate common AR struct
			if err := s.renderTemplate(tplDir, "admin_plat", "common_ar.tpl", tempDir, "apps/common/models_ar/"+modelNameLower+".go", td); err != nil {
				return nil, err
			}
			// multi-tenant: additionally generate generic functions in base_ar
			if isMultiTenant {
				if err := s.renderTemplate(tplDir, "admin_plat", "base_ar.tpl", tempDir, "apps/common/models_ar/base_ar/common_"+modelNameLower+".go", td); err != nil {
					return nil, err
				}
			}
		}
		if contains(codeTypes, CODE_TYPE_SERVICE) {
			if err := s.renderTemplate(tplDir, "admin_plat", "common_service.tpl", tempDir, "apps/common/services/"+menuModule+"/"+modelNameLower+".go", td); err != nil {
				return nil, err
			}
		}

		for _, appModule := range appModules {
			td.AppModule = appModule
			td.ApiUrlPrefix = "/" + appModule + "/" + menuModule + "-" + strings.ReplaceAll(tableName, "_", "-")
			td.ApiNamePrefix = strings.ToUpper(menuModule) + "_" + strings.ToUpper(strings.ReplaceAll(tableName, "_", "_"))

			if contains(codeTypes, CODE_TYPE_AR) {
				if err := s.renderTemplate(tplDir, "admin_plat", "ar.tpl", tempDir, "apps/"+appModule+"/models_ar/"+modelNameLower+"_ar.go", td); err != nil {
					return nil, err
				}
			}
			if contains(codeTypes, CODE_TYPE_SERVICE) {
				if err := s.renderTemplate(tplDir, "admin_plat", "service.tpl", tempDir, "apps/"+appModule+"/services/"+menuModule+"/"+modelNameLower+".go", td); err != nil {
					return nil, err
				}
			}
			if contains(codeTypes, CODE_TYPE_CONTROLLER) {
				if err := s.renderTemplate(tplDir, "admin_plat", "controller.tpl", tempDir, "apps/"+appModule+"/controllers/"+menuModule+"/"+modelNameLower+".go", td); err != nil {
					return nil, err
				}
			}
		}
	} else {
		// single-app: generate everything for one app
		appModule := appModules[0]
		td.AppModule = appModule
		td.ApiUrlPrefix = "/" + appModule + "/" + menuModule + "-" + strings.ReplaceAll(tableName, "_", "-")
		td.ApiNamePrefix = strings.ToUpper(menuModule) + "_" + strings.ToUpper(strings.ReplaceAll(tableName, "_", "_"))

		for _, codeType := range codeTypes {
			switch codeType {
			case CODE_TYPE_AR:
				if err := s.renderTemplate(tplDir, "admin_plat", "ar.tpl", tempDir, "apps/"+appModule+"/models_ar/"+modelNameLower+"_ar.go", td); err != nil {
					return nil, err
				}
			case CODE_TYPE_SERVICE:
				if err := s.renderTemplate(tplDir, "admin_plat", "service.tpl", tempDir, "apps/"+appModule+"/services/"+menuModule+"/"+modelNameLower+".go", td); err != nil {
					return nil, err
				}
			case CODE_TYPE_CONTROLLER:
				if err := s.renderTemplate(tplDir, "admin_plat", "controller.tpl", tempDir, "apps/"+appModule+"/controllers/"+menuModule+"/"+modelNameLower+".go", td); err != nil {
					return nil, err
				}
			}
		}
	}

	// frontend code
	if contains(codeTypes, CODE_TYPE_VIEW) && viewType == VIEW_TYPE_ELEMENT_PLUS {
		if err := s.renderTemplate(tplDir, "web/element-plus", "index.vue.tpl", tempDir, "web/"+modelNameLower+"/index.vue", td); err != nil {
			return nil, err
		}
		if err := s.renderTemplate(tplDir, "web/element-plus", "form.vue.tpl", tempDir, "web/"+modelNameLower+"/form.vue", td); err != nil {
			return nil, err
		}
		if err := s.renderTemplate(tplDir, "web/element-plus", "hook.tsx.tpl", tempDir, "web/"+modelNameLower+"/hook.tsx", td); err != nil {
			return nil, err
		}
		if err := s.renderTemplate(tplDir, "web/element-plus", "api.ts.tpl", tempDir, "web/"+modelNameLower+"/api.ts", td); err != nil {
			return nil, err
		}
	}

	ddlContent := s.generateDDL(genCodeDetail.TableName, columnConfigs)
	ddlPath := filepath.Join(tempDir, "sql", tableName+".sql")
	if err := os.MkdirAll(filepath.Dir(ddlPath), os.ModePerm); err != nil {
		return nil, err
	}
	if err := os.WriteFile(ddlPath, []byte(ddlContent), 0644); err != nil {
		return nil, err
	}

	dmlContent := s.generateMenuDML(genCodeDetail.TableName, reqDto.MenuName, td.ApiUrlPrefix)
	dmlPath := filepath.Join(tempDir, "sql", tableName+"_menu.sql")
	if err := os.WriteFile(dmlPath, []byte(dmlContent), 0644); err != nil {
		return nil, err
	}

	if err := s.createZip(tempDir, zipPath); err != nil {
		return nil, fmt.Errorf("创建zip失败: %v", err)
	}

	if err := os.RemoveAll(tempDir); err != nil {
		global.Log.Error("清理临时目录失败: %v", err)
	}

	result := map[string]string{
		"zipPath": modelNameLower + "-" + timestamp + ".zip",
		"zipName": modelNameLower + "-" + timestamp + ".zip",
	}
	return result, nil
}

func (s *GenerateCodeService) renderTemplate(tplDir, subDir, tplFile, tempDir, outFile string, td interface{}) error {
	tplPath := filepath.Join(tplDir, subDir, tplFile)
	if !helper.PathIsExist(tplPath) {
		return fmt.Errorf("模板文件不存在: %s", tplPath)
	}

	tmpl, err := template.ParseFiles(tplPath)
	if err != nil {
		return fmt.Errorf("解析模板失败 %s: %v", tplFile, err)
	}

	outPath := filepath.Join(tempDir, outFile)
	if err := os.MkdirAll(filepath.Dir(outPath), os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败 %s: %v", outFile, err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, td); err != nil {
		return fmt.Errorf("渲染模板失败 %s: %v", tplFile, err)
	}

	return nil
}

func (s *GenerateCodeService) generateDDL(tableName string, columns []ColumnConfig) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("-- DDL for table: %s\n", tableName))
	sb.WriteString("-- Generated by wen-beego code generator\n\n")
	sb.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", tableName))
	for i, col := range columns {
		comma := ","
		if i == len(columns)-1 {
			comma = ""
		}
		nullable := "NOT NULL"
		if !col.Required {
			nullable = ""
		}
		sb.WriteString(fmt.Sprintf("    %s %s %s%s\n", col.Name, col.PgType, nullable, comma))
	}
	sb.WriteString(");\n")
	return sb.String()
}

func (s *GenerateCodeService) generateMenuDML(tableName, menuName, apiUrlPrefix string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("-- Menu DML for: %s (%s)\n", tableName, menuName))
	sb.WriteString("-- Generated by wen-beego code generator\n\n")
	sb.WriteString(fmt.Sprintf("-- INSERT INTO plat_menu (id, name, path, ...) VALUES (...);\n"))
	sb.WriteString(fmt.Sprintf("-- API prefix: %s\n", apiUrlPrefix))
	sb.WriteString(fmt.Sprintf("-- Menu name: %s\n", menuName))
	return sb.String()
}

func (s *GenerateCodeService) createZip(srcDir, zipPath string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		writer, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		_, err = writer.Write(content)
		return err
	})
	return err
}

func snakeToPascal(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, "")
}

func snakeToCamel(s string) string {
	pascal := snakeToPascal(s)
	if len(pascal) == 0 {
		return ""
	}
	return strings.ToLower(pascal[:1]) + pascal[1:]
}

func getGoPgType(dbType string) (goType, pgType string) {
	switch strings.ToLower(dbType) {
	case "bpchar", "char":
		return "string", "bpchar(36)"
	case "varchar":
		return "string", "varchar(255)"
	case "text":
		return "string", "text"
	case "int2":
		return "int", "int2"
	case "int4", "integer":
		return "int", "int4"
	case "int8", "bigint":
		return "int64", "int8"
	case "float4":
		return "float32", "float4"
	case "float8":
		return "float64", "float8"
	case "bool":
		return "bool", "bool"
	case "timestamptz", "timestamp":
		return "time.Time", "timestamptz(6)"
	case "date":
		return "time.Time", "date"
	case "json", "jsonb":
		return "string", "jsonb"
	case "uuid":
		return "string", "uuid"
	default:
		return "string", "text"
	}
}

func hasColumn(columns []ColumnConfig, name string) bool {
	for _, col := range columns {
		if col.Name == name {
			return true
		}
	}
	return false
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
