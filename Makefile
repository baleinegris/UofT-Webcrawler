start-qdrant:
	@echo "Starting Qdrant server..."
	docker compose up -d qdrant
start-qdrant-interactor:
	@echo "Starting development environment..."
	cd packages/qdrant-interactor && \
		python main.py