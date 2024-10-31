# Go-templates [![License](https://img.shields.io/badge/License-MIT%202.0-green.svg)](https://github.com/thylong/go-templates/blob/main/01-simple-k8s-application/LICENSE)

This repository contains a list of go application templates to use with `$ gonew`.
The idea is reduce the burden of bootstrapping my Go applications.

## Requirements

```bash
# install gonew bin
go install golang.org/x/tools/cmd/gonew@latest
```

## Quickstart

```bash
gonew github.com/thylong/go-templates/01-simple-k8s-application example.com/simple-app
```

## Features

### Common

- [x] Makefile
- [x] README.md
- [x] docker & docker-compose
- [x] hardened k8s template
- [x] Simple go app including go.mod & cobra
- [x] Restrictive k8s Network Policy
- [x] LICENSE & related badge

## License

This library is licensed under MIT Full license text is available in [LICENSE](https://github.com/thylong/go-templates/blob/main/LICENSE).
