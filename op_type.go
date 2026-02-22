package golang

import (
	"bytes"
	"fmt"

	program "github.com/Opticode-Project/go-compiler/program"
)

func (g *Generator) op_type(buf *bytes.Buffer, node *program.IndexedNode, flags EvalFlags) error {
	var field program.NodeValue
	node.Fields(&field, 0)

	name, ok := g.LookUpStr(node.Id())
	if !ok {
		return fmt.Errorf("string with id %d is undefined", node.Id())
	}

	// keyword
	buf.Write(TokenType.Bytes())

	buf.Write(TokenSpace.Bytes())
	buf.Write(name)
	buf.Write(TokenSpace.Bytes())

	// Look for the field type
	def, ok := g.LookUpType(field.Type())
	if !ok {
		return fmt.Errorf("type with id %d is undefined", field.Type())
	}

	// Check if the field is a structure or a function type
	_type := def.TypeType()
	if _type != program.TypeStructureType && _type != program.TypeFunctionType {
		return fmt.Errorf("the operand must be either a structure or a function but got: %d", _type)
	}

	if _type == program.TypeFunctionType {
		buf.Write(TokenFunc.Bytes())
	}

	// Get the type and write it to the buffer
	if err := g.evalType(buf, def); err != nil {
		return err
	}

	return nil
}
