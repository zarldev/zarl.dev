run:
	@echo "Starting App..."
	@./dist/zarldev

# Path: Makefile
build:
	@echo "Building Assets..."
	@echo "Tailwind CSS..."
	@npx tailwindcss -i assets/css/app.css -o assets/css/apptw.css --minify
	@echo "Templ Templates..."
	@templ generate
	@go generate ./...
	@go build -o dist/zarldev ./cmd/zarldev.go

