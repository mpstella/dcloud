package gcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"net/http"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
)

const (
	location        = "australia-southeast1"
	scopes          = "https://www.googleapis.com/auth/cloud-platform"
	serviceEndpoint = "https://australia-southeast1-aiplatform.googleapis.com/v1beta1"
)

type NotebookClient struct {
	url   string
	token string
}

func NewNotebookClient(projectID string) NotebookClient {

	ctx := context.Background()

	creds, err := google.FindDefaultCredentials(ctx, scopes)
	if err != nil {
		logrus.Fatalf("Failed to obtain default credentials: %v", err)
	}

	token, err := creds.TokenSource.Token()
	if err != nil {
		logrus.Fatalf("Failed to get token: %v", err)
	}

	return NotebookClient{
		url:   fmt.Sprintf("%s/projects/%s/locations/%s/notebookRuntimeTemplates", serviceEndpoint, projectID, location),
		token: token.AccessToken,
	}
}

func (nc *NotebookClient) curl(method string, url string, payload io.Reader) []byte {

	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		logrus.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", nc.token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		logrus.Fatalf("Failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		logrus.Fatalf("Failed to read response body: %v", err)
	}

	if resp.StatusCode != 200 {
		logrus.Fatalf("Status returned: %d (%s)", resp.StatusCode, string(body))
	}

	logrus.Debugf("Response status: %s", resp.Status)
	return body
}

func (nc *NotebookClient) GetNotebookRuntimeTemplates() *ListNotebookRuntimeTemplatesResult {

	body := nc.curl("GET", nc.url, nil)

	var templates ListNotebookRuntimeTemplatesResult
	json.Unmarshal(body, &templates)
	return &templates
}

func (nc *NotebookClient) DeleteNotebookRuntimeTemplate(name string) {

	url := fmt.Sprintf("%s/%s", serviceEndpoint, name)
	logrus.Infof("Deleting: %s", url)
	nc.curl("DELETE", url, nil)
}

func (nc *NotebookClient) DeployNotebookRuntimeTemplate(template *NotebookRuntimeTemplate) {

	payload, err := json.Marshal(template)

	if err != nil {
		logrus.Fatal("Error creating JSON Payload", err)
	}

	nc.curl("POST", nc.url, bytes.NewBuffer(payload))
}
