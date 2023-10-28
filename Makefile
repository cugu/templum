.PHONY: install
install:
	@echo "Installing..."
	go install github.com/a-h/templ/cmd/templ@latest
	npm install -g tailwindcss@latest
	npm install -D tailwindcss@latest

.PHONY: install-dev
install-dev:
	@echo "Installing..."
	go install github.com/niklasfasching/templ/cmd/templ@latest
	npm install -g tailwindcss@latest postcss@latest autoprefixer@latest
	npm install -D tailwindcss@latest postcss@latest autoprefixer@latest
	npm install -g live-server@latest
	npm install -D live-server@latest
	go install github.com/cosmtrek/air@latest

.PHONY: generate
generate:
	@echo "Generating..."
	templ generate
	npx tailwindcss -i ./theme/default/in.css -o ./theme/default/style.css
	go run ./cmd/tempel/. --content content --output public

.PHONY: watch
watch:
	@echo "Watching..."
	air

.PHONY: serve
serve:
	@echo "Serving..."
	live-server --host=localhost public