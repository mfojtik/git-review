// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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

package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/mfojtik/git-review/pkg/util"
	"github.com/spf13/cobra"
)

type GitReviewOptions struct {
	GithubToken string
	GithubUser  string
	GithubOrg   string
	GithubRepo  string
}

var (
	options = &GitReviewOptions{}

	RootCmd = &cobra.Command{
		Use:   "git-review",
		Short: "Print all review comments for the currrent branch",
		Run: func(cmd *cobra.Command, args []string) {
			client := util.GithubClient(options.GithubToken)

			_, comments, err := util.GetPullRequestComments(client, options.GithubOrg, options.GithubRepo, options.GithubUser)
			if err != nil {
				log.Fatalf("failed to fetch pull request for current branch: %v", err)
			}

			//log.Printf("pull=%#+v", *pull)
			for _, c := range comments {
				// header
				fmt.Fprintf(os.Stdout, "@%s on %s:%d %s:\n", *c.User.Login, *c.Path, *c.OriginalPosition, humanize.Time(*c.CreatedAt))
				util.FancyDiff(*c.DiffHunk, os.Stdout)
				fmt.Fprintf(os.Stdout, "\n\n--> %s\n\n", strings.TrimSpace(*c.Body))
			}
		},
	}
)

func Execute() {
	RootCmd.Flags().StringVarP(&options.GithubToken, "github-token", "t", os.Getenv("GITHUB_API_KEY"), "Github API key")
	RootCmd.Flags().StringVarP(&options.GithubUser, "github-user", "u", os.Getenv("GITHUB_USERNAME"), "Github user name")
	RootCmd.Flags().StringVarP(&options.GithubOrg, "github-org", "o", "", "Github organization name")
	RootCmd.Flags().StringVarP(&options.GithubRepo, "github-repo", "r", "", "Github repository name")

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
