package nexus

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNexu2Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{\"totalCount\":60,\"from\":-1,\"count\":-1,\"tooManyResults\":false,\"collapsed\":false,\"repoDetails\":[{\"repositoryId\":\"releases\",\"repositoryName\":\"Releases\",\"repositoryContentClass\":\"maven2\",\"repositoryKind\":\"hosted\",\"repositoryPolicy\":\"RELEASE\",\"repositoryURL\":\"https://nexus/nexus/service/local/repositories/releases\"}],\"data\":[{\"groupId\":\"com.group.id\",\"artifactId\":\"app-server\",\"version\":\"1.0.0\",\"latestRelease\":\"1.0.0\",\"latestReleaseRepositoryId\":\"releases\",\"artifactHits\":[{\"repositoryId\":\"releases\",\"artifactLinks\":[{\"extension\":\"pom\"},{\"classifier\":\"el5-x86_64\",\"extension\":\"rpm\"},{\"classifier\":\"el7-x86_64\",\"extension\":\"rpm\"}]}]}]}")
	}))
	defer ts.Close()

	artifact := ArtifactInfo{
		ArtifactID: "app-server",
		GroupID:    "com.group.id",
		Packaging:  "rpm",
		Verbose:    true,
	}
	n := Nexus2{}
	result, err := n.findArtifact(&artifact, ts.URL)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expectedResult := "1.0.0"
	if result != expectedResult {
		t.Errorf("expected: %s actual: %s", expectedResult, result)
	}
}

func TestNexus2InvalidURL(t *testing.T) {

	artifact := ArtifactInfo{
		ArtifactID:   "app-server",
		GroupID:      "com.group.id",
		Packaging:    "rpm",
		RepositoryID: "releases",
	}
	n := Nexus2{}
	_, err := n.findArtifact(&artifact, "invalid_url")

	if err == nil {
		t.Error("should have an parser url error")
	}

}

func TestNexus2StatusError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(503)
	}))
	defer ts.Close()

	artifact := ArtifactInfo{
		ArtifactID:   "app-server",
		GroupID:      "com.group.id",
		Packaging:    "rpm",
		RepositoryID: "releases",
	}
	n := Nexus2{}
	_, err := n.findArtifact(&artifact, ts.URL)

	if err == nil {
		t.Error("should have a status error")
	}

}

func TestNexus2MissingJsonDataTag(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{\"totalCount\":60,\"from\":-1,\"count\":-1,\"tooManyResults\":false,\"collapsed\":false,\"repoDetails\":[{\"repositoryId\":\"releases\",\"repositoryName\":\"Releases\",\"repositoryContentClass\":\"maven2\",\"repositoryKind\":\"hosted\",\"repositoryPolicy\":\"RELEASE\",\"repositoryURL\":\"https://nexus/nexus/service/local/repositories/releases\"}],\"data\":[]}")
	}))
	defer ts.Close()

	artifact := ArtifactInfo{
		ArtifactID: "app-server",
		GroupID:    "com.group.id",
		Packaging:  "rpm",
	}
	n := Nexus2{}
	_, err := n.findArtifact(&artifact, ts.URL)

	if err == nil {
		t.Error("should have an error")
	}

}

func TestNexus2MissingJsonRepositoryID(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{\"totalCount\":60,\"from\":-1,\"count\":-1,\"tooManyResults\":false,\"collapsed\":false,\"repoDetails\":[{\"repositoryName\":\"Releases\",\"repositoryContentClass\":\"maven2\",\"repositoryKind\":\"hosted\",\"repositoryPolicy\":\"RELEASE\",\"repositoryURL\":\"https://nexus/nexus/service/local/repositories/releases\"}],\"data\":[]}")
	}))
	defer ts.Close()

	artifact := ArtifactInfo{
		ArtifactID: "app-server",
		GroupID:    "com.group.id",
		Packaging:  "rpm",
	}
	n := Nexus2{}
	_, err := n.findArtifact(&artifact, ts.URL)

	if err == nil {
		t.Errorf("should have an error")
	}
	errorMsg := fmt.Sprintf("%v", err)
	if !strings.Contains(errorMsg, "Could not get repository id") {
		t.Errorf("Error, it should contains expected: 'Could not get repository id', actual: '%v'", err)
	}
}

