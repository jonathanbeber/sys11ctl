package metakube

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// GetProjects returns a list of projects given a metakube.Client
func GetProjects(c Client) ([]Project, error) {
	resp, err := c.do("GET", "/api/v1/projects", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var projects []Project
	var buf bytes.Buffer

	buf.ReadFrom(resp.Body)
	if err := json.Unmarshal(buf.Bytes(), &projects); err != nil {
		return nil, fmt.Errorf("failed to unmarshall received JSON: %s", err.Error())
	}

	return projects, nil
}

// Project represents a project returned by metakube API
type Project struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Status string  `json:"status"`
	Owners []Owner `json:"owners"`
}

// Owner represents a project owner returned by metakube API
type Owner struct {
	Name string `json:"Name"`
}
