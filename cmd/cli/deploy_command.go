package cmd

import (
	"os"
	"path/filepath"
	"time"

	"github.com/mpstella/dcloud/pkg/gcp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var deploymentTimestampUTC string

func resolve(nc *gcp.NotebookClient, template *gcp.NotebookRuntimeTemplate) (gcp.TemplateComparison, *gcp.NotebookRuntimeTemplate) {

	// retrieving each time in case another template is deployed in the interim.
	templates := nc.GetNotebookRuntimeTemplates()

	for _, existing := range templates.NotebookRuntimeTemplates {

		result := template.ComparesTo(&existing)

		if result != gcp.DoesNotMatch {
			return result, &existing
		}
	}
	return gcp.DoesNotMatch, nil
}

var deployCmd = &cobra.Command{
	Use:   "deploy [project] [pathToTemplates]",
	Short: "Deploy NotebookRuntimeTemplates",

	Run: func(cmd *cobra.Command, args []string) {

		templates, err := os.ReadDir(templateDirectory)

		if err != nil {
			logrus.Fatalf("Error occurred reading directory %v", err)
		}

		nc := gcp.NewNotebookClient(projectID)

		for _, entry := range templates {

			if !entry.IsDir() {

				templateFile := filepath.Join(templateDirectory, entry.Name())

				logrus.Infof("Parsing template %s", templateFile)

				template := gcp.NewNotebookRuntimeTemplateFromFile(templateFile)

				cmp, existing := resolve(&nc, template)

				if cmp == gcp.Identical {
					logrus.Infof("Found existing template (%s) with same DisplayName and a md5 hash, skipping ..", *existing.DisplayName)
					continue
				}

				// add a bunch of labels to the resource so we can track deployments
				template.AddLabel("deployment_ts_utc", deploymentTimestampUTC)

				if val, ok := os.LookupEnv("GIT_SHA"); ok {
					template.AddLabel("git_sha", val)
				}

				if val, ok := os.LookupEnv("GITHUB_RUN_ID"); ok {
					template.AddLabel("git_run_id", val)
				}

				// deploy first as this does not impact any existing templates
				logrus.Infof("Attempting to deploy %s", templateFile)
				nc.DeployNotebookRuntimeTemplate(template)

				// now delete the duplicate
				if existing != nil {
					logrus.Infof("Found existing template (%s) with same DisplayName and a different md5 hash, will delete existing one ..", *existing.DisplayName)
					nc.DeleteNotebookRuntimeTemplate(*existing.Name)
				}
			}
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
