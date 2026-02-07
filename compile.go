package golang

import (
	program "github.com/Opticode-Project/go-compiler/program"
)

//go:generate bash ./generate-golang.sh
//go:generate bash ./generate-program.sh

// Compiles buffer into go files.
func Compile(buf *[]byte) ([]*GoFile, error) {
	app := program.GetRootAsApp(*buf, 0)

	// for i := range app.LutLength() {
	// 	var v program.StringEntry
	// 	app.Lut(&v, i)

	// 	log.Println(v.Key(), string(v.Value()))

	// 	ok := app.LutByKey(&v, v.Key())
	// 	if ok {
	// 		log.Println("yeah it ok")
	// 	}
	// }
	gen := NewGenerator(app, buf)

	for index := range gen.modulePath {
		gen.CompileModule(index)
	}

	// for i := range nodesLength {
	// 	log.Printf("index: %d", i)
	// 	var node program.Node
	// 	app.Nodes(&node, i)
	// 	gen.Eval(&node)
	// }

	//gen.PrintNodes()
	return gen.Export(gen.modulePath[0])
}
