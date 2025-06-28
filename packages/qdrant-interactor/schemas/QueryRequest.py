from pydantic import BaseModel

class QueryRequest(BaseModel):
    query: str
    collection_name: str
    limit: int = 10