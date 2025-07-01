from pydantic import BaseModel

class TextEmbeddingRequest(BaseModel):
    content: str
    url: str
    position: int
    collection_name: str
