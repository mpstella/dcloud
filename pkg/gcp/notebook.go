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

type ResponseError struct {
	Code    int
	Message string
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("Error response status (%d): %s", e.Code, e.Message)
}

func NewNotebookClient(projectID string) (*NotebookClient, error) {

	ctx := context.Background()

	creds, err := google.FindDefaultCredentials(ctx, scopes)
	if err != nil {
		return nil, err
	}

	token, err := creds.TokenSource.Token()
	if err != nil {
		return nil, err
	}

	return &NotebookClient{
		url:   fmt.Sprintf("%s/projects/%s/locations/%s/notebookRuntimeTemplates", serviceEndpoint, projectID, location),
		token: token.AccessToken,
	}, nil
}

func (nc *NotebookClient) curl(method string, url string, payload io.Reader) ([]byte, error) {

	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", nc.token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, &ResponseError{Code: resp.StatusCode, Message: string(body)}
	}

	return body, nil
}

func (nc *NotebookClient) GetNotebookRuntimeTemplates() (*ListNotebookRuntimeTemplatesResult, error) {

	body, err := nc.curl("GET", nc.url, nil)

	if err != nil {
		return nil, err
	}

	var templates ListNotebookRuntimeTemplatesResult
	err = json.Unmarshal(body, &templates)

	if err != nil {
		return nil, err
	}
	return &templates, nil
}

func (nc *NotebookClient) DeleteNotebookRuntimeTemplate(name string) error {

	url := fmt.Sprintf("%s/%s", serviceEndpoint, name)
	logrus.Infof("Deleting: %s", url)

	_, err := nc.curl("DELETE", url, nil)

	return err
}

func (nc *NotebookClient) DeployNotebookRuntimeTemplate(template *NotebookRuntimeTemplate) error {

	payload, err := json.Marshal(template)

	if err != nil {
		return err
	}
	_, err = nc.curl("POST", nc.url, bytes.NewBuffer(payload))
	return err
}
