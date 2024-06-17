package gcp

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	aiplatformpb "cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
)

const location = "australia-southeast1"
const endPoint = "australia-southeast1-aiplatform.googleapis.com:443"
const scopes = "https://www.googleapis.com/auth/cloud-platform"

type RuntimeTemplate struct {
	Name        string `json:"Name"`
	DisplayName string `json:"DisplayName"`
	Description string `json:"Description"`
	FileHash    string `json:"FileHash"`
	MachineType string `json:"MachineType"`
}

type CollabClient struct {
	projectID         string
	client            *aiplatform.NotebookClient
	existingTemplates map[string]RuntimeTemplate
	isInitialised     bool
}

func fullyQualifiedParent(projectID string) string {
	return fmt.Sprintf("projects/%s/locations/%s", projectID, location)
}

func NewCollabClient(projectID string) CollabClient {

	client, err := getClient()

	if err != nil {
		logrus.Fatal("Could not create Client", err)
	}

	return CollabClient{
		projectID:         projectID,
		client:            client,
		existingTemplates: make(map[string]RuntimeTemplate),
		isInitialised:     false,
	}
}

func getClient() (*aiplatform.NotebookClient, error) {

	ctx := context.Background()
	path := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	if path != "" {

		logrus.Infof("Logging onto GCP using credentials file: %s ...", path)

		return aiplatform.NewNotebookClient(ctx,
			option.WithCredentialsFile(path),
			option.WithEndpoint(endPoint),
		)
	}

	logrus.Info("Logging onto GCP using ADC ...")

	credentials, err := google.FindDefaultCredentials(ctx, scopes)
	if err != nil {
		logrus.Info("Could not obtain ADC credentials")
		return nil, err
	}
	return aiplatform.NewNotebookClient(ctx,
		option.WithCredentials(credentials),
		option.WithEndpoint(endPoint),
	)
}

func (cc *CollabClient) GetNotebookRuntimeTemplate(name string) (*aiplatformpb.NotebookRuntimeTemplate, error) {

	ctx := context.Background()

	req := &aiplatformpb.GetNotebookRuntimeTemplateRequest{
		Name: name,
	}

	return cc.client.GetNotebookRuntimeTemplate(ctx, req)
}

func (c *CollabClient) GetNotebookRuntimeTemplates() map[string]RuntimeTemplate {

	if c.isInitialised {
		return c.existingTemplates
	}

	logrus.Info("Retrieving existing deployed runtime templates ...")

	ctx := context.Background()

	// Define the request to list Notebook Runtime Templates
	req := &aiplatformpb.ListNotebookRuntimeTemplatesRequest{
		Parent: fullyQualifiedParent(c.projectID),
	}

	// List the Notebook Runtime Templates
	it := c.client.ListNotebookRuntimeTemplates(ctx, req)
	for {
		template, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			logrus.Fatalf("Failed to list Notebook Runtime Templates: %v", err)
		}

		c.existingTemplates[template.GetDisplayName()] = RuntimeTemplate{
			Name:        template.GetName(),
			DisplayName: template.GetDisplayName(),
			Description: template.GetDescription(),
			FileHash:    template.GetLabels()["md5"],
			MachineType: template.MachineSpec.GetMachineType(),
		}
	}
	return c.existingTemplates
}

func (c *CollabClient) Cleanup() {

	if c.client != nil {
		logrus.Info("Closing connection to GCP ...")
		c.client.Close()
	}
}

func (c *CollabClient) DeployNotebookRuntimeTemplate(templateFile string) {

	ctx := context.Background()

	data, err := os.ReadFile(templateFile)

	if err != nil {
		logrus.Fatalf("Error reading file %v\n", err)
	}

	var config aiplatformpb.NotebookRuntimeTemplate
	err = json.Unmarshal(data, &config)

	if err != nil {
		logrus.Fatalf("Error parsing JSON file: %v", err)
	}

	hash := md5.New()

	_, err = hash.Write(data)
	if err != nil {
		logrus.Fatalf("Failed to write data to hash: %v", err)
	}

	checksum := hex.EncodeToString(hash.Sum(nil))

	if existingTemplate, ok := c.GetNotebookRuntimeTemplates()[config.DisplayName]; ok {

		logrus.Info("A template already exists with this Display Name, will check for changes ...")

		if checksum == existingTemplate.FileHash {
			logrus.Infof("Template hash matches ('%s') skipping ...\n", checksum)
			return
		} else {
			logrus.Info("Will delete existing template and redeploy")
			c.DeleteNotebookRuntimeTemplate(existingTemplate.DisplayName)
		}
	}

	if config.Labels == nil {
		config.Labels = make(map[string]string)
	}

	config.Labels["md5"] = checksum

	// add some GITHUB goodness if running in CI/CD
	config.Labels["git_sha"] = os.Getenv("GITHUB_SHA")
	config.Labels["git_run_id"] = os.Getenv("GITHUB_RUN_ID")

	req := &aiplatformpb.CreateNotebookRuntimeTemplateRequest{
		Parent:                  fullyQualifiedParent(c.projectID),
		NotebookRuntimeTemplate: &config,
	}

	resp, err := c.client.CreateNotebookRuntimeTemplate(ctx, req)
	if err != nil {
		logrus.Fatalf("Failed to create Notebook Runtime Template: %v", err)
	}

	// add to cache to ensure uniqueness
	c.existingTemplates[config.DisplayName] = RuntimeTemplate{
		DisplayName: config.DisplayName,
		FileHash:    checksum,
	}

	logrus.Infof("Created Notebook Runtime Template: %v\n", resp)
}

func (c *CollabClient) DeleteNotebookRuntimeTemplate(name string) {

	ctx := context.Background()

	req := &aiplatformpb.DeleteNotebookRuntimeTemplateRequest{
		Name: name,
	}

	resp, err := c.client.DeleteNotebookRuntimeTemplate(ctx, req)
	if err != nil {
		logrus.Fatalf("Failed to delete Notebook Runtime Template: %v", err)
	}
	logrus.Infof("Deleted Notebook Runtime Template %+v", resp)

}
