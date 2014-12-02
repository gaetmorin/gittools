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

import "bytes"

type Task struct {
	Source
	Output bytes.Buffer
	Error  error
}

type Source struct {
	Dir    string
	Branch string
	URL    string
}

type sourcesByDir []Source

func (s sourcesByDir) Len() int           { return len(s) }
func (s sourcesByDir) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sourcesByDir) Less(i, j int) bool { return s[i].Dir < s[j].Dir }
