package golang

import (
	"bytes"
	"fmt"

	program "github.com/Opticode-Project/go-compiler/program"
)

func (g *Generator) op_inc(buf *bytes.Buffer, node *program.UnaryNode, flags EvalFlags) error {
	value := node.Value(nil)

	if value == nil {
		return fmt.Errorf("increment value cannot be nil")
	}

	if err := g.evalValue(buf, value, false); err != nil {
		return err
	}

	buf.Write(TokenSpace.Bytes())
	buf.Write(TokenIncrement.Bytes())
	return nil
}
