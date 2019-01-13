/*
Copyright 2019 Simon Symeonidis (psyomn)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/psyomn/phi"
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
