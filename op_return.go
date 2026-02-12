package golang

import (
	"bytes"
	"fmt"

	schema "github.com/Opticode-Project/go-compiler/golang"
	program "github.com/Opticode-Project/go-compiler/program"
)

func (g *Generator) op_return(buf *bytes.Buffer, node *program.IndexedNode, flags EvalFlags) error {
	length := node.FieldsLength()

	var params bytes.Buffer
	params.Grow(length * paramsGrowthModifer)

	for i := range length {
		var field program.NodeValue
		node.Fields(&field, i)

		if field.Flags()&uint32(schema.ValueFlagPointer) == 0 {
			return fmt.Errorf("func node fields can only be pointers")
		}

		target := g.GetNode(field.Value())
		if target == nil {
			return fmt.Errorf("attempt to access undefined node: %d", field.Value())
		}

		if params.Len() > 0 {
			params.Write(TokenComma.Bytes())
			params.Write(TokenSpace.Bytes())
		}

		if err := g.evalNode(&params, target, 0); err != nil {
			return err
		}
	}

	var declarationLength = TokenReturn.Len()
	buf.Grow(params.Len() + declarationLength)

	// Function declaration
	buf.Write(TokenReturn.Bytes())
	if params.Len() > 0 {
		buf.Write(TokenSpace.Bytes())
		buf.Write(params.Bytes())
	}

	return nil
}
