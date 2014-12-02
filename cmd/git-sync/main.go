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

// git sync keeps up with upstream changes
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gaetmorin/gittools/git"
)

var (
	dryRun    = flag.Bool("n", false, "only print the command(s), don't run them")
	allowPush = flag.Bool("push", false, "allow pushs to the remote repository")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	var p plan
	switch flag.NArg() {
	case 0:
		p.addSync(strings.TrimPrefix(p.Head(), "refs/heads/"))
	case 1:
		p.addSync(flag.Arg(0))
	default:
		fail("wrong number of arguments")
	}
	if *dryRun {
		p.print()
	} else {
		p.run()
	}
}

func run(c *exec.Cmd) {
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		fail("%s", err)
	}
}

type plan struct {
	git.Repo
	cmds []*exec.Cmd
}

func sync(name string) *plan {
	var p plan
	p.addSync(name)
	return &p
}

func (p *plan) print() {
	for _, c := range p.cmds {
		fmt.Println(strings.Join(c.Args, " "))
	}
}

func (p *plan) run() {
	for _, c := range p.cmds {
		run(c)
	}
}

func (p *plan) addGit(args ...string) {
	p.cmds = append(p.cmds, p.Git(args...))
}

func (p *plan) addSync(name string) {
	if p.Empty() {
		fail("empty repository")
	}

	if tag := "refs/tags/" + name; p.HasRef(tag) {
		p.addGit("checkout", tag)
		return
	}

	if branch := "refs/heads/" + name; p.HasRef(branch) {
		if p.Head() != branch {
			p.addGit("checkout", name)
		}
		upstream := p.UpstreamOf(branch)
		if upstream == "" {
			fail("%s: no upstream branch", branch)
		}
		if !p.HasRef(upstream) {
			fail("%s: upstream branch '%s' seems to have been deleted", branch, upstream)
		}
		parts := strings.SplitN(upstream, "/", 4)
		if len(parts) != 4 || parts[0] != "refs" || parts[1] != "remotes" || parts[3] != name {
			fail("%s: upstream branch '%s' is not of the form 'refs/remotes/<remote>/%s'", branch, upstream, name)
		}
		remote := parts[2]

		ahead, behind := p.AheadBehind(branch, upstream)
		switch {
		case ahead == 0 && behind > 0:
			p.addGit("merge", "--quiet", "--ff-only", upstream)
		case ahead > 0 && behind > 0:
			p.addGit("rebase", "--quiet", "-p", upstream)
		case ahead > 0 && behind == 0 && *allowPush:
			p.addGit("push", "--quiet", remote, "HEAD:"+branch)
		}
		return
	}

	for _, remote := range p.Remotes() {
		upstream := fmt.Sprintf("refs/heads/%s/%s", remote, name)
		if p.HasRef(upstream) {
			p.addGit("checkout", upstream, "-b", name)
			return
		}
	}

	fail("%s: no branch or tag found", name)
}

func fail(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s+"\n", args...)
	os.Exit(1)
}

func usage() {
	fmt.Println(`usage: git sync [options] [target]

Sync with 'target', which may be a branch (local or remote) or tag name. By
default, 'target' is the name of the current branch.

If HEAD is not aligned with 'target', switch to it first, creating a remote
tracking branch if necessary.

If target is a branch, keep up with upstream by:
  * rebasing on top of it if there are local, unpushed commits (with the -p
    option to preserve local merges);
  * merging if there are no local commits (fast-forward merge).

If the local branch is strictly ahead of upstream and the -push option is
specified, push local commits.

options:`)
	flag.PrintDefaults()
	os.Exit(1)
}
