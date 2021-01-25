package nexus

import (
	"fmt"
	"strings"
)

//ArtifactInfo Struct to Maven
type ArtifactInfo struct {
	ArtifactID      string
	GroupID         string
	Version         string
	Packaging       string
	ApplicationName string
	Classifier      string
	RepositoryID    string
	TargetFolder    string
	Verbose         bool
	NexusVersion    int
}
type commands interface {
	findArtifact(artifactInfo *ArtifactInfo, nexusHost string) (string, error)
	download(artifactInfo ArtifactInfo, nexusHost string) error
}

const errorMsg = "Could not find artifact\n\trepositoryId: %s\n\tgroupId: %s\n\tartifactId: %s\n\tpackaging: %s\n\tversion: %s\n\tclassifier: %s\n\tresponse: %+v\n"

//Nexus2 implementation
type Nexus2 struct{}

//Nexus3 implementation
type Nexus3 struct{}

var nexusImplementations = map[int]commands{
	2: Nexus2{},
	3: Nexus3{},
}

//DownloadArtifact search and download artifacts using nexus[2|3] implementation
func DownloadArtifact(artifactInfo ArtifactInfo, nexusHost string) error {
	implementation := nexusImplementations[artifactInfo.NexusVersion]
	if implementation == nil {
		return fmt.Errorf("Invalid nexus implementation %d, valid values are: %s", artifactInfo.NexusVersion, printKeys())
	}
	err := implementation.download(artifactInfo, nexusHost)
	if err != nil {
		return err
	}

	return nil
}

func printKeys() string {
	var keys []string
	for k := range nexusImplementations {
		keys = append(keys, fmt.Sprint(k))
	}

	return strings.Join(keys, ", ")
}
