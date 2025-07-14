.PHONY: install-dev
install-dev:
	@echo "Installing..."
	npm install -D live-server@latest
	go install github.com/bombsimon/wsl/v4/cmd...@v4.4.1
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.2.2

.PHONY: local_generate
local_generate:
	@echo "Generating..."
	tailwindcss --input ./theme/plain/static/in.css --output ./theme/plain/static/style.css
	go run ./cmd/templum/. --content content --output public --url "http://localhost:8080/"

.PHONY: serve
serve:
	@echo "Serving..."
	cd public && live-server --host=localhost --port=8080

.PHONY: fmt
fmt:
	@echo "Formatting..."
	go mod tidy
	wsl -fix ./... || true
	golangci-lint fmt ./...

.PHONE: lint
lint:
	@echo "Linting..."
	golangci-lint run -v ./...

.PHONY: test
test:
	@echo "Testing..."
	go test ./...