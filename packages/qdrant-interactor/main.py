from fastapi import FastAPI
from pydantic import BaseModel
from fastapi.middleware.cors import CORSMiddleware

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=['*'],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

class TextEmbeddingRequest(BaseModel):
    text: str
    url: str
    position: int


@app.post("/create_embedding")
async def create_embedding(request: TextEmbeddingRequest):
    """
    Given a text, url, and position, this endpoint creates an embedding for the text and adds it to the Qdrant database.
    """
    pass


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)