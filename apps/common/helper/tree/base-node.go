package tree

// BaseNode 基础树节点实现
type BaseNode struct {
	ID       interface{} `json:"id"`
	ParentID interface{} `json:"parentId"`
	Children []TreeNode  `json:"children,omitempty"`
}

func (n *BaseNode) GetID() interface{} {
	return n.ID
}

func (n *BaseNode) GetParentID() interface{} {
	return n.ParentID
}

func (n *BaseNode) GetChildren() []TreeNode {
	return n.Children
}

func (n *BaseNode) SetChildren(children []TreeNode) {
	n.Children = children
}
