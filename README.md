# Azure Logger

A simple way to have the standard log being piped to Azure's LogFiles/Application so it can be streamed by the Azure CLI.

## Installation

`go get github.com/dstpierre/azure-logger`

## How to use it

```go
package main

import (
    "log"
    logger "github.com/dstpierre/azure-logger"
)

func main() {
    err := logger.Start()
    if err != nil {
        // the log creation file failed 
    }
    defer logger.Stop()
    
    log.Println("This line will be streamed from azure site log tail sitename")
}
```

By using the standard `log` package, there's nothing to change if you were already using this. All output from the base logger
will be saved into the file, and optionally tail when using Azure's CLI command: `azure site log tail sitename`.

## What it is doing

1. It simply set the output of the `log` package to be a file in Azure's Application log directory.
2. When it start it will purge every log file older than 24 hours
3. The log file name format are this al_instanceid_date.log

Where `instanceid` is the 6 first characters of the environment variable `WEBSITE_INSTANCE_ID`

And `date` is the current date when the process was started. 

