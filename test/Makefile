sql:
	./cloud_sql_proxy -instances=mask-9999:asia-east2:health-insurance-special-pharmacy=tcp:5432

cmdpharmacy:
	DB_HOST=mask-9999:asia-east2:health-insurance-special-pharmacy \
	DB_PORT=5432 \
	DB_USER=postgres \
	DB_PASS=password \
	DB=mask \
	go run ../cmd/pharmacy/main.go


cmdfeedback:
	MASK_FEEDBACK_DB_HOST=mask-9999:asia-east2:health-insurance-special-pharmacy \
	MASK_FEEDBACK_DB_PORT=5432 \
	MASK_FEEDBACK_DB_USER=postgres \
	MASK_FEEDBACK_DB_PASS=password \
	MASK_FEEDBACK_DB=mask \
	go run ../cmd/feedback/main.go

press_test:
	ab -n 2000 -c 2 -p test.json -T 'application/json' https://mask-9999.appspot.com/api/pharmacies