from pydantic import BaseModel

class TextEmbeddingRequest(BaseModel):
    content: str
    title: str
    url: str
    position: int
    collection_name: str
