package golang

import (
	"fmt"

	program "github.com/Opticode-Project/go-compiler/program"
)

func (g *Generator) op_package(node *program.IndexedNode, flags EvalFlags) ([]byte, error) {
	g.StrLookupMutex.Lock()
	id, ok := g.LookUpStr(node.Id())
	g.StrLookupMutex.Unlock()
	if !ok {
		return nil, fmt.Errorf("string with id %d is undefined", node.Id())
	}
	return JoinBytes(TokenPackage.Bytes(), TokenSpace.Bytes(), id), nil
}
