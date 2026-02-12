package golang

import (
	program "github.com/Opticode-Project/go-compiler/program"
)

//go:generate bash ./generate-golang.sh
//go:generate bash ./generate-program.sh

// Compiles buffer into go files.
func Compile(buf *[]byte) ([]*GoFile, error) {
	app := program.GetRootAsApp(*buf, 0)

	gen := NewGenerator(app, buf)

	for index := range gen.modulePath {
		gen.CompileModule(index)
	}

	return gen.Export("main", gen.modulePath[0])
}
