package golang

import (
	"bytes"
	"fmt"

	schema "github.com/Opticode-Project/go-compiler/golang"
	program "github.com/Opticode-Project/go-compiler/program"
)

func (g *Generator) op_constValue(node *program.BinaryNode, flags EvalFlags) ([]byte, error) {
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

	g.StrLookupMutex.Lock()
	rightVal, ok := g.LookUpStr(uint32(right.Value()))
	g.StrLookupMutex.Unlock()
	if !ok {
		return nil, fmt.Errorf("string with id %d is undefined", left.Value())
	}
	buf.Write(JoinBytes(TokenSpace.Bytes(), TokenEqual.Bytes(), TokenSpace.Bytes()))
	if right.Flags()&uint32(schema.ValueFlagQuotation) != 0 {
		buf.Write(JoinBytes(TokenQuotation.Bytes(), rightVal, TokenQuotation.Bytes()))
	} else {
		buf.Write(rightVal)
	}

	return buf.Bytes(), nil
}
