package tests

import (
	"WenBeego/apps/common/helper/tree"
	"encoding/json"
	"fmt"
	"testing"
)

// 自定义节点结构
type Department struct {
	tree.BaseNode
	Name string `json:"name"`
	Code string `json:"code"`
}

type User struct {
	ID       int     `json:"id"`
	ParentID int     `json:"parentId"`
	Name     string  `json:"name"`
	Children []*User `json:"children,omitempty"`
}

func init() {
}

func TestTree(t *testing.T) {
	departments := []tree.TreeNode{
		&Department{
			BaseNode: tree.BaseNode{ID: 1, ParentID: 0},
			Name:     "总公司",
			Code:     "001",
		},
		&Department{
			BaseNode: tree.BaseNode{ID: 2, ParentID: 1},
			Name:     "技术部",
			Code:     "002",
		},
		&Department{
			BaseNode: tree.BaseNode{ID: 3, ParentID: 1},
			Name:     "市场部",
			Code:     "003",
		},
		&Department{
			BaseNode: tree.BaseNode{ID: 4, ParentID: 2},
			Name:     "前端组",
			Code:     "004",
		},
	}

	users := []*User{
		{ID: 1, ParentID: 0, Name: "管理员"},
		{ID: 2, ParentID: 1, Name: "用户组1"},
		{ID: 3, ParentID: 1, Name: "用户组2"},
		{ID: 4, ParentID: 2, Name: "用户A"},
		{ID: 5, ParentID: 2, Name: "用户B"},
	}

	test1(departments)
	test2(users)
	test3(departments)

	fmt.Println("=== end ===")
}
func test1(departments []tree.TreeNode) {
	// 示例1：使用TreeNode接口
	fmt.Println("=== 示例1：使用TreeNode接口 ===")

	treeData, err := tree.BuildTree(departments)
	if err != nil {
		panic(err)
	}

	jsonData, _ := json.MarshalIndent(treeData, "", "  ")
	fmt.Println(string(jsonData))
}

func test2(users []*User) {
	// 示例2：使用反射方式（不需要实现接口）
	fmt.Println("\n=== 示例2：使用反射方式 ===")

	treeData2, err := tree.ConvertSliceToTree(users, "ID", "ParentID", "Children")
	if err != nil {
		panic(err)
	}

	jsonData2, _ := json.MarshalIndent(treeData2, "", "  ")
	fmt.Println(string(jsonData2))

}
func test3(departments []tree.TreeNode) {

	// 示例3：使用配置选项
	fmt.Println("\n=== 示例3：使用配置选项 ===")

	converter := tree.NewTreeConverter().
		WithRootParentID(0).
		WithKeepOrder(true)

	treeData3, err := converter.ConvertToTree(departments)
	if err != nil {
		panic(err)
	}

	jsonData3, _ := json.MarshalIndent(treeData3, "", "  ")
	fmt.Println(string(jsonData3))
}
