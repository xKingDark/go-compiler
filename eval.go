package golang

import (
	"bytes"
	"fmt"
	"strconv"

	schema "github.com/Opticode-Project/go-compiler/golang"
	program "github.com/Opticode-Project/go-compiler/program"
	fb "github.com/google/flatbuffers/go"
)

type EvalFlags uint16

const (
	SeperatorTab EvalFlags = 1 << iota
	SeperatorSpace
)

func (g *Generator) Eval(node *program.Node, evalFlags EvalFlags) ([]byte, error) {
	var buf bytes.Buffer
	err := g.evalNode(&buf, node, evalFlags)
	return buf.Bytes(), err
}

func (g *Generator) evalNode(buf *bytes.Buffer, node *program.Node, evalFlags EvalFlags) error {
	if node == nil {
		return fmt.Errorf("node is nil")
	}

	var unionTable fb.Table
	if !node.Node(&unionTable) {
		return fmt.Errorf("failed to access union of node: %d", node.Id())
	}

	opcode := schema.Opcode(node.Opcode())

	switch node.NodeType() {
	case program.NodeUnionIndexedNode:
		n := new(program.IndexedNode)
		n.Init(unionTable.Bytes, unionTable.Pos)

		err := g.EvalIndexed(buf, opcode, n, evalFlags)
		if err != nil {
			return err
		}

	case program.NodeUnionBinaryNode:
		n := new(program.BinaryNode)
		n.Init(unionTable.Bytes, unionTable.Pos)

		err := g.EvalBinary(buf, opcode, n, evalFlags)
		if err != nil {
			return err
		}

	case program.NodeUnionUnaryNode:
		n := new(program.UnaryNode)
		n.Init(unionTable.Bytes, unionTable.Pos)

		err := g.EvalUnary(buf, opcode, n, evalFlags)
		if err != nil {
			return err
		}

	case program.NodeUnionNONE:
		return fmt.Errorf("node %d has no union payload", node.Id())

	default:
		return fmt.Errorf("unknown node union type %d", node.NodeType())
	}

	return nil
}

func (g *Generator) EvalIndexed(buf *bytes.Buffer, opcode schema.Opcode, node *program.IndexedNode, evalFlags EvalFlags) error {
	switch opcode {
	case schema.OpcodePackage:
		return g.op_package(buf, node, evalFlags)
	case schema.OpcodeImport:
		return g.op_import(buf, node, evalFlags)
	case schema.OpcodeConst:
		return g.op_const(buf, node, evalFlags)
	case schema.OpcodeVar:
		return g.op_var(buf, node, evalFlags)
	case schema.OpcodeIf:
		return g.op_if(buf, node, evalFlags)
	case schema.OpcodeFunc:
		return g.op_func(buf, node, evalFlags)
	case schema.OpcodeCall:
		return g.op_call(buf, node, evalFlags)
	case schema.OpcodeReturn:
		return g.op_return(buf, node, evalFlags)
	}

	return fmt.Errorf("invalid opcode on node with opcode of %s", opcode)
}

