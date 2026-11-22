package sqlite_dms_dto

type TestConnectionDto struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	FilePath string `json:"file_path"`
	Password string `json:"password"`
}

type SaveDataSourceDto struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	FilePath    string `json:"file_path"`
	Password    string `json:"password"`
	Description string `json:"description"`
}

type DeleteDataSourceDto struct {
	Ids []string `json:"ids"`
}

type DataSourceIdDto struct {
	DataSourceId string `json:"dataSourceId"`
}

type TableNameDto struct {
	DataSourceId string `json:"dataSourceId"`
	TableName    string `json:"tableName"`
}

type TableDataDto struct {
	DataSourceId string `json:"dataSourceId"`
	TableName    string `json:"tableName"`
	PageSize     int    `json:"pageSize"`
	CurrentPage  int    `json:"currentPage"`
}

type ExecuteSqlDto struct {
	DataSourceId string `json:"dataSourceId"`
	Sql          string `json:"sql"`
}

type DataSourceOptionVo struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	FilePath    string `json:"file_path"`
	Description string `json:"description"`
}

type TableColumnVo struct {
	Cid        int         `json:"cid"`
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	NotNull    bool        `json:"notnull"`
	DefaultVal interface{} `json:"dflt_value"`
	PrimaryKey bool        `json:"pk"`
}

type TableDataVo struct {
	Columns     []string        `json:"columns"`
	Rows        [][]interface{} `json:"rows"`
	Total       int64           `json:"total"`
	PageSize    int             `json:"pageSize"`
	CurrentPage int             `json:"currentPage"`
}

type ExecuteSqlResultVo struct {
	QueryType    string          `json:"queryType"`
	Columns      []string        `json:"columns"`
	Rows         [][]interface{} `json:"rows"`
	RowsAffected int64           `json:"rowsAffected"`
	LastInsertId int64           `json:"lastInsertId"`
	ElapsedMs    int64           `json:"elapsedMs"`
	Truncated    bool            `json:"truncated"`
}
