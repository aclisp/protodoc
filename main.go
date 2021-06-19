package main

import (
	"log"
	"os"
	"path/filepath"

	pp "github.com/yoheimuta/go-protoparser/v4"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Please specify a proto file.")
	}

	protoFileName := os.Args[1]
	reader, err := os.Open(protoFileName)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer reader.Close()

	parsed, err := pp.Parse(reader,
		pp.WithDebug(false),
		pp.WithPermissive(true),
		pp.WithFilename(filepath.Base(protoFileName)))
	if err != nil {
		log.Fatalf("%v", err)
	}

	protoFile := extract(parsed)
	protoFile.output()
}
