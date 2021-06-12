# bearpush

[![Go Report](https://img.shields.io/badge/go%20report-A-green.svg?style=flat)](https://goreportcard.com/report/github.com/Frixuu/BearPush) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Frixuu/BearPush) ![Lines of code](https://img.shields.io/tokei/lines/github/Frixuu/BearPush)

The poor man's CD pipeline.  
Bearpush supplies you with a web server, bring your own scripts and auth (optional)

## Download

Binaries for x86-64 Linux are available from the GitHub Actions.

## Build from source

```sh
go build .
```

## Usage

To use Bearpush, you need to configure at least one **product** - a type of entity Bearpush can process. You can scaffold the configuration with ```bearpush product new [product name]```.

The newly created configuration file will include a randomly generated static token. Feel free to change it. However, the config will not have any script attached yet. You can create one now.

Run ```bearpush``` to start the server.

The simplest way to upload your artifact is to use curl:

```sh
curl -X POST http://localhost:44344/v1/upload/your-product-name \
  -H "Content-Type: multipart/form-data" \
  -H "Authorization: Bearer your-token-here" \
  -F "artifact=@your-file.tar"
```
