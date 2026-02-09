package golang

import (
	"bytes"
	"fmt"

	schema "github.com/Opticode-Project/go-compiler/golang"
	program "github.com/Opticode-Project/go-compiler/program"
)

func (g *Generator) op_func(buf *bytes.Buffer, node *program.IndexedNode, flags EvalFlags) error {
	length := node.FieldsLength()

	var (
		funcId   []byte
		funcType *program.FunctionType
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

			funcName, ok := g.LookUpStr(def.Id())
			if !ok {
				return fmt.Errorf("string with id %d is undefined", def.Id())
			}

			funcId = funcName

			v, err := EvalType(def)
			if err != nil {
				return err
			}

			ft, ok := v.(*program.FunctionType)
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

	// Function declaration
	buf.Write(TokenFunc.Bytes())
	buf.Write(TokenSpace.Bytes())
	buf.Write(funcId)
	buf.Write(TokenParenLeft.Bytes())

	// Parameters
	if funcType != nil && funcType.ParamsLength() > 0 {
		err := g.writePairList(buf, funcType.ParamsLength(), funcType.Params)
		if err != nil {
			return err
		}
	}

	buf.Write(TokenParenRight.Bytes())

	// Return values
	if funcType != nil && funcType.ResultsLength() > 0 {
		buf.Write(TokenSpace.Bytes())

		// look at first result to decide parentheses
		var first program.Pair
		funcType.Results(&first, 0)

		name, ok := g.LookUpStr(first.Key())
		if !ok {
			return fmt.Errorf("string with id %d is undefined", first.Key())
		}

		needParens := funcType.ResultsLength() > 1 || len(name) > 0
		if needParens {
			buf.Write(TokenParenLeft.Bytes())
		}

		err := g.writePairList(buf, funcType.ResultsLength(), funcType.Results)
		if err != nil {
			return err
		}

		if needParens {
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

func (g *Generator) writePairList(buf *bytes.Buffer, listLength int, getPair func(obj *program.Pair, j int) bool) error {
	for i := range listLength {
		var p program.Pair
		getPair(&p, i)

		name, ok := g.LookUpStr(p.Key())
		if !ok {
			return fmt.Errorf("string with id %d is undefined", p.Key())
		}

		tdef, ok := g.LookUpType(p.Value())
		if !ok {
			return fmt.Errorf("type with id %d is undefined", p.Value())
		}

		typeStr, ok := g.LookUpStr(tdef.Id())
		if !ok {
			return fmt.Errorf("string with id %d is undefined", tdef.Id())
		}

		if i > 0 {
			buf.Write(TokenComma.Bytes())
			buf.Write(TokenSpace.Bytes())
		}

		if len(name) > 0 {
			buf.Write(name)
			buf.Write(TokenSpace.Bytes())
		}

		buf.Write(typeStr)
	}
	return nil
}
