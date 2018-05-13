package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/anz-bank/sysl-go-rest/pb"
	"github.com/golang/protobuf/proto"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		log.Fatal("Missing input file")
	}
	data, err := ioutil.ReadFile(args[0])
	if err != nil {
		log.Fatal(err)
	}
	module := &pb.Module{}
	err = proto.Unmarshal(data, module)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}
	fmt.Printf("Module: %v\n\n", module)
	for _, app := range module.GetApps() {
		for _, ep := range app.GetEndpoints() {
			fmt.Printf("app: %s, ep: %s, params: %v\n", app.GetName(), ep.GetName(), ep.GetParam())
		}
	}

	fmt.Printf("Finished successfully\n")

}
