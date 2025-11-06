package tree

// TreeNode 通用树节点接口
type TreeNode interface {
	GetID() interface{}
	GetParentID() interface{}
	GetChildren() []TreeNode
	SetChildren(children []TreeNode)
}

// 便捷函数

// BuildTree 快速构建树形结构（使用TreeNode接口）
func BuildTree(nodes []TreeNode) ([]TreeNode, error) {
	return NewTreeConverter().ConvertToTree(nodes)
}

// BuildTreeWithRootID 使用指定的根节点父ID构建树
func BuildTreeWithRootID(nodes []TreeNode, rootParentID interface{}) ([]TreeNode, error) {
	return NewTreeConverter().WithRootParentID(rootParentID).ConvertToTree(nodes)
}

// NewTreeConverter 创建树形转换器
func NewTreeConverter() *TreeConverter {
	return &TreeConverter{
		rootParentID: nil,
		keepOrder:    true,
	}
}
