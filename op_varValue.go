package golang

import (
	"bytes"
	"fmt"

	program "github.com/Opticode-Project/go-compiler/program"
)

var test = 10

func (g *Generator) op_varValue(buf *bytes.Buffer, node *program.BinaryNode, flags EvalFlags) error {
	// Get the left and right values
	left := node.Left(nil)
	right := node.Right(nil)

	if left == nil || right == nil {
		return fmt.Errorf("assignment operands cannot be nil")
	}

	// indentation
	if flags&SeperatorTab != 0 {
		buf.Write(TokenTab.Bytes())
	} else if flags&SeperatorSpace != 0 {
		buf.Write(TokenSpace.Bytes())
	}

	leftVal, ok := g.LookUpStr(uint32(left.Value()))
	if !ok {
		return fmt.Errorf("string with id %d is undefined", left.Value())
	}

	buf.Write(leftVal)

	buf.Write(TokenSpace.Bytes())
	buf.Write(TokenEqual.Bytes())
	buf.Write(TokenSpace.Bytes())

	return g.evalValue(buf, right, false)
}
