package golang

import (
	"bytes"

	program "github.com/Opticode-Project/go-compiler/program"
)

func (g *Generator) op_equal(node *program.BinaryNode, flags EvalFlags) ([]byte, error) {
	// Get the left and right values
	left := node.Left(nil)
	right := node.Right(nil)

	// New buffer for building the content
	var buf = new(bytes.Buffer)
	g.evalValue(buf, left)

	buf.Write(JoinBytes(TokenSpace.Bytes(), TokenCompare.Bytes(), TokenSpace.Bytes()))

	g.evalValue(buf, right)

	return buf.Bytes(), nil
}
