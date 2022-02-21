package version

import (
	"fmt"
	"runtime"
)

var (
	// application's name
	Name = "ci123chain"
	// application's version string
	Version = "<version>"
	// commit
	Commit = ""
	// build tags
	BuildTags = ""
)


type Info struct {
	Name 	string 	`json: "name" yaml:"name"`
	Version    string `json:"version" yaml:"version"`
	GitCommit  string `json:"commit" yaml:"commit"`
	BuildTags  string `json:"build_tags" yaml:"build_tags"`
	GoVersion  string `json:"go" yaml:"go"`
}

func NewInfo() Info {
	return Info{
		Name:       Name,
		Version:    Version,
		GitCommit:  Commit,
		BuildTags:  BuildTags,
		GoVersion:  fmt.Sprintf("go version %s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH),
	}
}

func GetVersion() string {
	return Version
}