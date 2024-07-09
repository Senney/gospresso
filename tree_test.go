package gospresso

import (
	"net/http"
	"testing"
)

func TestInsertRoute(t *testing.T) {
	tree := NewRouteTree()
	tree.root.Insert(mGET, "/root", http.NotFoundHandler())
	tree.root.Insert(mGET, "/root/foo", http.NotFoundHandler())
	tree.root.Insert(mGET, "/root/foo/1", http.NotFoundHandler())
	tree.root.Insert(mGET, "/", http.NotFoundHandler())
}

func TestFindRoute(t *testing.T) {
	tree := NewRouteTree()
	rootNode := tree.root.Insert(mGET, "/root", http.NotFoundHandler())
	rootFooNode := tree.root.Insert(mGET, "/root/foo", http.NotFoundHandler())
	rootFooOneNode := tree.root.Insert(mGET, "/root/foo/1", http.NotFoundHandler())

	tests := []struct {
		method uint
		path   string
		want   *routeTreeNode
	}{
		{mGET, "/root", rootNode},
		{mGET, "/root/foo", rootFooNode},
		{mGET, "/root/foo/1", rootFooOneNode},
		{mGET, "/404", nil},
		{mGET, "/root/foo/bar", nil},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			v := tree.root.Search(tt.method, tt.path)
			if v != tt.want {
				t.Fatalf("Expected node %p, found %p", tt.want, v)
			}
		})
	}
}

func FuzzSearch(f *testing.F) {
	tree := NewRouteTree()
	tree.root.Insert(mGET, "/root", http.NotFoundHandler())
	tree.root.Insert(mGET, "/root/foo", http.NotFoundHandler())
	tree.root.Insert(mGET, "/root/foo/1", http.NotFoundHandler())

	f.Add("/root")
	f.Add("/root/foo")

	f.Fuzz(func(t *testing.T, s string) {
		tree.root.Search(mGET, s)
	})
}

func FuzzInsert(f *testing.F) {
	handler := http.NotFoundHandler()
	f.Add("/root")
	f.Add("/root/foo")

	f.Fuzz(func(t *testing.T, s string) {
		tree := NewRouteTree()
		tree.root.Insert(mGET, s, handler)
	})
}
