package golang

import (
	"bytes"
	"fmt"

	schema "github.com/Opticode-Project/go-compiler/golang"
	program "github.com/Opticode-Project/go-compiler/program"
)

const (
	bodyGrowthModifer    = 32
	paramsGrowthModifer  = 16
	resultsGrowthModifer = 16
)

func (g *Generator) op_func(buf *bytes.Buffer, node *program.IndexedNode, flags EvalFlags) error {
	length := node.FieldsLength()

	var (
		funcId   []byte
		funcType *program.FuncType
	)

	var (
		params bytes.Buffer
		body   bytes.Buffer
	)

	for i := range length {
		var field program.NodeValue
		node.Fields(&field, i)

		if field.Flags()&uint32(schema.ValueFlagFuncMeta) != 0 {
			def, ok := g.LookUpType(field.Type())
			if !ok {
				return fmt.Errorf("type with id %d is undefined", field.Type())
			}

			funcId = def.Id()

			v, err := EvalType(def)
			if err != nil {
				return err
			}

			ft, ok := v.(*program.FuncType)
			if !ok {
				return fmt.Errorf("expected *program.FuncType but got %T", v)
			}

			funcType = ft
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
		case field.Flags()&uint32(schema.ValueFlagFuncParam) != 0:
			if params.Len() > 0 {
				params.Write(TokenComma.Bytes())
				params.Write(TokenSpace.Bytes())
			}

			if err := g.evalNode(&params, target, 0); err != nil {
				return err
			}

		case field.Flags()&uint32(schema.ValueFlagFuncBody) != 0:
			body.Write(TokenTab.Bytes())
			if err := g.evalNode(&body, target, 0); err != nil {
				return err
			}

			body.Write(TokenNewLine.Bytes())
		}
	}

	resultsLen := funcType.ResultsLength()

	// Function declaration
	buf.Write(TokenFunc.Bytes())
	buf.Write(TokenSpace.Bytes())
	buf.Write(funcId)
	buf.Write(TokenParenLeft.Bytes())
	buf.Write(params.Bytes())
	buf.Write(TokenParenRight.Bytes())

	// Return values
	if funcType != nil && resultsLen > 0 {
		buf.Write(TokenSpace.Bytes())
		if resultsLen > 1 {
			buf.Write(TokenParenLeft.Bytes())
		}

		for i := range resultsLen {
			var p program.Pair
			funcType.Results(&p, i)

			name, ok := g.LookUpStr(p.Key())
			if !ok {
				return fmt.Errorf("string with id %d is undefined", p.Key())
			}

			_type, ok := g.LookUpType(p.Value())
			if !ok {
				return fmt.Errorf("type with id %d is undefined", p.Value())
			}

			if i > 0 {
				buf.Write(TokenComma.Bytes())
				buf.Write(TokenSpace.Bytes())
			}

			if len(name) > 0 {
				buf.Write(name)
				buf.Write(TokenSpace.Bytes())
			}

			buf.Write(_type.Id())
		}

		if resultsLen > 1 {
			buf.Write(TokenParenRight.Bytes())
		}
	}

	// Body
	buf.Write(TokenSpace.Bytes())
	buf.Write(TokenBracesLeft.Bytes())
	buf.Write(TokenNewLine.Bytes())
	buf.Write(body.Bytes())
	buf.Write(TokenBracesRight.Bytes())

	return nil
}
