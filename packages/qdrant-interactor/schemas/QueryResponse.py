from pydantic import BaseModel

class QueryResult(BaseModel):
    id: str
    content: str
    source: str
    score: float

class QueryResponse(BaseModel):
    results: list[QueryResult]