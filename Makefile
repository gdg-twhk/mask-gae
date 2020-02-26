ACCOUNT = $(shell gcloud auth list --filter=status:ACTIVE --format='value(account)')
PROJECT = $(shell gcloud config list --format 'value(core.project)')
VERSION = $(shell git rev-parse --short HEAD)

all: help

## deploy_pharmacy [v=version-name]: deploy pharmacy service
deploy_pharmacy:
ifdef v
	gcloud app deploy --version ${v} --project ${PROJECT} -q cmd/pharmacy/app.yaml
else
	gcloud app deploy --version ${VERSION} --project ${PROJECT} -q cmd/pharmacy/app.yaml
endif

## deploy_docs [v=version-name]: deploy docs service
deploy_docs:
ifdef v
	gcloud app deploy --version ${v} --project ${PROJECT} -q cmd/docs/app.yaml
else
	gcloud app deploy --version ${VERSION} --project ${PROJECT} -q cmd/docs/app.yaml
endif

## deploy_feedback [v=version-name]: deploy feedback service
deploy_feedback:
ifdef v
	gcloud app deploy --version ${v} --project ${PROJECT} -q cmd/feedback/app.yaml
else
	gcloud app deploy --version ${VERSION} --project ${PROJECT} -q cmd/feedback/app.yaml
endif

## deploy_dispatch: deploy disptach
deploy_dispatch:
	gcloud app deploy --project ${PROJECT} -q dispatch.yaml

## build_swagger: generate swagger docs
build_swagger:
	swag init --dir .  --generalInfo ./cmd/docs/main.go --output ./cmd/docs/docs

.PHONY: all help

help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
	@echo ""
	@for var in $(helps); do \
		echo $$var; \
	done | column -t -s ':' |  sed -e 's/^/  /'