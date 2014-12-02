Work in progress

Intro
=====

Gittools is a collection of tools for working with Git.

 -  dep: execute commands on multiple repositories of a project
 -  git-divg: see the divergences from upstream
 -  git-sync: synchronize with upstream

This is how a typical fetch / integrate / push workflow unfolds:

    % git divg -u

Fetches changes and shows how the local repository diverges from upstream
(commits ahead/behind, local changes). Without -u, there is no fetch and
therefore no interaction with the network.

    % git sync

Rebases over upstream changes. The actual command depends on the context. Use
git sync -n to print the command instead of executing it. There is no output if
everything goes well.

    % git divg

Divg shows the also shows the commits ahead of upstream, so use it to see what
would be pushed, and finally validate with:

    % git sync -push

As is, git-divg and git-sync provide little value over the bare Git commands.
The dep tool is what makes them useful. Dep executes commands on multiple
repositories (dependencies) that compose a project. In the previous example,
replacing git by dep turns the single repository workflow into a multiple
repository workflow. `dep divg` shows the dependencies that diverges, and `dep
sync` synchronizes only these dependencies.

For now, a dep project is a directory containing a special Makefile, but this
is subject to change.

Install
=======

    % go get github.com/gaetmorin/gittools/cmd/...

Ensure $GOPATH/bin is added to your $PATH.


