package golang

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestCompile(t *testing.T) {
	file, err := os.ReadFile("nodes.opt")
	if err != nil {
		panic(err)
	}

	now := time.Now()

	gf, err := Compile(&file)
	if err != nil {
		panic(err)
	}

	log.Printf("Time elapse: %dms", time.Since(now).Milliseconds())

	for _, g := range gf {
		log.Println(string(*g.Content))
	}
	if len(os.Args) > 0 && os.Args[len(os.Args)-1] == "export-as-files" {
		for _, g := range gf {
			d := "./exports/" + g.Path + ".go"
			err := g.Write(d)
			if err != nil {
				panic(err)
			}
			log.Printf("Wrote to %s", d)
		}
	}
}
