.PHONY: install
install:
	@echo "Installing..."
	go install github.com/a-h/templ/cmd/templ@latest
	npm install tailwindcss@latest @tailwindcss/typography

.PHONY: install-dev
install-dev: install
	@echo "Installing..."
	npm install -D live-server@latest
	go install github.com/bombsimon/wsl/v4/cmd...@master
	go install mvdan.cc/gofumpt@latest
	go install github.com/daixiang0/gci@latest

.PHONY: local_generate
local_generate:
	@echo "Generating..."
	templ generate
	npx tailwindcss -i ./theme/plain/static/in.css -o ./theme/plain/static/style.css
	go run ./cmd/templum/. --content content --output public --url "http://localhost:8080/"

.PHONY: serve
serve:
	@echo "Serving..."
	live-server --host=localhost public

.PHONY: fmt
fmt:
	@echo "Formatting..."
	gci write -s standard -s default -s "prefix(github.com/cugu/templum)" .
	gofumpt -l -w .
	wsl -fix ./...

.PHONE: lint
lint:
	@echo "Linting..."
	golangci-lint run ./...

.PHONY: test
test:
	@echo "Testing..."
	go test ./...