func (g *Generator) EvalBinary(buf *bytes.Buffer, opcode schema.Opcode, node *program.BinaryNode, evalFlags EvalFlags) error {
	switch opcode {
	case schema.OpcodeImportValue:
		return g.op_importValue(buf, node, evalFlags)
	case schema.OpcodeConstValue:
		return g.op_constValue(buf, node, evalFlags)
	case schema.OpcodeVarValue:
		return g.op_varValue(buf, node, evalFlags)

	case schema.OpcodeEqual:
		return g.op_binary(buf, node, TokenCompare, evalFlags)
	case schema.OpcodeNotEqual:
		return g.op_binary(buf, node, TokenNotEqual, evalFlags)
	case schema.OpcodeLess:
		return g.op_binary(buf, node, TokenLess, evalFlags)
	case schema.OpcodeLessEqual:
		return g.op_binary(buf, node, TokenLessEqual, evalFlags)
	case schema.OpcodeGreater:
		return g.op_binary(buf, node, TokenGreater, evalFlags)
	case schema.OpcodeGreaterEqual:
		return g.op_binary(buf, node, TokenGreaterEqual, evalFlags)
	case schema.OpcodeAnd:
		return g.op_binary(buf, node, TokenAnd, evalFlags)
	case schema.OpcodeOr:
		return g.op_binary(buf, node, TokenOr, evalFlags)

	case schema.OpcodeAdd:
		return g.op_binary(buf, node, TokenPlus, evalFlags)
	case schema.OpcodeSub:
		return g.op_binary(buf, node, TokenMinus, evalFlags)
	case schema.OpcodeMul:
		return g.op_binary(buf, node, TokenStar, evalFlags)
	case schema.OpcodeDiv:
		return g.op_binary(buf, node, TokenSlash, evalFlags)
	case schema.OpcodeMod:
		return g.op_binary(buf, node, TokenModulus, evalFlags)
	case schema.OpcodeAssign:
		return g.op_binary(buf, node, TokenEqual, evalFlags)
	case schema.OpcodeAddAssign:
		return g.op_binary(buf, node, TokenAddAssign, evalFlags)
	case schema.OpcodeSubAssign:
		return g.op_binary(buf, node, TokenSubAssign, evalFlags)
	case schema.OpcodeMulAssign:
		return g.op_binary(buf, node, TokenMulAssign, evalFlags)
	case schema.OpcodeDivAssign:
		return g.op_binary(buf, node, TokenDivAssign, evalFlags)
	case schema.OpcodeModAssign:
		return g.op_binary(buf, node, TokenModAssign, evalFlags)

	case schema.OpcodeBitAndAssign:
		return g.op_binary(buf, node, TokenBitAndAssign, evalFlags)
	case schema.OpcodeBitOrAssign:
		return g.op_binary(buf, node, TokenBitOrAssign, evalFlags)
	case schema.OpcodeBitXorAssign:
		return g.op_binary(buf, node, TokenBitXorAssign, evalFlags)
	case schema.OpcodeBitClearAssign:
		return g.op_binary(buf, node, TokenBitClearAssign, evalFlags)
	case schema.OpcodeLeftShiftAssign:
		return g.op_binary(buf, node, TokenShiftLeftAssign, evalFlags)
	case schema.OpcodeRightShiftAssign:
		return g.op_binary(buf, node, TokenShiftRightAssign, evalFlags)
	case schema.OpcodeBitAnd:
		return g.op_binary(buf, node, TokenBitAnd, evalFlags)
	case schema.OpcodeBitOr:
		return g.op_binary(buf, node, TokenBitOr, evalFlags)
	case schema.OpcodeBitXor:
		return g.op_binary(buf, node, TokenBitXor, evalFlags)
	case schema.OpcodeBitClear:
		return g.op_binary(buf, node, TokenBitClear, evalFlags)
	case schema.OpcodeLeftShift:
		return g.op_binary(buf, node, TokenShiftLeft, evalFlags)
	case schema.OpcodeRightShift:
		return g.op_binary(buf, node, TokenShiftRight, evalFlags)
	}

	return fmt.Errorf("invalid opcode on node with opcode of %s", opcode)
}

func (g *Generator) EvalUnary(buf *bytes.Buffer, opcode schema.Opcode, node *program.UnaryNode, evalFlags EvalFlags) error {
	switch opcode {
	case schema.OpcodeNot:
		return g.op_unaryPrefix(buf, node, TokenNot, evalFlags)
	case schema.OpcodeDefer:
		return g.op_defer(buf, node, evalFlags)
	case schema.OpcodeGoRoutine:
		return g.op_goRoutine(buf, node, evalFlags)

	case schema.OpcodeInc:
		return g.op_unaryPostfix(buf, node, TokenIncrement, evalFlags)
	case schema.OpcodeDec:
		return g.op_unaryPostfix(buf, node, TokenDecrement, evalFlags)
	case schema.OpcodeAddrOf:
		return g.op_unaryPrefix(buf, node, TokenBitAnd, evalFlags)
	case schema.OpcodeDeref:
		return g.op_unaryPrefix(buf, node, TokenStar, evalFlags)
	}

	return fmt.Errorf("invalid opcode on node with opcode of %s", opcode)
}

func isConstOperator(op schema.Opcode) bool {
	// arithmetic
	if op >= schema.OpcodeAdd && op <= schema.OpcodeMod {
		return true
	}

	// comparisons
	if op >= schema.OpcodeEqual && op <= schema.OpcodeGreaterEqual {
		return true
	}

	// logical
	if op >= schema.OpcodeAnd && op <= schema.OpcodeNot {
		return true
	}

	// bitwise
	if op >= schema.OpcodeBitAnd && op <= schema.OpcodeRightShift {
		return true
	}

	return false
}

