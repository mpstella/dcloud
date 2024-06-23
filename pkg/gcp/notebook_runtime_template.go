package gcp

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	DoesNotMatch = iota
	Identical
	Different
)

type TemplateComparison int

func NewNotebookRuntimeTemplateFromFile(path string) *NotebookRuntimeTemplate {

	bytes, err := os.ReadFile(path)
	if err != nil {
		logrus.Fatal("Could not read file", err)
	}

	hash := md5.New()

	hash.Write(bytes)
	checksum := hex.EncodeToString(hash.Sum(nil))

	var template NotebookRuntimeTemplate

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

	template.AddLabel("md5", checksum)

	return &template
}

func (t *NotebookRuntimeTemplate) AddLabel(name string, value string) {
	if t.Labels == nil {
		labels := make(map[string]string)
		t.Labels = &labels
	}
	(*t.Labels)[name] = value
}

func (t *NotebookRuntimeTemplate) AddLabels(labels map[string]string) {
	if t.Labels != nil {
		for key, value := range *t.Labels {
			labels[key] = value
		}
	}
	t.Labels = &labels
}

func (t *NotebookRuntimeTemplate) ComparesTo(c *NotebookRuntimeTemplate) TemplateComparison {

	if *t.DisplayName == *c.DisplayName {

		// deployed template does not have any labels
		if c.Labels == nil {
			logrus.Warn("Found matching template with no labels attached")
			return Different
		}

		// check if a 'md5' k,v exists in the deployed template's labels
		if _, exists := (*c.Labels)["md5"]; exists {

			if (*t.Labels)["md5"] == (*c.Labels)["md5"] {
				return Identical
			} else {
				return Different
			}

		}
		logrus.Warn("Found matching template with no md5 label")
		return Different
	}
	return DoesNotMatch
}
