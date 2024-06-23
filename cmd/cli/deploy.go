package cmd

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/mpstella/dcloud/pkg/gcp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const (
	doesNotExist = iota
	existsAndIsIdentical
	existsButIsDifferent
)

type templateComparison int

func resolve(nc *gcp.NotebookClient, template gcp.NotebookRuntimeTemplate, checksum string) (templateComparison, *gcp.NotebookRuntimeTemplate) {

	templates := nc.GetNotebookRuntimeTemplates()

	for _, existing := range templates.NotebookRuntimeTemplates {

		if *existing.DisplayName == *template.DisplayName {
			if (*existing.Labels)["md5"] == checksum {
				return existsAndIsIdentical, &existing
			}
			return existsButIsDifferent, &existing
		}
	}
	return doesNotExist, nil
}

func readTemplateFile(path string) (gcp.NotebookRuntimeTemplate, string) {

	bytes, err := os.ReadFile(path)
	if err != nil {
		logrus.Fatal("Could not read file", err)
	}

	hash := md5.New()

	hash.Write(bytes)
	checksum := hex.EncodeToString(hash.Sum(nil))

	var template gcp.NotebookRuntimeTemplate

	ext := filepath.Ext(path)

	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(bytes, &template); err != nil {
			logrus.Fatal("Could not unmarshall JSON file", err)
		}
	case ".json":
		if err := json.Unmarshal(bytes, &template); err != nil {
			logrus.Fatal("Could not unmarshall JSON file", err)
		}
	default:
		logrus.Fatalf("Unsupported file extension '%s'", ext)
	}
	return template, checksum
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

				template, checksum := readTemplateFile(templateFile)
				resolveAction, existing := resolve(&nc, template, checksum)

				if resolveAction == existsAndIsIdentical {
					logrus.Infof("Found existing template (%s) with same DisplayName and md5 hash, skipping ..", *existing.DisplayName)
					continue
				}

				if resolveAction == existsButIsDifferent {
					logrus.Infof("Found existing template (%s) with same DisplayName and a different md5 hash, will delete existing one ..", *existing.DisplayName)
					nc.DeleteNotebookRuntimeTemplate(*existing.Name)
				}

				labels := map[string]string{
					"md5":        checksum,
					"git_sh":     os.Getenv("GITHUB_SHA"),
					"git_run_id": os.Getenv("GITHUB_RUN_ID"),
				}

				if template.Labels != nil {
					for key, value := range *template.Labels {
						labels[key] = value
					}
				}
				template.Labels = &labels

				logrus.Infof("Attempting to deploy %s", templateFile)
				nc.DeployNotebookRuntimeTemplate(&template)
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
