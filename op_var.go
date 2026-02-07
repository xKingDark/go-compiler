package golang

import (
	"bytes"
	"fmt"

	schema "github.com/Opticode-Project/go-compiler/golang"
	program "github.com/Opticode-Project/go-compiler/program"
)

func (g *Generator) op_var(buf *bytes.Buffer, node *program.IndexedNode, flags EvalFlags, isConst bool) error {
	length := node.FieldsLength()
	if length == 0 {
		return nil
	}

	multiline := length > 1
	separatorFlag := SeperatorSpace
	if multiline {
		separatorFlag = SeperatorTab
	}

	// keyword
	if isConst {
		buf.Write(TokenConst.Bytes())
	} else {
		buf.Write(TokenVar.Bytes())
	}

	if multiline {
		buf.Write(TokenSpace.Bytes())
		buf.Write(TokenParenLeft.Bytes())
	}

	for i := range length {
		var field program.NodeValue
		node.Fields(&field, i)

		if field.Flags()&uint32(schema.ValueFlagPointer) == 0 {
			return fmt.Errorf(
				"%s node fields must be pointers",
				map[bool]string{true: "const", false: "var"}[isConst],
			)
		}

		target := g.GetNode(field.Value())
		if target == nil {
			return fmt.Errorf("attempt to access undefined node: %d", field.Value())
		}

		if multiline {
			buf.Write(TokenNewLine.Bytes())
			buf.Write(TokenTab.Bytes())
		} else if i > 0 {
			buf.Write(TokenSpace.Bytes())
		}

		if err := g.evalNode(buf, target, separatorFlag); err != nil {
			return err
		}
	}

	if multiline {
		buf.Write(TokenNewLine.Bytes())
		buf.Write(TokenParenRight.Bytes())
	}

	return nil
}
