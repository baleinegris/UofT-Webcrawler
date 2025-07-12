start-qdrant:
	@echo "Starting Qdrant server..."
	docker compose up -d qdrant

dev-start-qdrant-interactor:
	@echo "Starting development environment..."
	cd packages/qdrant-interactor && \
		python main.py

dev-crawl:
	@echo "Starting web crawler development environment..."
	cd packages/web-crawler && \
		go run main.go

build-qdrant-interactor:
	@echo "Building Qdrant Interactor..."
	docker compose build --no-cache qdrant-interactor
	@echo "Qdrant Interactor built successfully ✅"

deploy-qdrant-interactor:
	@echo "Running Qdrant Interactor..."
	docker compose up -d qdrant-interactor
	@echo "Qdrant Interactor is running ✅"

destroy-qdrant-interactor:
	@echo "Stopping and removing Qdrant Interactor..."
	docker compose stop qdrant-interactor
	docker compose rm -f qdrant-interactor
	@echo "Removing Qdrant Interactor image..."
	IMAGE=$$(docker images --format '{{.Repository}}:{{.Tag}}' | grep qdrant-interactor || true); \
	if [ -n "$$IMAGE" ]; then docker rmi $$IMAGE; else echo "No qdrant-interactor image found."; fi

stop-qdrant:
	@echo "Stopping Qdrant server..."
	docker compose stop qdrant
	@echo "Qdrant server stopped ✅"

destroy-qdrant:
	@echo "Stopping and removing Qdrant server..."
	docker compose stop qdrant
	docker compose rm -f qdrant
	@echo "Qdrant server destroyed ✅"

restart-qdrant: destroy-qdrant start-qdrant
	@echo "Qdrant server restarted ✅"

stop-all:
	@echo "Stopping all services..."
	docker compose stop
	@echo "All services stopped ✅"