package golang

import (
	"bytes"
	"fmt"

	program "github.com/Opticode-Project/go-compiler/program"
)

func (g *Generator) op_unarySuffix(buf *bytes.Buffer, node *program.UnaryNode, t TokenKind, flags EvalFlags) error {
	value := node.Value(nil)
	if value == nil {
		return fmt.Errorf("unary operands cannot be nil")
	}

	if err := g.evalValue(buf, value, false); err != nil {
		return err
	}

	buf.Write(t.Bytes())
	return nil
}
