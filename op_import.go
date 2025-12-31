package golang

import (
	"bytes"
	"fmt"

	schema "github.com/Opticode-Project/go-compiler/golang"
	program "github.com/Opticode-Project/go-compiler/program"
)

func (g *Generator) op_import(node *program.IndexedNode, flags EvalFlags) ([]byte, error) {
	length := node.FieldsLength()
	buf := new(bytes.Buffer)
	var flag EvalFlags
	if length > 1 {
		flag |= SeperatorTab
	} else {
		flag |= SeperatorSpace
	}

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

			//! This preforms an allocation each cycle. (not good)
			buf.Write(JoinBytes(TokenSpace.Bytes(), out))
		} else {
			return nil, fmt.Errorf("import node fields can only be pointers")
		}
		if i > 0 {
			buf.WriteByte('\n')
		}
	}

	if length > 1 {
		//! messy and needs work
		return JoinBytes(TokenImport.Bytes(), TokenSpace.Bytes(), TokenParenLeft.Bytes(), TokenNewLine.Bytes(), buf.Bytes(), TokenParenRight.Bytes()), nil
	}
	//! includes the tab seperator (not intended)
	return JoinBytes(TokenImport.Bytes(), buf.Bytes()), nil
}
