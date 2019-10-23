package version

import (
	"fmt"
	"runtime"
)

var (
	// application's name
	Name = "ci123chain"
	// server binary name
	ServerName = "<cid>"
	// client binary name
	ClientName = "<cicli>"
	// application's version string
	Version = "1.0.0-beta"
	// commit
	Commit = ""
	// build tags
	BuildTags = ""
)


type Info struct {
	Name 	string `json: "name" yaml:"name"`
	ServerName string `json:"server_name" yaml:"server_name"`
	ClientName string `json:"client_name" yaml:"client_name"`
	Version    string `json:"version" yaml:"version"`
	GitCommit  string `json:"commit" yaml:"commit"`
	BuildTags  string `json:"build_tags" yaml:"build_tags"`
	GoVersion  string `json:"go" yaml:"go"`
}

func NewInfo() Info {
	return Info{
		Name:       Name,
		ServerName: ServerName,
		ClientName: ClientName,
		Version:    Version,
		GitCommit:  Commit,
		BuildTags:  BuildTags,
		GoVersion:  fmt.Sprintf("go version %s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH),
	}
}
