package nexus

import (
	"fmt"
	"log"
	"net/url"

	"github.com/sfragata/download-artifact/download"
)

const nexusDownloadPath = "/nexus/service/local/artifact/maven/content?r=%s&g=%s&a=%s&p=%s&v=%s&c=%s"

// Download donwload maven artifact
func (n Nexus2) download(artifactInfo ArtifactInfo, nexusHost string) error {
	validatedVersion, err := n.findArtifact(&artifactInfo, nexusHost)
	baseDownloadURL := nexusHost + nexusDownloadPath

	if err != nil {
		return err
	}

	downloadString := fmt.Sprintf(baseDownloadURL, artifactInfo.RepositoryID, artifactInfo.GroupID, artifactInfo.ArtifactID, artifactInfo.Packaging, validatedVersion, artifactInfo.Classifier)

	downloadURL, err := url.Parse(downloadString)

	if artifactInfo.Verbose {
		log.Printf("Download url: %s", downloadURL.String())
	}

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
