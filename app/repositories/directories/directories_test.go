package directories

import (
	"fmt"
	"os"
	Path "path/filepath"
	"runtime"
	"saymow/version-manager/app/pkg/errors"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

const PATH_SEPARATOR = string(Path.Separator)

func TestAddNode(t *testing.T) {
	dir := &Dir{
		Path:     Path.Join("home", "project"),
		Children: make(map[string]*Node),
	}

	dir.AddNode("a.txt", &Change{ChangeType: Modification, File: &File{Filepath: "home/project/a.txt"}})

	assert.Equal(t, dir.Path, Path.Join("home", "project"))
	assert.Equal(t, len(dir.Children), 1)
	assert.Equal(t, dir.Children["a.txt"].File, &File{Filepath: "home/project/a.txt"})
}

func TestAddNodeNestedPath(t *testing.T) {
	dir := &Dir{
		Path:     Path.Join("home", "project"),
		Children: make(map[string]*Node),
	}

	dir.AddNode("1.txt", &Change{ChangeType: Modification, File: &File{Filepath: "home/project/1.txt"}})
	dir.AddNode(fmt.Sprintf("a%s2.txt", PATH_SEPARATOR), &Change{ChangeType: Modification, File: &File{Filepath: "home/project/a/2.txt"}})
	dir.AddNode(fmt.Sprintf("a%s3.txt", PATH_SEPARATOR), &Change{ChangeType: Modification, File: &File{Filepath: "home/project/a/3.txt"}})
	dir.AddNode(fmt.Sprintf("a%sb%s4.txt", PATH_SEPARATOR, PATH_SEPARATOR), &Change{ChangeType: Modification, File: &File{Filepath: "home/project/a/b/4.txt"}})
	dir.AddNode(fmt.Sprintf("a%sb%s5.txt", PATH_SEPARATOR, PATH_SEPARATOR), &Change{ChangeType: Modification, File: &File{Filepath: "home/project/a/b/5.txt"}})
	dir.AddNode(fmt.Sprintf("a%sb%s6.txt", PATH_SEPARATOR, PATH_SEPARATOR), &Change{ChangeType: Modification, File: &File{Filepath: "home/project/a/b/6.txt"}})
	dir.AddNode(fmt.Sprintf("a%sc%s7.txt", PATH_SEPARATOR, PATH_SEPARATOR), &Change{ChangeType: Modification, File: &File{Filepath: "home/project/a/c/7.txt"}})
	dir.AddNode(fmt.Sprintf("a%sc%s8.txt", PATH_SEPARATOR, PATH_SEPARATOR), &Change{ChangeType: Modification, File: &File{Filepath: "home/project/a/c/8.txt"}})

	assert.Equal(t, dir.Path, Path.Join("home", "project"))
	assert.Equal(t, len(dir.Children), 2)
	assert.Equal(t, dir.Children["1.txt"].File, &File{Filepath: "home/project/1.txt"})
	assert.Equal(t, dir.Children["a"].NodeType, DirType)
	assert.Equal(t, dir.Children["a"].Dir.Path, Path.Join("home", "project", "a"))
	assert.Equal(t, len(dir.Children["a"].Dir.Children), 4)
	assert.Equal(t, dir.Children["a"].Dir.Children["2.txt"].File, &File{Filepath: "home/project/a/2.txt"})
	assert.Equal(t, dir.Children["a"].Dir.Children["3.txt"].File, &File{Filepath: "home/project/a/3.txt"})
	assert.Equal(t, dir.Children["a"].Dir.Children["b"].NodeType, DirType)
	assert.Equal(t, dir.Children["a"].Dir.Children["b"].Dir.Path, Path.Join("home", "project", "a", "b"))
	assert.Equal(t, len(dir.Children["a"].Dir.Children["b"].Dir.Children), 3)
	assert.Equal(t, dir.Children["a"].Dir.Children["b"].Dir.Children["4.txt"].File, &File{Filepath: "home/project/a/b/4.txt"})
	assert.Equal(t, dir.Children["a"].Dir.Children["b"].Dir.Children["5.txt"].File, &File{Filepath: "home/project/a/b/5.txt"})
	assert.Equal(t, dir.Children["a"].Dir.Children["b"].Dir.Children["6.txt"].File, &File{Filepath: "home/project/a/b/6.txt"})
	assert.Equal(t, dir.Children["a"].Dir.Children["c"].Dir.Path, Path.Join("home", "project", "a", "c"))
	assert.Equal(t, len(dir.Children["a"].Dir.Children["c"].Dir.Children), 2)
	assert.Equal(t, dir.Children["a"].Dir.Children["c"].Dir.Children["7.txt"].File, &File{Filepath: "home/project/a/c/7.txt"})
	assert.Equal(t, dir.Children["a"].Dir.Children["c"].Dir.Children["8.txt"].File, &File{Filepath: "home/project/a/c/8.txt"})
}

func TestAddNodeRemovalChanges(t *testing.T) {
	dir := &Dir{
		Path:     Path.Join("home", "project"),
		Children: make(map[string]*Node),
	}

	dir.AddNode("a.txt", &Change{ChangeType: Modification, File: &File{Filepath: "home/project/a.txt"}})
	dir.AddNode("b.txt", &Change{ChangeType: Modification, File: &File{Filepath: "home/project/b.txt"}})
	dir.AddNode("a.txt", &Change{ChangeType: Removal, Removal: &FileRemoval{Filepath: "home/project/a.txt"}})
	dir.AddNode("c.txt", &Change{ChangeType: Removal, Removal: &FileRemoval{Filepath: "home/project/c.txt"}})

	assert.Equal(t, dir.Path, Path.Join("home", "project"))
	assert.Equal(t, len(dir.Children), 1)
	assert.Equal(t, dir.Children["b.txt"].File, &File{Filepath: "home/project/b.txt"})
}

func TestAddNodeOverrideRemovalChanges(t *testing.T) {
	dir := &Dir{
		Path:     Path.Join("home", "project"),
		Children: make(map[string]*Node),
	}

	dir.AddNode("a.txt", &Change{ChangeType: Modification, File: &File{Filepath: "home/project/a.txt", ObjectName: "old-version"}})
	dir.AddNode("a.txt", &Change{ChangeType: Removal, Removal: &FileRemoval{Filepath: "home/project/a.txt"}})
	dir.AddNode("a.txt", &Change{ChangeType: Modification, File: &File{Filepath: "home/project/a.txt", ObjectName: "newer-version"}})

	assert.Equal(t, dir.Path, Path.Join("home", "project"))
	assert.Equal(t, len(dir.Children), 1)
	assert.Equal(t, dir.Children["a.txt"].File, &File{Filepath: "home/project/a.txt", ObjectName: "newer-version"})
}

func TestAddNodeRemovalChangesRemovesEmptyDir(t *testing.T) {
	dir := &Dir{
		Path: Path.Join("home", "project"),
		Children: map[string]*Node{
			"dir": {
				NodeType: DirType,
				Dir: &Dir{
					Path: Path.Join("home", "project", "dir"),
					Children: map[string]*Node{
						"a.txt": {
							NodeType: FileType,
							File: &File{
								"home/project/dir/a.txt",
								"object-a",
							},
						},
						"b.txt": {
							NodeType: FileType,
							File: &File{
								"home/project/dir/b.txt",
								"object-b",
							},
						},
					},
				},
			},
		},
	}

	dir.AddNode(Path.Join("dir", "a.txt"), &Change{ChangeType: Removal, Removal: &FileRemoval{Filepath: "home/project/dir/a.txt"}})
	assert.Equal(t, len(dir.Children), 1)
	assert.Equal(t, dir.Children["dir"].NodeType, DirType)
	dir.AddNode(Path.Join("dir", "b.txt"), &Change{ChangeType: Removal, Removal: &FileRemoval{Filepath: "home/project/dir/b.txt"}})
	assert.Equal(t, len(dir.Children), 0)
}

func TestAddNodeRemovalChangesNestedPath(t *testing.T) {
	dir := &Dir{
		Path:     Path.Join("home", "project"),
		Children: make(map[string]*Node),
	}

	dir.AddNode("1.txt", &Change{ChangeType: Modification, File: &File{Filepath: "home/project/1.txt"}})
	dir.AddNode(fmt.Sprintf("a%s2.txt", PATH_SEPARATOR), &Change{ChangeType: Modification, File: &File{Filepath: "home/project/a/2.txt"}})
	dir.AddNode(fmt.Sprintf("a%s3.txt", PATH_SEPARATOR), &Change{ChangeType: Modification, File: &File{Filepath: "home/project/a/3.txt"}})
	dir.AddNode(fmt.Sprintf("a%sb%s4.txt", PATH_SEPARATOR, PATH_SEPARATOR), &Change{ChangeType: Modification, File: &File{Filepath: "home/project/a/b/4.txt"}})
	dir.AddNode(fmt.Sprintf("a%sb%s5.txt", PATH_SEPARATOR, PATH_SEPARATOR), &Change{ChangeType: Modification, File: &File{Filepath: "home/project/a/b/5.txt"}})
	dir.AddNode(fmt.Sprintf("a%sb%s6.txt", PATH_SEPARATOR, PATH_SEPARATOR), &Change{ChangeType: Modification, File: &File{Filepath: "home/project/a/b/6.txt"}})
	dir.AddNode(fmt.Sprintf("a%sc%s7.txt", PATH_SEPARATOR, PATH_SEPARATOR), &Change{ChangeType: Modification, File: &File{Filepath: "home/project/a/c/7.txt"}})
	dir.AddNode(fmt.Sprintf("a%sc%s8.txt", PATH_SEPARATOR, PATH_SEPARATOR), &Change{ChangeType: Modification, File: &File{Filepath: "home/project/a/c/8.txt"}})

	assert.Equal(t, dir.Path, Path.Join("home", "project"))
	assert.Equal(t, len(dir.Children), 2)
	assert.Equal(t, dir.Children["1.txt"].File, &File{Filepath: "home/project/1.txt"})
	assert.Equal(t, dir.Children["a"].NodeType, DirType)
	assert.Equal(t, len(dir.Children["a"].Dir.Children), 4)
	assert.Equal(t, dir.Children["a"].Dir.Children["2.txt"].File, &File{Filepath: "home/project/a/2.txt"})
	assert.Equal(t, dir.Children["a"].Dir.Children["3.txt"].File, &File{Filepath: "home/project/a/3.txt"})
	assert.Equal(t, dir.Children["a"].Dir.Children["b"].NodeType, DirType)
	assert.Equal(t, dir.Children["a"].Dir.Children["b"].Dir.Path, Path.Join("home", "project", "a", "b"))
	assert.Equal(t, len(dir.Children["a"].Dir.Children["b"].Dir.Children), 3)
	assert.Equal(t, dir.Children["a"].Dir.Children["b"].Dir.Children["4.txt"].File, &File{Filepath: "home/project/a/b/4.txt"})
	assert.Equal(t, dir.Children["a"].Dir.Children["b"].Dir.Children["5.txt"].File, &File{Filepath: "home/project/a/b/5.txt"})
	assert.Equal(t, dir.Children["a"].Dir.Children["b"].Dir.Children["6.txt"].File, &File{Filepath: "home/project/a/b/6.txt"})
	assert.Equal(t, dir.Children["a"].Dir.Children["b"].NodeType, DirType)
	assert.Equal(t, dir.Children["a"].Dir.Children["c"].Dir.Path, Path.Join("home", "project", "a", "c"))
	assert.Equal(t, len(dir.Children["a"].Dir.Children["c"].Dir.Children), 2)
	assert.Equal(t, dir.Children["a"].Dir.Children["c"].Dir.Children["7.txt"].File, &File{Filepath: "home/project/a/c/7.txt"})
	assert.Equal(t, dir.Children["a"].Dir.Children["c"].Dir.Children["8.txt"].File, &File{Filepath: "home/project/a/c/8.txt"})

	dir.AddNode(fmt.Sprintf("a%s3.txt", PATH_SEPARATOR), &Change{ChangeType: Removal, Removal: &FileRemoval{Filepath: "home/project/a/3.txt"}})
	dir.AddNode(fmt.Sprintf("a%sb%s5.txt", PATH_SEPARATOR, PATH_SEPARATOR), &Change{ChangeType: Removal, Removal: &FileRemoval{Filepath: "home/project/a/b/5.txt"}})
	dir.AddNode(fmt.Sprintf("a%sc%s8.txt", PATH_SEPARATOR, PATH_SEPARATOR), &Change{ChangeType: Removal, Removal: &FileRemoval{Filepath: "home/project/a/c/8.txt"}})

	assert.Equal(t, len(dir.Children), 2)
	assert.Equal(t, dir.Children["1.txt"].File, &File{Filepath: "home/project/1.txt"})
	assert.Equal(t, dir.Children["a"].NodeType, DirType)
	assert.Equal(t, len(dir.Children["a"].Dir.Children), 3)
	assert.Equal(t, dir.Children["a"].Dir.Children["2.txt"].File, &File{Filepath: "home/project/a/2.txt"})
	_, ok := dir.Children["a"].Dir.Children["3.txt"]
	assert.False(t, ok)
	assert.Equal(t, dir.Children["a"].Dir.Children["b"].NodeType, DirType)
	assert.Equal(t, len(dir.Children["a"].Dir.Children["b"].Dir.Children), 2)
	assert.Equal(t, dir.Children["a"].Dir.Children["b"].Dir.Children["4.txt"].File, &File{Filepath: "home/project/a/b/4.txt"})
	_, ok = dir.Children["a"].Dir.Children["b"].Dir.Children["5.txt"]
	assert.False(t, ok)
	assert.Equal(t, dir.Children["a"].Dir.Children["b"].Dir.Children["6.txt"].File, &File{Filepath: "home/project/a/b/6.txt"})
	assert.Equal(t, len(dir.Children["a"].Dir.Children["c"].Dir.Children), 1)
	assert.Equal(t, dir.Children["a"].Dir.Children["c"].Dir.Children["7.txt"].File, &File{Filepath: "home/project/a/c/7.txt"})
	_, ok = dir.Children["a"].Dir.Children["c"].Dir.Children["8.txt"]
	assert.False(t, ok)

	dir.AddNode(fmt.Sprintf("a%sc%s8.txt", PATH_SEPARATOR, PATH_SEPARATOR), &Change{ChangeType: Modification, File: &File{Filepath: "home/project/a/c/8.txt", ObjectName: "newer-version"}})

	assert.Equal(t, len(dir.Children), 2)
	assert.Equal(t, dir.Children["1.txt"].File, &File{Filepath: "home/project/1.txt"})
	assert.Equal(t, dir.Children["a"].NodeType, DirType)
	assert.Equal(t, len(dir.Children["a"].Dir.Children), 3)
	assert.Equal(t, dir.Children["a"].Dir.Children["2.txt"].File, &File{Filepath: "home/project/a/2.txt"})
	_, ok = dir.Children["a"].Dir.Children["3.txt"]
	assert.False(t, ok)
	assert.Equal(t, dir.Children["a"].Dir.Children["b"].NodeType, DirType)
	assert.Equal(t, len(dir.Children["a"].Dir.Children["b"].Dir.Children), 2)
	assert.Equal(t, dir.Children["a"].Dir.Children["b"].Dir.Children["4.txt"].File, &File{Filepath: "home/project/a/b/4.txt"})
	_, ok = dir.Children["a"].Dir.Children["b"].Dir.Children["5.txt"]
	assert.False(t, ok)
	assert.Equal(t, dir.Children["a"].Dir.Children["b"].Dir.Children["6.txt"].File, &File{Filepath: "home/project/a/b/6.txt"})
	assert.Equal(t, len(dir.Children["a"].Dir.Children["c"].Dir.Children), 2)
	assert.Equal(t, dir.Children["a"].Dir.Children["c"].Dir.Children["7.txt"].File, &File{Filepath: "home/project/a/c/7.txt"})
	assert.Equal(t, dir.Children["a"].Dir.Children["c"].Dir.Children["8.txt"].File, &File{Filepath: "home/project/a/c/8.txt", ObjectName: "newer-version"})
}

func TestFindNode(t *testing.T) {
	dir := &Dir{
		Path: Path.Join("home", "project"),
		Children: map[string]*Node{
			"a.txt": {
				NodeType: FileType,
				File: &File{
					"home/project/a.txt",
					"object-a",
				},
			},
			"b.txt": {
				NodeType: FileType,
				File: &File{
					"home/project/b.txt",
					"object-b",
				},
			},
		},
	}

	assert.Equal(t, dir.FindNode("").NodeType, DirType)
	assert.Equal(t, dir.FindNode("a.txt").NodeType, FileType)
	assert.Equal(t, dir.FindNode("a.txt").File, &File{"home/project/a.txt", "object-a"})
	assert.Equal(t, dir.FindNode("b.txt").NodeType, FileType)
	assert.Equal(t, dir.FindNode("b.txt").File, &File{"home/project/b.txt", "object-b"})
}

func TestFindNodeNestedPath(t *testing.T) {
	dir := &Dir{
		Children: map[string]*Node{
			"a.txt": {
				NodeType: FileType,
				File: &File{
					"home/project/a.txt",
					"object-a",
				},
			},
			"b.txt": {
				NodeType: FileType,
				File: &File{
					"home/project/b.txt",
					"object-b",
				},
			},
			"subdir": {
				NodeType: DirType,
				Dir: &Dir{
					Children: map[string]*Node{
						"a.txt": {
							NodeType: FileType,
							File: &File{
								"home/project/subdir/a.txt",
								"object-subdir-a",
							},
						},
						"c.txt": {
							NodeType: FileType,
							File: &File{
								"home/project/subdir/c.txt",
								"object-subdir-c",
							},
						},
						"nested-subdir": {
							NodeType: DirType,
							Dir: &Dir{
								Children: map[string]*Node{
									"b.txt": {
										NodeType: FileType,
										File: &File{
											"home/project/subdir/nested-subdir/b.txt",
											"object-subdir-nested-subdir-b",
										},
									},
									"d.txt": {
										NodeType: FileType,
										File: &File{
											"home/project/subdir/nested-subdir/d.txt",
											"object-subdir-nested-subdir-d",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// File nodes

	assert.Equal(t, dir.FindNode("").NodeType, DirType)
	assert.Equal(t, dir.FindNode("a.txt").NodeType, FileType)
	assert.Equal(t, dir.FindNode("a.txt").File, &File{"home/project/a.txt", "object-a"})
	assert.Equal(t, dir.FindNode("b.txt").NodeType, FileType)
	assert.Equal(t, dir.FindNode("b.txt").File, &File{"home/project/b.txt", "object-b"})
	assert.Equal(t, dir.FindNode("subdir").NodeType, DirType)
	assert.Equal(t, dir.FindNode(fmt.Sprintf("subdir%sa.txt", PATH_SEPARATOR)).NodeType, FileType)
	assert.Equal(t, dir.FindNode(fmt.Sprintf("subdir%sa.txt", PATH_SEPARATOR)).File, &File{"home/project/subdir/a.txt", "object-subdir-a"})
	assert.Equal(t, dir.FindNode(fmt.Sprintf("subdir%sc.txt", PATH_SEPARATOR)).NodeType, FileType)
	assert.Equal(t, dir.FindNode(fmt.Sprintf("subdir%sc.txt", PATH_SEPARATOR)).File, &File{"home/project/subdir/c.txt", "object-subdir-c"})
	assert.Equal(t, dir.FindNode(Path.Join("subdir", "nested-subdir")).NodeType, DirType)
	assert.Equal(t, dir.FindNode(fmt.Sprintf("subdir%snested-subdir%sb.txt", PATH_SEPARATOR, PATH_SEPARATOR)).NodeType, FileType)
	assert.Equal(t, dir.FindNode(fmt.Sprintf("subdir%snested-subdir%sb.txt", PATH_SEPARATOR, PATH_SEPARATOR)).File, &File{"home/project/subdir/nested-subdir/b.txt", "object-subdir-nested-subdir-b"})
	assert.Equal(t, dir.FindNode(fmt.Sprintf("subdir%snested-subdir%sd.txt", PATH_SEPARATOR, PATH_SEPARATOR)).NodeType, FileType)
	assert.Equal(t, dir.FindNode(fmt.Sprintf("subdir%snested-subdir%sd.txt", PATH_SEPARATOR, PATH_SEPARATOR)).File, &File{"home/project/subdir/nested-subdir/d.txt", "object-subdir-nested-subdir-d"})

}

func TestAllCollectFiles(t *testing.T) {
	dir := &Dir{
		Children: map[string]*Node{
			"a.txt": {
				NodeType: FileType,
				File: &File{
					"home/project/a.txt",
					"1",
				},
			},
			"b.txt": {
				NodeType: FileType,
				File: &File{
					"home/project/b.txt",
					"2",
				},
			},
			"subdir": {
				NodeType: DirType,
				Dir: &Dir{
					Children: map[string]*Node{
						"a.txt": {
							NodeType: FileType,
							File: &File{
								"home/project/subdir/a.txt",
								"3",
							},
						},
						"c.txt": {
							NodeType: FileType,
							File: &File{
								"home/project/subdir/c.txt",
								"4",
							},
						},
						"nested-subdir": {
							NodeType: DirType,
							Dir: &Dir{
								Children: map[string]*Node{
									"b.txt": {
										NodeType: FileType,
										File: &File{
											"home/project/subdir/nested-subdir/b.txt",
											"5",
										},
									},
									"d.txt": {
										NodeType: FileType,
										File: &File{
											"home/project/subdir/nested-subdir/d.txt",
											"6",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	received := dir.CollectAllFiles()

	sort.Slice(received, func(i, j int) bool { return received[i].ObjectName < received[j].ObjectName })
	assert.Equal(t,
		received,
		[]*File{
			{
				"home/project/a.txt",
				"1",
			},
			{
				"home/project/b.txt",
				"2",
			},
			{
				"home/project/subdir/a.txt",
				"3",
			},
			{
				"home/project/subdir/c.txt",
				"4",
			},
			{
				"home/project/subdir/nested-subdir/b.txt",
				"5",
			},
			{
				"home/project/subdir/nested-subdir/d.txt",
				"6",
			},
		},
	)
}

// Since it's we cannot rely on the sequence the of the map iteration, the test becomes
// hard. This is the reason there are only file nodes on the last dir node.
// This teste ensure that the nodes are indeed collected in pre order.
func TestPreOrderTraversal(t *testing.T) {
	dir := &Dir{
		Path: "",
		Children: map[string]*Node{
			"subdir": {
				NodeType: DirType,
				Dir: &Dir{
					Path: "subdir",
					Children: map[string]*Node{
						"nested-subdir": {
							NodeType: DirType,
							Dir: &Dir{
								Path: Path.Join("subdir", "nested-subdir"),
								Children: map[string]*Node{
									"b.txt": {
										NodeType: FileType,
										File: &File{
											"home/project/subdir/nested-subdir/b.txt",
											"5",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	received := dir.PreOrderTraversal()

	assert.EqualValues(t, len(received), 4)
	assert.EqualValues(t, received[0].NodeType, DirType)
	assert.EqualValues(t, received[0].Dir.Path, "")
	assert.EqualValues(t, received[1].NodeType, DirType)
	assert.EqualValues(t, received[1].Dir.Path, "subdir")
	assert.EqualValues(t, received[2].NodeType, DirType)
	assert.EqualValues(t, received[2].Dir.Path, Path.Join("subdir", "nested-subdir"))
	assert.EqualValues(t,
		received[3],
		&Node{
			NodeType: FileType,
			File: &File{
				"home/project/subdir/nested-subdir/b.txt",
				"5",
			},
		},
	)
}

func TestMerge(t *testing.T) {
	// unmodified
	{
		dir := &Dir{
			Path: Path.Join(getOsRoot(), "home", "project"),
			Children: map[string]*Node{
				"a.txt": {
					NodeType: FileType,
					File: &File{
						"home/project/a.txt",
						"object-a",
					},
				},
				"b.txt": {
					NodeType: FileType,
					File: &File{
						"home/project/b.txt",
						"object-b",
					},
				},
			},
		}

		dir.Merge(&Dir{
			Path: Path.Join(getOsRoot(), "home"),
			Children: map[string]*Node{
				"a.txt": {
					NodeType: FileType,
					File: &File{
						"home/a.txt",
						"object-a",
					},
				},
				"b.txt": {
					NodeType: FileType,
					File: &File{
						"home/b.txt",
						"object-b",
					},
				},
			},
		})

		assert.Equal(t, dir.Path, Path.Join(getOsRoot(), "home", "project"))
		assert.Equal(t, len(dir.Children), 2)
		assert.Equal(t, dir.Children["a.txt"].File, &File{Filepath: "home/project/a.txt", ObjectName: "object-a"})
		assert.Equal(t, dir.Children["b.txt"].File, &File{Filepath: "home/project/b.txt", ObjectName: "object-b"})
	}

	// add nodes
	{
		dir := &Dir{
			Path: Path.Join(getOsRoot(), "home", "project"),
			Children: map[string]*Node{
				"a.txt": {
					NodeType: FileType,
					File: &File{
						Path.Join(getOsRoot(), "home", "project", "a.txt"),
						"object-a",
					},
				},
				"b.txt": {
					NodeType: FileType,
					File: &File{
						Path.Join(getOsRoot(), "home", "project", "b.txt"),
						"object-b",
					},
				},
			},
		}

		dir.Merge(&Dir{
			Path: Path.Join(getOsRoot(), "home"),
			Children: map[string]*Node{
				"a.txt": {
					NodeType: FileType,
					File: &File{
						Path.Join(getOsRoot(), "home", "a.txt"),
						"object-a",
					},
				},
				"b.txt": {
					NodeType: FileType,
					File: &File{
						Path.Join(getOsRoot(), "home", "b.txt"),
						"object-b",
					},
				},
				"project": {
					NodeType: DirType,
					Dir: &Dir{
						Path: Path.Join(getOsRoot(), "home", "project"),
						Children: map[string]*Node{
							"c.txt": {
								NodeType: FileType,
								File: &File{
									Path.Join(getOsRoot(), "home", "project", "c.txt"),
									"object-c",
								},
							},
							"d.txt": {
								NodeType: FileType,
								File: &File{
									Path.Join(getOsRoot(), "home", "project", "d.txt"),
									"object-d",
								},
							},
							"folder": {
								NodeType: DirType,
								Dir: &Dir{
									Path: Path.Join(getOsRoot(), "home", "project", "folder"),
									Children: map[string]*Node{
										"a.txt": {
											NodeType: FileType,
											File: &File{
												Path.Join(getOsRoot(), "home", "project", "folder", "a.txt"),
												"object-a",
											},
										},
										"b.txt": {
											NodeType: FileType,
											File: &File{
												Path.Join(getOsRoot(), "home", "project", "folder", "b.txt"),
												"object-b",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		})

		assert.Equal(t, dir.Path, Path.Join(getOsRoot(), "home", "project"))
		assert.Equal(t, len(dir.Children), 5)
		assert.Equal(t, dir.Children["a.txt"].File, &File{Filepath: Path.Join(getOsRoot(), "home", "project", "a.txt"), ObjectName: "object-a"})
		assert.Equal(t, dir.Children["b.txt"].File, &File{Filepath: Path.Join(getOsRoot(), "home", "project", "b.txt"), ObjectName: "object-b"})
		assert.Equal(t, dir.Children["c.txt"].File, &File{Filepath: Path.Join(getOsRoot(), "home", "project", "c.txt"), ObjectName: "object-c"})
		assert.Equal(t, dir.Children["d.txt"].File, &File{Filepath: Path.Join(getOsRoot(), "home", "project", "d.txt"), ObjectName: "object-d"})
		assert.Equal(t, dir.Children["folder"].NodeType, DirType)
		assert.Equal(t, dir.Children["folder"].Dir.Path, Path.Join(getOsRoot(), "home", "project", "folder"))
		assert.Equal(t, len(dir.Children["folder"].Dir.Children), 2)
		assert.Equal(t, dir.Children["folder"].Dir.Children["a.txt"].File, &File{Filepath: Path.Join(getOsRoot(), "home", "project", "folder", "a.txt"), ObjectName: "object-a"})
		assert.Equal(t, dir.Children["folder"].Dir.Children["b.txt"].File, &File{Filepath: Path.Join(getOsRoot(), "home", "project", "folder", "b.txt"), ObjectName: "object-b"})
	}

	// add and override nodes
	{
		dir := &Dir{
			Path: Path.Join(getOsRoot(), "home", "project"),
			Children: map[string]*Node{
				"a.txt": {
					NodeType: FileType,
					File: &File{
						Path.Join(getOsRoot(), "home", "project", "a.txt"),
						"object-a",
					},
				},
				"b.txt": {
					NodeType: FileType,
					File: &File{
						Path.Join(getOsRoot(), "home", "project", "b.txt"),
						"object-b",
					},
				},
			},
		}

		dir.Merge(&Dir{
			Path: Path.Join(getOsRoot(), "home", "project"),
			Children: map[string]*Node{
				"a.txt": {
					NodeType: FileType,
					File: &File{
						Path.Join(getOsRoot(), "home", "project", "a.txt"),
						"object-a-overridden",
					},
				},
				"b.txt": {
					NodeType: FileType,
					File: &File{
						Path.Join(getOsRoot(), "home", "project", "b.txt"),
						"object-b-overridden",
					},
				},
				"folder": {
					NodeType: DirType,
					Dir: &Dir{
						Path: Path.Join(getOsRoot(), "home", "project", "folder"),
						Children: map[string]*Node{
							"a.txt": {
								NodeType: FileType,
								File: &File{
									Path.Join(getOsRoot(), "home", "project", "folder", "a.txt"),
									"object-a",
								},
							},
							"b.txt": {
								NodeType: FileType,
								File: &File{
									Path.Join(getOsRoot(), "home", "project", "folder", "b.txt"),
									"object-b",
								},
							},
						},
					},
				},
			},
		})

		assert.Equal(t, dir.Path, Path.Join(getOsRoot(), "home", "project"))
		assert.Equal(t, len(dir.Children), 3)
		assert.Equal(t, dir.Children["a.txt"].File, &File{Filepath: Path.Join(getOsRoot(), "home", "project", "a.txt"), ObjectName: "object-a-overridden"})
		assert.Equal(t, dir.Children["b.txt"].File, &File{Filepath: Path.Join(getOsRoot(), "home", "project", "b.txt"), ObjectName: "object-b-overridden"})
		assert.Equal(t, dir.Children["folder"].NodeType, DirType)
		assert.Equal(t, dir.Children["folder"].Dir.Path, Path.Join(getOsRoot(), "home", "project", "folder"))
		assert.Equal(t, len(dir.Children["folder"].Dir.Children), 2)
		assert.Equal(t, dir.Children["folder"].Dir.Children["a.txt"].File, &File{Filepath: Path.Join(getOsRoot(), "home", "project", "folder", "a.txt"), ObjectName: "object-a"})
		assert.Equal(t, dir.Children["folder"].Dir.Children["b.txt"].File, &File{Filepath: Path.Join(getOsRoot(), "home", "project", "folder", "b.txt"), ObjectName: "object-b"})
	}
}

func getOsRoot() string {
	if runtime.GOOS == "windows" {
		return "C:\\"
	}

	return "/"
}

func TestNormalizePath(t *testing.T) {
	base, err := os.Getwd()
	errors.Check(err)

	dir := &Dir{
		Path:     base,
		Children: make(map[string]*Node),
	}

	filepath, err := dir.NormalizePath(Path.Join(base, "1.txt"))
	assert.Nil(t, err)
	assert.Equal(t, filepath, "1.txt")

	filepath, err = dir.NormalizePath(Path.Join(base, "a", "1.txt"))
	assert.Nil(t, err)
	assert.Equal(t, filepath, Path.Join("a", "1.txt"))

	filepath, err = dir.NormalizePath(Path.Join(base, "a", "b", "1.txt"))
	assert.Nil(t, err)
	assert.Equal(t, filepath, Path.Join("a", "b", "1.txt"))

	filepath, err = dir.NormalizePath(Path.Join(base, "folder"))
	assert.Nil(t, err)
	assert.Equal(t, filepath, "folder")

	filepath, err = dir.NormalizePath(Path.Join(base, "a", "b"))
	assert.Nil(t, err)
	assert.Equal(t, filepath, Path.Join("a", "b"))

	filepath, err = dir.NormalizePath(base)
	assert.Nil(t, err)
	assert.Equal(t, filepath, "")

	filepath, err = dir.NormalizePath(Path.Join(getOsRoot(), "a", "b"))
	assert.Error(t, err, "invalid path.")
	assert.Equal(t, filepath, "")

	filepath, err = dir.NormalizePath(Path.Join(getOsRoot(), "home", "a"))
	assert.Error(t, err, "invalid path.")
	assert.Equal(t, filepath, "")

	filepath, err = dir.NormalizePath(Path.Join(getOsRoot(), "home", "projectads", "a"))
	assert.Error(t, err, "invalid path.")
	assert.Equal(t, filepath, "")
}

func TestAbsPath(t *testing.T) {
	base, err := os.Getwd()
	errors.Check(err)

	dir := &Dir{
		Path:     base,
		Children: make(map[string]*Node),
	}

	filepath, err := dir.AbsPath(Path.Join(base, "1.txt"))
	assert.Nil(t, err)
	assert.Equal(t, filepath, Path.Join(base, "1.txt"))

	filepath, err = dir.AbsPath(Path.Join(base, "a", "1.txt"))
	assert.Nil(t, err)
	assert.Equal(t, filepath, Path.Join(base, "a", "1.txt"))

	filepath, err = dir.AbsPath("1.txt")
	assert.Nil(t, err)
	assert.Equal(t, filepath, Path.Join(base, "1.txt"))

	filepath, err = dir.AbsPath(Path.Join("a", "1.txt"))
	assert.Nil(t, err)
	assert.Equal(t, filepath, Path.Join(base, "a", "1.txt"))

	filepath, err = dir.AbsPath(Path.Join(getOsRoot(), "a", "b"))
	assert.Error(t, err, "invalid path.")
	assert.Equal(t, filepath, "")

	filepath, err = dir.AbsPath(Path.Join(getOsRoot(), "home", "a"))
	assert.Error(t, err, "invalid path.")
	assert.Equal(t, filepath, "")

	filepath, err = dir.AbsPath(Path.Join(getOsRoot(), "home", "projectads", "a"))
	assert.Error(t, err, "invalid path.")
	assert.Equal(t, filepath, "")
}

func TestChangeConflicts(t *testing.T) {
	// Distinct Filepaths

	change := &Change{
		ChangeType: Creation,
		File: &File{
			Filepath:   "a",
			ObjectName: "",
		},
	}
	assert.False(
		t, change.Conflicts(&Change{
			ChangeType: Removal,
			Removal: &FileRemoval{
				Filepath: "b",
			},
		}),
	)

	change = &Change{
		ChangeType: Removal,
		Removal: &FileRemoval{
			Filepath: "b",
		},
	}
	assert.False(
		t, change.Conflicts(
			&Change{
				ChangeType: Creation,
				File: &File{
					Filepath:   "a",
					ObjectName: "",
				},
			},
		),
	)

	// Both are removals

	change = &Change{
		ChangeType: Removal,
		Removal: &FileRemoval{
			Filepath: "a",
		},
	}
	assert.False(
		t, change.Conflicts(&Change{
			ChangeType: Removal,
			Removal: &FileRemoval{
				Filepath: "a",
			},
		}),
	)

	// Equal file hash

	change = &Change{
		ChangeType: Creation,
		File: &File{
			Filepath:   "a",
			ObjectName: "hash",
		},
	}
	assert.False(
		t, change.Conflicts(&Change{
			ChangeType: Modification,
			File: &File{
				Filepath:   "a",
				ObjectName: "hash",
			},
		}),
	)

	change = &Change{
		ChangeType: Modification,
		File: &File{
			Filepath:   "a",
			ObjectName: "hash",
		},
	}
	assert.False(
		t, change.Conflicts(&Change{
			ChangeType: Creation,
			File: &File{
				Filepath:   "a",
				ObjectName: "hash",
			},
		}),
	)

	// Distinct file hash

	change = &Change{
		ChangeType: Modification,
		File: &File{
			Filepath:   "a",
			ObjectName: "hash",
		},
	}
	assert.True(
		t, change.Conflicts(&Change{
			ChangeType: Removal,
			Removal: &FileRemoval{
				Filepath: "a",
			},
		}),
	)

	change = &Change{
		ChangeType: Removal,
		Removal: &FileRemoval{
			Filepath: "a",
		},
	}
	assert.True(
		t,
		change.Conflicts(&Change{
			ChangeType: Modification,
			File: &File{
				Filepath:   "a",
				ObjectName: "hash",
			},
		}),
	)

	change = &Change{
		ChangeType: Creation,
		File: &File{
			Filepath:   "a",
			ObjectName: "hash",
		},
	}
	assert.True(
		t, change.Conflicts(&Change{
			ChangeType: Removal,
			Removal: &FileRemoval{
				Filepath: "a",
			},
		}),
	)

	change = &Change{
		ChangeType: Removal,
		Removal: &FileRemoval{
			Filepath: "a",
		},
	}
	assert.True(
		t,
		change.Conflicts(&Change{
			ChangeType: Creation,
			File: &File{
				Filepath:   "a",
				ObjectName: "hash",
			},
		}),
	)

	change = &Change{
		ChangeType: Creation,
		File: &File{
			Filepath:   "a",
			ObjectName: "different-hash",
		},
	}
	assert.True(
		t,
		change.Conflicts(&Change{
			ChangeType: Creation,
			File: &File{
				Filepath:   "a",
				ObjectName: "hash",
			},
		}),
	)

	change = &Change{
		ChangeType: Modification,
		File: &File{
			Filepath:   "a",
			ObjectName: "different-hash",
		},
	}
	assert.True(
		t,
		change.Conflicts(&Change{
			ChangeType: Modification,
			File: &File{
				Filepath:   "a",
				ObjectName: "hash",
			},
		}),
	)
}
