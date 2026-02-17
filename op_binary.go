package golang

import (
	"bytes"
	"fmt"

	program "github.com/Opticode-Project/go-compiler/program"
)

func (g *Generator) op_binary(buf *bytes.Buffer, node *program.BinaryNode, t TokenKind, flags EvalFlags) error {
	// Get the left and right values
	left := node.Left(nil)
	right := node.Right(nil)

	if left == nil || right == nil {
		return fmt.Errorf("equal value operands cannot be nil")
	}

	if err := g.evalValue(buf, left, false); err != nil {
		return err
	}

	buf.Write(TokenSpace.Bytes())
	buf.Write(t.Bytes())
	buf.Write(TokenSpace.Bytes())

	return g.evalValue(buf, right, false)
}
