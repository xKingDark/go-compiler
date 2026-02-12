package golang

import (
	"bytes"
	"fmt"

	program "github.com/Opticode-Project/go-compiler/program"
)

func (g *Generator) op_deref(buf *bytes.Buffer, node *program.UnaryNode, flags EvalFlags) error {
	value := node.Value(nil)

	if value == nil {
		return fmt.Errorf("dereference value cannot be nil")
	}

	buf.Write(TokenStar.Bytes())
	if err := g.evalValue(buf, value, false); err != nil {
		return err
	}

	return nil
}
