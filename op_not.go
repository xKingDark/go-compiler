package golang

import (
	"bytes"
	"fmt"

	program "github.com/Opticode-Project/go-compiler/program"
)

func (g *Generator) op_not(buf *bytes.Buffer, node *program.UnaryNode, flags EvalFlags) error {
	value := node.Value(nil)

	if value == nil {
		return fmt.Errorf("not value cannot be nil")
	}

	buf.Write(TokenNot.Bytes())
	if err := g.evalValue(buf, value, false); err != nil {
		return err
	}

	return nil
}
