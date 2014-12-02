/*
Copyright 2014 Gaëtan Morin

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

// Dep executes commands on many repositories that compose a project. For now,
// the layout is very specific: it is extracted from a Makefile that contains
// "git clone" commands.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
)

const defaultCmd = "git"

var (
	command = flag.String("c", defaultCmd, "name of the command to run")
	jobs    = flag.Int("j", 100, "number of parallel jobs")
	short   = flag.Bool("s", false, "print in short, machine-friendly format")
)

var root = Source{
	Dir:    ".",
	Branch: "HEAD",
}

func main() {
	flag.Parse()

	sources, err := findSources()
	sort.Sort(sourcesByDir(sources))
	if err != nil {
		fmt.Println("no source")
		os.Exit(1)
	}

	if flag.NArg() == 0 && *command == defaultCmd {
		list(sources)
	} else {
		dep(sources)
	}
}

func dep(sources []Source) {
	var tasks []Task
	sources = append([]Source{root}, sources...)
	go indicateProgress(len(sources))
	tasks = append(tasks, run(sources, *jobs)...)
	waitProgress()

	if *short {
		printShort(tasks)
	} else {
		print(tasks)
	}
}

func print(tasks []Task) {
	for _, t := range tasks {
		if t.Error == nil && t.Output.Len() == 0 {
			continue
		}
		if t.Error != nil {
			fmt.Printf("%s %s\n", t.Dir, t.Error)
		} else {
			fmt.Printf("%s\n", t.Dir)
		}
		s := bufio.NewScanner(&t.Output)
		for s.Scan() {
			fmt.Printf("\t%s\n", trimcr(s.Text()))
		}
		// TODO check s.Err()
	}
}

func printShort(tasks []Task) {
	w := newTabWriter()
	for _, t := range tasks {
		if t.Error != nil {
			fmt.Fprintf(w, "%s\t%s\n", t.Dir, t.Error)
		}
		s := bufio.NewScanner(&t.Output)
		for s.Scan() {
			fmt.Fprintf(w, "%s\t%s\n", t.Dir, trimcr(s.Text()))
		}
		// TODO check s.Err()
	}
	w.Flush()
}

func list(sources []Source) {
	w := newTabWriter()
	for _, s := range sources {
		fmt.Fprintf(w, "%s\t%s\t%s\n", s.Dir, s.Branch, s.URL)
	}
	w.Flush()
}

func newTabWriter() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 10, 0, 1, ' ', 0)
}
