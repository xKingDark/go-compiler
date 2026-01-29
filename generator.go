package golang

import (
	"bytes"
	"log"
	"os"
	"sync"

	schema "github.com/Opticode-Project/go-compiler/golang"
	program "github.com/Opticode-Project/go-compiler/program"
)

type GoFile struct {
	Path    string
	Content *[]byte
	nodes   *map[int64]*DeserializeNode
}

func (gf *GoFile) Write(s string) error {
	return os.WriteFile(s, *gf.Content, 0)
}

type DeserializeNode struct {
	Flags   schema.NodeFlag // Node flags
	Module  uint8           // Index of the nodes' module
	Content *[]byte         // Evaluated node content
	Span    [2]uint32       // File line span
}

type Generator struct {
	program *program.App
	buf     *[]byte

	nodeOffsets    map[int64]int
	StrLookupMutex sync.Mutex
	modulePath     map[uint8][]int

	// nodeId -> deserialized node
	nodes      map[int64]*DeserializeNode
	nodesMutex sync.Mutex
}

func NewGenerator(app *program.App, buf *[]byte) *Generator {
	var nodesLength = app.NodesLength()
	var nodeOffset = make(map[int64]int, nodesLength)
	var path = []int{0}
	for i := range nodesLength {
		var n program.Node
		app.Nodes(&n, i)

		log.Println(n.Id(), schema.Opcode(n.Opcode()), schema.NodeFlag(n.Flags()), n.Next())

		if i == int(path[len(path)-1]) {
			path = append(path, int(n.Next()))
		}

		nodeOffset[n.Id()] = i
	}

	modulePath := make(map[uint8][]int)
	modulePath[0] = path
	log.Printf("Path: %v", path)
	return &Generator{
		program:     app,
		buf:         buf,
		nodeOffsets: nodeOffset,
		modulePath:  modulePath,
		nodes:       make(map[int64]*DeserializeNode, nodesLength),
	}
}

func (g *Generator) LookUpStr(i uint32) ([]byte, bool) {
	var str program.StringEntry
	ok := g.program.LutByKey(&str, i)
	if !ok {
		return nil, ok
	}
	return str.Value(), ok
}

func (g *Generator) LookUpType(t uint32) (*program.TypeDef, bool) {
	var _type program.TypeEntry
	ok := g.program.TypesByKey(&_type, t)
	if !ok {
		return nil, ok
	}
	return _type.Value(nil), true
}

func (g *Generator) GetNode(id int64) *program.Node {
	i, ok := g.nodeOffsets[id]
	if !ok {
		return nil
	}
	var node program.Node
	g.program.Nodes(&node, i)

	return &node
}

func (g *Generator) Write(id int64, flags schema.NodeFlag, module uint8, content *[]byte) {
	g.nodes[id] = &DeserializeNode{
		Flags:   flags,
		Module:  module,
		Content: content,

		Span: [2]uint32{},
	}
}

func (g *Generator) PrintNodes() {
	for k, v := range g.nodes {
		log.Println(k, string(*v.Content))
	}
}

func (g *Generator) Export(p []int) ([]*GoFile, error) {
	var out bytes.Buffer
	for _, id := range p {
		if n, ok := g.nodes[int64(id)]; ok {
			out.Grow(len(*n.Content))
			out.Write(*n.Content)
			out.WriteRune('\n')
		}
	}

	o := out.Bytes()

	t := &GoFile{
		Content: &o,
		nodes:   &g.nodes,
	}

	return []*GoFile{t}, nil
}
