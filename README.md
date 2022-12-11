# Go GitHub Versioning Utility

A library to help to version your Go binaries with GitHub releases.

## Usage

### Step 1
Make a `version.go` file and populate it with the following information:

```go
package utils

import (
	ghvu "github.com/BRUHItsABunny/go-ghvu"
	"strings"
)

const none string = ""

// ldflags
var (
	AppVersion      = "v0.0.1"
	BuildTime       = none
	GitCommit       = none
	GitRef          = none
	GitRepo         = "https://github.com/USER/REPO/"
	CurrentVersion  *ghvu.Version
	CurrentCodeBase *ghvu.CodeBase
)

func init() {
	CurrentVersion = ghvu.NewVersionOrDefault(AppVersion, GitCommit, GitRef, BuildTime)
	repoSplit := strings.Split(GitRepo, "/")
	CurrentCodeBase = ghvu.NewCodeBase(nil, repoSplit[len(repoSplit)-3], repoSplit[len(repoSplit)-2])
}
```

### Step 2
Print the current version and check for a newer function like this:

```go
package main

import "fmt"

func main(){
	currentPrompt := utils.CurrentCodeBase.PromptCurrentVersion(utils.CurrentVersion)
	latestVersion, err := utils.CurrentCodeBase.GetLatestVersion(context.Background(), nil)
	if err != nil {
		panic(fmt.Errorf("utils.CurrentCodeBase.GetLatestVersion: %w", err))
	}
	isOutdated, latestPrompt := utils.CurrentCodeBase.PromptLatestVersion(utils.CurrentVersion, latestVersion)
	fmt.Println(currentPrompt.Output)
	if isOutdated {
		fmt.Println(latestPrompt.Output)
		fmt.Println(fmt.Sprintf("You can find more here:\n%s\n", latestPrompt.UpdateURL))
	}
}
```

### Step 3
Build your executable with your ldflags

`go build main.go -ldflags="-X 'utils.AppVersion=v0.0.1' -X 'utils.GitCommit=somehash' -X 'utils.GitRef=sometag'"`

### Step 4
Tag and release your code in accordance with your version you built with

## Examples

### 1. Canary Replay
You can find the whole repository [here](https://github.com/BRUHItsABunny/canary-replay).
* Step 1: [version.go](https://github.com/BRUHItsABunny/canary-replay/blob/main/utils/version.go)
* Step 2: [main.go](https://github.com/BRUHItsABunny/canary-replay/blob/abbf10c740bae9082fab9a97336c66b72c91d589/main.go#L74)
* Step 3+4: [build with GitHub actions](https://github.com/BRUHItsABunny/canary-replay/blob/abbf10c740bae9082fab9a97336c66b72c91d589/.github/workflows/tag.yaml)
