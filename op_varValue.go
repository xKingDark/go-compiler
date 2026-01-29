package golang

import (
	"bytes"
	"fmt"

	program "github.com/Opticode-Project/go-compiler/program"
)

var test = 10

func (g *Generator) op_varValue(node *program.BinaryNode, flags EvalFlags) ([]byte, error) {
	// Get the left and right values
	left := node.Left(nil)
	right := node.Right(nil)

	// New buffer for building the content
	var buf = new(bytes.Buffer)

	if flags&SeperatorTab != 0 {
		buf.WriteByte('	')
	} else if flags&SeperatorSpace != 0 {
		buf.WriteByte(' ')
	}

	g.StrLookupMutex.Lock()
	leftVal, ok := g.LookUpStr(uint32(left.Value()))
	g.StrLookupMutex.Unlock()
	if !ok {
		return nil, fmt.Errorf("string with id %d is undefined", left.Value())
	}
	buf.Write(leftVal)

	buf.Write(JoinBytes(TokenSpace.Bytes(), TokenEqual.Bytes(), TokenSpace.Bytes()))

	g.evalValue(buf, right)

	return buf.Bytes(), nil
}
