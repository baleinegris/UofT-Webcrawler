from fastapi import FastAPI
from pydantic import BaseModel
from chatbot import startChatbot, queryChatbot
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import StreamingResponse


class QueryRequest(BaseModel):
    query: str

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=['*'],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.on_event("startup")
async def startup_event():
    """Initialize the chatbot on startup."""
    try:
        startChatbot()
    except Exception as e:
        print(f"Error starting chatbot: {e}")

@app.post("/query")
async def query(request: QueryRequest):
    return StreamingResponse(
        queryChatbot(request.query),
        media_type="text/event-stream",
        )

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8001)