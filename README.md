DAW gcloud for NotebookRuntimeTemplates

## CLI Options
```sh
> go run dcloud
Given we are waiting on Google this is our dodgy gcloud helper

Usage:
  dcloud [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  delete      Delete an existing NotebookRuntimeTemplate by DisplayName
  deploy      Deploy NotebookRuntimeTemplates
  generate    Generates a sample JSON Template
  help        Help about any command
  list        Retrieve existing NotebookRuntimeTemplates

Flags:
  -h, --help             help for dcloud
  -p, --project string   GCP Project Name
  -t, --toggle           Help message for toggle

Use "dcloud [command] --help" for more information about a command.
```

## Deploy a notebook runtime template
```sh
$> go run dcloud deploy -p gamma-priceline-playground -t templates
Attempting to deploy templates/sample.json
Created Notebook Runtime Template: &{0x14000178168}
```

## Get deployed notebook runtime templates
```sh
$> go run dcloud list -h
Retrieve existing NotebookRuntimeTemplates

Usage:
  dcloud list [project] [flags]

Flags:
  -h, --help   help for list

Global Flags:
  -p, --project string   GCP Project Name

$> go run dcloud list -p gamma-priceline-playground
'Test 123 this is my request' -> projects/1019340507365/locations/australia-southeast1/notebookRuntimeTemplates/7327778806196862976
```

## Delete deployed notebook runtime templates

```sh
$> go run dcloud delete -p gamma-priceline-playground -n "Test 123 this is my request"
Found the template: projects/1019340507365/locations/australia-southeast1/notebookRuntimeTemplates/7327778806196862976
Deleted Notebook Runtime Template &{0x140001a76e0}
                                
```
