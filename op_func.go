package golang

import (
	"bytes"
	"fmt"

	schema "github.com/Opticode-Project/go-compiler/golang"
	program "github.com/Opticode-Project/go-compiler/program"
)

func (g *Generator) op_func(buf *bytes.Buffer, node *program.IndexedNode, flags EvalFlags) error {
	length := node.FieldsLength()

	var body bytes.Buffer
	body.Grow(length * 32)

	for i := range length {
		var field program.NodeValue
		node.Fields(&field, i)

		if field.Flags()&uint32(schema.ValueFlagFuncMeta) != 0 {
			def, ok := g.LookUpType(field.Type())
			if !ok {
				return fmt.Errorf("type with id %d is undefined", field.Type())
			}

			if err := g.evalType(buf, def); err != nil {
				return err
			}
			continue
		}

		if field.Flags()&uint32(schema.ValueFlagPointer) == 0 {
			return fmt.Errorf("func node fields can only be pointers")
		}

		target := g.GetNode(field.Value())
		if target == nil {
			return fmt.Errorf("attempt to access undefined node: %d", field.Value())
		}

		switch {
		case field.Flags()&uint32(schema.ValueFlagFuncBody) != 0:
			body.Write(TokenTab.Bytes())
			if err := g.evalNode(&body, target, 0); err != nil {
				return err
			}

			body.Write(TokenNewLine.Bytes())
		}
	}

	// Body
	buf.Write(TokenSpace.Bytes())
	buf.Write(TokenBraceLeft.Bytes())
	buf.Write(TokenNewLine.Bytes())
	buf.Write(body.Bytes())
	buf.Write(TokenBraceRight.Bytes())

	return nil
}
