package util

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func GithubClient(token string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	return github.NewClient(oauth2.NewClient(oauth2.NoContext, ts))
}

func GetPullRequestComments(client *github.Client, org, name, user string) (*github.PullRequest, []*github.PullRequestComment, error) {
	result, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return nil, nil, fmt.Errorf("git rev-parse --abbrev-ref HEAD failed: %s (%v)", string(result), err)
	}

	listOpts := &github.PullRequestListOptions{Head: user + ":" + strings.TrimSpace(string(result))}
	pulls, _, err := client.PullRequests.List(org, name, listOpts)
	if err != nil {
		return nil, nil, err
	}

	switch {
	case len(pulls) == 0:
		return nil, nil, fmt.Errorf("no pull request found for branch %q", listOpts.Head)
	case len(pulls) > 1:
		return nil, nil, fmt.Errorf("multiple pull requests found for branch %q", listOpts.Head)
	}

	pullCommentOpt := &github.PullRequestListCommentsOptions{}
	comments, _, err := client.PullRequests.ListComments(org, name, *pulls[0].Number, pullCommentOpt)
	if err != nil {
		return nil, nil, err
	}

	return pulls[0], comments, nil
}

func FancyDiff(in string, out io.Writer) {
	cmd := exec.Command("diff-so-fancy")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return
	}
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, in)
	}()
	output, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	out.Write(output)
}
