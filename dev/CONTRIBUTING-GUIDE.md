## Prerequisite
- [Go 1.22.X](https://go.dev/doc/install)
- [Docker](https://docs.docker.com/engine/install/)
- [make](https://askubuntu.com/questions/161104/how-do-i-install-make)

## How to Contribute

### Build Dev Image
Dev Image contains everything that are needed to develop golib, from go goimports and golangci-lint.
```bash
make dev
```
Above command will build `golib-dev` docker image and run it in a container, and open interactive session to it.

### Edit Code
Root directory is mounted to the golib-dev. Any chances happening in your local will also reflect to the docker container.

### Format and Lint
Run following command *inside interactive session* to format and lint your newly added code.
```bash
make format
make lint
```

### Test
Run following command *inside interactive session* to test your newly added code.
```bash
make test
```
