package githubVersionChecker

import (
	"encoding/json"
	"fmt"
	"time"
)

func UnmarshalRawCommit(data []byte) (RawCommit, error) {
	var r RawCommit
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *RawCommit) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type RawCommit struct {
	SHA         string      `json:"sha"`
	NodeID      string      `json:"node_id"`
	Commit      *CommitData `json:"commit"`
	URL         string      `json:"url"`
	HTMLURL     string      `json:"html_url"`
	CommentsURL string      `json:"comments_url"`
}

type CommitData struct {
	Author       *RawCommitIdentity `json:"author"`
	Committer    *RawCommitIdentity `json:"committer"`
	Message      string             `json:"message"`
	Tree         *Tree              `json:"tree"`
	URL          string             `json:"url"`
	CommentCount int                `json:"comment_count"`
	Verification *Verification      `json:"verification"`
}

type RawCommitIdentity struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}

type auxRawCommitIdentity struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Date  string `json:"date"`
}

func (x *RawCommitIdentity) MarshalJSON() ([]byte, error) {
	aux := &auxRawCommitIdentity{
		Name:  x.Name,
		Email: x.Email,
		Date:  x.Date.Format(time.RFC3339),
	}
	return json.Marshal(aux)
}

func (x *RawCommitIdentity) UnmarshalJSON(data []byte) error {
	aux := &auxRawCommitIdentity{}
	err := json.Unmarshal(data, aux)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	x.Name = aux.Name
	x.Email = aux.Email
	x.Date, err = time.Parse(time.RFC3339, aux.Date)
	if err != nil {
		return fmt.Errorf("time.Parse: %w", err)
	}
	return nil
}

type Tree struct {
	SHA string `json:"sha"`
	URL string `json:"url"`
}

type Verification struct {
	Verified  bool   `json:"verified"`
	Reason    string `json:"reason"`
	Signature string `json:"signature"`
	Payload   string `json:"payload"`
}
