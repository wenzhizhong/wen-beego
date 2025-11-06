package tree

import (
	"fmt"
	"reflect"
)

// TreeConverter 树形转换器
type TreeConverter struct {
	// 根节点的父ID值，默认为nil/0/""等
	rootParentID interface{}
	// 是否保持原始顺序
	keepOrder bool
}

// WithRootParentID 设置根节点的父ID值
func (tc *TreeConverter) WithRootParentID(rootParentID interface{}) *TreeConverter {
	tc.rootParentID = rootParentID
	return tc
}

// WithKeepOrder 设置是否保持原始顺序
func (tc *TreeConverter) WithKeepOrder(keepOrder bool) *TreeConverter {
	tc.keepOrder = keepOrder
	return tc
}

// ConvertToTree 将数组转换为树形结构
func (tc *TreeConverter) ConvertToTree(nodes []TreeNode) ([]TreeNode, error) {
	if len(nodes) == 0 {
		return []TreeNode{}, nil
	}

	// 创建ID到节点的映射
	nodeMap := make(map[interface{}]TreeNode)
	// 创建父ID到子节点列表的映射
	childrenMap := make(map[interface{}][]TreeNode)

	// 第一次遍历：建立映射关系
	for _, node := range nodes {
		nodeID := node.GetID()
		parentID := node.GetParentID()

		// 检查ID是否重复
		if _, exists := nodeMap[nodeID]; exists {
			return nil, fmt.Errorf("duplicate node ID: %v", nodeID)
		}

		nodeMap[nodeID] = node

		// 将节点添加到父节点的子节点列表中
		childrenMap[parentID] = append(childrenMap[parentID], node)
	}

	// 第二次遍历：构建树形结构
	var roots []TreeNode
	processed := make(map[interface{}]bool)

	for _, node := range nodes {
		if tc.isRootNode(node) {
			roots = append(roots, node)
		}
		tc.buildTree(node, childrenMap, processed)
	}

	// 如果指定了保持顺序，按照原始顺序排序根节点
	if tc.keepOrder {
		roots = tc.sortRootsByOriginalOrder(roots, nodes)
	}

	return roots, nil
}

// isRootNode 判断是否为根节点
func (tc *TreeConverter) isRootNode(node TreeNode) bool {
	parentID := node.GetParentID()

	// 如果设置了特定的根节点父ID
	if tc.rootParentID != nil {
		return reflect.DeepEqual(parentID, tc.rootParentID)
	}

	// 默认情况：父ID为nil、0、""等被认为是根节点
	return parentID == nil || parentID == 0 || parentID == "" ||
		(reflect.ValueOf(parentID).Kind() == reflect.String && parentID == "")
}

// buildTree 递归构建树
func (tc *TreeConverter) buildTree(node TreeNode, childrenMap map[interface{}][]TreeNode, processed map[interface{}]bool) {
	nodeID := node.GetID()

	// 如果已经处理过，直接返回
	if processed[nodeID] {
		return
	}

	processed[nodeID] = true

	// 获取子节点
	children := childrenMap[nodeID]

	// 递归处理子节点
	for _, child := range children {
		tc.buildTree(child, childrenMap, processed)
	}

	// 设置子节点
	if tc.keepOrder {
		children = tc.sortChildrenByOriginalOrder(children, childrenMap[nodeID])
	}
	node.SetChildren(children)
}

// sortRootsByOriginalOrder 按照原始顺序排序根节点
func (tc *TreeConverter) sortRootsByOriginalOrder(roots []TreeNode, originalNodes []TreeNode) []TreeNode {
	if len(roots) <= 1 {
		return roots
	}

	positionMap := make(map[interface{}]int)
	for i, node := range originalNodes {
		positionMap[node.GetID()] = i
	}

	// 使用稳定的排序算法
	for i := 0; i < len(roots)-1; i++ {
		for j := 0; j < len(roots)-i-1; j++ {
			pos1 := positionMap[roots[j].GetID()]
			pos2 := positionMap[roots[j+1].GetID()]
			if pos1 > pos2 {
				roots[j], roots[j+1] = roots[j+1], roots[j]
			}
		}
	}

	return roots
}

// sortChildrenByOriginalOrder 按照原始顺序排序子节点
func (tc *TreeConverter) sortChildrenByOriginalOrder(children []TreeNode, originalOrder []TreeNode) []TreeNode {
	if len(children) <= 1 {
		return children
	}

	positionMap := make(map[interface{}]int)
	for i, node := range originalOrder {
		positionMap[node.GetID()] = i
	}

	// 使用稳定的排序算法
	for i := 0; i < len(children)-1; i++ {
		for j := 0; j < len(children)-i-1; j++ {
			pos1 := positionMap[children[j].GetID()]
			pos2 := positionMap[children[j+1].GetID()]
			if pos1 > pos2 {
				children[j], children[j+1] = children[j+1], children[j]
			}
		}
	}

	return children
}

