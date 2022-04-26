package dragon

import (
	"strings"
)

func parsePath(path string) []string {
	p := strings.Split(path, "/")

	parts := make([]string, 0)
	for _, item := range p {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

type regexNode struct {
	path    string
	part    string
	regex   string
	isMatch bool
}

const (
	closeProtected        uint8 = ^(uint8(1))
	closeRunFree          uint8 = ^(uint8(2))
	closeSpacingDetection uint8 = ^(uint8(4))
	closeCompeteStatus    uint8 = ^(uint8(8))
)

func regexParsePath(path string) (string, []regexNode) {
	var (
		state = struct {
			pathLen                   int    // 路径长度
			bracketIndex              uint64 // 保存括号下标
			unprotectedPartitionIndex uint64 // 保存非保护模式下的分隔符/下标
			protectedPartitionIndex   uint64 // 保存保护模式下的分隔符/下标
			// 第1位     进行{检测，保护机制 protected;
			// 第2位     轮空控制 runFree;
			// 第3位     跨区检测 spacingDetection;
			// 第4位		竞态检测 competeStatus;
			controlFlag uint8
		}{pathLen: len(path)}
		parts = make([]regexNode, 0, 5)
	)

	// 去除最后一项的'/'
	if path[state.pathLen-1] != '/' {
		path = path + "/"
	}

	// 进行路径遍历，跳过第一个'/'
	for i := 1; i < len(path); i++ {
		// 文件类型检测
		// 检测到文件类型后续所有的操作都停止进行
		if path[i] == '*' {
			parts = append(parts, regexNode{
				part: path[i:],
			})
			break
		}

		// 加速无效字段匹配速度
		if path[i] != '/' && path[i] != '{' && path[i] != '}' {
			continue
		}

		// 得到正则表达式后，后续所有内容全部轮空直到'/'
		if state.controlFlag&2 != 0 {
			if path[i] == '/' {
				state.controlFlag &= closeRunFree // 关闭轮空模式
				state.unprotectedPartitionIndex = uint64(i)
			}
			continue
		}

		// 保护模式状态检测
		// 保护模式开启后，将不在检测'/'字段，直至遇到跨区括号碰撞或者'}'结束
		if state.controlFlag&1 != 0 {
			switch {
			// 保留最近截点，保证碰撞后再拆分到最近截点位置停止
			case path[i] == '/':
				state.controlFlag &= closeSpacingDetection // 解除跨区检测
				state.protectedPartitionIndex = uint64(i)
			// 正则表达式匹配成功
			case path[i] == '}':
				r := path[state.bracketIndex+1 : i]
				parts = append(parts, regexNode{
					part:  path[state.unprotectedPartitionIndex+1 : state.bracketIndex],
					regex: r,
					// 判断正则表达式是否为空
					isMatch: len(r) != 0,
				})
				state.controlFlag &= closeProtected // 关闭保护模式
				state.controlFlag |= 2              // 开启轮空模式
			// 跨区情况下的括号碰撞
			case path[i] == '{' && state.controlFlag&4 != 0:
				state.controlFlag &= closeProtected
				state.protectedPartitionIndex = uint64(i)
				i = int(state.unprotectedPartitionIndex)
				state.controlFlag |= 8 // 开启竞态检测
			}
		} else {
			// 非保护模式下匹配到'/'，保存非正则路径段信息
			if path[i] == '/' {
				parts = append(parts, regexNode{
					part:    path[state.unprotectedPartitionIndex+1 : i],
					isMatch: false,
				})
				// 非保模分隔符下标保存
				state.unprotectedPartitionIndex = uint64(i)
			} else if path[i] == '{' {
				if state.controlFlag&8 != 0 {
					continue
				}
				state.controlFlag &= closeCompeteStatus
				state.controlFlag |= 1         // 开启保护模式
				state.controlFlag |= 4         // 开启跨区检测
				state.bracketIndex = uint64(i) // 存留保护模式开启位置
			}
		}
	}

	// 保护模式下的结尾路由段可能会被忽略，循环结束后进行判断
	if state.controlFlag&1 != 0 {
		parts = append(parts, regexNode{
			part: path[state.unprotectedPartitionIndex+1 : state.pathLen],
		})
	}

	return path, parts
}

type node struct {
	path     string
	part     string
	regex    string
	division uint32
	children []*node
	isWild   bool
	isMatch  bool
}

/* 普通匹配 */
// 匹配单节点 - 插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *node) insert(path string, parts []string, height int) {
	if len(parts) == height {
		n.path = path
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}
	child.insert(path, parts, height+1)
}

// 匹配多节点 - 查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.path == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}

/* 正则路径匹配 */
func (n *node) regexMatchChild(part regexNode) *node {
	for _, child := range n.children {
		if child.isMatch && ((child.isWild && child.regex == part.regex) || (child.part == part.part && child.regex == part.regex)) {
			return child
		}
		if !child.isMatch && (child.part == part.part || child.isWild) {
			return child
		}
	}
	return nil
}

func (n *node) regexInsert(path string, parts []regexNode, height int) {
	if len(parts) == height {
		n.path = path
		return
	}

	part := parts[height]
	child := n.regexMatchChild(part)
	if child == nil {
		child = &node{
			part:    part.part,
			regex:   part.regex,
			isMatch: part.isMatch,
		}
		if len(child.part) > 0 {
			child.isWild = child.part[0] == ':' || child.part[0] == '*'
		}
		if child.isMatch {
			n.children = append(n.children, child)
		} else {
			t := n.children
			n.children = append(t[0:n.division], child)
			n.children = append(n.children, t[n.division:]...)
			n.division++
		}
	}
	child.regexInsert(path, parts, height+1)
}

func (n *node) regexMatchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.isMatch && regexVerify(child.regex, part) {
			nodes = append(nodes, child)
		}
		if !child.isMatch && (child.part == part || child.isWild) {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) regexSearch(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.path == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.regexMatchChildren(part)

	for _, child := range children {
		result := child.regexSearch(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
