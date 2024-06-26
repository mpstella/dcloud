package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mpstella/dcloud/pkg/gcp"

	"github.com/spf13/cobra"
)

var (
	projectID         string
	templateName      string
	templateDirectory string
	version           string = "v0.1.6"
)

func prettyPrinter(arg interface{}) {
	prettyString, _ := json.MarshalIndent(arg, "", "  ")
	fmt.Printf("%s\n", string(prettyString))
}

var rootCmd = &cobra.Command{
	Use:   "dcloud",
	Short: "DAW gcloud makeshift utility",
	Long:  `Given we are waiting on Google this is our dodgy gcloud helper`,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version number of dcloud",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\n", version)
	},
}

var listCmd = &cobra.Command{
	Use:   "list [project]",
	Short: "Retrieve existing NotebookRuntimeTemplates",

	Run: func(cmd *cobra.Command, args []string) {

		nc, err := gcp.NewNotebookClient(projectID)

		if err != nil {
			log.Fatalf("Could not get client %v\n", err)
		}

		existingTemplates, err := nc.GetNotebookRuntimeTemplates()

		if err != nil {
			log.Fatalf("Failed to retreive templates: %v\n", err)
		}

		prettyPrinter(existingTemplates.NotebookRuntimeTemplates)

	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete an existing NotebookRuntimeTemplate",

	Run: func(cmd *cobra.Command, args []string) {

		nc, err := gcp.NewNotebookClient(projectID)

		if err != nil {
			log.Fatal(err)
		}

		err = nc.DeleteNotebookRuntimeTemplate(templateName)
		if err != nil {
			log.Fatal(err)
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

	listCmd.PersistentFlags().StringVar(&projectID, "project", "", "GCP Project Name")
	listCmd.MarkPersistentFlagRequired("project")

	deleteCmd.PersistentFlags().StringVar(&templateName, "name", "", "Name of the template")
	deleteCmd.MarkPersistentFlagRequired("name")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// add all the commands here
	rootCmd.AddCommand(versionCmd, listCmd, deleteCmd)
}
