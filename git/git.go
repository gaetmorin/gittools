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

// git command-line wrappers
//
// THIS INTERFACE IS EXPERIMENTAL AND IS LIKELY TO CHANGE.
package git

import (
	"fmt"
	"os/exec"
	"strings"
)

type Repo string

func (r Repo) HasUntracked() bool {
	return test(r.Git("ls-files", "--other", "--error-unmatch", "--exclude-standard", "."))
}

func (r Repo) Clean() bool {
	return test(r.Git("diff", "--quiet", "HEAD"))
}

func (r Repo) WorkdirClean() bool {
	return test(r.Git("diff", "--quiet"))
}

func (r Repo) IndexClean() bool {
	return test(r.Git("diff", "--cached", "--quiet", "HEAD"))
}

func (r Repo) Empty() bool {
	return !test(r.Git("rev-parse", "HEAD"))
}

func (r Repo) HasRef(ref string) bool {
	return test(r.Git("show-ref", "--verify", "--quiet", ref))
}

func (r Repo) UpstreamOf(ref string) string {
	return subst(r.Git("for-each-ref", "--format=%(upstream)", ref))
}

func (r Repo) Abbrev(ref string) string {
	return subst(r.Git("rev-parse", "--abbrev-ref", ref))
}

func (r Repo) Ref(name string) string {
	return subst(r.Git("rev-parse", "--symbolic-full-name", name))
}

func (r Repo) ShortSHA1(ref string) string {
	return subst(r.Git("rev-parse", "--short", ref))
}

func (r Repo) MergeBaseSHA1(left, right string) string {
	return subst(r.Git("merge-base", left, right))
}

func (r Repo) Head() string {
	return r.Ref("HEAD")
}

func (r Repo) AheadBehind(local, upstream string) (ahead, behind int) {
	rangespec := local + "..." + upstream
	counts := subst(r.Git("rev-list", "--count", "--left-right", rangespec))
	_, err := fmt.Sscanln(counts, &ahead, &behind)
	if err != nil {
		return -1, -1
	}
	return ahead, behind
}

func (r Repo) Remotes() []string {
	return strings.Fields(subst(r.Git("remote")))
}

func (r Repo) Git(args ...string) *exec.Cmd {
	cmd := exec.Command("git", args...)
	cmd.Dir = string(r)
	return cmd
}

func test(cmd *exec.Cmd) bool {
	err := cmd.Run()
	return err == nil
}

func subst(cmd *exec.Cmd) string {
	b, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimRight(string(b), "\n\r")
}
