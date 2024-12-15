# Go-templates [![License](https://img.shields.io/badge/License-MIT%202.0-green.svg)](https://github.com/thylong/go-templates/blob/main/01-simple-k8s-application/LICENSE) [![Go Monorepo CI](https://github.com/thylong/go-templates/actions/workflows/go-monorepo.yml/badge.svg?branch=main)](https://github.com/thylong/go-templates/actions/workflows/go-monorepo.yml)


This repository contains a list of go application templates to use with `$ gonew`.
The idea is reduce the burden of bootstrapping my Go applications.

## Template list

01. Simple k8s application [![Go Report Card](https://goreportcard.com/badge/github.com/thylong/gonew-templates/01-simple-k8s-application)](https://goreportcard.com/report/github.com/thylong/gonew-templates/01-simple-k8s-application)
02. Simple k8s fiber app [![Go Report Card](https://goreportcard.com/badge/github.com/thylong/go-templates/02-simple-k8s-fiber-app)](https://goreportcard.com/report/github.com/thylong/go-templates/02-simple-k8s-fiber-app)
03. k8s fiber sqlc [![Go Report Card](https://goreportcard.com/badge/github.com/thylong/go-templates/03-k8s-fiber-sqlc)](https://goreportcard.com/report/github.com/thylong/go-templates/03-k8s-fiber-sqlc)

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
