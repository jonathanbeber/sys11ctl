package metakube

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type mockClient struct {
	doFn func(string, string, io.Reader) (*http.Response, error)
}

func (c mockClient) do(method, url string, body io.Reader) (*http.Response, error) {
	return c.doFn(method, url, body)
}

func TestGet(t *testing.T) {
	var calledMethod, calledPath string
	client := mockClient{
		doFn: func(method string, url string, body io.Reader) (*http.Response, error) {
			calledMethod = method
			calledPath = url
			return &http.Response{
				Body: ioutil.NopCloser(strings.NewReader(`
					[
						{
							"id": "abcdefgh123",
							"name": "test1",
							"creationTimestamp": "2019-04-03T12:00:00Z",
							"status": "Active",
							"owners": [
								{
								"name": "a@test.com",
								"creationTimestamp": "0001-01-01T00:00:00Z",
								"email": "a@test.com"
								}
							]
						}
					]`)),
			}, nil
		},
	}
	projects, err := GetProjects(client)
	if err != nil {
		t.Fatalf("Err not expected: %v", err)
	}

	if calledMethod != "GET" {
		t.Fatalf("Expected 'GET' method, found '%s'", calledMethod)
	}
	if calledPath != "/api/v1/projects" {
		t.Fatalf("Expected '/api/v1/projects' path, found '%s'", calledPath)
	}

	if len(projects) != 1 {
		t.Fatalf("Expected 1 project, found '%d'", len(projects))
	}
	if projects[0].ID != "abcdefgh123" {
		t.Fatalf("expected id '%s', found '%s'", "abcdefgh123", projects[0].ID)
	}
	if projects[0].Name != "test1" {
		t.Fatalf("expected name '%s', found '%s'", "test1", projects[0].Name)
	}
	if projects[0].Status != "Active" {
		t.Fatalf("expected status '%s', found '%s'", "Active", projects[0].Status)
	}
	if len(projects[0].Owners) != 1 {
		t.Fatalf("Expected 1 project owner, found '%d'", len(projects[0].Owners))
	}
	if projects[0].Owners[0].Name != "a@test.com" {
		t.Fatalf("expected owner name '%s', found '%s'", "a@test.com", projects[0].Owners[0].Name)
	}
}

func TestGetErrorOnClient(t *testing.T) {
	client := mockClient{
		doFn: func(method string, url string, body io.Reader) (*http.Response, error) {
			return nil, errors.New("Expected error")
		},
	}
	_, err := GetProjects(client)
	if err == nil {
		t.Fatalf("Expected error 'Expected error', none found")
	}
	if err.Error() != "Expected error" {
		t.Fatalf("Expected error 'Expected error', found '%s'", err.Error())
	}
}

func TestGetErrorOnJson(t *testing.T) {
	var calledMethod, calledPath string
	client := mockClient{
		doFn: func(method string, url string, body io.Reader) (*http.Response, error) {
			calledMethod = method
			calledPath = url
			return &http.Response{
				Body: ioutil.NopCloser(strings.NewReader(`
					"id": "anInvalidJSON",
					"name": "test1",
				`)),
			}, nil
		},
	}
	_, err := GetProjects(client)
	if err == nil {
		t.Fatalf("Expected error 'Expected error', none found")
	}
	if !strings.HasPrefix(err.Error(), "failed to unmarshall received JSON") {
		t.Fatalf("Expected error 'failed to unmarshall received JSON', found '%s'", err.Error())
	}

	if calledMethod != "GET" {
		t.Fatalf("Expected 'GET' method, found '%s'", calledMethod)
	}
	if calledPath != "/api/v1/projects" {
		t.Fatalf("Expected '/api/v1/projects' path, found '%s'", calledPath)
	}
}
