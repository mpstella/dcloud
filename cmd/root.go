package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mpstella/dcloud/pkg/gcp"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var (
	projectID         string
	templateName      string
	templateDirectory string
	dryRun            bool
	silentMode        bool
)

func initConfig() {
	if silentMode {
		logrus.SetLevel(logrus.WarnLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func prettyPrinter(arg interface{}) {

	prettyString, _ := json.MarshalIndent(arg, "", "  ")

	if silentMode {
		fmt.Println(string(prettyString))
	} else {
		logrus.Infof("%s\n", string(prettyString))
	}
}

var rootCmd = &cobra.Command{
	Use:   "dcloud",
	Short: "DAW gcloud makeshift utility",
	Long:  `Given we are waiting on Google this is our dodgy gcloud helper`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
	},
}

var listCmd = &cobra.Command{
	Use:   "list [project]",
	Short: "Retrieve existing NotebookRuntimeTemplates",

	Run: func(cmd *cobra.Command, args []string) {

		cc := gcp.NewCollabClient(projectID)
		defer cc.Cleanup()

		existingTemplates := cc.GetNotebookRuntimeTemplates()

		for _, template := range existingTemplates {
			prettyPrinter(template)
		}
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete an existing NotebookRuntimeTemplate",

	Run: func(cmd *cobra.Command, args []string) {

		cc := gcp.NewCollabClient(projectID)
		defer cc.Cleanup()

		cc.DeleteNotebookRuntimeTemplate(templateName)
	},
}

var deployCmd = &cobra.Command{
	Use:   "deploy [project] [pathToTemplates]",
	Short: "Deploy NotebookRuntimeTemplates",

	Run: func(cmd *cobra.Command, args []string) {

		templates, err := os.ReadDir(templateDirectory)

		if err != nil {
			logrus.Fatalf("Error occurred reading directory %v\n", err)
		}

		cc := gcp.NewCollabClient(projectID)
		defer cc.Cleanup()

		for _, entry := range templates {
			if !entry.IsDir() {

				templateFile := filepath.Join(templateDirectory, entry.Name())
				logrus.Infof("Attempting to deploy %s\n", templateFile)
				cc.DeployNotebookRuntimeTemplate(templateFile, dryRun)
			}
		}
	},
}

var exportCmd = &cobra.Command{
	Use:   "export [name]",
	Short: "Export an existing NotebookRutimeTemplate",

	Run: func(cmd *cobra.Command, args []string) {

		cc := gcp.NewCollabClient(projectID)
		defer cc.Cleanup()

		template, err := cc.GetNotebookRuntimeTemplate(templateName)

		if err != nil {
			logrus.Fatal("Failed to retrieve RuntimeTemplate", err)
		}

		prettyPrinter(template)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().BoolVar(&silentMode, "silent", false, "Minimise output to stdout")

	listCmd.PersistentFlags().StringVar(&projectID, "project", "", "GCP Project Name")
	listCmd.MarkPersistentFlagRequired("project")

	deployCmd.PersistentFlags().StringVar(&projectID, "project", "", "GCP Project Name")
	deployCmd.PersistentFlags().StringVar(&templateDirectory, "templates", "", "Directory where templates are located")
	deployCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Run the command in dry-run mode")

	deployCmd.MarkPersistentFlagRequired("project")
	deployCmd.MarkPersistentFlagRequired("templates")

	deleteCmd.PersistentFlags().StringVar(&templateName, "name", "", "Name of the template")
	deleteCmd.MarkPersistentFlagRequired("name")

	exportCmd.PersistentFlags().StringVar(&templateName, "name", "", "Name of the template")
	exportCmd.MarkPersistentFlagRequired("name")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// add all the commands here
	rootCmd.AddCommand(deployCmd, listCmd, deleteCmd, exportCmd)
}