func (g *Generator) isConstValue(v *program.NodeValue) bool {
	if v == nil {
		return false
	}

	// literal value
	if v.Flags()&uint32(schema.ValueFlagPointer) == 0 {
		return true
	}

	// pointer -> must recurse
	node := g.GetNode(v.Value())
	if node == nil {
		return false
	}

	return g.isConstantExpression(node)
}

func (g *Generator) isConstantExpression(node *program.Node) bool {
	if node == nil {
		return false
	}

	op := schema.Opcode(node.Opcode())

	// only allow const-safe operators
	if !isConstOperator(op) {
		return false
	}

	var unionTable fb.Table
	if !node.Node(&unionTable) {
		return false
	}

	switch node.NodeType() {
	case program.NodeUnionUnaryNode:
		n := new(program.UnaryNode)
		n.Init(unionTable.Bytes, unionTable.Pos)

		return g.isConstValue(n.Value(nil))

	case program.NodeUnionBinaryNode:
		n := new(program.BinaryNode)
		n.Init(unionTable.Bytes, unionTable.Pos)

		return g.isConstValue(n.Left(nil)) &&
			g.isConstValue(n.Right(nil))

	default:
		return false
	}
}

func (g *Generator) evalValue(buf *bytes.Buffer, nodeValue *program.NodeValue, isConst bool) error {
	if nodeValue == nil {
		return fmt.Errorf("node value is null")
	}

	if nodeValue.Flags()&uint32(schema.ValueFlagPointer) != 0 {
		if isConst {
			return fmt.Errorf("const value cannot reference runtime expression")
		}

		node := g.GetNode(nodeValue.Value())
		if node == nil {
			return fmt.Errorf("attempt to access undefined node: %d", nodeValue.Value())
		}

		err := g.evalNode(buf, node, 0)
		if err != nil {
			return err
		}
		return nil
	}

	value, ok := g.LookUpStr(uint32(nodeValue.Value()))
	if !ok {
		return fmt.Errorf("string with id %d is undefined", nodeValue.Value())
	}

	if nodeValue.Flags()&uint32(schema.ValueFlagQuotation) != 0 {
		buf.Write(TokenQuotation.Bytes())
		buf.Write(value)
		buf.Write(TokenQuotation.Bytes())
	} else {
		buf.Write(value)
	}

	return nil
}

func EvalType(t *program.TypeDef) (any, error) {
	var unionTable fb.Table
	if !t.Type(&unionTable) {
		return nil, fmt.Errorf("failed to access union of type: %d", t.Id())
	}

	switch t.TypeType() {
	case program.TypeFunctionType:
		ptr := new(program.FunctionType)
		ptr.Init(unionTable.Bytes, unionTable.Pos)

		return ptr, nil
	case program.TypePointerType:
		ptr := new(program.PointerType)
		ptr.Init(unionTable.Bytes, unionTable.Pos)

		return ptr, nil
	case program.TypeMapType:
		ptr := new(program.MapType)
		ptr.Init(unionTable.Bytes, unionTable.Pos)

		return ptr, nil
	case program.TypeArrayType:
		ptr := new(program.ArrayType)
		ptr.Init(unionTable.Bytes, unionTable.Pos)

		return ptr, nil
	case program.TypeTupleType:
		ptr := new(program.TupleType)
		ptr.Init(unionTable.Bytes, unionTable.Pos)

		return ptr, nil
	case program.TypeStructureType:
		ptr := new(program.StructureType)
		ptr.Init(unionTable.Bytes, unionTable.Pos)

		return ptr, nil

	default:
		return nil, fmt.Errorf("unknown type kind: %d", t.TypeType())
	}
}

