migrate-create:
	migrate create -ext sql -seq course_ecommerce_db

run:
	@go run cmd/api/main.go

watch:
	air --build.cmd "go build -o bin/api cmd/api/main.go" --build.bin "./bin/api"