// ConvertSliceToTree 通用切片转树形结构（使用反射，不需要实现TreeNode接口）
func ConvertSliceToTree(slice interface{}, idField, parentIdField, childrenField string) (interface{}, error) {
	sliceValue := reflect.ValueOf(slice)
	if sliceValue.Kind() != reflect.Slice {
		return nil, fmt.Errorf("input must be a slice")
	}

	length := sliceValue.Len()
	if length == 0 {
		return []interface{}{}, nil
	}

	// // 获取切片元素类型
	// elemType := sliceValue.Type().Elem()
	// isPtr := elemType.Kind() == reflect.Ptr
	// var valueType reflect.Type
	// if isPtr {
	// 	valueType = elemType.Elem()
	// } else {
	// 	valueType = elemType
	// }

	// 创建映射
	nodeMap := make(map[interface{}]reflect.Value)
	childrenMap := make(map[interface{}][]reflect.Value)

	// 第一次遍历：建立映射关系
	for i := 0; i < length; i++ {
		item := sliceValue.Index(i)

		// 获取实际的值（如果是指针，则获取指向的值）
		var itemValue reflect.Value
		if item.Kind() == reflect.Ptr {
			if item.IsNil() {
				continue // 跳过nil指针
			}
			itemValue = item.Elem()
		} else {
			itemValue = item
		}

		id := getFieldValue(itemValue, idField)
		parentId := getFieldValue(itemValue, parentIdField)

		nodeMap[id.Interface()] = item // 存储原始项（可能是指针或值）

		childrenMap[parentId.Interface()] = append(childrenMap[parentId.Interface()], item)
	}

	// 构建树形结构
	var roots []reflect.Value
	processed := make(map[interface{}]bool)

	for i := 0; i < length; i++ {
		item := sliceValue.Index(i)

		// 获取实际的值
		var itemValue reflect.Value
		if item.Kind() == reflect.Ptr {
			if item.IsNil() {
				continue
			}
			itemValue = item.Elem()
		} else {
			itemValue = item
		}

		parentId := getFieldValue(itemValue, parentIdField)

		// 判断是否为根节点
		if isZeroValue(parentId) {
			roots = append(roots, item)
		}

		// 递归构建树
		buildTreeWithReflection(item, childrenMap, processed, childrenField)
	}

	// 转换回切片
	result := reflect.MakeSlice(sliceValue.Type(), len(roots), len(roots))
	for i, root := range roots {
		result.Index(i).Set(root)
	}

	return result.Interface(), nil
}

// buildTreeWithReflection 递归构建树（反射版本）
func buildTreeWithReflection(node reflect.Value, childrenMap map[interface{}][]reflect.Value, processed map[interface{}]bool, childrenField string) {
	// 获取节点的实际值
	var nodeValue reflect.Value
	if node.Kind() == reflect.Ptr {
		if node.IsNil() {
			return
		}
		nodeValue = node.Elem()
	} else {
		nodeValue = node
	}

	id := getFieldValue(nodeValue, "ID") // 假设ID字段名为"ID"
	nodeID := id.Interface()

	// 如果已经处理过，直接返回
	if processed[nodeID] {
		return
	}

	processed[nodeID] = true

	// 获取子节点
	children := childrenMap[nodeID]

	// 递归处理子节点
	for _, child := range children {
		buildTreeWithReflection(child, childrenMap, processed, childrenField)
	}

	// 设置子节点
	childrenFieldValue := nodeValue.FieldByName(childrenField)
	if childrenFieldValue.IsValid() && childrenFieldValue.CanSet() {
		// 创建正确类型的切片
		childrenSliceType := childrenFieldValue.Type()
		childrenSlice := reflect.MakeSlice(childrenSliceType, len(children), len(children))

		for j, child := range children {
			// 处理类型匹配
			targetType := childrenSliceType.Elem()
			sourceType := child.Type()

			if targetType == sourceType {
				// 类型完全匹配，直接设置
				childrenSlice.Index(j).Set(child)
			} else if targetType.Kind() == reflect.Ptr && sourceType.Kind() != reflect.Ptr {
				// 目标是指针，源是值：创建新指针
				newPtr := reflect.New(sourceType)
				newPtr.Elem().Set(child)
				childrenSlice.Index(j).Set(newPtr)
			} else if targetType.Kind() != reflect.Ptr && sourceType.Kind() == reflect.Ptr {
				// 目标是值，源是指针：获取指针的值
				if !child.IsNil() {
					childrenSlice.Index(j).Set(child.Elem())
				} else {
					// 如果指针是nil，创建零值
					zeroVal := reflect.New(targetType).Elem()
					childrenSlice.Index(j).Set(zeroVal)
				}
			} else {
				// 其他类型不匹配情况
				childrenSlice.Index(j).Set(child)
			}
		}

		childrenFieldValue.Set(childrenSlice)
	}
}

// 辅助函数：获取字段值（处理指针）
func getFieldValue(v reflect.Value, fieldName string) reflect.Value {
	// 如果v是指针，获取其指向的值
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		// 尝试通过json tag查找
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			fieldType := t.Field(i)
			jsonTag := fieldType.Tag.Get("json")
			if jsonTag == fieldName {
				field = v.Field(i)
				break
			}
		}
	}
	return field
}

// isZeroValue 判断是否为零值
func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Invalid:
		return true
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Slice, reflect.Map, reflect.Array:
		return v.Len() == 0
	default:
		return false
	}
}
