build:
	@go build -o bin/ravan .

# should implement full interactive mode first
# run: build
# 	@./bin/ravan

tidy:
	@go mod tidy

test:
	@go test ./... -v