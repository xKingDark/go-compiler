package golang

import (
	"bytes"
	"fmt"

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
		return g.op_equal(buf, node, evalFlags)
	case schema.OpcodeNotEqual:
		return g.op_notEqual(buf, node, evalFlags)
	case schema.OpcodeLess:
		return g.op_less(buf, node, evalFlags)
	case schema.OpcodeLessEqual:
		return g.op_lessEqual(buf, node, evalFlags)
	case schema.OpcodeGreater:
		return g.op_less(buf, node, evalFlags)
	case schema.OpcodeGreaterEqual:
		return g.op_greaterEqual(buf, node, evalFlags)
	case schema.OpcodeAnd:
		return g.op_and(buf, node, evalFlags)
	case schema.OpcodeOr:
		return g.op_or(buf, node, evalFlags)

	case schema.OpcodeAdd:
		return g.op_add(buf, node, evalFlags)
	case schema.OpcodeSub:
		return g.op_sub(buf, node, evalFlags)
	case schema.OpcodeMul:
		return g.op_mul(buf, node, evalFlags)
	case schema.OpcodeDiv:
		return g.op_div(buf, node, evalFlags)
	case schema.OpcodeMod:
		return g.op_mod(buf, node, evalFlags)
	case schema.OpcodeAssign:
		return g.op_assign(buf, node, evalFlags)
	case schema.OpcodeAddAssign:
		return g.op_addAssign(buf, node, evalFlags)
	case schema.OpcodeSubAssign:
		return g.op_subAssign(buf, node, evalFlags)
	case schema.OpcodeMulAssign:
		return g.op_mulAssign(buf, node, evalFlags)
	case schema.OpcodeDivAssign:
		return g.op_divAssign(buf, node, evalFlags)
	case schema.OpcodeModAssign:
		return g.op_modAssign(buf, node, evalFlags)

	case schema.OpcodeBitAndAssign:
		return g.op_bitAndAssign(buf, node, evalFlags)
	case schema.OpcodeBitOrAssign:
		return g.op_bitOrAssign(buf, node, evalFlags)
	case schema.OpcodeBitXorAssign:
		return g.op_bitXorAssign(buf, node, evalFlags)
	case schema.OpcodeBitClearAssign:
		return g.op_bitClearAssign(buf, node, evalFlags)
	case schema.OpcodeLeftShiftAssign:
		return g.op_leftShiftAssign(buf, node, evalFlags)
	case schema.OpcodeRightShiftAssign:
		return g.op_rightShiftAssign(buf, node, evalFlags)
	case schema.OpcodeBitAnd:
		return g.op_bitAnd(buf, node, evalFlags)
	case schema.OpcodeBitOr:
		return g.op_bitOr(buf, node, evalFlags)
	case schema.OpcodeBitXor:
		return g.op_bitXor(buf, node, evalFlags)
	case schema.OpcodeBitClear:
		return g.op_bitClear(buf, node, evalFlags)
	case schema.OpcodeLeftShift:
		return g.op_leftShift(buf, node, evalFlags)
	case schema.OpcodeRightShift:
		return g.op_rightShift(buf, node, evalFlags)
	}

	return fmt.Errorf("invalid opcode on node with opcode of %s", opcode)
}

func (g *Generator) EvalUnary(buf *bytes.Buffer, opcode schema.Opcode, node *program.UnaryNode, evalFlags EvalFlags) error {
	switch opcode {
	case schema.OpcodeNot:
		return g.op_not(buf, node, evalFlags)
	case schema.OpcodeDefer:
		return g.op_defer(buf, node, evalFlags)
	case schema.OpcodeGoRoutine:
		return g.op_goRoutine(buf, node, evalFlags)

	case schema.OpcodeInc:
		return g.op_inc(buf, node, evalFlags)
	case schema.OpcodeDec:
		return g.op_dec(buf, node, evalFlags)
	case schema.OpcodeAddrOf:
		return g.op_addrOf(buf, node, evalFlags)
	case schema.OpcodeDeref:
		return g.op_deref(buf, node, evalFlags)
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

	default:
		return nil, fmt.Errorf("unknown type kind: %d", t.TypeType())
	}
}
