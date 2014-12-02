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
	"fmt"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	progress     = make(chan bool)
	progressDone = make(chan bool)
)

func indicateProgress(n int) {
	msg := ""
	running := true
	display := false
	i := 0
	delay := time.After(1 * time.Second)
	for running {
		if display {
			msg = fmt.Sprintf("%d / %d", i, n)
			fmt.Fprint(os.Stderr, msg)
		}
		select {
		case _, running = <-progress:
			i++
		case <-delay:
			display = true
		}
		if display {
			fmt.Fprintf(os.Stderr, "\r%s\r", strings.Repeat(" ", utf8.RuneCountInString(msg)))
		}
	}
	close(progressDone)
}

func waitProgress() {
	close(progress)
	<-progressDone
}
