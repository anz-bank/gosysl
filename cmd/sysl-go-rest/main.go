package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/anz-bank/gosysl"
	"github.com/anz-bank/gosysl/pb"
	"github.com/golang/protobuf/proto"
)

func main() {
	fmt.Println("sysl-go-rest started")
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		log.Fatal("Usage: sysl-go-rest <INPUT.pb> <OUTPUT_DIR>")
	}
	data, err := ioutil.ReadFile(args[0])
	if err != nil {
		log.Fatal(err)
	}
	module := &pb.Module{}
	err = proto.Unmarshal(data, module)
	if err != nil {
		log.Fatal("Unmarshaling error: ", err)
	}
	outDir := args[1]
	os.MkdirAll(outDir, os.ModePerm)
	if _, err := os.Stat(outDir); err != nil {
		log.Fatal("Cannot access output directory, error: ", err)
	}
	pkg := gosysl.GetPackage(outDir)
	result, err := gosysl.Generate(module, pkg)
	if err != nil {
		log.Fatal("Code generation error: ", err)
	}
	fmt.Printf("Main.Result: %v\n", result)
	fmt.Printf("Finished successfully\n")

}
