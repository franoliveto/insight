// Copyright 2025 Francisco Oliveto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/franoliveto/insight"
)

func main() {
	log.SetFlags(0)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s system name\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}
	system := flag.Arg(0)
	name := flag.Arg(1)

	client := insight.NewClient("", nil)

	p, err := client.GetPackage(system, name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("System: %s\nPackage: %s\n", p.PackageKey.System, p.PackageKey.Name)
	fmt.Println("Versions:")
	for _, v := range p.Versions {
		fmt.Printf("%s\n", v.VersionKey.Version)
	}
}
