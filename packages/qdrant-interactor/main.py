from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from schemas.TextEmbeddingRequest import TextEmbeddingRequest
from embedding_agent import startInteractor, addDocument

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=['*'],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.post("/add_embedding")
async def add_embedding(request: TextEmbeddingRequest):
    """
    Given a text, url, and position, this endpoint creates an embedding for the text and adds it to the Qdrant database.
    """
    status = addDocument(content=request.content, source=request.url, collection_name="test_collection")
    if status is True:
        return {"message": "Document added successfully."}
    else:
        return {"message": "Failed to add document."}



if __name__ == "__main__":
    import uvicorn
    startInteractor(collection_name="test_collection")
    uvicorn.run(app, host="0.0.0.0", port=8000)