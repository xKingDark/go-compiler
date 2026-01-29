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

func (g *Generator) op_func(node *program.IndexedNode, flags EvalFlags) ([]byte, error) {
	length := node.FieldsLength()

	buf := new(bytes.Buffer)

	var funcId = make([]byte, 32)
	var funcType *program.FuncType

	params := new(bytes.Buffer)
	params.Grow(length * paramsGrowthModifer)

	body := new(bytes.Buffer)
	body.Grow(length * bodyGrowthModifer)

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

			if field.Flags()&uint32(schema.ValueFlagFuncParam) != 0 {
				if params.Len() > 0 {
					params.WriteByte(',')
					params.WriteByte(' ')
				}
				params.Write(out)
			} else if field.Flags()&uint32(schema.ValueFlagFuncBody) != 0 {
				body.WriteByte('	')
				body.Write(out)
				body.WriteByte('\n')
			}
		} else if field.Flags()&uint32(schema.ValueFlagFuncMeta) != 0 {
			continue
			def, ok := g.LookUpType(field.Type())
			if !ok {
				return nil, fmt.Errorf("type with id, %d is undefined", field.Type())
			}

			funcId = def.Id()

			if v := EvalType(def); v != nil {
				if v, ok := v.(program.FuncType); ok {
					funcType = &v
				}
			}
		}
	}
	resultsLen := funcType.ResultsLength()

	var declarationLength = TokenFunc.Len() + (resultsLen * resultsGrowthModifer) + 8
	buf.Grow(params.Len() + body.Len() + declarationLength)

	buf.Write(TokenFunc.Bytes())
	buf.WriteByte(' ')
	buf.Write(funcId)
	buf.WriteByte('(')
	buf.Write(params.Bytes())
	buf.WriteByte(')')

	if resultsLen > 0 {
		buf.WriteByte(' ')
		if resultsLen > 1 {
			buf.WriteByte('(')
		}

		for i := range resultsLen {
			var p program.Pair
			funcType.Results(&p, i)

			var ok bool
			var name = make([]byte, 32)
			g.StrLookupMutex.Lock()
			name, ok = g.LookUpStr(p.Key())
			g.StrLookupMutex.Unlock()
			if !ok {
				return nil, fmt.Errorf("string with id %d is undefined", node.Id())
			}
			var _type = make([]byte, 45)
			g.StrLookupMutex.Lock()
			_type, ok = g.LookUpStr(p.Value())
			g.StrLookupMutex.Unlock()
			if !ok {
				return nil, fmt.Errorf("string with id %d is undefined", node.Id())
			}

			if i > 0 {
				buf.WriteByte(',')
				buf.WriteByte(' ')
			}
			if len(name) > 0 {
				buf.Write(name)
				buf.WriteByte(' ')
			}

			buf.Write(_type)
		}

		if resultsLen > 1 {
			buf.WriteByte(')')
		}
	}
	buf.WriteByte(' ')
	buf.WriteByte('{')
	buf.WriteByte('\n')
	buf.Write(body.Bytes())
	buf.WriteByte('}')

	return buf.Bytes(), nil
}
