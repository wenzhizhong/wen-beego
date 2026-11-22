package systemtools

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"WenBeego/apps/admin_plat/models"
	"WenBeego/apps/admin_plat/models_ar"
	dto "WenBeego/apps/common/dto_vo"
	"WenBeego/apps/common/dto_vo/page_dto"
	"WenBeego/apps/common/dto_vo/sqlite_dms_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"

	sqlite "gosqlite.org"
	"gosqlite.org/vfs/crypto"

	"gorm.io/gorm"
)

type SqliteDmsService struct {
	SqliteDataSourceModel models.SqliteDataSource
	SqliteDataSourceAr    models_ar.SqliteDataSourceAr
}

const (
	sqliteSecretSeed   = "wen-beego-sqlite-dms-secret"
	maxExecuteRowLimit = 5000
)

var queryKeywords = map[string]bool{
	"SELECT":  true,
	"PRAGMA":  true,
	"EXPLAIN": true,
	"WITH":    true,
	"VALUES":  true,
}

func (s *SqliteDmsService) GetDataSourceList(reqDto page_dto.SqliteDataSourceListReqDto) (*dto.RespDataListDto, error) {
	data, count, err := s.SqliteDataSourceAr.GetList(reqDto)
	if err != nil {
		return nil, err
	}
	res := &dto.RespDataListDto{
		List:        data,
		Total:       count,
		PageSize:    reqDto.PageSize,
		CurrentPage: reqDto.CurrentPage,
	}
	return res, nil
}

func (s *SqliteDmsService) GetAllDataSources() ([]sqlite_dms_dto.DataSourceOptionVo, error) {
	list, err := s.SqliteDataSourceAr.GetAll()
	if err != nil {
		return nil, err
	}
	data := make([]sqlite_dms_dto.DataSourceOptionVo, 0, len(list))
	for _, item := range list {
		data = append(data, sqlite_dms_dto.DataSourceOptionVo{
			Id:          item.Id,
			Name:        item.Name,
			FilePath:    item.FilePath,
			Description: item.Description,
		})
	}
	return data, nil
}

func (s *SqliteDmsService) TestConnection(reqDto sqlite_dms_dto.TestConnectionDto) (map[string]interface{}, error) {
	if strings.TrimSpace(reqDto.FilePath) == "" {
		return nil, fmt.Errorf("数据文件路径不能为空")
	}
	encryptedPassword := ""
	if reqDto.Password != "" {
		value, err := encryptSecret(reqDto.Password)
		if err != nil {
			return nil, err
		}
		encryptedPassword = value
	} else if reqDto.Id != "" {
		oldData, err := s.SqliteDataSourceAr.GetById(reqDto.Id)
		if err == nil {
			encryptedPassword = oldData.Password
		}
	}
	if err := ensureParentDir(reqDto.FilePath); err != nil {
		return nil, err
	}
	db, err := openSqliteDatabase(reqDto.FilePath, encryptedPassword)
	if err != nil {
		return nil, friendlyOpenError(err, reqDto.Password != "")
	}
	defer db.Close()

	var tableCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table'`).Scan(&tableCount); err != nil {
		return nil, friendlyOpenError(err, reqDto.Password != "")
	}

	return map[string]interface{}{
		"success":     true,
		"encrypted":   encryptedPassword != "",
		"file_path":   reqDto.FilePath,
		"table_count": tableCount,
		"message":     "连接成功",
	}, nil
}

