package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/integrii/flaggy"
	"github.com/sfragata/download-artifact/nexus"
	"github.com/sfragata/download-artifact/utils"
)

// These variables will be replaced by real values when do gorelease
var (
	version = "none"
	date    string
	commit  string
)

func main() {

	info := fmt.Sprintf(
		"%s\nDate: %s\nCommit: %s\nOS: %s\nArch: %s",
		version,
		date,
		commit,
		runtime.GOOS,
		runtime.GOARCH,
	)

	flaggy.SetName("download-artifact")
	flaggy.SetDescription("Utility to download artifacts hosted on Nexus using Lucene and nexus rest api")
	flaggy.SetVersion(info)

	var artifact string
	flaggy.String(&artifact, "a", "artifact-id", "Maven artifact id")

	var group string
	flaggy.String(&group, "g", "group-id", "Maven group id")

	var version string
	flaggy.String(&version, "v", "artifact-version", "Artifact version")

	var packaging string
	flaggy.String(&packaging, "p", "packaging", "Type of packaging (ex. pom, jar, war etc)")

	var applicationName string
	flaggy.String(&applicationName, "n", "appName", "Name to be used as filename when dowloading")

	var classifier string
	flaggy.String(&classifier, "c", "classifier", "Artifact classifier")

	var targetFolder = "."
	flaggy.String(&targetFolder, "t", "target", "Target folder")

	var nexusHost string
	flaggy.String(&nexusHost, "H", "host", "Base nexus url")

	var repositoryID string
	flaggy.String(&repositoryID, "r", "repository", "Nexus repository id")

	var verbose bool
	flaggy.Bool(&verbose, "V", "verbose", "Verbose mode")

	var nexusVersion = 3
	flaggy.Int(&nexusVersion, "nv", "nexus-version", "Nexus version")

	flaggy.Parse()

	if utils.IsEmpty(artifact) || utils.IsEmpty(group) || utils.IsEmpty(nexusHost) {
		flaggy.ShowHelpAndExit("group-id, artifact-id and nexus host are mandatory")
	}

	artifactInfo := nexus.ArtifactInfo{
		ArtifactID:      artifact,
		GroupID:         group,
		Version:         version,
		Packaging:       packaging,
		Classifier:      classifier,
		ApplicationName: applicationName,
		RepositoryID:    repositoryID,
		TargetFolder:    targetFolder,
		Verbose:         verbose,
		NexusVersion:    nexusVersion,
	}

	err := nexus.DownloadArtifact(artifactInfo, nexusHost)

	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}
}
