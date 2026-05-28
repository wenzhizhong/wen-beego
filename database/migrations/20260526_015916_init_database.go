package main

import (
	"os"
	"regexp"
	"strings"

	"github.com/beego/beego/v2/client/orm/migration"
	googleUuid "github.com/google/uuid"
)

// DO NOT MODIFY
type InitDatabase_20260526_015916 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &InitDatabase_20260526_015916{}
	m.Created = "20260526_015916"

	migration.Register("InitDatabase_20260526_015916", m)
}

// Run the migrations
func (m *InitDatabase_20260526_015916) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	ddl, err := os.ReadFile("../../doc/数据库表结构.sql")
	dml, err2 := os.ReadFile("../../doc/数据库表基本数据.sql")
	if err != nil {
		panic(err)
	}
	if err2 != nil {
		panic(err2)
	}
	if ddl == nil || dml == nil {
		panic("ddl or dml is nil")
	}

	newDml := m.replaceDmlUuid(string(dml))
	executeSql := string(ddl) + "\n;\n"
	executeSql += newDml
	m.SQL(executeSql)

}

// Reverse the migrations
func (m *InitDatabase_20260526_015916) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	ddl, err := os.ReadFile("../../doc/数据库表结构.sql")
	if err != nil {
		panic(err)
	}
	if ddl != nil {
		viewNames := m.filterViewNameFromDdl(string(ddl))
		tableNames := m.filterTableNameFromDdl(string(ddl))

		executeSql := ""
		for _, viewName := range viewNames {
			executeSql += "DROP VIEW IF EXISTS " + viewName + ";\n"
		}
		for tableName := range tableNames {
			executeSql += "DROP TABLE IF EXISTS " + tableName + ";\n"
		}
		m.SQL(executeSql)
	}

}

func (m *InitDatabase_20260526_015916) filterViewNameFromDdl(dml string) []string {
	modeRegexp := regexp.MustCompile(`CREATE OR REPLACE VIEW\s*"?[a-zA-Z0-9_-]+"?\."?[a-zA-Z0-9_-]+"?`)
	tmpViewArr := modeRegexp.FindAllString(string(dml), -1)
	viewNames := []string{}
	for _, tmpView := range tmpViewArr {
		tmpView = strings.ReplaceAll(tmpView, "CREATE OR REPLACE VIEW", "")
		tmpView = strings.Trim(tmpView, " \r\n\t")
		viewNames = append(viewNames, tmpView)
	}
	return viewNames
}
func (m *InitDatabase_20260526_015916) filterTableNameFromDdl(dml string) map[string]int {
	tableNameRegexp := regexp.MustCompile(`CREATE\s+TABLE\s+"?[a-zA-Z0-9_-]+"?\."?[a-zA-Z0-9_-]+"?\s+\(`)
	tmpNameArr := tableNameRegexp.FindAllString(dml, -1)

	replaceStrSlice := []string{"CREATE", "TABLE", strings.ToLower("CREATE"), strings.ToLower("TABLE"), "("}
	nameMap := map[string]int{}
	for _, tmpName := range tmpNameArr {
		for _, replaceStr := range replaceStrSlice {
			tmpName = strings.ReplaceAll(tmpName, replaceStr, "")
		}
		tmpName = strings.Trim(tmpName, " \r\n\t")
		nameMap[tmpName] = 1
	}
	return nameMap
}
func (m *InitDatabase_20260526_015916) replaceDmlUuid(oldStr string) string {
	if oldStr == "" {
		return ""
	}
	uuidMaps := make(map[string]string)
	uuidRegexp := regexp.MustCompile(`([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})`)
	newSql := uuidRegexp.ReplaceAllStringFunc(oldStr, func(s string) string {
		if newID, ok := uuidMaps[s]; ok {
			return newID
		}
		uuid, err := googleUuid.NewV7()
		if err != nil {
			panic(err)
		}
		newID := uuid.String()
		uuidMaps[s] = newID
		return newID
	})
	return newSql
}
