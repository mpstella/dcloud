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

var deploymentTimestampUTC string

type localTemplate struct {
	path                   string
	template               *gcp.NotebookRuntimeTemplate
	matchingRemoteTemplate *gcp.NotebookRuntimeTemplate
}

type deploymentActions struct {
	toBeDeployed []*localTemplate
	toBeDeleted  []*localTemplate
}

func sortItOut(nc *gcp.NotebookClient, localTemplates []*localTemplate) (*deploymentActions, error) {

	existingTemplates, err := nc.GetNotebookRuntimeTemplates()

	if err != nil {
		return nil, err
	}

	actions := deploymentActions{}

	for _, lt := range localTemplates {

		matchedTemplate, comparisonResult := existingTemplates.Compare(lt.template)

		if comparisonResult == gcp.Identical {
			fmt.Printf("Template '%s' matches '%s' - skipping\n", lt.path, *matchedTemplate.Name)
			continue
		}

		// if we get here we either have a new template or need to 'modify' an existing.
		actions.toBeDeployed = append(actions.toBeDeployed, lt)

		if comparisonResult == gcp.Different {
			fmt.Printf("Template '%s' matches '%s' but has changed - marking for future delete\n", lt.path, *matchedTemplate.Name)
			lt.matchingRemoteTemplate = matchedTemplate
			actions.toBeDeleted = append(actions.toBeDeleted, lt)
		}
	}
	return &actions, nil
}

func deployTemplate(nc *gcp.NotebookClient, template *gcp.NotebookRuntimeTemplate) error {

	// add a bunch of labels to the resource so we can track deployments
	template.AddLabel("deployment_ts_utc", deploymentTimestampUTC)

	if val, ok := os.LookupEnv("GIT_SHA"); ok {
		template.AddLabel("git_sha", val)
	}

	if val, ok := os.LookupEnv("GITHUB_RUN_ID"); ok {
		template.AddLabel("git_run_id", val)
	}

	err := nc.DeployNotebookRuntimeTemplate(template)

	if err != nil {
		return err
	}
	return nil
}

var deployCmd = &cobra.Command{
	Use:   "deploy [project] [pathToTemplates]",
	Short: "Deploy NotebookRuntimeTemplates",

	Run: func(cmd *cobra.Command, args []string) {

		templates, err := os.ReadDir(templateDirectory)

		if err != nil {
			log.Fatal(fmt.Errorf("error occurred reading directory %v", err))
		}

		var notebookRuntimeTemplates = make([]*localTemplate, len(templates))

		// let's read everything first in case we get an error
		for i, entry := range templates {

			templateFile := filepath.Join(templateDirectory, entry.Name())

			fmt.Printf("Reading template: %s\n", templateFile)

			notebookRuntimeTemplates[i] = &localTemplate{
				path:     templateFile,
				template: gcp.NewNotebookRuntimeTemplateFromFile(templateFile),
			}
		}

		nc, err := gcp.NewNotebookClient(projectID)

		if err != nil {
			log.Fatal(fmt.Errorf("could not createa a client %v", err))
		}

		actions, err := sortItOut(nc, notebookRuntimeTemplates)

		if err != nil {
			log.Fatal(fmt.Errorf("error in comparison %v", err))
		}

		fmt.Printf("Deploy Count: %d\nDelete Count: %d\n", len(actions.toBeDeployed), len(actions.toBeDeleted))

		for _, d := range actions.toBeDeployed {
			fmt.Printf("Deploying template: %s\n", d.path)
			deployTemplate(nc, d.template)
		}

		for _, d := range actions.toBeDeleted {
			existingName := *d.matchingRemoteTemplate.Name
			fmt.Printf("Deleting matched template (%s) -> %s\n", d.path, existingName)
			nc.DeleteNotebookRuntimeTemplate(existingName)
		}
	},
}

func init() {

	now := time.Now().UTC()
	deploymentTimestampUTC = now.Format("20060102_150405")

	deployCmd.PersistentFlags().StringVar(&projectID, "project", "", "GCP Project Name")
	deployCmd.PersistentFlags().StringVar(&templateDirectory, "templates", "", "Directory where templates are located")

	deployCmd.MarkPersistentFlagRequired("project")
	deployCmd.MarkPersistentFlagRequired("templates")

	rootCmd.AddCommand(deployCmd)
}
