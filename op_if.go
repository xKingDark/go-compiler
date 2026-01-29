package golang

import (
	"bytes"
	"fmt"

	schema "github.com/Opticode-Project/go-compiler/golang"
	program "github.com/Opticode-Project/go-compiler/program"
)

func (g *Generator) op_if(node *program.IndexedNode, flags EvalFlags) ([]byte, error) {
	length := node.FieldsLength()

	buf := new(bytes.Buffer)
	condition := new(bytes.Buffer)
	body := new(bytes.Buffer)
	_else := new(bytes.Buffer)

	for i := range length {
		var field program.NodeValue
		node.Fields(&field, i)

		//log.Printf("Field: %d, %s", field.Value(), schema.ValueFlag(field.Flags()))
		if field.Flags()&uint32(schema.ValueFlagPointer) != 0 {
			node := g.GetNode(field.Value())
			if node == nil {
				return nil, fmt.Errorf("attempt to access undefined node: %d", field.Value())
			}
			out, err := g.Eval(node, 0)
			if err != nil {
				return nil, err
			}

			if field.Flags()&uint32(schema.ValueFlagIfConditon) != 0 {
				condition.Write(out)
			} else if field.Flags()&uint32(schema.ValueFlagIfBody) != 0 {
				body.WriteByte('	')
				body.Write(out)
				body.WriteByte('\n')
			} else if field.Flags()&uint32(schema.ValueFlagIfElse) != 0 {
				body.WriteByte('	')
				_else.Write(out)
				body.WriteByte('\n')
			}
		} else {
			return nil, fmt.Errorf("func node fields can only be pointers")
		}
	}
	buf.Grow(condition.Len() + body.Len() + _else.Len())

	buf.Write(TokenIf.Bytes())
	buf.WriteByte(' ')
	buf.Write(condition.Bytes())
	buf.WriteByte(' ')
	buf.WriteByte('{')
	buf.WriteByte('\n')
	buf.Write(body.Bytes())
	buf.WriteByte('}')

	if _else.Len() > 0 {
		buf.WriteByte(' ')
		buf.Write(TokenElse.Bytes())
		// WARN - unable to create `if else`
		buf.WriteByte('{')
		buf.WriteByte('\n')
		buf.Write(_else.Bytes())
		buf.WriteByte('}')
	}

	return buf.Bytes(), nil
}
