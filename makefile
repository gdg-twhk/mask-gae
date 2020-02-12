ACCOUNT = cage.chung@gmail.com
PROJECT = mask-9999
# VERSION = 4c833b1makm
VERSION = 4

NODE_BIN = $(shell npm bin)

set_config:
	gcloud config set account $(ACCOUNT)
	gcloud config set project $(PROJECT)

run:
	dev_appserver.py \
	dispatch.yaml \
	ownership/app.yaml \
	frontend/app.yaml \
	endpoints/app.yaml \
	--skip_sdk_update_check=yes \
	--host 0.0.0.0 \
	--enable_sendmail=yes

update_frontend:
	gcloud app deploy --version $(VERSION) frontend/app.yaml --project $(PROJECT) --promote -q

# update_endpoints:
# 	gcloud app deploy --version $(VERSION) endpoints/app.yaml --project $(PROJECT) --promote -q
#
# update_ownership:
# 	gcloud app deploy --version=$(VERSION) ownership/app.yaml --project $(PROJECT) --promote -q
#
update_pharmacy:
	gcloud app deploy --version=0-1 pharmacy/app.yaml --project $(PROJECT) --promote -q


update_dispatch:
	gcloud app deploy --version=$(VERSION)  dispatch.yaml --project $(PROJECT) -q
#
# update: update_frontend update_endpoints update_ownership update_dispatch
# update: update_frontend update_endpoints



sql:
	./cloud_sql_proxy -instances=mask-9999:asia-east2:health-insurance-special-pharmacy=tcp:5432

cmdpharmacy:
	MADK_PHARMACY_DB_HOST=localhost MADK_PHARMACY_DB_PORT=5433 go run pharmacy/main.go