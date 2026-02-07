package golang

import (
	"bytes"
	"fmt"

	schema "github.com/Opticode-Project/go-compiler/golang"
	program "github.com/Opticode-Project/go-compiler/program"
)

func (g *Generator) op_if(buf *bytes.Buffer, node *program.IndexedNode, flags EvalFlags) error {
	var (
		condition bytes.Buffer
		thenBody  bytes.Buffer
		elseBody  bytes.Buffer
	)
	var elseIfNode *program.Node

	length := node.FieldsLength()
	for i := range length {
		var field program.NodeValue
		node.Fields(&field, i)

		if field.Flags()&uint32(schema.ValueFlagPointer) == 0 {
			return fmt.Errorf("if node fields can only be pointers")
		}

		node := g.GetNode(field.Value())
		if node == nil {
			return fmt.Errorf("attempt to access undefined node: %d", field.Value())
		}

		switch {
		case field.Flags()&uint32(schema.ValueFlagIfConditon) != 0:
			if err := g.evalNode(&condition, node, 0); err != nil {
				return err
			}

		case field.Flags()&uint32(schema.ValueFlagIfBody) != 0:
			thenBody.Write(TokenTab.Bytes())
			if err := g.evalNode(&thenBody, node, 0); err != nil {
				return err
			}
			thenBody.Write(TokenNewLine.Bytes())

		case field.Flags()&uint32(schema.ValueFlagIfElse) != 0:
			if schema.Opcode(node.Opcode()) == schema.OpcodeIf {
				elseIfNode = node
			} else {
				elseBody.Write(TokenTab.Bytes())
				if err := g.evalNode(&elseBody, node, 0); err != nil {
					return err
				}
				elseBody.Write(TokenNewLine.Bytes())
			}
		}
	}

	buf.Write(TokenIf.Bytes())
	buf.Write(TokenSpace.Bytes())
	buf.Write(condition.Bytes())
	buf.Write(TokenSpace.Bytes())

	buf.Write(TokenBracesLeft.Bytes())
	buf.Write(TokenNewLine.Bytes())
	buf.Write(thenBody.Bytes())
	buf.Write(TokenBracesRight.Bytes())

	if elseIfNode != nil {
		buf.Write(TokenSpace.Bytes())
		buf.Write(TokenElse.Bytes())
		buf.Write(TokenSpace.Bytes())

		if err := g.evalNode(buf, elseIfNode, flags); err != nil {
			return err
		}
	} else if elseBody.Len() > 0 {
		buf.Write(TokenSpace.Bytes())
		buf.Write(TokenElse.Bytes())
		buf.Write(TokenSpace.Bytes())

		buf.Write(TokenBracesLeft.Bytes())
		buf.Write(TokenNewLine.Bytes())
		buf.Write(elseBody.Bytes())
		buf.Write(TokenBracesRight.Bytes())
	}

	return nil
}
