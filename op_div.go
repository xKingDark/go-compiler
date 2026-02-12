package golang

import (
	"bytes"
	"fmt"

	program "github.com/Opticode-Project/go-compiler/program"
)

func (g *Generator) op_div(buf *bytes.Buffer, node *program.BinaryNode, flags EvalFlags) error {
	// Get the left and right values
	left := node.Left(nil)
	right := node.Right(nil)

	if left == nil || right == nil {
		return fmt.Errorf("division value operands cannot be nil")
	}

	if err := g.evalValue(buf, left, false); err != nil {
		return err
	}

	buf.Write(TokenSpace.Bytes())
	buf.Write(TokenSlash.Bytes())
	buf.Write(TokenSpace.Bytes())

	if err := g.evalValue(buf, right, false); err != nil {
		return err
	}

	return nil
}
