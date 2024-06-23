package cmd

import (
	"encoding/json"
	"os"

	"github.com/mpstella/dcloud/pkg/gcp"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var (
	projectID         string
	templateName      string
	templateDirectory string
)

func prettyPrinter(arg interface{}) {
	prettyString, _ := json.MarshalIndent(arg, "", "  ")
	logrus.Infof("%s", string(prettyString))
}

var rootCmd = &cobra.Command{
	Use:   "dcloud",
	Short: "DAW gcloud makeshift utility",
	Long:  `Given we are waiting on Google this is our dodgy gcloud helper`,
}

var listCmd = &cobra.Command{
	Use:   "list [project]",
	Short: "Retrieve existing NotebookRuntimeTemplates",

	Run: func(cmd *cobra.Command, args []string) {

		nc := gcp.NewNotebookClient(projectID)

		existingTemplates := nc.GetNotebookRuntimeTemplates()

		for _, template := range existingTemplates.NotebookRuntimeTemplates {
			prettyPrinter(template)
		}
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete an existing NotebookRuntimeTemplate",

	Run: func(cmd *cobra.Command, args []string) {

		nc := gcp.NewNotebookClient(projectID)
		nc.DeleteNotebookRuntimeTemplate(templateName)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	listCmd.PersistentFlags().StringVar(&projectID, "project", "", "GCP Project Name")
	listCmd.MarkPersistentFlagRequired("project")

	deleteCmd.PersistentFlags().StringVar(&templateName, "name", "", "Name of the template")
	deleteCmd.MarkPersistentFlagRequired("name")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// add all the commands here
	rootCmd.AddCommand(listCmd, deleteCmd)
}