func (s *SqliteDmsService) SaveDataSource(baseParamDto dto.BaseParamDto, reqDto sqlite_dms_dto.SaveDataSourceDto) error {
	if strings.TrimSpace(reqDto.Name) == "" {
		return fmt.Errorf("数据源名称不能为空")
	}
	if strings.TrimSpace(reqDto.FilePath) == "" {
		return fmt.Errorf("数据文件路径不能为空")
	}

	now := time.Now()
	data := models.SqliteDataSource{
		Id:          reqDto.Id,
		Name:        strings.TrimSpace(reqDto.Name),
		FilePath:    strings.TrimSpace(reqDto.FilePath),
		Description: reqDto.Description,
		UpdateTime:  &now,
	}

	if err := ensureParentDir(data.FilePath); err != nil {
		return err
	}

	encryptedPassword, err := encryptSecret(reqDto.Password)
	if err != nil {
		return err
	}
	data.Password = encryptedPassword

	if reqDto.Id != "" {
		oldData, err := s.SqliteDataSourceAr.GetById(reqDto.Id)
		if err != nil {
			return fmt.Errorf("数据源不存在: %v", err)
		}
		if reqDto.Password == "" {
			data.Password = oldData.Password
		}
	} else {
		id, err := helper.GetUuid()
		if err != nil {
			return err
		}
		data.Id = id
		data.CreateTime = &now
	}

	if _, err := testOpen(data.FilePath, data.Password); err != nil {
		return err
	}

	return global.GetWriteDb().Transaction(func(tx *gorm.DB) error {
		if reqDto.Id != "" {
			return s.SqliteDataSourceAr.Update(tx, &data)
		}
		return s.SqliteDataSourceAr.Insert(tx, &data)
	})
}

func (s *SqliteDmsService) DeleteDataSources(ids []string) error {
	if len(ids) == 0 {
		return fmt.Errorf("ids不能为空")
	}
	return global.GetWriteDb().Transaction(func(tx *gorm.DB) error {
		return s.SqliteDataSourceAr.Delete(tx, ids)
	})
}

func (s *SqliteDmsService) GetTables(dataSourceId string) ([]string, error) {
	db, err := s.openByDataSourceId(dataSourceId)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT name FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%' ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := make([]string, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		data = append(data, name)
	}
	return data, rows.Err()
}

func (s *SqliteDmsService) GetTableSchema(dataSourceId, tableName string) ([]sqlite_dms_dto.TableColumnVo, error) {
	if strings.TrimSpace(tableName) == "" {
		return nil, fmt.Errorf("表名不能为空")
	}
	db, err := s.openByDataSourceId(dataSourceId)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	exists, err := tableExists(db.DB, tableName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("表不存在: %s", tableName)
	}

	quoted, err := quoteIdentifier(tableName)
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("PRAGMA table_info(" + quoted + ")")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := make([]sqlite_dms_dto.TableColumnVo, 0)
	for rows.Next() {
		var (
			cid        int
			name       string
			colType    string
			notNull    int
			defaultV   sql.NullString
			primaryKey int
		)
		if err := rows.Scan(&cid, &name, &colType, &notNull, &defaultV, &primaryKey); err != nil {
			return nil, err
		}
		item := sqlite_dms_dto.TableColumnVo{
			Cid:        cid,
			Name:       name,
			Type:       colType,
			NotNull:    notNull == 1,
			PrimaryKey: primaryKey > 0,
		}
		if defaultV.Valid {
			item.DefaultVal = defaultV.String
		}
		data = append(data, item)
	}
	return data, rows.Err()
}

func (s *SqliteDmsService) GetTableData(reqDto sqlite_dms_dto.TableDataDto) (*sqlite_dms_dto.TableDataVo, error) {
	if strings.TrimSpace(reqDto.TableName) == "" {
		return nil, fmt.Errorf("表名不能为空")
	}
	pageSize := reqDto.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 500 {
		pageSize = 500
	}
	currentPage := reqDto.CurrentPage
	if currentPage <= 0 {
		currentPage = 1
	}
	offset := (currentPage - 1) * pageSize

	db, err := s.openByDataSourceId(reqDto.DataSourceId)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	exists, err := tableExists(db.DB, reqDto.TableName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("表不存在: %s", reqDto.TableName)
	}

	quoted, err := quoteIdentifier(reqDto.TableName)
	if err != nil {
		return nil, err
	}

	var total int64
	if err := db.QueryRow("SELECT COUNT(*) FROM " + quoted).Scan(&total); err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT * FROM "+quoted+" LIMIT ? OFFSET ?", pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, rowData, err := scanRows(rows)
	if err != nil {
		return nil, err
	}

	return &sqlite_dms_dto.TableDataVo{
		Columns:     columns,
		Rows:        rowData,
		Total:       total,
		PageSize:    pageSize,
		CurrentPage: currentPage,
	}, nil
}

