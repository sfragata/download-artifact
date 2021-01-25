package nexus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/sfragata/download-artifact/utils"
)

const nexus3SearchPath = "/service/rest/v1/search/assets?repository=%s&group=%s&name=%s&maven.extension=%s&version=%s&maven.classifier=%s"

type artifactNexus3Response struct {
	Items []struct {
		DownloadURL  string    `json:"downloadUrl"`
		Repository   string    `json:"repository"`
		LastModified time.Time `json:"lastModified"`
	} `json:"items"`
}

// FindArtifact maven artifact and return its downloadURL
func (n Nexus3) findArtifact(artifactInfo *ArtifactInfo, nexusHost string) (string, error) {
	baseSearchURL := nexusHost + nexus3SearchPath

	searchString := fmt.Sprintf(baseSearchURL, artifactInfo.RepositoryID, artifactInfo.GroupID, artifactInfo.ArtifactID, artifactInfo.Packaging, artifactInfo.Version, artifactInfo.Classifier)

	searchURL, err := url.Parse(searchString)
	if err != nil {
		return "", err
	}

	if artifactInfo.Verbose {
		log.Printf("Search url: %s", searchURL.String())
	}

	request, err := http.NewRequest("GET", searchURL.String(), nil)
	if err != nil {
		return "", err
	}

	request.Header.Set("Accept", "application/json")

	httpClient := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	resp, err := httpClient.Do(request)

	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if artifactInfo.Verbose {
		log.Printf("response: %s\n", body)
	}

	artifactResponse := artifactNexus3Response{}
	err = json.Unmarshal(body, &artifactResponse)
	if err != nil {
		return "", fmt.Errorf("Invalid JSON: %v", err)
	}

	if len(artifactResponse.Items) == 0 {
		return "", fmt.Errorf(errorMsg, artifactInfo.RepositoryID, artifactInfo.GroupID, artifactInfo.ArtifactID, artifactInfo.Packaging, artifactInfo.Version, artifactInfo.Classifier, artifactResponse)
	}

	//fetch first element of slice
	var downloadURL string
	if !utils.IsEmpty(artifactResponse.Items[0].DownloadURL) {
		downloadURL = artifactResponse.Items[0].DownloadURL
	}

	if utils.IsEmpty(downloadURL) {
		return "", fmt.Errorf(errorMsg, artifactInfo.RepositoryID, artifactInfo.GroupID, artifactInfo.ArtifactID, artifactInfo.Packaging, artifactInfo.Version, artifactInfo.Classifier, artifactResponse.Items)
	}

	return downloadURL, nil
}
