.PHONY: install
install:
	@echo "Installing..."
	go install github.com/a-h/templ/cmd/templ@v0.3.906
	npm install tailwindcss@latest @tailwindcss/typography

.PHONY: install-dev
install-dev: install
	@echo "Installing..."
	npm install -D live-server@latest

.PHONY: install_golangci_lint
install_golangci_lint:
	@echo "Installing..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.11.4
	@echo "Done Installing."

.PHONY: local_generate
local_generate:
	@echo "Generating..."
	templ generate
	npx tailwindcss -i ./theme/plain/static/in.css -o ./theme/plain/static/style.css
	go run ./cmd/templum/. --content content --output public --url "http://localhost:8080/"

.PHONY: serve
serve:
	@echo "Serving..."
	cd public && live-server --host=localhost --port=8080

.PHONY: fmt
fmt:
	@echo "Formatting..."
	golangci-lint fmt ./...
	templ fmt .

.PHONE: lint
lint:
	@echo "Linting..."
	golangci-lint run -v ./...

.PHONY: test
test:
	@echo "Testing..."
	go test ./...