package gee

import "strings"

// 静态路由时，路径与handle一一映射，由map记录

// 通过trie树，可以匹配后缀不同的url，且trie相对map对key存储更节省空间

// node 将parttern信息-url存储于叶子节点，part用于提取共同前缀，即part用作key，pattern用作value。
// 匹配某个路由时，遍历node的part，查询一个node，parts--，直到parts为空匹配成功，读取/写入 当前node的pattern

// pattern不止存在叶节点，任一节点若含有pattern，也代表止于该节点的url
type node struct {
	pattern  string
	part     string
	children []*node
	isWild   bool
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// pattern指一个url
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	chidren := n.matchChildren(part)

	for _, child := range chidren {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
