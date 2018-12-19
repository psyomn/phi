package main

import (
	"flag"
	"fmt"
	"os"

	"local/psyomn/phi"
)

var (
	inDir  = ""
	outDir = ""
)

func init() {
	flag.StringVar(&inDir, "dir", inDir, "input directory")
	flag.StringVar(&outDir, "output", outDir, "output directory")
}

// Print usage
func usage() {
	fmt.Println("phi-store -dir <in-dir> -output <out-dir>")
}

func main() {
	flag.Parse()

	if inDir == "" || outDir == "" {
		usage()
		os.Exit(1)
	}

	phi.SortByModTime(inDir, outDir)
}