func (s *SqliteDmsService) ExecuteSql(reqDto sqlite_dms_dto.ExecuteSqlDto) (*sqlite_dms_dto.ExecuteSqlResultVo, error) {
	sqlText := strings.TrimSpace(reqDto.Sql)
	if sqlText == "" {
		return nil, fmt.Errorf("SQL 不能为空")
	}

	db, err := s.openByDataSourceId(reqDto.DataSourceId)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	start := time.Now()
	keyword := firstKeyword(sqlText)
	result := &sqlite_dms_dto.ExecuteSqlResultVo{
		Columns: make([]string, 0),
		Rows:    make([][]interface{}, 0),
	}

	if queryKeywords[keyword] {
		rows, err := db.Query(sqlText)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		columns, rowData, err := scanRowsLimit(rows, maxExecuteRowLimit)
		if err != nil {
			return nil, err
		}
		result.QueryType = "query"
		result.Columns = columns
		result.Rows = rowData
		result.Truncated = len(rowData) >= maxExecuteRowLimit
		result.ElapsedMs = time.Since(start).Milliseconds()
		return result, nil
	}

	execResult, err := db.Exec(sqlText)
	if err != nil {
		return nil, err
	}
	result.QueryType = "exec"
	result.RowsAffected, _ = execResult.RowsAffected()
	result.LastInsertId, _ = execResult.LastInsertId()
	result.ElapsedMs = time.Since(start).Milliseconds()
	return result, nil
}

func (s *SqliteDmsService) openByDataSourceId(dataSourceId string) (*sqlite.DB, error) {
	if strings.TrimSpace(dataSourceId) == "" {
		return nil, fmt.Errorf("数据源ID不能为空")
	}
	dataSource, err := s.SqliteDataSourceAr.GetById(dataSourceId)
	if err != nil {
		return nil, fmt.Errorf("数据源不存在: %v", err)
	}
	if strings.TrimSpace(dataSource.FilePath) == "" {
		return nil, fmt.Errorf("数据源未配置文件路径")
	}
	if err := ensureParentDir(dataSource.FilePath); err != nil {
		return nil, err
	}
	db, err := openSqliteDatabase(dataSource.FilePath, dataSource.Password)
	if err != nil {
		return nil, friendlyOpenError(err, dataSource.Password != "")
	}
	return db, nil
}

func testOpen(filePath, encryptedPassword string) (map[string]interface{}, error) {
	db, err := openSqliteDatabase(filePath, encryptedPassword)
	if err != nil {
		return nil, friendlyOpenError(err, encryptedPassword != "")
	}
	defer db.Close()

	var tableCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table'`).Scan(&tableCount); err != nil {
		return nil, friendlyOpenError(err, encryptedPassword != "")
	}
	return map[string]interface{}{
		"table_count": tableCount,
	}, nil
}

func openSqliteDatabase(filePath, encryptedPassword string) (*sqlite.DB, error) {
	cfg := sqlite.Config{
		Path: filePath,
		Mode: sqlite.ModeReadWriteCreate,
		Pragmas: sqlite.Pragmas{
			BusyTimeout: 5 * time.Second,
			ForeignKeys: true,
		},
	}
	if encryptedPassword == "" {
		return sqlite.Open(cfg)
	}

	password, err := decryptSecret(encryptedPassword)
	if err != nil {
		return nil, fmt.Errorf("解密数据源密码失败: %v", err)
	}
	// key, err := deriveCryptoKey(password, filePath)
	// if err != nil {
	// 	return nil, err
	// }
	return crypto.Open(cfg, crypto.Options{Key: []byte(password)})
}

