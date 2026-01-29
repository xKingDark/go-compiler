package golang

import (
	"log"
	"sync"

	schema "github.com/Opticode-Project/go-compiler/golang"
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

	const maxRoutines = 5 // maximum allowed concurrent goroutines

	sem := make(chan struct{}, maxRoutines) // semaphore channel
	var wg sync.WaitGroup

	for _, i := range gen.modulePath[0] {
		if i == -1 {
			break
		}
		sem <- struct{}{} // acquire a "slot" (pauses if full)
		wg.Add(1)

		go func(index int) {
			defer wg.Done()
			defer func() { <-sem }() // release the slot when done

			var node program.Node
			app.Nodes(&node, index)

			buf, err := gen.Eval(&node, 0)
			if err != nil {
				log.Println(err)
			}
			gen.nodesMutex.Lock()
			gen.Write(node.Id(), schema.NodeFlag(node.Flags()), 0, &buf)
			gen.nodesMutex.Unlock()
		}(i)
	}

	wg.Wait()

	// for i := range nodesLength {
	// 	log.Printf("index: %d", i)
	// 	var node program.Node
	// 	app.Nodes(&node, i)
	// 	gen.Eval(&node)
	// }

	//gen.PrintNodes()
	return gen.Export(gen.modulePath[0])
}
