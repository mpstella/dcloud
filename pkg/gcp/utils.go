package gcp

func (list *ListNotebookRuntimeTemplatesResult) ToMap() *map[string]*NotebookRuntimeTemplate {

	existingTemplates := make(map[string]*NotebookRuntimeTemplate)

	// create a map of existing templates
	for _, x := range list.NotebookRuntimeTemplates {
		existingTemplates[*x.DisplayName] = &x
	}
	return &existingTemplates
}

func (list *ListNotebookRuntimeTemplatesResult) Compare(t *NotebookRuntimeTemplate) (*NotebookRuntimeTemplate, TemplateComparison) {

	for _, existing := range list.NotebookRuntimeTemplates {
		if *t.DisplayName == *existing.DisplayName {
			return &existing, existing.ComparesTo(t)
		}
	}
	return nil, DoesNotMatch
}
