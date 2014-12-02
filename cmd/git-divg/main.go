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

// git divg shows how two git references diverge
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/gaetmorin/gittools/git"
)

var (
	update        = flag.Bool("u", false, "fetch first")
	showUntracked = flag.Bool("a", false, "show untracked files")
)

// repo is the current directory
var repo git.Repo

func main() {
	flag.Usage = usage
	flag.Parse()

	if *update {
		fgit(os.Stdout, "fetch")
	}

	switch flag.NArg() {
	case 0:
		divgUpstream(repo.Head())
	case 1:
		divgUpstream(repo.Ref(flag.Arg(0)))
	case 2:
		l, r := flag.Arg(0), flag.Arg(1)
		divg(repo.Ref(l), repo.Ref(r))
	default:
		fail("wrong number of arguments\n")
	}
}

func divgUpstream(ref string) {
	checkRef(ref)
	remote := repo.UpstreamOf(ref)
	if remote == "" {
		fail("%s: no upstream branch\n", ref)
	}
	if !repo.HasRef(remote) {
		fail("%s: upstream branch '%s' seems to have been deleted\n", ref, remote)
	}
	divg(ref, remote)
}

func divg(local, remote string) {
	checkRef(local)
	checkRef(remote)

	ahead, behind := repo.AheadBehind(local, remote)
	if ahead < 0 || behind < 0 {
		fail("error computing ahead/behind commit counts")
	}

	rangespec := local + "..." + remote
	abbrev := repo.Abbrev(remote)
	var b bytes.Buffer
	if ahead > 0 {
		fmt.Fprintf(&b, "%d commit%s ahead of %s:\n", ahead, s(ahead), abbrev)
		fgit(&b, "log", "--pretty=%m  %h %s", "--cherry-mark", "--no-merges", "--left-only", rangespec)
	}
	if behind > 0 {
		fmt.Fprintf(&b, "%d commit%s behind %s:\n", behind, s(behind), abbrev)
		fgit(&b, "log", "--pretty=%m  %h %s", "--cherry-mark", "--no-merges", "--right-only", rangespec)
	}
	if local == repo.Head() && (!repo.Clean() || (*showUntracked && repo.HasUntracked())) {
		fmt.Fprintln(&b, "status:")
		fgit(&b, "status", "--short")
	}

	fmt.Print(b.String())
}

func checkRef(ref string) {
	if !repo.HasRef(ref) {
		fail("%s: no such reference\n", ref)
	}
}

func s(x int) string {
	if x > 1 {
		return "s"
	}
	return ""
}

func fgit(w io.Writer, args ...string) {
	cmd := repo.Git(args...)
	cmd.Stdout = w
	err := cmd.Run()
	failOn(err, cmd)
}

func failOn(err error, cmd *exec.Cmd) {
	if err != nil {
		fail("error running %s: %s\n", strings.Join(cmd.Args, " "), err)
	}
}

func fail(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s, args...)
	os.Exit(1)
}

func usage() {
	fmt.Println(`usage: git divg [options] [local [base]]

Show the commits in 'local' that are ahead of or behind 'base'. 'local' and
'base' can be commit ids or references. By default, 'base' is the upstream of
'local', and 'local' is the current branch.

If 'local' is the current branch and there are uncommited changes, show a short
status.

options:`)
	flag.PrintDefaults()
	os.Exit(1)
}
