# hcp - copy via http

> Transfer files between devices via HTTP

## Usage

### Start server

```
docker run -p 8080:8080 hellodhlyn/hcp
```

### Transfer files

You can upload and download files using HTTP clients such as curl or wget.

For example:

```sh
# Upload
curl -F 'file=@/path/to/file' http://<server-address>/<key>

# Download
curl -o <filename> -f http://<server-address>/<key>
```
