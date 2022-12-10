package ghvu

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/go-version"
	"testing"
	"time"
)

func mustVersion(versionStr string, delta int) *Version {
	versionObj, _ := version.NewVersion(versionStr)
	date := time.Now().AddDate(0, 0, delta)
	return &Version{
		Version:   versionObj,
		Commit:    date.String(),
		Ref:       "refs/tags/" + date.String(),
		BuildTime: date,
	}
}

// Not actually tests, more like sanity checks

func TestCodeBase_GetLatestVersion(t *testing.T) {
	base := NewCodeBase(nil, "BRUHItsABunny", "Premiumize-File-Sync")
	v, err := base.GetLatestVersion(context.Background(), nil)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(spew.Sdump(v))
}

func TestCodeBase_PromptCurrentVersion(t *testing.T) {
	base := NewCodeBase(nil, "BRUHItsABunny", "Premiumize-File-Sync")
	v, err := base.GetLatestVersion(context.Background(), nil)
	if err != nil {
		t.Error(err)
	}
	prompt := base.PromptCurrentVersion(v)
	fmt.Println(prompt.Output)
}

func TestCodeBase_PromptLatestVersion(t *testing.T) {
	base := NewCodeBase(nil, "BRUHItsABunny", "Premiumize-File-Sync")
	v, err := base.GetLatestVersion(context.Background(), nil)
	if err != nil {
		t.Error(err)
	}
	v2 := mustVersion("v1.0.0", -1)
	v3 := mustVersion("v10.0.0", 1)
	ok, prompt := base.PromptLatestVersion(v2, v)
	if !ok {
		t.Error("should be true when comparing v2")
	}
	fmt.Println(prompt.Output)
	fmt.Println(fmt.Sprintf("You can find more here:\n%s\n", prompt.UpdateURL))
	fmt.Println("====================================")

	ok, prompt = base.PromptLatestVersion(v3, v)
	if ok {
		t.Error("should be false when comparing v3")
	}
	fmt.Println(prompt.Output)
}
