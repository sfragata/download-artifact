package nexus

import (
	"log"
	"net/url"
	"strings"

	"github.com/sfragata/download-artifact/download"
	"github.com/sfragata/download-artifact/utils"
)

// Download donwload maven artifact
func (n Nexus3) download(artifactInfo ArtifactInfo, nexusHost string) error {
	downloadString, err := n.findArtifact(&artifactInfo, nexusHost)

	if err != nil {
		return err
	}

	downloadURL, err := url.Parse(downloadString)

	if artifactInfo.Verbose {
		log.Printf("Download url: %s", downloadURL.String())
	}

	if err != nil {
		return err
	}

	var filename string
	if utils.IsEmpty(artifactInfo.ApplicationName) {
		lastSlash := strings.LastIndex(downloadString, "/")
		filename = downloadString[lastSlash+1:]
	}

	downloadOptions := download.Options{
		URL:        *downloadURL,
		Filename:   filename,
		FolderPath: artifactInfo.TargetFolder,
	}

	err = download.GetFile(downloadOptions)

	if err != nil {
		return err
	}

	return nil
}
