package revision_plate

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func runServer(t *testing.T, h http.Handler, method string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(method, "/site/sha", nil)
	if err != nil {
		t.Fatal(err)
	}
	h.ServeHTTP(w, req)
	return w
}

func createRevisionFile(t *testing.T, path string, rev string) {
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	_, err = file.WriteString(rev)
	if err != nil {
		t.Fatal(err)
	}
}

func removeRevisionFile(t *testing.T, path string) {
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
}

func TestGetRevision(t *testing.T) {
	createRevisionFile(t, "REVISION", "deadbeef")
	defer removeRevisionFile(t, "REVISION")

	recorder := runServer(t, New(""), "GET")
	if recorder.Code != 200 {
		t.Errorf("Expected status code 200, but got %d", recorder.Code)
	}
	body := recorder.Body.String()
	if body != "deadbeef" {
		t.Errorf("Expected response body 'deadbeef', but got %s", body)
	}
}

func TestHeadRevision(t *testing.T) {
	createRevisionFile(t, "REVISION", "deadbeef")
	defer removeRevisionFile(t, "REVISION")

	recorder := runServer(t, New(""), "HEAD")
	if recorder.Code != 200 {
		t.Errorf("Expected status code 200, but got %d", recorder.Code)
	}
	body := recorder.Body.String()
	if body != "" {
		t.Errorf("Expected empty response body, but got %s", body)
	}
}

func TestGetRevisionWithoutFile(t *testing.T) {
	recorder := runServer(t, New(""), "GET")
	if recorder.Code != 404 {
		t.Errorf("Expected status code 404, but got %d", recorder.Code)
	}
	body := recorder.Body.String()
	if body != "REVISION_FILE_NOT_FOUND" {
		t.Errorf("Expected response body 'REVISION_FILE_NOT_FOUND', but got %s", body)
	}
}

func TestHeadRevisionWithRemovedFile(t *testing.T) {
	createRevisionFile(t, "REVISION", "deadbeef")
	defer removeRevisionFile(t, "REVISION")
	h := New("")

	r1 := runServer(t, h, "HEAD")
	if r1.Code != 200 {
		t.Errorf("Expected status code 200, but got %d", r1.Code)
	}
	b1 := r1.Body.String()
	if b1 != "" {
		t.Errorf("Expected empty esponse body, but got %s", b1)
	}

	removeRevisionFile(t, "REVISION")

	r2 := runServer(t, h, "HEAD")
	if r2.Code != 404 {
		t.Errorf("Expected status code 404, but got %d", r2.Code)
	}
	b2 := r2.Body.String()
	if b2 != "" {
		t.Errorf("Expected empty esponse body, but got %s", b2)
	}
}

func TestHeadRevisionWithoutFile(t *testing.T) {
	recorder := runServer(t, New(""), "HEAD")
	if recorder.Code != 404 {
		t.Errorf("Expected status code 404, but got %d", recorder.Code)
	}
	body := recorder.Body.String()
	if body != "" {
		t.Errorf("Expected empty response body, but got %s", body)
	}
}

func TestGetRevisionWithRemovedFile(t *testing.T) {
	createRevisionFile(t, "REVISION", "deadbeef")
	defer removeRevisionFile(t, "REVISION")
	h := New("")

	r1 := runServer(t, h, "GET")
	if r1.Code != 200 {
		t.Errorf("Expected status code 200, but got %d", r1.Code)
	}
	b1 := r1.Body.String()
	if b1 != "deadbeef" {
		t.Errorf("Expected esponse body 'deadbeef', but got %s", b1)
	}

	removeRevisionFile(t, "REVISION")

	r2 := runServer(t, h, "GET")
	if r2.Code != 404 {
		t.Errorf("Expected status code 404, but got %d", r2.Code)
	}
	b2 := r2.Body.String()
	if b2 != "REVISION_FILE_REMOVED" {
		t.Errorf("Expected esponse body 'REVISION_FILE_REMOVED', but got %s", b2)
	}
}

func TestGetRevisionWithCustomPath(t *testing.T) {
	createRevisionFile(t, "site-sha", "cafebabe")
	defer removeRevisionFile(t, "site-sha")

	r1 := runServer(t, New(""), "GET")
	if r1.Code != 404 {
		t.Errorf("Expected status code 404, but got %d", r1.Code)
	}

	r2 := runServer(t, New("site-sha"), "GET")
	if r2.Code != 200 {
		t.Errorf("Expected status code 200, but got %d", r2.Code)
	}

	body := r2.Body.String()
	if body != "cafebabe" {
		t.Errorf("Expected response body 'cafebabe', but got %s", body)
	}
}

func TestGetRevisionWithUpdatedFile(t *testing.T) {
	createRevisionFile(t, "REVISION", "deadbeef")
	defer removeRevisionFile(t, "REVISION")

	h := New("")
	r1 := runServer(t, h, "GET")
	if r1.Code != 200 {
		t.Errorf("Expected status code 200, but got %d", r1.Code)
	}
	b1 := r1.Body.String()
	if b1 != "deadbeef" {
		t.Errorf("Expected response body 'deadbeef', but got %s", b1)
	}

	createRevisionFile(t, "REVISION", "cafebabe")

	r2 := runServer(t, h, "GET")
	if r2.Code != 200 {
		t.Errorf("Expected status code 200, but got %d", r2.Code)
	}
	b2 := r2.Body.String()
	if b2 != "deadbeef" {
		t.Errorf("Expected response body 'deadbeef', but got %s", b2)
	}
}

func TestGetRevisionWithFileCreatedAfterInitialize(t *testing.T) {
	h := New("")
	r1 := runServer(t, h, "GET")
	if r1.Code != 404 {
		t.Errorf("Expected status code 404, but got %d", r1.Code)
	}
	b1 := r1.Body.String()
	if b1 != "REVISION_FILE_NOT_FOUND" {
		t.Errorf("Expected response body 'REVISION_FILE_NOT_FOUND', but got %s", b1)
	}

	createRevisionFile(t, "REVISION", "deadbeef")
	defer removeRevisionFile(t, "REVISION")

	r2 := runServer(t, h, "GET")
	if r2.Code != 404 {
		t.Errorf("Expected status code 404, but got %d", r2.Code)
	}
	b2 := r2.Body.String()
	if b2 != "REVISION_FILE_NOT_FOUND" {
		t.Errorf("Expected response body 'REVISION_FILE_NOT_FOUND', but got %s", b2)
	}
}
