package actions

import (
	"context"

	"github.com/google/go-github/v37/github"
	"golang.org/x/oauth2"
)

const (
	GithubPRBaseBranchMerge = "master"
	GithubPRState           = "open"
	GithubMaxPRPerQuery     = 100
)

func getGithubClient(accessToken string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)

	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func GetPullRequests(p params) ([]*github.PullRequest, error) {
	client := getGithubClient(p.accessToken)
	pullRequests, _, err := client.PullRequests.List(context.Background(), p.owner, p.repo, &github.PullRequestListOptions{
		State: GithubPRState,
		Base:  GithubPRBaseBranchMerge,
		ListOptions: github.ListOptions{
			PerPage: GithubMaxPRPerQuery,
		},
	})

	return pullRequests, err
}

func CommentPR(p params, pr *github.PullRequest, message string) error {
	client := getGithubClient(p.accessToken)
	review, _, err := client.PullRequests.CreateReview(context.Background(), p.owner, p.repo, pr.GetNumber(), &github.PullRequestReviewRequest{
		Body: ptrStr(message),
	})
	if err != nil {
		return err
	}
	_, _, err = client.PullRequests.SubmitReview(context.Background(), p.owner, p.repo, pr.GetNumber(), review.GetID(), &github.PullRequestReviewRequest{
		Event: ptrStr("COMMENT"),
	})
	return err
}

func MergePR(p params, pr *github.PullRequest) error {
	client := getGithubClient(p.accessToken)
	_, _, err := client.PullRequests.Merge(context.Background(), p.owner, p.repo, pr.GetNumber(), "", &github.PullRequestOptions{})
	return err
}

func getIssue(p params, issueNumber int) (*github.Issue, error) {
	client := getGithubClient(p.accessToken)
	issue, _, err := client.Issues.Get(context.Background(), p.owner, p.repo, issueNumber)
	return issue, err
}

func getLatestIssueCommentOfIssue(p params, issue *github.Issue) (*github.IssueComment, error) {
	client := getGithubClient(p.accessToken)
	comments, _, err := client.Issues.ListComments(context.Background(), p.owner, p.repo, issue.GetNumber(), &github.IssueListCommentsOptions{
		ListOptions: github.ListOptions{
			Page:    issue.GetComments(),
			PerPage: 1,
		},
	})
	if err != nil {
		return nil, err
	}

	return comments[0], nil
}

func addIssueComment(p params, issueNumber int, content string) error {
	client := getGithubClient(p.accessToken)
	_, _, err := client.Issues.CreateComment(context.Background(), p.owner, p.repo, issueNumber, &github.IssueComment{
		Body: &content,
	})
	return err
}
