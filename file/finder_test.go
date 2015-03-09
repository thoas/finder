package file

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type Node struct {
	name    string
	entries []*Node // nil if the entry is a file
	mark    int
}

var dir, err = filepath.Abs(filepath.Dir(os.Args[0]))

func In(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

var tree = &Node{
	filepath.Join(dir, "testdata"),
	[]*Node{
		{"a", nil, 0},
		{"b", []*Node{}, 0},
		{"c", nil, 0},
		{
			"d",
			[]*Node{
				{"x", nil, 0},
				{"y", []*Node{}, 0},
				{
					"z",
					[]*Node{
						{"u", nil, 0},
						{"v", nil, 0},
					},
					0,
				},
			},
			0,
		},
	},
	0,
}

func walkTree(n *Node, path string, f func(path string, n *Node)) {
	f(path, n)
	for _, e := range n.entries {
		walkTree(e, filepath.Join(path, e.name), f)
	}
}

func makeTree(t *testing.T) {
	walkTree(tree, tree.name, func(path string, n *Node) {
		if n.entries == nil {
			fd, err := os.Create(path)
			if err != nil {
				t.Errorf("makeTree: %v", err)
				return
			}
			fd.Close()
		} else {
			os.Mkdir(path, 0770)
		}
	})
}

func withTree(t *testing.T, f func(t *testing.T)) {
	makeTree(t)

	f(t)

	if err := os.RemoveAll(tree.name); err != nil {
		t.Errorf("removeTree: %v", err)
	}
}

func TestList(t *testing.T) {
	withTree(t, func(t *testing.T) {
		finder := New([]string{tree.name})
		paths := finder.List([]string{})

		for _, filename := range []string{"testdata/a", "testdata/b"} {
			if In(filepath.Join(dir, filename), paths) == false {
				t.Error(fmt.Sprintf("%s does not belong to paths", filename))
			}
		}
	})
}

func TestListWithPatterns(t *testing.T) {
	withTree(t, func(t *testing.T) {
		finder := New([]string{tree.name})
		paths := finder.List([]string{".*testdata.*"})

		for _, filename := range []string{"testdata/a", "testdata/b"} {
			if In(filepath.Join(dir, filename), paths) == true {
				t.Error(fmt.Sprintf("%s should not belong to paths", filename))
			}
		}
	})
}
