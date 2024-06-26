package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/mpstella/dcloud/pkg/gcp"

	"github.com/spf13/cobra"
)

var (
	maximumConcurrentThreads int
	deploymentTimestampUTC   string
	notebookClient           *gcp.NotebookClient
)

func processFile(path string) {

	fmt.Printf("Parsing template: %s\n", path)

	template := gcp.NewNotebookRuntimeTemplateFromFile(path)

	templates, err := notebookClient.GetNotebookRuntimeTemplates()

	var templateToDelete *gcp.NotebookRuntimeTemplate

	for _, existing := range templates.NotebookRuntimeTemplates {

		comparisonResult := template.ComparesTo(&existing)

		if comparisonResult == gcp.Identical {
			fmt.Println("Found existing template with same DisplayName and a md5 hash, skipping ..")
			return
		}

		if comparisonResult == gcp.Different {
			fmt.Println("Found existing template with same DisplayName and a different md5 hash, will delete existing one post deployment..")
			templateToDelete = &existing
			break
		}
	}
	// if we get to here we are ready to deploy

	// add a bunch of labels to the resource so we can track deployments
	template.AddLabel("deployment_ts_utc", deploymentTimestampUTC)

	if val, ok := os.LookupEnv("GIT_SHA"); ok {
		template.AddLabel("git_sha", val)
	}

	if val, ok := os.LookupEnv("GITHUB_RUN_ID"); ok {
		template.AddLabel("git_run_id", val)
	}

	// deploy first as this does not impact any existing templates
	fmt.Println("Deploying template.")
	err = notebookClient.DeployNotebookRuntimeTemplate(template)

	if err != nil {
		log.Fatalf("Deploy failed %v\n", err)
	}

	if templateToDelete != nil {
		fmt.Printf("Deleting template: %s\n", *templateToDelete.Name)
		err := notebookClient.DeleteNotebookRuntimeTemplate(*templateToDelete.Name)

		if err != nil {
			log.Fatalf("Deletion failed: %v\n", err)
		}
	}
	fmt.Println("Processed template.")
}

var deployCmd = &cobra.Command{
	Use:   "deploy [project] [pathToTemplates]",
	Short: "Deploy NotebookRuntimeTemplates",

	Run: func(cmd *cobra.Command, args []string) {

		templates, err := os.ReadDir(templateDirectory)

		if err != nil {
			log.Fatal(fmt.Errorf("error occurred reading directory %v", err))
		}

		notebookClient, err = gcp.NewNotebookClient(projectID)

		if err != nil {
			log.Fatal(fmt.Errorf("could not createa a client %v", err))
		}

		for _, entry := range templates {

			if !entry.IsDir() {
				processFile(filepath.Join(templateDirectory, entry.Name()))
			}
		}
	},
}

func init() {

	now := time.Now().UTC()
	deploymentTimestampUTC = now.Format("20060102_150405")

	deployCmd.PersistentFlags().StringVar(&projectID, "project", "", "GCP Project Name")
	deployCmd.PersistentFlags().StringVar(&templateDirectory, "templates", "", "Directory where templates are located")
	deployCmd.PersistentFlags().IntVar(&maximumConcurrentThreads, "threads", 1, "Number of concurrent threads")

	deployCmd.MarkPersistentFlagRequired("project")
	deployCmd.MarkPersistentFlagRequired("templates")

	rootCmd.AddCommand(deployCmd)
}
