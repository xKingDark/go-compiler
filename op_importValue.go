package golang

import (
	"bytes"
	"fmt"

	program "github.com/Opticode-Project/go-compiler/program"
)

func (g *Generator) op_importValue(node *program.BinaryNode, flags EvalFlags) ([]byte, error) {
	left := node.Left(nil)
	right := node.Right(nil)

	var buf = new(bytes.Buffer)
	if flags&SeperatorTab != 0 {
		buf.WriteByte('	')
	} else if flags&SeperatorSpace != 0 {
		buf.WriteByte(' ')
	}
	// write import alias
	g.StrLookupMutex.Lock()
	leftValue, ok := g.LookUpStr(uint32(left.Value()))
	g.StrLookupMutex.Unlock()
	if !ok {
		return nil, fmt.Errorf("string of with id %d is undefined", left.Value())
	}
	buf.Write(leftValue)
	// write seperator
	buf.WriteByte(' ')
	// write package path
	g.StrLookupMutex.Lock()
	rightValue, ok := g.LookUpStr(uint32(right.Value()))
	g.StrLookupMutex.Unlock()
	if !ok {
		return nil, fmt.Errorf("string of with id %d is undefined", right.Value())
	}
	buf.Write(JoinBytes(TokenQuotation.Bytes(), rightValue, TokenQuotation.Bytes()))

	return buf.Bytes(), nil
}
