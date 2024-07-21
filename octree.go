package sixel

import "image/color"

const MAX_DEPTH = 8

type Color struct {
	R, G, B int
}

type Node struct {
	Level, ColorCount, ColorIndex int
	RGB                           Color
	IsLeaf                        bool
	Next                          [8]*Node
}

var (
	Size        int
	OctreeDepth int
	Octree      *Node
)

func NewOctree(depth int) *Node {
	is_leaf := depth == OctreeDepth
	if is_leaf {
		Size++
	}
	return &Node{
		Level:  depth,
		IsLeaf: is_leaf,
	}
}

func GetBitAt(num int, i int) int {
	return (num >> i) & 1
}

func Branch(rgb Color, depth int) int {
	i := MAX_DEPTH - depth
	return GetBitAt(rgb.R, i)*4 + GetBitAt(rgb.G, i)*2 + GetBitAt(rgb.B, i)
}

func InsertTree(tree *Node, rgb Color, depth int) *Node {
	if tree == nil {
		tree = NewOctree(depth)
	}
	if tree.IsLeaf {
		tree.ColorCount++
		AddColors(tree.RGB, rgb)
	} else {
		InsertTree(tree.Next[Branch(rgb, depth)], rgb, depth+1)
	}
}

func ReduceTree() {
	tree := GetReducible(tree)
	sum := Color{}
	children := 0
	for _, next := range tree.Next {
		if next == nil {
			continue
		}
		children++
		AddColors(sum, next.RGB)

	}
	tree.RGB = sum
	tree.IsLeaf = true
	Size -= children + 1
}

func GenerateOctree(pixels []color.Color, k int) {
	for _, p := range pixels {
		r, g, b, _ := p.RGBA()
		rgb := Color{R: int(r >> 8), G: int(g >> 8), B: int(b >> 8)}
		Octree = InsertTree(Octree, rgb, 1)
		for Size > k {
			ReduceTree()
		}
	}
}

func GetPalette() {

}