func (g *Generator) evalType(buf *bytes.Buffer, t *program.TypeDef) error {
	if t.TypeType() == program.TypeNONE {
		name, ok := g.LookUpStr(t.Base())
		if !ok {
			return fmt.Errorf("string with id %d is undefined", t.Base())
		}

		buf.Write(name)
		return nil
	}

	v, err := EvalType(t)
	if err != nil {
		return err
	}

	switch ty := v.(type) {
	case *program.PointerType:
		buf.Write(TokenStar.Bytes())
		elem, ok := g.LookUpType(ty.Elem())
		if !ok {
			return fmt.Errorf("type with id %d is undefined", ty.Elem())
		}

		return g.evalType(buf, elem)

	case *program.MapType:
		buf.Write(TokenMap.Bytes())
		buf.Write(TokenBracketLeft.Bytes())

		key, ok := g.LookUpType(ty.Key())
		if !ok {
			return fmt.Errorf("type with id %d is undefined", ty.Key())
		}

		if err := g.evalType(buf, key); err != nil {
			return err
		}

		buf.Write(TokenBracketRight.Bytes())

		value, ok := g.LookUpType(ty.Value())
		if !ok {
			return fmt.Errorf("type with id %d is undefined", ty.Value())
		}

		return g.evalType(buf, value)

	case *program.ArrayType:
		buf.Write(TokenBracketLeft.Bytes())
		buf.WriteString(strconv.Itoa(int(ty.Size())))
		buf.Write(TokenBracketRight.Bytes())

		elem, ok := g.LookUpType(ty.Elem())
		if !ok {
			return fmt.Errorf("type with id %d is undefined", ty.Elem())
		}

		return g.evalType(buf, elem)

	case *program.TupleType:
		buf.Write(TokenParenLeft.Bytes())
		for i := 0; i < ty.ElemLength(); i++ {
			if i > 0 {
				buf.Write(TokenComma.Bytes())
				buf.Write(TokenSpace.Bytes())
			}

			elem, ok := g.LookUpType(ty.Elem(i))
			if !ok {
				return fmt.Errorf("type with id %d is undefined", ty.Elem(i))
			}

			if err := g.evalType(buf, elem); err != nil {
				return err
			}
		}
		buf.Write(TokenParenRight.Bytes())
		return nil

	case *program.StructureType:
		buf.Write(TokenStruct.Bytes())
		buf.Write(TokenSpace.Bytes())

		buf.Write(TokenBraceLeft.Bytes())

		buf.Write(TokenNewLine.Bytes())
		for i := 0; i < ty.FieldsLength(); i++ {
			buf.Write(TokenTab.Bytes())

			var f program.StructureField
			ty.Fields(&f, i)

			name, ok := g.LookUpStr(f.Name())
			if !ok {
				return fmt.Errorf("string with id %d is undefined", f.Name())
			}
			buf.Write(name)
			buf.Write(TokenSpace.Bytes())

			def, ok := g.LookUpType(f.Type())
			if !ok {
				return fmt.Errorf("type with id %d is undefined", f.Type())
			}

			if err := g.evalType(buf, def); err != nil {
				return err
			}

			buf.Write(TokenNewLine.Bytes())
		}

		buf.Write(TokenBraceRight.Bytes())
		return nil

	case *program.FunctionType:
		// Function declaration
		buf.Write(TokenFunc.Bytes())

		funcName, ok := g.LookUpStr(t.Id())
		if !ok {
			return fmt.Errorf("string with id %d is undefined", t.Id())
		}

		buf.Write(TokenSpace.Bytes())
		buf.Write(funcName)

		// Parameters
		buf.Write(TokenParenLeft.Bytes())
		if ty.ParamsLength() > 0 {
			err := g.writePairList(buf, ty.ParamsLength(), ty.Params)
			if err != nil {
				return err
			}
		}
		buf.Write(TokenParenRight.Bytes())

		// Return values
		if ty.ResultsLength() > 0 {
			buf.Write(TokenSpace.Bytes())

			// look at first result to decide parentheses
			var first program.Pair
			ty.Results(&first, 0)

			name, ok := g.LookUpStr(first.Key())
			if !ok {
				return fmt.Errorf("string with id %d is undefined", first.Key())
			}

			needParens := ty.ResultsLength() > 1 || len(name) > 0
			if needParens {
				buf.Write(TokenParenLeft.Bytes())
			}

			err := g.writePairList(buf, ty.ResultsLength(), ty.Results)
			if err != nil {
				return err
			}

			if needParens {
				buf.Write(TokenParenRight.Bytes())
			}
		}
		return nil

	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

func (g *Generator) writePairList(buf *bytes.Buffer, listLength int, getPair func(obj *program.Pair, j int) bool) error {
	for i := range listLength {
		if i > 0 {
			buf.Write(TokenComma.Bytes())
			buf.Write(TokenSpace.Bytes())
		}

		var p program.Pair
		getPair(&p, i)

		name, ok := g.LookUpStr(p.Key())
		if !ok {
			return fmt.Errorf("string with id %d is undefined", p.Key())
		}

		if len(name) > 0 {
			buf.Write(name)
			buf.Write(TokenSpace.Bytes())
		}

		def, ok := g.LookUpType(p.Value())
		if !ok {
			return fmt.Errorf("type with id %d is undefined", p.Value())
		}

		if err := g.evalType(buf, def); err != nil {
			return err
		}
	}
	return nil
}
