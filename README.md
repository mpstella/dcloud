DAW gcloud for NotebookRuntimeTemplates

## CLI Options
```text
> ./dcloud --help
Given we are waiting on Google this is our dodgy gcloud helper

Usage:
  dcloud [command]

Available Commands:
  delete      Delete an existing NotebookRuntimeTemplate
  deploy      Deploy NotebookRuntimeTemplates
  help        Help about any command
  list        Retrieve existing NotebookRuntimeTemplates
  version     Show version number of dcloud

Flags:
  -h, --help     help for dcloud
  -t, --toggle   Help message for toggle

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
  -h, --help               help for deploy
      --project string     GCP Project Name
      --templates string   Directory where templates are located
      --threads int        Number of concurrent threads (default 1)
```

### example
```text
$> ./dcloud deploy --project XXXXX --templates templates --threads 10
INFO[0000] (templates/sample3.json) Parsing template    
INFO[0000] (templates/sample2.yaml) Parsing template    
INFO[0000] (templates/sample1.json) Parsing template    
INFO[0000] (templates/sample3.json) Deploying template. 
INFO[0000] (templates/sample2.yaml) Deploying template. 
INFO[0000] (templates/sample1.json) Deploying template. 
INFO[0001] (templates/sample3.json) Processed template. 
INFO[0001] (templates/sample2.yaml) Processed template. 
INFO[0001] (templates/sample1.json) Processed template.   
```



### idempotent
```text
$> ./dcloud deploy --project XXXXX --templates templates --threads 10
INFO[0000] (templates/sample1.json) Parsing template    
INFO[0000] (templates/sample2.yaml) Parsing template    
INFO[0000] (templates/sample3.json) Parsing template    
INFO[0000] (templates/sample1.json) Found existing template with same DisplayName and a md5 hash, skipping .. 
INFO[0000] (templates/sample2.yaml) Found existing template with same DisplayName and a md5 hash, skipping .. 
INFO[0000] (templates/sample3.json) Found existing template with same DisplayName and a md5 hash, skipping .. 
```

## Get deployed notebook runtime templates
### help
```text
$> ./dcloud list --help
Retrieve existing NotebookRuntimeTemplates

Usage:
  dcloud list [project] [flags]

Flags:
  -h, --help             help for list
      --project string   GCP Project Name
```
### example
```text
$>  ./dcloud list  --project XXXXX
INFO[0000] {
  "name": "projects/1019340507365/locations/australia-southeast1/notebookRuntimeTemplates/2933954419743522816",
  "displayName": "This is an example of a runtime template [sample1.json]",
  "description": "Deployed from sample1.json",
  "machineSpec": {
    "machineType": "e2-standard-2"
  },
  "dataPersistentDiskSpec": {
    "diskType": "pd-standard",
    "diskSizeGb": "10"
  },
  "networkSpec": {
    "enableInternetAccess": true,
    "network": "projects/1019340507365/global/networks/default"
  },
  "etag": "AMEw9yOL_SUIaJb8cHdsiC3i-0gwovQ6Wph6kjILvMJX_L-XJjlzHuWutNTXOpcqjiZT",
  "labels": {
    "deployment_ts_utc": "20240625_053002",
    "env": "dev2",
    "md5": "bd1be799f147d7f3cee1cac98fba3066",
    "source": "sample1"
  },
  "idleShutdownConfig": {
    "idleTimeout": "600s"
  },
  "eucConfig": {},
  "createTime": "2024-06-25T05:30:04.856134Z",
  "updateTime": "2024-06-25T05:30:04.856134Z",
  "notebookRuntimeType": "USER_DEFINED"
} 
INFO[0000] {
  "name": "projects/1019340507365/locations/australia-southeast1/notebookRuntimeTemplates/4075616925281943552",
  "displayName": "This is an example of a runtime template [sample2.yaml]",
  "description": "Deployed from sample2.yaml",
  "machineSpec": {
    "machineType": "e2-standard-4"
  },
  "dataPersistentDiskSpec": {
    "diskType": "pd-standard",
    "diskSizeGb": "10"
  },
  "networkSpec": {
    "enableInternetAccess": true,
    "network": "projects/1019340507365/global/networks/default"
  },
  "etag": "AMEw9yPXdMyvHsDHrwQi2D5zByAuKra0TtXfUgPJbsxfha9JLtc_HZ3TA8g68NxKcsYP",
  "labels": {
    "deployment_ts_utc": "20240625_053002",
    "md5": "2b881d4beafacc8cd0d9fcb9420b82fe"
  },
  "idleShutdownConfig": {
    "idleTimeout": "600s"
  },
  "eucConfig": {},
  "createTime": "2024-06-25T05:30:04.836839Z",
  "updateTime": "2024-06-25T05:30:04.836839Z",
  "notebookRuntimeType": "USER_DEFINED"
} 
INFO[0000] {
  "name": "projects/1019340507365/locations/australia-southeast1/notebookRuntimeTemplates/7676244827364655104",
  "displayName": "This is an example of a runtime template [sample3.json]",
  "description": "Deployed from sample3.json",
  "machineSpec": {
    "machineType": "e2-standard-4"
  },
  "dataPersistentDiskSpec": {
    "diskType": "pd-standard",
    "diskSizeGb": "10"
  },
  "networkSpec": {
    "enableInternetAccess": true,
    "network": "projects/1019340507365/global/networks/default"
  },
  "etag": "AMEw9yMHlv6Q12FCIiGMzRnQkwscZcdGK0_dVlHelcxD4rSTvTveR2y95umxRwWcvinA",
  "labels": {
    "deployment_ts_utc": "20240625_053002",
    "md5": "10b06b42701ada152dc131ff0148fce0"
  },
  "idleShutdownConfig": {
    "idleTimeout": "600s"
  },
  "eucConfig": {},
  "createTime": "2024-06-25T05:30:04.618378Z",
  "updateTime": "2024-06-25T05:30:04.618378Z",
  "notebookRuntimeType": "USER_DEFINED"
} 
```

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
```

## example
```text
$> ./dcloud delete --name "projects/1019340507365/locations/australia-southeast1/notebookRuntimeTemplates/7676244827364655104"     
```
