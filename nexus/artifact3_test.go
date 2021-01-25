package nexus

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNexus3Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{\"items\":[{\"downloadUrl\":\"success\",\"path\":\"path\",\"id\":\"id\",\"repository\":\"maven-central\",\"format\":\"maven2\",\"checksum\":{\"sha1\":\"sha1\",\"sha256\":\"sha256\",\"sha512\":\"sha512\",\"md5\":\"md5\"},\"contentType\":\"application/java-archive\",\"lastModified\":\"2006-03-14T05:31:30.000+00:00\",\"maven2\":{\"extension\":\"jar\",\"groupId\":\"stax\",\"artifactId\":\"stax-api\",\"version\":\"1.0.1\"}}],\"continuationToken\":null}")
	}))
	defer ts.Close()

	artifact := ArtifactInfo{
		ArtifactID: "app-server",
		GroupID:    "com.group.id",
		Packaging:  "rpm",
		Verbose:    true,
	}
	n := Nexus3{}
	result, err := n.findArtifact(&artifact, ts.URL)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expectedResult := "success"
	if result != expectedResult {
		t.Errorf("expected: %s actual: %s", expectedResult, result)
	}
}

func TestNexus3InvalidURL(t *testing.T) {

	artifact := ArtifactInfo{
		ArtifactID:   "app-server",
		GroupID:      "com.group.id",
		Packaging:    "rpm",
		RepositoryID: "releases",
	}
	n := Nexus3{}
	_, err := n.findArtifact(&artifact, "invalid_url")

	if err == nil {
		t.Error("should have an parser url error")
	}

}

func TestNexus3StatusError(t *testing.T) {
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
	n := Nexus3{}
	_, err := n.findArtifact(&artifact, ts.URL)

	if err == nil {
		t.Error("should have a status error")
	}

}

func TestNexus3MissingJsonItemsTag(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{\"items\":[],\"continuationToken\":null}")
	}))
	defer ts.Close()

	artifact := ArtifactInfo{
		ArtifactID: "app-server",
		GroupID:    "com.group.id",
		Packaging:  "rpm",
	}
	n := Nexus3{}
	_, err := n.findArtifact(&artifact, ts.URL)

	if err == nil {
		t.Error("should have an error")
	}
	errorMsg := fmt.Sprintf("%v", err)
	if !strings.Contains(errorMsg, "Could not find artifact") {
		t.Errorf("Error, it should contains expected: 'Could not find artifact', actual: '%v'", err)
	}

}

func TestNexus3MissingJsonDownloadURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{\"items\":[{\"path\":\"path\",\"id\":\"id\",\"repository\":\"maven-central\",\"format\":\"maven2\",\"checksum\":{\"sha1\":\"sha1\",\"sha256\":\"sha256\",\"sha512\":\"sha512\",\"md5\":\"md5\"},\"contentType\":\"application/java-archive\",\"lastModified\":\"2006-03-14T05:31:30.000+00:00\",\"maven2\":{\"extension\":\"jar\",\"groupId\":\"stax\",\"artifactId\":\"stax-api\",\"version\":\"1.0.1\"}}],\"continuationToken\":null}")
	}))
	defer ts.Close()

	artifact := ArtifactInfo{
		ArtifactID: "app-server",
		GroupID:    "com.group.id",
		Packaging:  "rpm",
	}
	n := Nexus3{}
	_, err := n.findArtifact(&artifact, ts.URL)

	if err == nil {
		t.Errorf("should have an error")
	}
	errorMsg := fmt.Sprintf("%v", err)
	if !strings.Contains(errorMsg, "Could not find artifact") {
		t.Errorf("Error, it should contains expected: 'Could not get repository id', actual: '%v'", err)
	}
}

func TestNexus3NullJsonDownloadURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{\"items\":[{\"downloadUrl\":null,\"path\":\"path\",\"id\":\"id\",\"repository\":\"maven-central\",\"format\":\"maven2\",\"checksum\":{\"sha1\":\"sha1\",\"sha256\":\"sha256\",\"sha512\":\"sha512\",\"md5\":\"md5\"},\"contentType\":\"application/java-archive\",\"lastModified\":\"2006-03-14T05:31:30.000+00:00\",\"maven2\":{\"extension\":\"jar\",\"groupId\":\"stax\",\"artifactId\":\"stax-api\",\"version\":\"1.0.1\"}}],\"continuationToken\":null}")
	}))
	defer ts.Close()

	artifact := ArtifactInfo{
		ArtifactID: "app-server",
		GroupID:    "com.group.id",
		Packaging:  "rpm",
	}
	n := Nexus3{}
	_, err := n.findArtifact(&artifact, ts.URL)

	if err == nil {
		t.Errorf("should have an error")
	}
	errorMsg := fmt.Sprintf("%v", err)
	if !strings.Contains(errorMsg, "Could not find artifact") {
		t.Errorf("Error, it should contains expected: 'Could not get repository id', actual: '%v'", err)
	}
}

func TestNexus3EmptyJsonDownloadURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{\"items\":[{\"downloadUrl\":\"\",\"path\":\"path\",\"id\":\"id\",\"repository\":\"maven-central\",\"format\":\"maven2\",\"checksum\":{\"sha1\":\"sha1\",\"sha256\":\"sha256\",\"sha512\":\"sha512\",\"md5\":\"md5\"},\"contentType\":\"application/java-archive\",\"lastModified\":\"2006-03-14T05:31:30.000+00:00\",\"maven2\":{\"extension\":\"jar\",\"groupId\":\"stax\",\"artifactId\":\"stax-api\",\"version\":\"1.0.1\"}}],\"continuationToken\":null}")
	}))
	defer ts.Close()

	artifact := ArtifactInfo{
		ArtifactID: "app-server",
		GroupID:    "com.group.id",
		Packaging:  "rpm",
	}
	n := Nexus3{}
	_, err := n.findArtifact(&artifact, ts.URL)

	if err == nil {
		t.Errorf("should have an error")
	}
	errorMsg := fmt.Sprintf("%v", err)
	if !strings.Contains(errorMsg, "Could not find artifact") {
		t.Errorf("Error, it should contains expected: 'Could not get repository id', actual: '%v'", err)
	}
}

func TestNexus3InvalidJson(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{\"items\": [{\"downloadUrl\": \"sucess\",")
	}))
	defer ts.Close()

	artifact := ArtifactInfo{
		ArtifactID: "app-server",
		GroupID:    "com.group.id",
		Packaging:  "rpm",
	}
	n := Nexus3{}
	_, err := n.findArtifact(&artifact, ts.URL)

	if err == nil {
		t.Errorf("should have an error")
	}
	errorMsg := fmt.Sprintf("%v", err)
	if !strings.Contains(errorMsg, "Invalid JSON") {
		t.Errorf("Error, it should contains expected: 'Invalid JSON', actual: '%v'", err)
	}
}
