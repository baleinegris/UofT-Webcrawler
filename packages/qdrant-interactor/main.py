from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from schemas.TextEmbeddingRequest import TextEmbeddingRequest
from schemas.QueryRequest import QueryRequest
from schemas.QueryResponse import QueryResponse
from embedding_agent import startInteractor, addDocument, queryDatabase

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=['*'],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.on_event("startup")
def on_startup():
    startInteractor(collection_name="test_collection")

@app.post("/add_embedding")
async def add_embedding(request: TextEmbeddingRequest):
    """
    Given a text, url, and position, this endpoint creates an embedding for the text and adds it to the Qdrant database.
    """
    print(f"DEBUG: Received request - content: {request.content[:50]}..., title: '{request.title}', url: {request.url}")
    status = addDocument(content=request.content, title=request.title, source=request.url, collection_name=request.collection_name)
    if status is True:
        return {"message": "Document added successfully."}
    else:
        return {"message": "Failed to add document."}

@app.post("/query")
async def query(request: QueryRequest) -> QueryResponse:
    """
    Given a query string, this endpoint retrieves relevant documents from the Qdrant database.
    """
    results: QueryResponse = queryDatabase(
        query=request.query,
        collection_name=request.collection_name,
        limit=request.limit
    )
    return results

if __name__ == "__main__":
    import uvicorn
    startInteractor(collection_name="test_collection")
    uvicorn.run(app, host="0.0.0.0", port=8080)