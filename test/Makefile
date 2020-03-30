sql:
	./cloud_sql_proxy -instances=mask-9999:asia-east2:health-insurance-special-pharmacy=tcp:5432

cmdpharmacy:
	MASK_PHARMACY_PROJECT_ID=mask-9999 \
	MASK_PHARMACY_LOCATION_ID=asia-east2 \
	MASK_PHARMACY_QUEUE_ID=sync-points-queue2 \
	MASK_PHARMACY_BUCKET_ID=mask-9999-pharmacies \
	MASK_PHARMACY_POINTS_OBJECT_NAME=/points.json \
	MASK_PHARMACY_DB_HOST=localhost \
	MASK_PHARMACY_DB_PORT=5433 \
	MASK_PHARMACY_DB_USER=postgres \
	MASK_PHARMACY_DB_PASS=password \
	MASK_PHARMACY_DB=mask \
	go run ../cmd/pharmacy/main.go


cmdfeedback:
	MASK_FEEDBACK_DB_HOST=localhost \
	MASK_FEEDBACK_DB_PORT=5433 \
	MASK_FEEDBACK_DB_USER=postgres \
	MASK_FEEDBACK_DB_PASS=password \
	MASK_FEEDBACK_DB=mask \
	go run ../cmd/feedback/main.go
	go run ../cmd/feedback/main.go
