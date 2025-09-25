# Demo Application

This is a simple demo application to test the idpbuilder OCI build and push commands.

## Contents

- `Dockerfile` - Defines the container image
- `hello.txt` - Sample text file
- `README.md` - This file

## Testing

This directory will be built into an OCI image using:
```bash
idpbuilder build --context . --tag hello-world:v1
```

Then pushed to the Gitea registry using:
```bash
idpbuilder push hello-world:v1
```