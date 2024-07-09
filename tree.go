package gospresso

import (
	"net/http"
	"sort"
)

type RouteTree struct {
	root *routeTreeNode
}

type routeTreeNode struct {
	edges  nodes
	prefix string
	label  byte

	Handler http.Handler
}

type nodes []*routeTreeNode

func (ns nodes) Sort()              { sort.Sort(ns) }
func (ns nodes) Len() int           { return len(ns) }
func (ns nodes) Less(i, j int) bool { return ns[i].label < ns[j].label }
func (ns nodes) Swap(i, j int)      { ns[i], ns[j] = ns[j], ns[i] }

func NewRouteTree() *RouteTree {
	return &RouteTree{
		root: &routeTreeNode{edges: make([]*routeTreeNode, 0)},
	}
}

func (n *routeTreeNode) Search(method uint, pattern string) *routeTreeNode {
	if n == nil {
		panic("RouteTree has not been initialized")
	}

	search := pattern

	for {
		if len(search) == 0 {
			return nil
		}

		label := search[0]
		matchingEdge := n.getEdge(label)

		if matchingEdge == nil {
			return nil
		}

		n = matchingEdge

		if n.prefix == search {
			return n
		}

		search = search[longestPrefix(search, n.prefix):]
	}
}

func (n *routeTreeNode) Insert(method uint, pattern string, handler http.Handler) *routeTreeNode {
	if n == nil {
		panic("RouteTree has not been initialized")
	}

	var parent *routeTreeNode
	search := pattern

	for {
		if len(search) == 0 {
			// n.setEndpoint(method, handler, pattern)
			return n
		}

		// Use the first char of search as our index. As we traverse the tree
		// the "first char" will move through the search string.
		label := search[0]
		parent = n

		// Find the edge which corresponds to the first char, if there is one.
		n = n.getEdge(label)

		// If not, create a child for this string as we've exhausted the tree.
		if n == nil {
			child := &routeTreeNode{label: label, prefix: search}
			hn := parent.addChild(child)
			hn.setHandler(handler)
			return hn
		}

		commonPrefix := longestPrefix(search, n.prefix)

		// Continue the search down the tree, as our search string contains the full
		// prefix encapsulated by this node.
		// e.g. search = /get/foo/1
		//      prefix = /get/
		if commonPrefix == len(n.prefix) {
			search = search[commonPrefix:]
			continue
		}

		// Otherwise, we've identified a new parent prefix
		// e.g. search = /gerbil
		//      n.label = /
		//      n.prefix = /get
		//      new child
		//         prefix = /ge
		//         label  = /
		//         children = [{ prefix: /get, label: t }, { prefix: /gerbil, label: r }]
		child := &routeTreeNode{
			label:  label,
			prefix: search[:commonPrefix],
		}

		// Updates the parent to replace the child for '/' with the new parent child for the old nodes.
		parent.replaceChild(label, child)

		// Update the old node with new label and prefix info, excluding the common prefix.
		n.label = n.prefix[commonPrefix]
		n.prefix = n.prefix[commonPrefix:]
		child.addChild(n)

		// Add the new node for our new child
		newNode := &routeTreeNode{
			label:  search[0],
			prefix: search,
		}
		newNode.setHandler(handler)
		child.addChild(newNode)
	}
}

func (n *routeTreeNode) Walk(visitor func(node *routeTreeNode)) {
	if n == nil {
		return
	}

	var stack nodes = make(nodes, 1)
	stack[0] = n

	for len(stack) != 0 {
		top := stack[0]
		stack = append(stack, top.edges...)

		visitor(top)

		stack = stack[1:]
	}
}

func (n *routeTreeNode) Len() int {
	if n == nil {
		return 0
	}

	i := 0
	n.Walk(func(node *routeTreeNode) {
		if node.Handler != nil {
			i++
		}
	})

	return i
}

func (n *routeTreeNode) getEdge(label byte) *routeTreeNode {
	for i := 0; i < len(n.edges); i++ {
		if n.edges[i].label == label {
			return n.edges[i]
		}
	}

	return nil
}

func (n *routeTreeNode) replaceChild(label byte, child *routeTreeNode) {
	for i := 0; i < len(n.edges); i++ {
		if n.edges[i].label == label {
			n.edges[i] = child
			return
		}
	}

	panic("Could not find child to replace.")
}

func (n *routeTreeNode) addChild(child *routeTreeNode) *routeTreeNode {
	n.edges = append(n.edges, child)
	n.edges.Sort()

	return child
}

func (n *routeTreeNode) setHandler(handler http.Handler) {
	n.Handler = handler
}

func longestPrefix(key string, prefix string) int {
	i := 0
	for ; i < len(prefix) && i < len(key); i++ {
		if key[i] != prefix[i] {
			break
		}
	}

	return i
}
