package gcp

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"fmt"

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
		log.Fatal(fmt.Errorf("Could not read file: %v", err))
	}

	hash := md5.New()

	hash.Write(bytes)
	checksum := hex.EncodeToString(hash.Sum(nil))

	var template NotebookRuntimeTemplate

	ext := filepath.Ext(path)

	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(bytes, &template); err != nil {
			log.Fatal(fmt.Errorf("could not unmarshall YAML file: %v", err))
		}
	case ".json":
		if err := json.Unmarshal(bytes, &template); err != nil {
			log.Fatal(fmt.Errorf("could not unmarshall JSON file: %v", err))
		}
	default:
		log.Fatalf(fmt.Sprintf("Unsupported file extension '%s'", ext))
	}

	template.AddLabel("md5", checksum)

	return &template
}

func (t *NotebookRuntimeTemplate) AddLabel(key string, value string) {

	if t.Labels == nil {
		labels := make(map[string]string)
		t.Labels = &labels
	}
	(*t.Labels)[key] = value
}

func (t *NotebookRuntimeTemplate) ComparesTo(c *NotebookRuntimeTemplate) TemplateComparison {

	if *t.DisplayName == *c.DisplayName {

		// deployed template does not have any labels
		if c.Labels == nil {

			fmt.Println("Warning: Found matching deployed template with no labels attached")
			return Different
		}

		// check if a 'md5' key exists in the deployed template's labels
		if _, exists := (*c.Labels)["md5"]; exists {

			if (*t.Labels)["md5"] == (*c.Labels)["md5"] {
				return Identical
			} else {
				return Different
			}

		}
		fmt.Println("Warning: Found matching template with no md5 label")
		return Different
	}
	return DoesNotMatch
}
