start-qdrant:
	@echo "Starting Qdrant server..."
	docker compose up -d qdrant

dev-start-qdrant-interactor:
	@echo "Starting development environment..."
	cd packages/qdrant-interactor && \
    	. ./.venv/bin/activate && \
		python main.py

dev-start-chatbot:
	@echo "Starting chatbot development environment..."
	cd packages/chatbot && \
    	. ./.venv/bin/activate && \
		python main.py

dev-start-web-client:
	@echo "Starting web client development environment..."
	cd packages/web-client && \
		npm install && \
		npm run dev

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

dev-start-all:
	@echo "Starting all development environments in parallel..."
	@trap 'echo "Stopping all services..."; kill 0; exit' INT; \
	(cd packages/qdrant-interactor && . ./.venv/bin/activate && python main.py) & \
	PID1=$$!; \
	(cd packages/chatbot && . ./.venv/bin/activate && python main.py) & \
	PID2=$$!; \
	(cd packages/web-client && npm install && npm run dev) & \
	PID3=$$!; \
	echo "All services started in background ✅"; \
	echo "PIDs: Qdrant-Interactor=$$PID1, Chatbot=$$PID2, Web-Client=$$PID3"; \
	echo "Press Ctrl+C to stop all services"; \
	wait

dev-stop-all:
	@echo "Stopping all development services..."
	@pkill -f "python.*qdrant-interactor.*main.py" || echo "Qdrant Interactor not running"
	@pkill -f "python.*chatbot.*main.py" || echo "Chatbot not running"
	@pkill -f "npm.*run.*dev" || echo "Web Client not running"
	@pkill -f "node.*vite" || echo "Vite dev server not running"
	@echo "All development services stopped ✅"