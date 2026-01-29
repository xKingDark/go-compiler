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
}
