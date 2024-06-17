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
	"google.golang.org/grpc/grpclog"

	"github.com/spf13/cobra"
)

var (
	projectID         string
	templateName      string
	templateDirectory string
	dryRun            bool
	silentMode        bool
	debugMode         bool
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

func initConfig() {

	logrus.SetFormatter(&CustomFormatter{})

	if silentMode {
		logrus.SetLevel(logrus.WarnLevel)
	} else if debugMode {
		logrus.SetLevel(logrus.DebugLevel)
		grpclog.SetLoggerV2(grpclog.NewLoggerV2WithVerbosity(os.Stdout, os.Stderr, os.Stderr, 99))
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func prettyPrinter(arg interface{}) {

	prettyString, _ := json.MarshalIndent(arg, "", "  ")

	if silentMode {
		fmt.Println(string(prettyString))
	} else {
		logrus.Infof("%s", string(prettyString))
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
			logrus.Fatalf("Error occurred reading directory %v", err)
		}

		cc := gcp.NewCollabClient(projectID)
		defer cc.Cleanup()

		for _, entry := range templates {
			if !entry.IsDir() {

				templateFile := filepath.Join(templateDirectory, entry.Name())
				logrus.Infof("Attempting to deploy %s", templateFile)
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
	rootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "Enable debugging of application")

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
