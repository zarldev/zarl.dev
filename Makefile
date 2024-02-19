run:
	@./dist/zarldotdev
	@brave http://localhost:8080

build:
	@echo "===========[ Building ZarlDev ]==========="
	@echo "===========[ Tailwind CSS ]==========="
	@npx tailwindcss -i assets/css/style.css -o assets/css/app.css --minify
	@echo "===========[ Templ HTML Templates ]==========="
	@echo "Templ Templates..."
	@templ generate
	@echo "===========[ Go Generate ]==========="
	@go generate ./...
	@echo "===========[ Binary ]==========="
	@go build -o dist/zarldotdev ./cmd/zarldotdev.go
	@ls -lah dist/zarldotdev | awk '{print "Location:" $$9, "Size:" $$5}' | column -t
	@echo "===========[ DONE ]==========="

docker-build:
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

docker-publish : docker-build
	@echo "===========[ Publishing ZarlDev Docker Image ]==========="
	@echo "Publishing Docker Image..."
	@docker tag zarldotdev:latest ghcr.io/zarldev/zarldotdev:latest
	@docker push ghcr.io/zarldev/zarldotdev:latest
	@echo "===========[ DONE ]==========="

