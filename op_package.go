package golang

import (
	"bytes"
	"fmt"

	program "github.com/Opticode-Project/go-compiler/program"
)

func (g *Generator) op_package(buf *bytes.Buffer, node *program.IndexedNode, flags EvalFlags) error {
	id, ok := g.LookUpStr(node.Id())
	if !ok {
		return fmt.Errorf("string with id %d is undefined", node.Id())
	}

	buf.Write(TokenPackage.Bytes())
	buf.Write(TokenSpace.Bytes())
	buf.Write(id)
	return nil
}
