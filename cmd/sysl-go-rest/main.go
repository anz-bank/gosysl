package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

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

	s := reflect.ValueOf(&result).Elem()
	for i := 0; i < s.NumField(); i++ {
		content := s.Field(i).Interface().(string)
		basename := strings.ToLower(s.Type().Field(i).Name)
		filename := filepath.Join(outDir, basename+".go")
		err = ioutil.WriteFile(filename, []byte(content), 0644)
		if err != nil {
			log.Fatal("Cannot write file ", filename)
		}
	}
	fmt.Printf("Finished successfully\n")
}
