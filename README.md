DAW gcloud for NotebookRuntimeTemplates

## CLI Options
```text
> ./dcloud --help
Given we are waiting on Google this is our dodgy gcloud helper

Usage:
  dcloud [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  delete      Delete an existing NotebookRuntimeTemplate
  deploy      Deploy NotebookRuntimeTemplates
  export      Export an existing NotebookRutimeTemplate
  help        Help about any command
  list        Retrieve existing NotebookRuntimeTemplates

Flags:
      --dry-run   Run the command in dry-run mode
  -h, --help      help for dcloud
      --silent    Minimise output to stdout
  -t, --toggle    Help message for toggle

Use "dcloud [command] --help" for more information about a command.
```

## Credentials

This application will use the magic of GCP [Application Default Credentials](https://cloud.google.com/docs/authentication/application-default-credentials) to log you in.

## Deploy a notebook runtime template
### help
```text
$> ./dcloud deploy --help
Deploy NotebookRuntimeTemplates

Usage:
  dcloud deploy [project] [pathToTemplates] [flags]

Flags:
      --dry-run            Run the command in dry-run mode
  -h, --help               help for deploy
      --project string     GCP Project Name
      --templates string   Directory where templates are located

Global Flags:
      --silent   Minimise output to stdout
```

### example
```text
$> ./dcloud deploy --project gamma-priceline-playground --templates templates 
INFO[0000] Logging onto GCP using ADC ...               
INFO[0000] Attempting to deploy templates/sample.json   
INFO[0000] Retrieving existing deployed runtime templates ... 
INFO[0001] Created Notebook Runtime Template: &{0x14000596810} 
INFO[0001] Attempting to deploy templates/sample2.json  
INFO[0001] Retrieving existing deployed runtime templates ... 
INFO[0002] Created Notebook Runtime Template: &{0x140000e6720} 
INFO[0002] Closing connection to GCP ...    
```

### dry-run
```text
$> ./dcloud deploy --project gamma-priceline-playground --templates templates --dry-run
INFO[0000] Attempting to deploy templates/sample.json   
INFO[0000] Retrieving existing deployed runtime templates ... 
INFO[0000] This is a dry-run, however, the template would be deployed 
INFO[0000] Attempting to deploy templates/sample2.json  
INFO[0000] Retrieving existing deployed runtime templates ... 
INFO[0000] A template already exists with this Display Name, will check for changes ... 
INFO[0000] This is a dry-run, however, the template would be deleted as the hashes do not match 
INFO[0000] This is a dry-run, however, the template would be deployed 
INFO[0000] Attempting to deploy templates/sample2_duplicate.json 
INFO[0000] Retrieving existing deployed runtime templates ... 
INFO[0000] A template already exists with this Display Name, will check for changes ... 
INFO[0000] Template hash matches ('56f6d374c5d9b04304bc38a069ebbf84') skipping ... 
INFO[0000] Closing connection to GCP ...     
```

### idempotent
```text
$> ./dcloud deploy --project gamma-priceline-playground --templates templates
INFO[0000] Logging onto GCP using ADC ...               
INFO[0000] Attempting to deploy templates/sample.json   
INFO[0000] Retrieving existing deployed runtime templates ... 
INFO[0000] A template already exists with this Display Name, will check for changes ... 
INFO[0000] Template hash matches ('18879ff2f8ee028deff132b953c534b4') skipping ... 
INFO[0000] Attempting to deploy templates/sample2.json  
INFO[0000] Retrieving existing deployed runtime templates ... 
INFO[0000] A template already exists with this Display Name, will check for changes ... 
INFO[0000] Template hash matches ('56f6d374c5d9b04304bc38a069ebbf84') skipping ... 
INFO[0000] Closing connection to GCP ...  
```

## Get deployed notebook runtime templates
### help
```text
$> ./dcloud list -h                           
Retrieve existing NotebookRuntimeTemplates

Usage:
  dcloud list [project] [flags]

Flags:
  -h, --help             help for list
      --project string   GCP Project Name

Global Flags:
      --dry-run   Run the command in dry-run mode
      --silent    Minimise output to stdout
```
### example
```text
$>  ./dcloud list --project gamma-priceline-playground
INFO[0000] Logging onto GCP using ADC ...               
INFO[0000] Retrieving existing deployed runtime templates ... 
INFO[0000] {
  "Name": "projects/1019340507365/locations/australia-southeast1/notebookRuntimeTemplates/3241888044264980480",
  "DisplayName": "This is an another example of a runtime template",
  "Description": "Deployed from sample2.json",
  "FileHash": "56f6d374c5d9b04304bc38a069ebbf84",
  "MachineType": "e2-standard-4"
} 
INFO[0000] {
  "Name": "projects/1019340507365/locations/australia-southeast1/notebookRuntimeTemplates/2570851699786776576",
  "DisplayName": "This is an example of a runtime template",
  "Description": "Deployed from sample.json",
  "FileHash": "18879ff2f8ee028deff132b953c534b4",
  "MachineType": "e2-standard-2"
} 
INFO[0000] Closing connection to GCP ... 
```
## Export an existing notebook runtime template

### help
```text
$> ./dcloud export --help
Export an existing NotebookRutimeTemplate

Usage:
  dcloud export [name] [flags]

Flags:
  -h, --help          help for export
      --name string   Name of the template

Global Flags:
      --dry-run   Run the command in dry-run mode
      --silent    Minimise output to stdout
```

### example
```text
$> ./dcloud export --name "projects/1019340507365/locations/australia-southeast1/notebookRuntimeTemplates/2570851699786776576"
INFO[0000] Logging onto GCP using ADC ...               
INFO[0000] {
  "name": "projects/1019340507365/locations/australia-southeast1/notebookRuntimeTemplates/2570851699786776576",
  "display_name": "This is an example of a runtime template",
  "description": "Deployed from sample.json",
  "machine_spec": {
    "machine_type": "e2-standard-2"
  },
  "data_persistent_disk_spec": {
    "disk_type": "pd-standard",
    "disk_size_gb": 10
  },
  "network_spec": {
    "enable_internet_access": true,
    "network": "projects/1019340507365/global/networks/default"
  },
  "etag": "AMEw9yOC7ouezSgCrI0ZeAj-_BOAVFGRnmg3sq8G1sWIE2UHAAAx6ktVVS_3XUzdRgfU",
  "labels": {
    "git_run_id": "",
    "git_sha": "",
    "md5": "18879ff2f8ee028deff132b953c534b4",
    "source": "sample"
  },
  "idle_shutdown_config": {
    "idle_timeout": {
      "seconds": 600
    }
  },
  "euc_config": {},
  "create_time": {
    "seconds": 1718587516,
    "nanos": 808801000
  },
  "update_time": {
    "seconds": 1718587516,
    "nanos": 808801000
  },
  "notebook_runtime_type": 1
} 
INFO[0000] Closing connection to GCP ...   
```

## Delete existing notebook runtime templates

## help
```text
$/ ./dcloud delete --help
Delete an existing NotebookRuntimeTemplate

Usage:
  dcloud delete [name] [flags]

Flags:
  -h, --help          help for delete
      --name string   Name of the template

Global Flags:
      --dry-run   Run the command in dry-run mode
      --silent    Minimise output to stdout
```

## example
```text
$> ./dcloud delete  --name "projects/1019340507365/locations/australia-southeast1/notebookRuntimeTemplates/2570851699786776576"
INFO[0000] Logging onto GCP using ADC ...               
INFO[0000] Deleted Notebook Runtime Template &{0x14000130af8}
INFO[0000] Closing connection to GCP ...              
```
