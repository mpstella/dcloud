package cmd

import (
	"dcloud/pkg/gcp"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var projectID string
var displayName string
var templateDirectory string

var rootCmd = &cobra.Command{
	Use:   "dcloud",
	Short: "DAW gcloud makeshift utility",
	Long:  `Given we are waiting on Google this is our dodgy gcloud helper`,
}

var listCmd = &cobra.Command{
	Use:   "list [project]",
	Short: "Retrieve existing NotebookRuntimeTemplates",

	Run: func(cmd *cobra.Command, args []string) {

		cc := gcp.NewCollabClient(projectID)
		defer cc.Cleanup()

		existingTemplates := cc.GetNotebookRuntimeTemplates()

		for _, template := range existingTemplates {
			prettyString, err := json.MarshalIndent(template, "", "  ")
			if err != nil {
				fmt.Printf("%s\n", template)
			} else {
				fmt.Printf("%s\n", string(prettyString))
			}
		}
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [project] [name]",
	Short: "Delete an existing NotebookRuntimeTemplate by DisplayName",

	Run: func(cmd *cobra.Command, args []string) {

		cc := gcp.NewCollabClient(projectID)
		defer cc.Cleanup()

		cc.DeleteNotebookRuntimeTemplate(displayName)
	},
}

var deployCmd = &cobra.Command{
	Use:   "deploy [project]",
	Short: "Deploy NotebookRuntimeTemplates",

	Run: func(cmd *cobra.Command, args []string) {

		templates, err := os.ReadDir(templateDirectory)

		if err != nil {
			log.Fatalf("Error occurred reading directory %v\n", err)
		}

		cc := gcp.NewCollabClient(projectID)
		defer cc.Cleanup()

		for _, entry := range templates {
			if !entry.IsDir() {

				templateFile := filepath.Join(templateDirectory, entry.Name())
				fmt.Printf("Attempting to deploy %s\n", templateFile)
				cc.DeployNotebookRuntimeTemplate(templateFile)
			}
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().StringVarP(&projectID, "project", "p", "", "GCP Project Name")
	rootCmd.MarkPersistentFlagRequired("project")

	rootCmd.AddCommand(listCmd)

	deployCmd.PersistentFlags().StringVarP(&templateDirectory, "templates", "t", "", "Directory where templates are located")
	deployCmd.MarkPersistentFlagRequired("templates")
	rootCmd.AddCommand(deployCmd)

	deleteCmd.PersistentFlags().StringVarP(&displayName, "name", "n", "", "Display Name of the template")
	deleteCmd.MarkPersistentFlagRequired("name")
	rootCmd.AddCommand(deleteCmd)

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
