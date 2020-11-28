package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/sfragata/download-artifact/nexus"
	"github.com/sfragata/download-artifact/utils"
)

func main() {
	artifact := flag.String("artifactId", "", "maven artifact id")
	group := flag.String("groupId", "", "maven group id")
	version := flag.String("version", "", "artifact version")
	packaging := flag.String("packaging", "war", "packaging")
	applicationName := flag.String("appName", "", "name to be used when dowloading")
	classifier := flag.String("classifier", "", "artifact classifier")
	targetFolder := flag.String("target", "/tmp", "target folder")
	repositoryType := "snapshots"
	nexusHost := flag.String("nexus", "", "base nexus url")

	flag.Parse()

	if utils.IsEmpty(*artifact) || utils.IsEmpty(*group) || utils.IsEmpty(*nexusHost) {
		log.Println("groupId, artifactId and nexusHost are mandatory")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if !strings.ContainsAny(*version, "SNAPSHOT") && len(*version) > 0 {
		repositoryType = "releases"
	}

	artifactInfo := nexus.ArtifactInfo{
		ArtifactID:      *artifact,
		GroupID:         *group,
		Version:         *version,
		Packaging:       *packaging,
		Classifier:      *classifier,
		ApplicationName: *applicationName,
		RepositoryType:  repositoryType,
		TargetFolder:    *targetFolder,
	}
	err := nexus.Download(artifactInfo, *nexusHost)

	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}
}
