package githubVersionChecker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-github/v48/github"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	ErrNoTags         = errors.New("repository has no tags")
	ErrCommitNotFound = errors.New("unable to find this commit")
)

type CodeBase struct {
	Client     *github.Client
	Owner      string
	Repository string
}

func NewCodeBase(hClient *http.Client, owner, repo string) *CodeBase {
	return &CodeBase{
		Client:     github.NewClient(hClient),
		Owner:      owner,
		Repository: repo,
	}
}

func (c *CodeBase) GetLatestVersion(ctx context.Context, opts *github.ListOptions) (*Version, error) {
	tags, _, err := c.Client.Repositories.ListTags(ctx, c.Owner, c.Repository, opts)
	if err != nil {
		return nil, fmt.Errorf("c.Client.Repositories.ListTags: %w", err)
	}

	if len(tags) < 1 {
		return nil, ErrNoTags
	}

	latestTag := tags[0]
	// Get the tag and parse its RAW response to get fields the GitHub library in question doesn't expose
	rawResp, err := c.Client.Client().Get(latestTag.GetCommit().GetURL())
	if err != nil {
		return nil, fmt.Errorf("c.Client.Client().Get: %w", err)
	}
	rawRespBytes, err := io.ReadAll(rawResp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}
	rawCommit := &RawCommit{}
	err = json.Unmarshal(rawRespBytes, rawCommit)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	latestVersion, err := NewVersion(latestTag.GetName(), latestTag.GetCommit().GetSHA(), rawCommit.Commit.Committer.Date)
	if err != nil {
		return nil, err
	}

	return latestVersion, nil
}

func (c *CodeBase) PromptCurrentVersion(version *Version) *Prompt {
	outputBuilder := strings.Builder{}
	outputBuilder.WriteString(fmt.Sprintf("Current version data for %s\n", c.Repository))
	outputBuilder.WriteString(fmt.Sprintf("Version:\t%s\n", version.Version.Original()))
	outputBuilder.WriteString(fmt.Sprintf("Git commit:\t%s\n", version.Commit))
	outputBuilder.WriteString(fmt.Sprintf("Git Ref:\t%s\n", version.Ref))
	outputBuilder.WriteString(fmt.Sprintf("Date:\t%s\n", version.BuildTime.Format(time.UnixDate)))
	repoURL := fmt.Sprintf("https://github.com/%s/%s", c.Owner, c.Repository)
	return &Prompt{Output: outputBuilder.String(), RepositoryURL: repoURL, UpdateURL: repoURL + "/releases/" + version.Version.Original()}
}

func (c *CodeBase) PromptLatestVersion(current, latest *Version) (bool, *Prompt) {
	outputBuilder := strings.Builder{}
	outdated := current.Version.LessThan(latest.Version)
	repoURL := fmt.Sprintf("https://github.com/%s/%s", c.Owner, c.Repository)
	updateURL := repoURL + "/releases/" + current.Version.Original()
	if outdated {
		outputBuilder.WriteString(fmt.Sprintf("Latest version data for %s\n", c.Repository))
		outputBuilder.WriteString(fmt.Sprintf("Version:\t%s\n", latest.Version.Original()))
		outputBuilder.WriteString(fmt.Sprintf("Git commit:\t%s\n", latest.Commit))
		outputBuilder.WriteString(fmt.Sprintf("Git Ref:\t%s\n", latest.Ref))
		outputBuilder.WriteString(fmt.Sprintf("Date:\t%s\n", latest.BuildTime.Format(time.UnixDate)))
		updateURL = repoURL + "/releases/" + latest.Version.Original()
	} else {
		outputBuilder.WriteString("You are using the latest version!\n")
	}

	return outdated, &Prompt{Output: outputBuilder.String(), RepositoryURL: repoURL, UpdateURL: updateURL}
}
