.PHONY: install
install:
	@echo "Installing..."
	go install github.com/a-h/templ/cmd/templ@latest
	npm install -g tailwindcss@latest
	npm install -D tailwindcss@latest @tailwindcss/typography

.PHONY: install-dev
install-dev: install
	@echo "Installing..."
	npm install -g postcss@latest autoprefixer@latest
	npm install -D postcss@latest autoprefixer@latest
	npm install -g live-server@latest
	npm install -D live-server@latest
	go install github.com/cosmtrek/air@latest
	go install github.com/bombsimon/wsl/v4/cmd...@master
	go install mvdan.cc/gofumpt@latest
	go install github.com/daixiang0/gci@latest

.PHONY: generate
generate:
	@echo "Generating..."
	templ generate
	npx tailwindcss -i ./theme/plain/static/in.css -o ./theme/plain/static/style.css
	go run ./cmd/tempel/. --content content --output public

.PHONY: watch
watch:
	@echo "Watching..."
	air

.PHONY: serve
serve:
	@echo "Serving..."
	live-server --host=localhost public

.PHONY: fmt
fmt:
	@echo "Formatting..."
	gci write -s standard -s default -s "prefix(github.com/cugu/tempel)" .
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