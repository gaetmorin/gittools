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

package main

import (
	"bufio"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// TODO(@gaetmorin): this part is horribly project specific.

var re = regexp.MustCompile(`git\s+clone\s+(?P<url>[^\s]+)\s+-b\s+(?P<branch>[^\s]+)\s+(?P<dir>[^\s]+)`)

func findSources() ([]Source, error) {
	return sourcesFromMakefile(detectDependencyfile())
}

func detectDependencyfile() (makefile, target string) {
	for _, x := range []struct {
		makefile, target string
	}{
		{"Dependencyfile", "all"},
		{"Depfile.mk", "clone"},
		{"Makefile", "sources"},
	} {
		if _, err := os.Stat(x.makefile); err == nil {
			return x.makefile, x.target
		}
	}
	return "", ""
}

func sourcesFromMakefile(makefile, target string) ([]Source, error) {
	cmd := exec.Command("make", "-Bnf", makefile, target)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(stdout)
	var sources []Source
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !re.MatchString(line) {
			continue
		}
		sources = append(sources, Source{
			URL:    re.ReplaceAllString(line, "$url"),
			Branch: re.ReplaceAllString(line, "$branch"),
			Dir:    re.ReplaceAllString(line, "$dir"),
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if err := cmd.Wait(); err != nil {
		return nil, err
	}
	return sources, nil
}
