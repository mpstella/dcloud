DAW gcloud for NotebookRuntimeTemplates

## CLI Options
```sh
> ./dcloud --help
Given we are waiting on Google this is our dodgy gcloud helper

Usage:
  dcloud [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  delete      Delete an existing NotebookRuntimeTemplate by DisplayName
  deploy      Deploy NotebookRuntimeTemplates
  help        Help about any command
  list        Retrieve existing NotebookRuntimeTemplates

Flags:
  -h, --help             help for dcloud
  -p, --project string   GCP Project Name
  -t, --toggle           Help message for toggle

Use "dcloud [command] --help" for more information about a command.
```

## Deploy a notebook runtime template
### help
```sh
$> ./dcloud deploy -h                                        

Deploy NotebookRuntimeTemplates

Usage:
  dcloud deploy [project] [flags]

Flags:
  -h, --help               help for deploy
  -t, --templates string   Directory where templates are located

Global Flags:
  -p, --project string   GCP Project Name
```

### example
```sh
$> ./dcloud deploy -p gamma-priceline-playground -t templates

Logging onto GCP using ADC ...
Attempting to deploy templates/sample.json
Retrieving existing deployed runtime templates ...
Could not find any existing runtime templates
Created Notebook Runtime Template: &{0x14000502810}
Attempting to deploy templates/sample2.json
Retrieving existing deployed runtime templates ...
Found an existing runtime template
Created Notebook Runtime Template: &{0x14000300078}
Closing connection to GCP ...
```

### idempotent
```sh
$> ./dcloud deploy -p gamma-priceline-playground -t templates

Logging onto GCP using ADC ...
Attempting to deploy templates/sample.json
Retrieving existing deployed runtime templates ...
Found 2 existing runtime templates
A template already exists with this Display Name, skipping ...
Template hash matches ('18879ff2f8ee028deff132b953c534b4')  skipping ...
Attempting to deploy templates/sample2.json
Retrieving existing deployed runtime templates ...
Found 2 existing runtime templates
A template already exists with this Display Name, skipping ...
Template hash matches ('56f6d374c5d9b04304bc38a069ebbf84')  skipping ...
Closing connection to GCP ...
```

## Get deployed notebook runtime templates
### help
```sh
$> ./dcloud list -h                           
Retrieve existing NotebookRuntimeTemplates

Usage:
  dcloud list [project] [flags]

Flags:
  -h, --help   help for list

Global Flags:
  -p, --project string   GCP Project Name
```

### example
```sh
$>  ./dcloud list -p gamma-priceline-playground 
Logging onto GCP using ADC ...
Retrieving existing deployed runtime templates ...
Found 2 existing runtime templates
{
  "Name": "projects/1019340507365/locations/australia-southeast1/notebookRuntimeTemplates/5680587242486104064",
  "DisplayName": "This is an example of a runtime template",
  "Description": "Deployed from sample.json",
  "FileHash": "18879ff2f8ee028deff132b953c534b4",
  "MachineType": "e2-standard-2"
}
{
  "Name": "projects/1019340507365/locations/australia-southeast1/notebookRuntimeTemplates/4971270301175250944",
  "DisplayName": "This is an another example of a runtime template",
  "Description": "Deployed from sample2.json",
  "FileHash": "56f6d374c5d9b04304bc38a069ebbf84",
  "MachineType": "e2-standard-4"
}
Closing connection to GCP ...
```

## Delete deployed notebook runtime templates

```sh
$> ./clouddelete -p gamma-priceline-playground -n 'This is an another example of a runtime template'

Logging onto GCP using ADC ...
Retrieving existing deployed runtime templates ...
Found 2 existing runtime templates
Found template: projects/1019340507365/locations/australia-southeast1/notebookRuntimeTemplates/8475070811269496832
Deleted Notebook Runtime Template &{0x140000e4f18}

Closing connection to GCP ...
                                
```
