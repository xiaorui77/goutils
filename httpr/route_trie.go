package httpr

import "strings"

// trie node
type node struct {
	pattern  string
	part     string
	children []*node
	isWild   bool
}

func (n *node) matchChild(part string) *node {
	for _, c := range n.children {
		if c.part == part || c.isWild {
			return c
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	var nodes []*node
	for _, c := range n.children {
		if c.part == part || c.isWild {
			nodes = append(nodes, c)
		}
	}
	return nodes
}

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
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
