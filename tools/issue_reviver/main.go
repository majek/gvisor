// Copyright 2019 The gVisor Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"os"

	"gvisor.dev/gvisor/tools/issue_reviver/github"
	"gvisor.dev/gvisor/tools/issue_reviver/reviver"
)

var (
	owner string
	repo  string
	token string
	path  string
)

// Keep the options simple for now. Supports only a single path and repo.
func init() {
	flag.StringVar(&owner, "owner", "google", "Github project org/owner to look for issues")
	flag.StringVar(&repo, "repo", "gvisor", "Github repo to look for issues")
	flag.StringVar(&token, "oauth-token", "", "OAUTH token to be used as credential to github")
	flag.StringVar(&path, "path", "", "Path to scan for TODOs")
}

func main() {
	flag.Parse()

	// Check for mandatory parameters.
	if len(owner) == 0 {
		fmt.Println("missing --owner option.")
		flag.Usage()
		os.Exit(1)
	}
	if len(repo) == 0 {
		fmt.Println("missing --repo option.")
		flag.Usage()
		os.Exit(1)
	}
	if len(path) == 0 {
		fmt.Println("missing --path option.")
		flag.Usage()
		os.Exit(1)
	}

	bugger, err := github.NewBugger(token, owner, repo)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	rev := reviver.New([]string{path}, []reviver.Bugger{bugger})
	if errs := rev.Run(); len(errs) > 0 {
		fmt.Fprintf(os.Stderr, "Encountered %d errors:\n", len(errs))
		for _, err := range errs {
			fmt.Fprintf(os.Stderr, "\t%v\n", err)
		}
		os.Exit(1)
	}
}
