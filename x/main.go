// Copyright 2025 Francisco Oliveto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/franoliveto/insights"
)

func doVersion(c *insights.Client, system, name, version string) error {
	var v *insights.Version
	v, err := c.GetVersion(system, name, version)
	if err != nil {
		return err
	}
	fmt.Println(*v)
	return nil
}

func doPackage(c *insights.Client, system, name string) error {
	var p *insights.Package
	p, err := c.GetPackage(system, name)
	if err != nil {
		return err
	}
	fmt.Println(*p)
	return nil
}

func main() {
	log.SetFlags(0)
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "Usage: x command [args]")
		os.Exit(1)
	}

	client := insights.NewClient()

	switch cmd := flag.Arg(0); cmd {
	case "package":
		if flag.NArg() < 3 {
			fmt.Fprintln(os.Stderr, "usage: x package system name")
			os.Exit(1)
		}
		system := flag.Arg(1)
		name := flag.Arg(2)
		if err := doPackage(client, system, name); err != nil {
			log.Fatal(err)
		}
	case "version":
		if flag.NArg() < 4 {
			fmt.Fprintln(os.Stderr, "usage: x version system name version")
			os.Exit(1)
		}
		system := flag.Arg(1)
		name := flag.Arg(2)
		version := flag.Arg(3)
		if err := doVersion(client, system, name, version); err != nil {
			log.Fatal(err)
		}
	case "dependencies":
		if flag.NArg() < 4 {
			fmt.Fprintln(os.Stderr, "usage: x dependencies system name version")
			os.Exit(1)
		}
		system := flag.Arg(1)
		name := flag.Arg(2)
		version := flag.Arg(3)
		d, err := client.GetDependencies(insights.VersionKey{System: system, Name: name, Version: version})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(*d)
	case "project":
		if flag.NArg() < 2 {
			fmt.Fprintln(os.Stderr, "usage: x project id")
			os.Exit(1)
		}
		p, err := client.GetProject(flag.Arg(1))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(*p)
	}

}