func deriveCryptoKey(password, filePath string) ([]byte, error) {
	sum := sha256.Sum256([]byte(filepath.Clean(filePath)))
	salt := sum[:crypto.MinSaltLen]
	return crypto.DeriveKey([]byte(password), salt, crypto.Adiantum)
}

func encryptSecret(plain string) (string, error) {
	if plain == "" {
		return "", nil
	}
	block, err := aes.NewCipher(secretAesKey())
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	sealed := gcm.Seal(nonce, nonce, []byte(plain), nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

func decryptSecret(cipherText string) (string, error) {
	if cipherText == "" {
		return "", nil
	}
	raw, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(secretAesKey())
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(raw) < gcm.NonceSize() {
		return "", errors.New("密码密文格式错误")
	}
	nonce, data := raw[:gcm.NonceSize()], raw[gcm.NonceSize():]
	plain, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

func secretAesKey() []byte {
	sum := sha256.Sum256([]byte(sqliteSecretSeed))
	return sum[:]
}

func ensureParentDir(filePath string) error {
	parent := filepath.Dir(filePath)
	if parent == "." || parent == "" {
		return nil
	}
	return helper.MkdirAll(parent)
}

func friendlyOpenError(err error, encrypted bool) error {
	if err == nil {
		return nil
	}
	msg := strings.ToLower(err.Error())
	if encrypted && (strings.Contains(msg, "not a database") || strings.Contains(msg, "malformed") || strings.Contains(msg, "invalid")) {
		return fmt.Errorf("打开加密数据库失败，请检查密码是否正确: %v", err)
	}
	return err
}

func tableExists(db *sql.DB, tableName string) (bool, error) {
	var name string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type IN ('table', 'view') AND name = ?`, tableName).Scan(&name)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func quoteIdentifier(name string) (string, error) {
	if strings.TrimSpace(name) == "" {
		return "", fmt.Errorf("标识符不能为空")
	}
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`, nil
}

func firstKeyword(sqlText string) string {
	text := strings.TrimSpace(sqlText)
	for {
		if strings.HasPrefix(text, "--") {
			if idx := strings.IndexAny(text, "\r\n"); idx >= 0 {
				text = strings.TrimSpace(text[idx+1:])
				continue
			}
			return ""
		}
		if strings.HasPrefix(text, "/*") {
			if idx := strings.Index(text, "*/"); idx >= 0 {
				text = strings.TrimSpace(text[idx+2:])
				continue
			}
			return ""
		}
		break
	}
	fields := strings.Fields(text)
	if len(fields) == 0 {
		return ""
	}
	return strings.ToUpper(strings.Trim(fields[0], "("))
}

func scanRows(rows *sql.Rows) ([]string, [][]interface{}, error) {
	return scanRowsLimit(rows, 0)
}

func scanRowsLimit(rows *sql.Rows, limit int) ([]string, [][]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}
	data := make([][]interface{}, 0)
	values := make([]interface{}, len(columns))
	pointers := make([]interface{}, len(columns))
	for i := range values {
		pointers[i] = &values[i]
	}
	for rows.Next() {
		if err := rows.Scan(pointers...); err != nil {
			return nil, nil, err
		}
		row := make([]interface{}, len(columns))
		for i, value := range values {
			row[i] = normalizeValue(value)
		}
		data = append(data, row)
		if limit > 0 && len(data) >= limit {
			break
		}
	}
	return columns, data, rows.Err()
}

func normalizeValue(value interface{}) interface{} {
	switch v := value.(type) {
	case []byte:
		return string(v)
	case time.Time:
		return v.Format("2006-01-02 15:04:05")
	default:
		return v
	}
}
