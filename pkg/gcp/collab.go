package gcp

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/durationpb"
	"log"
	"os"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	aiplatformpb "cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
)

const location = "australia-southeast1"
const endPoint = "australia-southeast1-aiplatform.googleapis.com:443"
const scopes = "https://www.googleapis.com/auth/cloud-platform"

func fullyQualifiedParent(projectID string) string {
	return fmt.Sprintf("projects/%s/locations/%s", projectID, location)
}

func NewClientUsingFile(path string) *aiplatform.NotebookClient {
	ctx := context.Background()
	client, err := aiplatform.NewNotebookClient(ctx,
		option.WithCredentialsFile(path),
		option.WithEndpoint(endPoint),
	)
	if err != nil {
		log.Fatal("Could not create Client", err)
	}
	return client
}

func NewClientUsingEnv() *aiplatform.NotebookClient {
	return NewClientUsingFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
}

func NewClientUsingADC() *aiplatform.NotebookClient {
	ctx := context.Background()
	credentials, err := google.FindDefaultCredentials(ctx, scopes)
	if err != nil {
		panic("Could not obtain ADC credentials")
	}
	client, err := aiplatform.NewNotebookClient(ctx,
		option.WithCredentials(credentials),
		option.WithEndpoint(endPoint),
	)
	if err != nil {
		log.Fatal("Could not create Client", err)
	}
	return client
}

func GetNotebookRuntimeTemplates(client *aiplatform.NotebookClient, projectID string) (map[string]string, error) {

	var existingTemplates = make(map[string]string)

	ctx := context.Background()

	// Define the request to list Notebook Runtime Templates
	req := &aiplatformpb.ListNotebookRuntimeTemplatesRequest{
		Parent: fullyQualifiedParent(projectID),
	}

	// List the Notebook Runtime Templates
	it := client.ListNotebookRuntimeTemplates(ctx, req)
	for {
		template, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to list Notebook Runtime Templates: %v", err)
		}
		existingTemplates[template.GetDisplayName()] = template.GetName()
	}
	return existingTemplates, nil
}

func DeployNotebookRuntimeTemplate(client *aiplatform.NotebookClient, projectID string, templateFile string, ensureUnique bool) {
	ctx := context.Background()

	var config aiplatformpb.NotebookRuntimeTemplate

	data, err := os.ReadFile(templateFile)

	if err != nil {
		log.Fatalf("Error reading file %v\n", err)
	}

	err = json.Unmarshal(data, &config)

	if err != nil {
		log.Fatalf("Error parsing JSON file: %v", err)
	}

	if ensureUnique {
		if GetNotebookRuntimeTemplateNameByDisplayName(client, projectID, config.DisplayName) != "" {
			fmt.Printf("A template already exists with this Display Name, skipping ...\n")
			return
		}
	}

	req := &aiplatformpb.CreateNotebookRuntimeTemplateRequest{
		Parent:                  fullyQualifiedParent(projectID),
		NotebookRuntimeTemplate: &config,
		//NotebookRuntimeTemplateId: "my-notebook-runtime-template",
	}
	resp, err := client.CreateNotebookRuntimeTemplate(ctx, req)
	if err != nil {
		log.Fatalf("Failed to create Notebook Runtime Template: %v", err)
	}
	fmt.Printf("Created Notebook Runtime Template: %v\n", resp)
}

func GenerateSampleTemplate() {

	machineSpec := &aiplatformpb.MachineSpec{
		MachineType:     "e2-standard-2",
		AcceleratorType: aiplatformpb.AcceleratorType_ACCELERATOR_TYPE_UNSPECIFIED,
	}

	networkSpec := &aiplatformpb.NetworkSpec{
		EnableInternetAccess: true,
		Subnetwork:           "subnetwork",
		Network:              "network",
	}

	duration := &durationpb.Duration{
		Seconds: 600,
	}

	idleShutdownConfig := &aiplatformpb.NotebookIdleShutdownConfig{
		IdleTimeout:          duration,
		IdleShutdownDisabled: false,
	}

	persistentDiskSpec := &aiplatformpb.PersistentDiskSpec{
		DiskType:   "pd-standard",
		DiskSizeGb: 10,
	}

	runtimeTemplate := &aiplatformpb.NotebookRuntimeTemplate{
		DisplayName:            "Test 123 this is my request",
		Description:            "This is a test template deployed by dcloud",
		IsDefault:              false,
		IdleShutdownConfig:     idleShutdownConfig,
		MachineSpec:            machineSpec,
		NetworkSpec:            networkSpec,
		DataPersistentDiskSpec: persistentDiskSpec,
	}

	jsonData, err := json.Marshal(runtimeTemplate)
	if err != nil {
		fmt.Println("Error converting to JSON", err)
	}
	fmt.Println(string(jsonData))
}

func GetNotebookRuntimeTemplateNameByDisplayName(client *aiplatform.NotebookClient, projectID string, displayName string) string {

	existingTemplates, err := GetNotebookRuntimeTemplates(client, projectID)

	if err != nil {
		panic("Couldn't get list of existing NotebookRuntimeTemplate")
	}

	if value, ok := existingTemplates[displayName]; ok {
		return value
	}
	return ""

}

func DeleteNotebookRuntimeTemplate(client *aiplatform.NotebookClient, templateId string) error {

	ctx := context.Background()

	req := &aiplatformpb.DeleteNotebookRuntimeTemplateRequest{
		Name: templateId,
	}

	resp, err := client.DeleteNotebookRuntimeTemplate(ctx, req)
	if err != nil {
		log.Fatalf("Failed to delete Notebook Runtime Template: %v", err)
	}

	fmt.Printf("Deleted Notebook Runtime Template %v\n\n", resp)
	return err
}
