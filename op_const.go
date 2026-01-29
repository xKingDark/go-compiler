package golang

import (
	"bytes"
	"fmt"

	schema "github.com/Opticode-Project/go-compiler/golang"
	program "github.com/Opticode-Project/go-compiler/program"
)

func (g *Generator) op_const(node *program.IndexedNode, flags EvalFlags) ([]byte, error) {
	length := node.FieldsLength()
	var flag EvalFlags
	if length > 1 {
		flag |= SeperatorTab
	} else {
		flag |= SeperatorSpace
	}
	buf := new(bytes.Buffer)
	for i := range length {
		var field program.NodeValue
		node.Fields(&field, i)
		if i > 0 {
			buf.WriteByte('\n')
		}
		//log.Printf("Field: %d, %s", field.Value(), schema.ValueFlag(field.Flags()))
		if field.Flags()&uint32(schema.ValueFlagPointer) != 0 {
			// is pointer
			node := g.GetNode(field.Value())
			if node == nil {
				return nil, fmt.Errorf("attempt to access undefined node: %d", field.Value())
			}
			out, err := g.Eval(node, flag)
			if err != nil {
				return nil, err
			}

			buf.Write(out)
		} else {
			return nil, fmt.Errorf("const node fields can only be pointers")
		}
		if i > 0 {
			buf.WriteByte('\n')
		}
	}

	if length > 1 {
		//! messy and needs work
		return JoinBytes(TokenConst.Bytes(), TokenSpace.Bytes(), TokenParenLeft.Bytes(), TokenNewLine.Bytes(), buf.Bytes(), TokenParenRight.Bytes()), nil
	}
	return JoinBytes(TokenConst.Bytes(), buf.Bytes()), nil
}