func TestNexus2SuccessSnapshot(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{\"totalCount\":60,\"from\":-1,\"count\":-1,\"tooManyResults\":false,\"collapsed\":false,\"repoDetails\":[{\"repositoryId\":\"snapshots\",\"repositoryName\":\"snapshots\",\"repositoryContentClass\":\"maven2\",\"repositoryKind\":\"hosted\",\"repositoryPolicy\":\"RELEASE\",\"repositoryURL\":\"https://nexus/nexus/service/local/repositories/releases\"}],\"data\":[{\"groupId\":\"com.group.id\",\"artifactId\":\"app-server\",\"version\":\"1.0.0\",\"latestSnapshot\":\"1.0.0-SNAPSHOT\",\"latestSnapshotsRepositoryId\":\"snapshots\",\"artifactHits\":[{\"repositoryId\":\"snapshots\",\"artifactLinks\":[{\"extension\":\"pom\"},{\"classifier\":\"el5-x86_64\",\"extension\":\"rpm\"},{\"classifier\":\"el7-x86_64\",\"extension\":\"rpm\"}]}]}]}")
	}))
	defer ts.Close()

	artifact := ArtifactInfo{
		ArtifactID: "app-server",
		GroupID:    "com.group.id",
		Packaging:  "rpm",
	}
	n := Nexus2{}
	result, err := n.findArtifact(&artifact, ts.URL)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expectedResult := "1.0.0-SNAPSHOT"
	if result != expectedResult {
		t.Errorf("expected: %s actual: %s", expectedResult, result)
	}
}

func TestNexus2MissingLatestSnapshotAndRelease(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{\"totalCount\":60,\"from\":-1,\"count\":-1,\"tooManyResults\":false,\"collapsed\":false,\"repoDetails\":[{\"repositoryId\":\"snapshots\",\"repositoryName\":\"snapshots\",\"repositoryContentClass\":\"maven2\",\"repositoryKind\":\"hosted\",\"repositoryPolicy\":\"RELEASE\",\"repositoryURL\":\"https://nexus/nexus/service/local/repositories/releases\"}],\"data\":[{\"groupId\":\"com.group.id\",\"artifactId\":\"app-server\",\"version\":\"1.0.0\",\"latestSnapshotsRepositoryId\":\"snapshots\",\"artifactHits\":[{\"repositoryId\":\"snapshots\",\"artifactLinks\":[{\"extension\":\"pom\"},{\"classifier\":\"el5-x86_64\",\"extension\":\"rpm\"},{\"classifier\":\"el7-x86_64\",\"extension\":\"rpm\"}]}]}]}")
	}))
	defer ts.Close()

	artifact := ArtifactInfo{
		ArtifactID: "app-server",
		GroupID:    "com.group.id",
		Packaging:  "rpm",
	}
	n := Nexus2{}
	_, err := n.findArtifact(&artifact, ts.URL)

	if err == nil {
		t.Errorf("should have an error")
	}
	errorMsg := fmt.Sprintf("%v", err)
	if !strings.Contains(errorMsg, "Could not find artifact") {
		t.Errorf("Error, it should contains expected: 'Could not find artifact', actual: '%v'", err)
	}
}

func TestNexus2InvalidJson(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{\"totalCount\":60,\"")
	}))
	defer ts.Close()

	artifact := ArtifactInfo{
		ArtifactID: "app-server",
		GroupID:    "com.group.id",
		Packaging:  "rpm",
	}
	n := Nexus2{}
	_, err := n.findArtifact(&artifact, ts.URL)

	if err == nil {
		t.Errorf("should have an error")
	}
	errorMsg := fmt.Sprintf("%v", err)
	if !strings.Contains(errorMsg, "Invalid JSON") {
		t.Errorf("Error, it should contains expected: 'Invalid JSON', actual: '%v'", err)
	}
}
