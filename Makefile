migrate-create:
	migrate create -ext sql -seq course_ecommerce_db

run:
	@go run cmd/api/main.go