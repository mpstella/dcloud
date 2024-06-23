package cmd

import (
	"os"
	"path/filepath"

	"github.com/mpstella/dcloud/pkg/gcp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

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

				if existing != nil {
					logrus.Infof("Found existing template (%s) with same DisplayName and a different md5 hash, will delete existing one ..", *existing.DisplayName)
					nc.DeleteNotebookRuntimeTemplate(*existing.Name)
				}

				template.AddLabel("git_sha", os.Getenv("GITHUB_SHA"))
				template.AddLabel("git_run_id", os.Getenv("GITHUB_RUN_ID"))

				logrus.Infof("Attempting to deploy %s", templateFile)
				nc.DeployNotebookRuntimeTemplate(template)
			}
		}
	},
}

func init() {

	deployCmd.PersistentFlags().StringVar(&projectID, "project", "", "GCP Project Name")
	deployCmd.PersistentFlags().StringVar(&templateDirectory, "templates", "", "Directory where templates are located")

	deployCmd.MarkPersistentFlagRequired("project")
	deployCmd.MarkPersistentFlagRequired("templates")

	rootCmd.AddCommand(deployCmd)
}
