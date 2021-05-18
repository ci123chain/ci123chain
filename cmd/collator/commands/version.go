package commands

var (
	// Version defines the application version (defined at compile time)
	Version = ""
	// Commit defines the application commit hash (defined at compile time)
	Commit = ""
	// SDKCommit defines the CosmosSDK commit hash (defined at compile time)
	SDKCommit = ""
)

type versionInfo struct {
	Version   string `json:"version" yaml:"version"`
	Commit    string `json:"commit" yaml:"commit"`
	Chain string `json:"chain" yaml:"chain"`
	Go        string `json:"go" yaml:"go"`
}
