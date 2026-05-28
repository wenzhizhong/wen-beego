package tests

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	googleUuid "github.com/google/uuid"
)

func TestMigrate(t *testing.T) {
	// filterViewNameFromDdl()
	filterTableNameFromDdl()
	// replaceDmlUuid()
}
func filterViewNameFromDdl() {
	sqlStr, err := os.ReadFile("doc/数据库表结构.sql")
	if err != nil {
		panic(err)
	}

	modeRegexp := regexp.MustCompile(`CREATE OR REPLACE VIEW\s*"?[a-zA-Z0-9_-]+"?\."?[a-zA-Z0-9_-]+"?`)
	tmpViewArr := modeRegexp.FindAllString(string(sqlStr), -1)
	viewNames := []string{}
	for _, tmpView := range tmpViewArr {
		tmpView = strings.ReplaceAll(tmpView, "CREATE OR REPLACE VIEW", "")
		tmpView = strings.Trim(tmpView, " \r\n\t")
		viewNames = append(viewNames, tmpView)
	}
	fmt.Println(viewNames)
}
func filterTableNameFromDdl() {
	sqlStr, err := os.ReadFile("doc/数据库表结构.sql")
	if err != nil {
		panic(err)
	}
	tableNameRegexp := regexp.MustCompile(`CREATE\s+TABLE\s+"?[a-zA-Z0-9_-]+"?\."?[a-zA-Z0-9_-]+"?\s+\(`)
	tmpNameArr := tableNameRegexp.FindAllString(string(sqlStr), -1)

	replaceStrSlice := []string{"CREATE", "TABLE", strings.ToLower("CREATE"), strings.ToLower("TABLE"), "("}
	nameMap := map[string]int{}
	for _, tmpName := range tmpNameArr {
		for _, replaceStr := range replaceStrSlice {
			tmpName = strings.ReplaceAll(tmpName, replaceStr, "")
		}
		tmpName = strings.Trim(tmpName, " \r\n\t")
		nameMap[tmpName] = 1
	}
	fmt.Println("nameMap=", nameMap)
}
func replaceDmlUuid() {
	sqlStr, err := os.ReadFile("doc/数据库表基本数据.sql")
	if err != nil {
		panic(err)
	}
	oldStr := string(sqlStr)
	uuidMaps := make(map[string]string)
	uuidRegexp := regexp.MustCompile(`([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})`)
	newStr := uuidRegexp.ReplaceAllStringFunc(oldStr, func(s string) string {
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
	fmt.Println("\n\n\nres=", newStr)
	fmt.Println("\n\n\noldStr=", oldStr)
	fmt.Println("==============")
}
