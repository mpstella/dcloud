package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/mpstella/dcloud/pkg/gcp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	maximumConcurrentThreads int
	deploymentTimestampUTC   string
	notebookClient           *gcp.NotebookClient
)

func processFile(path string, wg *sync.WaitGroup, ch chan<- string, errCh chan<- error, sem chan struct{}) {

	defer wg.Done()
	defer func() { <-sem }() // Release the spot in the semaphore when the goroutine completes

	logrus.Infof("(%s) Parsing template", path)

	template := gcp.NewNotebookRuntimeTemplateFromFile(path)

	templates, err := notebookClient.GetNotebookRuntimeTemplates()

	if err != nil {
		errCh <- err
		return
	}

	var templateToDelete *gcp.NotebookRuntimeTemplate

	for _, existing := range templates.NotebookRuntimeTemplates {

		comparisonResult := template.ComparesTo(&existing)

		if comparisonResult == gcp.Identical {
			ch <- fmt.Sprintf("(%s) Found existing template with same DisplayName and a md5 hash, skipping ..", path)
			return
		}

		if comparisonResult == gcp.Different {
			logrus.Infof("(%s) Found existing template with same DisplayName and a different md5 hash, will delete existing one post deployment..", path)
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
	logrus.Infof("(%s) Deploying template.", path)
	err = notebookClient.DeployNotebookRuntimeTemplate(template)

	if err != nil {
		errCh <- err
		return
	}

	if templateToDelete != nil {
		logrus.Infof("(%s) Deleting template: %s", path, *templateToDelete.Name)
		err := notebookClient.DeleteNotebookRuntimeTemplate(*templateToDelete.Name)
		if err != nil {
			errCh <- err
			return
		}
	}
	ch <- fmt.Sprintf("(%s) Processed template.", path)
}

var deployCmd = &cobra.Command{
	Use:   "deploy [project] [pathToTemplates]",
	Short: "Deploy NotebookRuntimeTemplates",

	Run: func(cmd *cobra.Command, args []string) {

		templates, err := os.ReadDir(templateDirectory)

		if err != nil {
			logrus.Fatalf("Error occurred reading directory %v", err)
		}

		notebookClient, err = gcp.NewNotebookClient(projectID)

		if err != nil {
			logrus.Fatal(err)
		}

		var wg sync.WaitGroup
		contentCh := make(chan string, len(templates))
		errorCh := make(chan error, len(templates))
		sem := make(chan struct{}, maximumConcurrentThreads)

		for _, entry := range templates {

			if !entry.IsDir() {
				wg.Add(1)
				sem <- struct{}{} // Acquire a spot in the semaphore
				go processFile(filepath.Join(templateDirectory, entry.Name()), &wg, contentCh, errorCh, sem)
			}
		}

		go func() {
			wg.Wait()
			close(contentCh)
			close(errorCh)
		}()

		for content := range contentCh {
			logrus.Info(content)
		}

		// Handle any errors
		for err := range errorCh {
			logrus.Error(err)
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
