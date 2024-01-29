run:
	@./dist/zarldotdev

# Path: Makefile
build:
	@echo "===========[ Building ZarlDev ]==========="
	@echo "===========[ Tailwind CSS ]==========="
	@npx tailwindcss -i assets/css/app.css -o assets/css/apptw.css --minify
	@echo "===========[ Templ HTML Templates ]==========="
	@echo "Templ Templates..."
	@templ generate
	@echo "===========[ Go Generate ]==========="
	@go generate ./...
	@echo "===========[ Binary ]==========="
	@go build -o dist/zarldotdev ./cmd/zarldotdev.go
	@ls -lah dist/zarldotdev | awk '{print "Location:" $$9, "Size:" $$5}' | column -t	
	@echo "===========[ DONE ]==========="

docker:
	@echo "===========[ Building ZarlDev Docker Image ]==========="
	@echo "Building Docker Image..."
	@docker build -t zarldotdev:latest . --no-cache
	@echo "===========[ DONE ]==========="

docker-run:
	@echo "===========[ Running ZarlDev Docker Image ]==========="
	@echo "Running Docker Image..."
	@docker compose up -d 
	@brave http://localhost:8080
	@echo "===========[ DONE ]==========="

docker-stop:
	@echo "===========[ Stopping ZarlDev Docker Image ]==========="
	@echo "Stopping Docker Image..."
	@docker compose down
	@echo "===========[ DONE ]==========="
