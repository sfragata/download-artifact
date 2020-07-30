package nexus

import (
	"fmt"
	"net/url"

	"github.com/download-artifact/download"
)

//ArtifactInfo Struct to Maven
type ArtifactInfo struct {
	ArtifactID      string
	GroupID         string
	Version         string
	Packaging       string
	ApplicationName string
	Classifier      string
	RepositoryType  string
	TargetFolder    string
}

const nexusDownloadPath = "/nexus/service/local/artifact/maven/content?r=%s&g=%s&a=%s&p=%s&v=%s&c=%s"
const errorMsg = "Could not find artifact\n\trepositoryType: %s\n\tgroupId: %s\n\tartifactId: %s\n\tpackaging: %s\n\tversion: %s\n\tclassifier: %s\n\tresponse: %+v\n"

// Download donwload maven artifact
func Download(artifactInfo ArtifactInfo, nexusHost string) error {
	validatedVersion, err := FindArtifact(artifactInfo, nexusHost)
	baseDownloadURL := nexusHost + nexusDownloadPath

	if err != nil {
		return err
	}

	downloadString := fmt.Sprintf(baseDownloadURL, artifactInfo.RepositoryType, artifactInfo.GroupID, artifactInfo.ArtifactID, artifactInfo.Packaging, validatedVersion, artifactInfo.Classifier)

	downloadURL, err := url.Parse(downloadString)

	if err != nil {
		return err
	}

	downloadOptions := download.Options{
		URL:           *downloadURL,
		Filename:      artifactInfo.ApplicationName,
		FileExtension: artifactInfo.Packaging,
		FolderPath:    artifactInfo.TargetFolder,
	}

	err = download.GetFile(downloadOptions)

	if err != nil {
		return err
	}

	return nil
}
