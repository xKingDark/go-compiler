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
	unionTable := new(fb.Table)

	if node.Node(unionTable) {
		nodeType := node.NodeType()
		out := []byte{} // WARN - Should be bytes.Buffer
		switch nodeType {
		case program.NodeUnionIndexedNode:
			type1 := new(program.IndexedNode)
			type1.Init(unionTable.Bytes, unionTable.Pos)
			var err error
			out, err = g.EvalIndexed(schema.Opcode(node.Opcode()), type1, evalFlags)
			if err != nil {
				return nil, err
			}
		case program.NodeUnionBinaryNode:
			type2 := new(program.BinaryNode)
			type2.Init(unionTable.Bytes, unionTable.Pos)
			var err error
			out, err = g.EvalBinary(schema.Opcode(node.Opcode()), type2, evalFlags)
			if err != nil {
				return nil, err
			}
		case program.NodeUnionUnaryNode:
			type3 := new(program.UnaryNode)
			type3.Init(unionTable.Bytes, unionTable.Pos)
			var err error
			out, err = g.EvalUnary(schema.Opcode(node.Opcode()), type3, evalFlags)
			if err != nil {
				return nil, err
			}
		case program.NodeUnionNONE:
			return nil, fmt.Errorf("failed to determine node type of node: %d", node.Id())
		}
		//log.Printf("node id: %d opcode: %s", node.Id(), schema.Opcode(node.Opcode()))
		return out, nil
	}
	return nil, fmt.Errorf("failed to access union of node: %d", node.Id())
}

func (g *Generator) EvalIndexed(opcode schema.Opcode, node *program.IndexedNode, evalFlags EvalFlags) ([]byte, error) {
	switch opcode {
	case schema.OpcodePackage:
		return g.op_package(node, evalFlags)
	case schema.OpcodeImport:
		return g.op_import(node, evalFlags)
	case schema.OpcodeConst:
		return g.op_const(node, evalFlags)
	case schema.OpcodeVar:
		return g.op_var(node, evalFlags)
	case schema.OpcodeIf:
		return g.op_if(node, evalFlags)
	}
	return nil, fmt.Errorf("invalid opcode on node with opcode of %s", opcode)
}

func (g *Generator) EvalBinary(opcode schema.Opcode, node *program.BinaryNode, evalFlags EvalFlags) ([]byte, error) {
	switch opcode {
	case schema.OpcodeImportValue:
		return g.op_importValue(node, evalFlags)
	case schema.OpcodeConstValue:
		return g.op_constValue(node, evalFlags)
	case schema.OpcodeVarValue:
		return g.op_varValue(node, evalFlags)
	case schema.OpcodeEqual:
		return g.op_equal(node, evalFlags)
	}
	return nil, fmt.Errorf("invalid opcode on node with opcode of %s", opcode)
}

func (g *Generator) EvalUnary(opcode schema.Opcode, node *program.UnaryNode, evalFlags EvalFlags) ([]byte, error) {
	switch opcode {
	}
	return nil, fmt.Errorf("invalid opcode on node with opcode of %s", opcode)
}

func (g *Generator) evalValue(buf *bytes.Buffer, nodeValue *program.NodeValue) error {
	if nodeValue.Flags()&uint32(schema.ValueFlagPointer) != 0 {
		node := g.GetNode(nodeValue.Value())
		if node == nil {
			return fmt.Errorf("attempt to access undefined node: %d", nodeValue.Value())
		}
		out, err := g.Eval(node, 0)
		if err != nil {
			return err
		}
		buf.Write(out)
	} else {
		g.StrLookupMutex.Lock()
		value, ok := g.LookUpStr(uint32(nodeValue.Value()))
		g.StrLookupMutex.Unlock()
		if !ok {
			return fmt.Errorf("string with id %d is undefined", nodeValue.Value())
		}
		if nodeValue.Flags()&uint32(schema.ValueFlagQuotation) != 0 {
			buf.Write(JoinBytes(TokenQuotation.Bytes(), value, TokenQuotation.Bytes()))
		} else {
			buf.Write(value)
		}
	}
	return nil
}

func EvalType(t *program.TypeDef) any {
	unionTable := new(fb.Table)

	if t.Type(unionTable) {
		goType := t.TypeType()
		switch goType {
		case program.TypePointerType:
			ptr := new(program.PointerType)
			ptr.Init(unionTable.Bytes, unionTable.Pos)
			return ptr
		}
	}
	return nil
}
