sql:
	./cloud_sql_proxy -instances=mask-9999:asia-east2:health-insurance-special-pharmacy=tcp:5432

cmdpharmacy:
	DB_HOST=127.0.0.1 \
	DB_PORT=5432 \
	DB_USER=postgres \
	DB_PASS=password \
	DB=mask \
	go run ../cmd/pharmacy/main.go


cmdfeedback:
	MASK_FEEDBACK_DB_HOST=localhost \
	MASK_FEEDBACK_DB_PORT=5433 \
	MASK_FEEDBACK_DB_USER=postgres \
	MASK_FEEDBACK_DB_PASS=password \
	MASK_FEEDBACK_DB=mask \
	go run ../cmd/feedback/main.go
	go run ../cmd/feedback/main.go


cmdsync:
	PROJECT_ID=mask-9999 \
	BUCKET_ID=mask-9999-pharmacies \
	POINTS_OBJECT_NAME=points.json \
	DB_HOST=127.0.0.1 \
	DB_PORT=5432 \
	DB_USER=postgres \
	DB_PASS=password \
	DB=mask \
	go run ../cmd/sync/main.go
