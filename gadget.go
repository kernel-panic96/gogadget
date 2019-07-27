package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/kernel-panic96/gogadget/find"
)

func main() {
	filename := flag.String("f", "", "The file which should be searched for sub tests.")
	debug := flag.Bool("debug", false, "will print the AST for easier verification")
	_ = debug
	flag.Parse()

	if *filename == "" {
		log.Fatalf("file must be specified with the -f flag")
	}
	tests, err := find.AllCalls(find.Call{
		ImportPath: "testing",
		TypeName:   "T",
		MethodName: "Run",
	}).InFile(*filename)

	if err != nil {
		log.Fatal(err)
	}
	for _, t := range tests {
		fmt.Println(t)
	}
}
