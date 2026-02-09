package golang

import (
	"bytes"
	"fmt"

	schema "github.com/Opticode-Project/go-compiler/golang"
	program "github.com/Opticode-Project/go-compiler/program"
)

func (g *Generator) op_constValue(buf *bytes.Buffer, node *program.BinaryNode, flags EvalFlags) error {
	// Get the left and right values
	left := node.Left(nil)
	right := node.Right(nil)

	if left == nil || right == nil {
		return fmt.Errorf("assignment operands cannot be nil")
	}

	// const semantic check
	if right.Flags()&uint32(schema.ValueFlagPointer) != 0 {
		target := g.GetNode(right.Value())
		if target == nil {
			return fmt.Errorf("undefined node: %d", right.Value())
		}

		if !g.isConstantExpression(target) {
			return fmt.Errorf("const value must be a constant expression")
		}
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

	return g.evalValue(buf, right, true)
}
