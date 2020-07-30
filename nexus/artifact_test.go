package nexus

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{\"totalCount\":60,\"from\":-1,\"count\":-1,\"tooManyResults\":false,\"collapsed\":false,\"repoDetails\":[{\"repositoryId\":\"releases\",\"repositoryName\":\"Releases\",\"repositoryContentClass\":\"maven2\",\"repositoryKind\":\"hosted\",\"repositoryPolicy\":\"RELEASE\",\"repositoryURL\":\"https://nexus/nexus/service/local/repositories/releases\"}],\"data\":[{\"groupId\":\"com.group.id\",\"artifactId\":\"app-server\",\"version\":\"1.0.0\",\"latestRelease\":\"1.0.0\",\"latestReleaseRepositoryId\":\"releases\",\"artifactHits\":[{\"repositoryId\":\"releases\",\"artifactLinks\":[{\"extension\":\"pom\"},{\"classifier\":\"el5-x86_64\",\"extension\":\"rpm\"},{\"classifier\":\"el7-x86_64\",\"extension\":\"rpm\"}]}]}]}")
	}))
	defer ts.Close()

	artifact := ArtifactInfo{
		ArtifactID:     "app-server",
		GroupID:        "com.group.id",
		Packaging:      "rpm",
		RepositoryType: "releases",
	}

	result, err := FindArtifact(artifact, ts.URL)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expectedResult := "1.0.0"
	if result != expectedResult {
		t.Errorf("expected: %s actual: %s", expectedResult, result)
	}
}

func TestInvalidURL(t *testing.T) {

	artifact := ArtifactInfo{
		ArtifactID:     "app-server",
		GroupID:        "com.group.id",
		Packaging:      "rpm",
		RepositoryType: "releases",
	}

	_, err := FindArtifact(artifact, "invalid_url")

	if err == nil {
		t.Error("should have an parser url error")
	}

}

func TestStaturError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(503)
	}))
	defer ts.Close()

	artifact := ArtifactInfo{
		ArtifactID:     "app-server",
		GroupID:        "com.group.id",
		Packaging:      "rpm",
		RepositoryType: "releases",
	}

	_, err := FindArtifact(artifact, ts.URL)

	if err == nil {
		t.Errorf("should have a status error")
	}

}
