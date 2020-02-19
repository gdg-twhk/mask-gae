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


.PHONY: all help

help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
	@echo ""
	@for var in $(helps); do \
		echo $$var; \
	done | column -t -s ':' |  sed -e 's/^/  /'