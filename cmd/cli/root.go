package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mpstella/dcloud/pkg/gcp"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var (
	projectID         string
	templateName      string
	templateDirectory string
)

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Format the time
	timestamp := time.Now().Format("2006/01/02 15:04:05")

	// Create a buffer to write the formatted log entry
	var b bytes.Buffer

	// Write the formatted log entry
	b.WriteString(fmt.Sprintf("%s %s: [%s] %s\n", timestamp, strings.ToUpper(entry.Level.String()), "collab", entry.Message))

	return b.Bytes(), nil
}

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
				logrus.Infof("Attempting to deploy %s", templateFile)
				nc.DeployNotebookRuntimeTemplateFromFile(templateFile)
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

	listCmd.PersistentFlags().StringVar(&projectID, "project", "", "GCP Project Name")
	listCmd.MarkPersistentFlagRequired("project")

	deployCmd.PersistentFlags().StringVar(&projectID, "project", "", "GCP Project Name")
	deployCmd.PersistentFlags().StringVar(&templateDirectory, "templates", "", "Directory where templates are located")

	deployCmd.MarkPersistentFlagRequired("project")
	deployCmd.MarkPersistentFlagRequired("templates")

	deleteCmd.PersistentFlags().StringVar(&templateName, "name", "", "Name of the template")
	deleteCmd.MarkPersistentFlagRequired("name")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// add all the commands here
	rootCmd.AddCommand(deployCmd, listCmd, deleteCmd)
}
