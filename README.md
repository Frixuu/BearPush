# bearpush

Utility software for pushing artifacts.

## Build

```sh
go build .
```

## Usage

```sh
curl -X POST http://localhost:44344/v1/upload/your-product-name \
  -H "Content-Type: multipart/form-data" \
  -H "Authorization: Bearer your-token-here" \
  -F "artifact=@your-file.tar"
```
