# UofT-Webcrawler

<p>This project aims to crawl UofT websites, and vector embed chunks of text using Qdrant. The project will have a web user interface and a REST API allowing users to search the database for relevant information. (eg: user searches "Bayesian Networks" and gets classes and programs related to Bayesian Networks).</p>
<p>The project operates with a microservices architecture, with packages Dockerized and orchestrated with Docker Compose. All commands are handled with GNU Make.</p>

## 📁 Project Structure

```
UofT-Webcrawler/
├── docker-compose.yml          # Orchestrates all services
├── Makefile                    # Build and deployment commands
├── README.md                   # Project documentation
├── go.mod                      # Go module dependencies
├── go.sum                      # Go module checksums
├── packages/
│   ├── qdrant-interactor/      # Python FastAPI service for Qdrant operations
│   │   ├── Dockerfile          # Container configuration
│   │   ├── main.py             # FastAPI application entry point
│   │   ├── embedding_agent.py  # Qdrant client and vector operations
│   │   ├── requirements.txt    # Python dependencies
│   │   └── schemas/            # Pydantic models for API requests/responses
│   │       ├── QueryRequest.py
│   │       ├── QueryResponse.py
│   │       └── TextEmbeddingRequest.py
│   └── web-crawler/            # Go-based web crawler service
│       └── main.go             # Web crawler implementation
└── qdrant_storage/             # Persistent storage for Qdrant database
    ├── raft_state.json
    ├── aliases/
    └── collections/
```

## 🏛️ Services Architecture

### 1. 🗄️ Qdrant Vector Database
- **Image**: `qdrant/qdrant:latest`
- **Purpose**: Stores and indexes vector embeddings for semantic search
- **Ports**: 
  - `6333`: HTTP API for vector operations
  - `6334`: gRPC API (alternative interface)
- **Storage**: Persistent volume for database state

### 2. 🐍 Qdrant Interactor (Python FastAPI)
- **Purpose**: REST API service for vector embedding and search operations
- **Language**: Python 3.11.8
- **Framework**: FastAPI with automatic OpenAPI documentation
- **Port**: `8000`
- **Dependencies**: Qdrant, Uvicorn, FastAPI, Pydantic

### 3. 🕷️ Web Crawler (Go)
- **Purpose**: Crawls UofT websites and extracts content for embedding
- **Language**: Go
- **Status**: 🚧 In development

## 🚀 Qdrant Interactor API

The Qdrant Interactor service provides a REST API for managing vector embeddings and performing semantic searches.

### 🔗 Endpoints

#### 📝 `POST /add_embedding`
Adds a new document to the vector database.

**Request Body:**
```json
{
  "content": "Document content to be embedded",
  "url": "https://source-url.com",
  "position": 100
}
```

**Response:**
- `200`: ✅ Document successfully added
- `500`: ❌ Error adding document

#### 🔍 `GET /query`
Performs semantic search against the vector database.

**Request Body:**
```json
{
  "query": "search terms",
  "collection_name": "test_collection",
  "limit": 10
}
```

**Response:**
```json
{
  "results": [
    {
      "id": "unique-document-id",
      "content": "matching document content",
      "source": "https://source-url.com",
      "score": 0.95
    }
  ]
}
```

### 🧠 How Vector Embeddings Work

1. **📥 Document Ingestion**: When content is submitted via `/add_embedding`:
   - The text content is processed by the embedding model (`BAAI/bge-small-en`)
   - A vector representation is generated from the text
   - The vector is stored in Qdrant with metadata (source URL, content)

2. **🔍 Semantic Search**: When a query is submitted via `/query`:
   - The query text is converted to a vector using the same embedding model
   - Qdrant performs similarity search against stored vectors
   - Results are ranked by cosine similarity score
   - Top matching documents are returned with their metadata

3. **💾 Vector Storage**: 
   - Uses dense vector embeddings for semantic understanding
   - Supports exact and approximate nearest neighbor search
   - Persistent storage ensures data survives container restarts

## 💻 Development Commands

### 📋 Prerequisites
- 🐳 Docker and Docker Compose
- 🔨 GNU Make
- 🐍 Python 3.11+ (for local development)
- 🐹 Go 1.19+ (for web crawler development)

### ⚡ Quick Start

1. **🚀 Start Qdrant Database**:
   ```bash
   make start-qdrant
   ```

2. **🔧 Build and Deploy Qdrant Interactor**:
   ```bash
   make build-qdrant-interactor
   make deploy-qdrant-interactor
   ```

3. **🌐 Access Services**:
   - 📊 Qdrant Dashboard: `http://localhost:6333/dashboard`
   - 📖 API Documentation: `http://localhost:8000/docs`

### 🛠️ Development Workflow

**💻 Local Development (Python)**:
```bash
make dev-start-qdrant-interactor
```

**🐳 Docker Development**:
```bash
# Build with no cache
make build-qdrant-interactor

# Deploy container
make deploy-qdrant-interactor

# Clean rebuild
make destroy-qdrant-interactor
make build-qdrant-interactor
make deploy-qdrant-interactor
```

### Available Make Targets

- 🚀 `start-qdrant`: Start Qdrant database service
- 💻 `dev-start-qdrant-interactor`: Run interactor locally for development
- 🔧 `build-qdrant-interactor`: Build Docker image for interactor
- 📦 `deploy-qdrant-interactor`: Deploy interactor container
- 🗑️ `destroy-qdrant-interactor`: Stop and remove interactor container and image

## 🗺️ Next Steps

- [ ] 🕷️ Complete Go web crawler implementation
- [ ] 🌐 Add web user interface
- [ ] 🔐 Implement authentication and rate limiting
- [ ] 📊 Add monitoring and logging
- [ ] ☁️ Deploy to production environment