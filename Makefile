.PHONY: generate
generate:
	@echo "Generating..."
	templ generate
	npx tailwindcss -i ./theme/default/in.css -o ./theme/default/out.css
	go run ./cmd/tempel/. --content content --output public

.PHONY: watch
watch:
	@echo "Watching..."
	air

.PHONY: serve
serve:
	@echo "Serving..."
	live-server --host=localhost public