from fastapi import FastAPI
from pydantic import BaseModel
from chatbot import startChatbot, queryChatbot

class QueryRequest(BaseModel):
    query: str

app = FastAPI()

@app.on_event("startup")
async def startup_event():
    """Initialize the chatbot on startup."""
    try:
        startChatbot()
    except Exception as e:
        print(f"Error starting chatbot: {e}")

@app.get("/query")
async def query(request: QueryRequest):
    response = queryChatbot(request.query)
    return response

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)