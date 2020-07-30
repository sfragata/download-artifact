package nexus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/download-artifact/utils"
)

const nexusSearchPath = "/nexus/service/local/lucene/search?repositoryId=%s&g=%s&a=%s&p=%s&v=%s&c=%s"

// FindArtifact maven artifact and return its version
func FindArtifact(artifactInfo ArtifactInfo, nexusHost string) (string, error) {
	baseSearchURL := nexusHost + nexusSearchPath

	searchString := fmt.Sprintf(baseSearchURL, artifactInfo.RepositoryType, artifactInfo.GroupID, artifactInfo.ArtifactID, artifactInfo.Packaging, artifactInfo.Version, artifactInfo.Classifier)

	searchURL, err := url.Parse(searchString)
	if err != nil {
		return "", err
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

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	datas := result["data"].([]interface{})

	if len(datas) == 0 {
		return "", fmt.Errorf(errorMsg, artifactInfo.RepositoryType, artifactInfo.GroupID, artifactInfo.ArtifactID, artifactInfo.Packaging, artifactInfo.Version, artifactInfo.Classifier, result)
	}

	//fetch first element of slice
	firstData := datas[0].(map[string]interface{})

	var latestVersion string
	if artifactInfo.RepositoryType == "releases" {
		latestVersion = firstData["latestRelease"].(string)
	} else {
		latestVersion = firstData["latestSnapshot"].(string)
	}

	if utils.IsEmpty(latestVersion) {
		return "", fmt.Errorf(errorMsg, artifactInfo.RepositoryType, artifactInfo.GroupID, artifactInfo.ArtifactID, artifactInfo.Packaging, artifactInfo.Version, artifactInfo.Classifier, datas)
	}

	return latestVersion, nil
}
