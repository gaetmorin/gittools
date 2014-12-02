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
	"flag"
	"os/exec"
	"strings"
	"sync"
)

func run(src []Source, jobs int) []Task {
	if jobs <= 0 {
		jobs = len(src)
	}
	tasks := prepare(src)
	queue, done := spawn(jobs)
	dispatch(queue, tasks)
	<-done
	return tasks
}

func prepare(src []Source) []Task {
	var tasks []Task
	for _, s := range src {
		tasks = append(tasks, Task{Source: s})
	}
	return tasks
}

func dispatch(queue chan<- func(), tasks []Task) {
	for i := range tasks {
		t := &tasks[i]
		cmd := exec.Command(*command, flag.Args()...)
		cmd.Dir = t.Dir
		cmd.Stdout = &t.Output
		cmd.Stderr = &t.Output
		queue <- func() {
			t.Error = cmd.Run()
			progress <- true
		}
	}
	close(queue)
}

func spawn(n int) (queue chan<- func(), done <-chan bool) {
	var wg sync.WaitGroup
	q := make(chan func())
	d := make(chan bool)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go work(&wg, q)
	}
	go wait(&wg, d)
	return q, d
}

func work(wg *sync.WaitGroup, queue <-chan func()) {
	defer wg.Done()
	for f := range queue {
		f()
	}
}

func wait(wg *sync.WaitGroup, done chan<- bool) {
	wg.Wait()
	close(done)
}

func trimcr(s string) string {
	i := strings.LastIndex(s, "\r")
	if i < 0 {
		return s
	}
	return s[i+1:]
}
