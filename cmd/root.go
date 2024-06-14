package cmd

import (
	"dcloud/pkg/gcp"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var projectID string
var displayName string
var templateDirectory string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dcloud",
	Short: "DAW gcloud",
	Long:  `Given we are waiting on Google this is our dodgy gcloud helper`,
}

var listCmd = &cobra.Command{
	Use:   "list [project]",
	Short: "Retrieve existing NotebookRuntimeTemplates",

	Run: func(cmd *cobra.Command, args []string) {

		client := gcp.NewClientUsingADC()
		defer client.Close()

		existing, err := gcp.GetNotebookRuntimeTemplates(client, projectID)

		if err != nil {
			panic("Couldn't retrieve existing templates")
		}
		for displayName, name := range existing {
			fmt.Printf("'%s' -> %s\n", displayName, name)
		}
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [project] [name]",
	Short: "Delete an existing NotebookRuntimeTemplate by DisplayName",

	Run: func(cmd *cobra.Command, args []string) {

		client := gcp.NewClientUsingADC()
		defer client.Close()

		templateId := gcp.GetNotebookRuntimeTemplateNameByDisplayName(client, projectID, displayName)
		fmt.Printf("Found the template: %s\n", templateId)

		if templateId != "" {
			gcp.DeleteNotebookRuntimeTemplate(client, templateId)
		}
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

		client := gcp.NewClientUsingADC()
		defer client.Close()

		for _, entry := range templates {
			if !entry.IsDir() {

				templateFile := filepath.Join(templateDirectory, entry.Name())
				fmt.Printf("Attempting to deploy %s\n", templateFile)
				gcp.DeployNotebookRuntimeTemplate(client, projectID, templateFile, true)
			}
		}
		//gcp.CreateNotebookRuntimeTemplate(client, projectID, true)
	},
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates a sample JSON Template",

	Run: func(cmd *cobra.Command, args []string) {
		gcp.GenerateSampleTemplate()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
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

	rootCmd.AddCommand(generateCmd)

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
