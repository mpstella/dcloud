package gcp

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"

	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
	"gopkg.in/yaml.v3"
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

type TemplateComparison int

const (
	DoesNotExist = iota
	ExistsAndIsIdentical
	ExistsButIsDifferent
)

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

func (nc *NotebookClient) curl(method string, url string, payload io.Reader) (string, []byte) {

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

	// Read and print the response
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		logrus.Fatalf("Failed to read response body: %v", err)
	}

	logrus.Debugf("Response status: %s", resp.Status)

	if resp.Status != "200 OK" {
		logrus.Error(string(body))
	}

	return resp.Status, body
}

func (nc *NotebookClient) resolve(template NotebookRuntimeTemplate, checksum string) (TemplateComparison, *NotebookRuntimeTemplate) {

	templates := nc.GetNotebookRuntimeTemplates()

	for _, existing := range templates.NotebookRuntimeTemplates {

		if *existing.DisplayName == *template.DisplayName {
			if (*existing.Labels)["md5"] == checksum {
				return ExistsAndIsIdentical, &existing
			}
			return ExistsButIsDifferent, &existing
		}
	}
	return DoesNotExist, nil
}

func (nc *NotebookClient) GetNotebookRuntimeTemplates() ListNotebookRuntimeTemplatesResult {

	status, body := nc.curl("GET", nc.url, nil)

	if status != "200 OK" {
		logrus.Fatalf("Status returned: %s (%s)", status, string(body))
	}
	var templates ListNotebookRuntimeTemplatesResult

	json.Unmarshal(body, &templates)

	return templates
}

func (nc *NotebookClient) DeleteNotebookRuntimeTemplate(name string) {

	url := fmt.Sprintf("%s/%s", serviceEndpoint, name)
	logrus.Infof("Deleting: %s", url)
	nc.curl("DELETE", url, nil)
}

func (nc *NotebookClient) DeployNotebookRuntimeTemplateFromFile(path string) {

	template, checksum := readTemplateFile(path)

	resolveAction, existing := nc.resolve(template, checksum)

	if resolveAction == ExistsAndIsIdentical {
		logrus.Infof("Found existing template (%s) with same DisplayName and md5 hash, skipping ..", *existing.DisplayName)
		return
	}

	if resolveAction == ExistsButIsDifferent {
		logrus.Infof("Found existing template (%s) with same DisplayName and a different md5 hash, will delete existing one ..", *existing.DisplayName)
		nc.DeleteNotebookRuntimeTemplate(*existing.Name)
	}

	labels := map[string]string{
		"md5":        checksum,
		"git_sh":     os.Getenv("GITHUB_SHA"),
		"git_run_id": os.Getenv("GITHUB_RUN_ID"),
	}

	if template.Labels != nil {
		for key, value := range *template.Labels {
			labels[key] = value
		}
	}
	template.Labels = &labels

	payload, err := json.Marshal(template)

	if err != nil {
		logrus.Fatal("Error creating JSON Payload", err)
	}

	nc.curl("POST", nc.url, bytes.NewBuffer(payload))
}

func readTemplateFile(path string) (NotebookRuntimeTemplate, string) {

	bytes, err := os.ReadFile(path)
	if err != nil {
		logrus.Fatal("Could not read file", err)
	}

	hash := md5.New()

	hash.Write(bytes)
	checksum := hex.EncodeToString(hash.Sum(nil))

	var template NotebookRuntimeTemplate

	ext := filepath.Ext(path)

	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(bytes, &template); err != nil {
			logrus.Fatal("Could not unmarshall JSON file", err)
		}
	case ".json":
		if err := json.Unmarshal(bytes, &template); err != nil {
			logrus.Fatal("Could not unmarshall JSON file", err)
		}
	default:
		logrus.Fatalf("Unsupported file extension '%s'", ext)
	}
	return template, checksum
}
