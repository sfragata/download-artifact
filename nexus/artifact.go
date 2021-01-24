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

const nexusSearchPath = "/nexus/service/local/lucene/search?repositoryId=%s&g=%s&a=%s&p=%s&v=%s&c=%s"

type artifactResponse struct {
	TotalCount       uint               `json:"totalCount"`
	Data             []data             `json:"data"`
	RepositoryDetail []repositoryDetail `json:"repoDetails"`
}

type data struct {
	LatestRelease  string `json:"latestRelease"`
	LatestSnapshot string `json:"latestSnapshot"`
}

type repositoryDetail struct {
	RepositoryID string `json:"repositoryId"`
}

// FindArtifact maven artifact and return its version
func FindArtifact(artifactInfo *ArtifactInfo, nexusHost string) (string, error) {
	baseSearchURL := nexusHost + nexusSearchPath

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

	artifactResponse := artifactResponse{}
	err = json.Unmarshal(body, &artifactResponse)
	if err != nil {
		return "", fmt.Errorf("Invalid JSON: %v", err)
	}

	if utils.IsEmpty(artifactInfo.RepositoryID) {
		if len(artifactResponse.RepositoryDetail) > 0 && !utils.IsEmpty(artifactResponse.RepositoryDetail[0].RepositoryID) {
			artifactInfo.RepositoryID = artifactResponse.RepositoryDetail[0].RepositoryID
		} else {
			return "", fmt.Errorf("Could not get repository id: %+v", artifactResponse.RepositoryDetail)
		}
	}
	if len(artifactResponse.Data) == 0 {
		return "", fmt.Errorf(errorMsg, artifactInfo.RepositoryID, artifactInfo.GroupID, artifactInfo.ArtifactID, artifactInfo.Packaging, artifactInfo.Version, artifactInfo.Classifier, artifactResponse)
	}

	//fetch first element of slice
	var latestVersion string
	if !utils.IsEmpty(artifactResponse.Data[0].LatestRelease) {
		latestVersion = artifactResponse.Data[0].LatestRelease
	} else if !utils.IsEmpty(artifactResponse.Data[0].LatestSnapshot) {
		latestVersion = artifactResponse.Data[0].LatestSnapshot
	}

	if utils.IsEmpty(latestVersion) {
		return "", fmt.Errorf(errorMsg, artifactInfo.RepositoryID, artifactInfo.GroupID, artifactInfo.ArtifactID, artifactInfo.Packaging, artifactInfo.Version, artifactInfo.Classifier, artifactResponse.Data)
	}

	return latestVersion, nil
}
