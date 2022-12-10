package ghvu

import (
	"fmt"
	"github.com/hashicorp/go-version"
	"time"
)

type Version struct {
	Version   *version.Version
	Commit    string
	Ref       string
	BuildTime time.Time
}

func DefaultVersion() *Version {
	versionObj, _ := version.NewVersion("v0.0.1")
	return &Version{
		Version:   versionObj,
		Commit:    "",
		Ref:       "refs/tags/v0.0.1",
		BuildTime: time.Time{},
	}
}

func NewVersion(versionStr, commit string, buildTime time.Time) (*Version, error) {
	versionObj, err := version.NewVersion(versionStr)
	if err != nil {
		return nil, fmt.Errorf("version.NewVersion: %w", err)
	}

	return &Version{
		Version:   versionObj,
		Commit:    commit,
		Ref:       "refs/tags/" + versionStr,
		BuildTime: buildTime,
	}, nil
}
