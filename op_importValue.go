package golang

import (
	"bytes"
	"fmt"

	program "github.com/Opticode-Project/go-compiler/program"
)

func (g *Generator) op_importValue(buf *bytes.Buffer, node *program.BinaryNode, flags EvalFlags) error {
	left := node.Left(nil)
	right := node.Right(nil)

	if left == nil || right == nil {
		return fmt.Errorf("import value operands cannot be nil")
	}

	// indentation
	if flags&SeperatorTab != 0 {
		buf.Write(TokenTab.Bytes())
	} else if flags&SeperatorSpace != 0 {
		buf.Write(TokenSpace.Bytes())
	}

	// write import alias
	leftValue, ok := g.LookUpStr(uint32(left.Value()))
	if !ok {
		return fmt.Errorf("string with id %d is undefined", left.Value())
	}
	buf.Write(leftValue)

	// write seperator
	if len(leftValue) > 0 {
		buf.Write(TokenSpace.Bytes())
	}

	// quoted package path
	rightValue, ok := g.LookUpStr(uint32(right.Value()))
	if !ok {
		return fmt.Errorf("string with id %d is undefined", right.Value())
	}

	buf.Write(TokenQuotation.Bytes())
	buf.Write(rightValue)
	buf.Write(TokenQuotation.Bytes())

	return nil
}
