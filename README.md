# bearpush

[![Go Report](https://img.shields.io/badge/go%20report-A-green.svg?style=flat)](https://goreportcard.com/report/github.com/Frixuu/BearPush) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Frixuu/BearPush) ![Lines of code](https://img.shields.io/tokei/lines/github/Frixuu/BearPush)

The poor man's CD pipeline.  
Bearpush supplies you with a web server, bring your own scripts and auth (optional)

## Download

Binaries for x86-64 Linux are available from the GitHub Actions.

## Build from source

```sh
go build ./cmd/bearpush/
```

## Usage

To use Bearpush, you need to configure at least one **product** - a type of entity Bearpush can process. You can scaffold the configuration with ```bearpush product new [product name]```.

The newly created configuration file will include a randomly generated token ```static-value```. Feel free to change it. This token has to be provided in your request in order for it to be processed.  

To do something with the uploaded file, you have to create a ```process-script```. By default nothing is associated with your product.
The simplest shell script could look like this:

```sh
#!/bin/bash

mkdir -p /app/foo
cp "${ARTIFACT_PATH}" /app/foo/my-file.txt
```

Note that Bearpush provides you with the ```ARTIFACT_PATH``` environment variable.
Update your config file to include the **absolute** path to the script: ```process-script: '/home/foo/bar.sh'```.  

Run ```bearpush``` to start the server.

The simplest way to upload your artifact is to use curl:

```sh
curl -X POST http://localhost:44344/v1/upload/your-product-name \
  -H "Content-Type: multipart/form-data" \
  -H "Authorization: Bearer your-token-here" \
  -F "artifact=@your-file.tar"
```
