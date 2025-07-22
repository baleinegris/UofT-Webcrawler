from langchain_openai import ChatOpenAI
from langchain.schema import HumanMessage, AIMessage, SystemMessage
import os
from dotenv import load_dotenv

load_dotenv()
chatbot = None  # Global variable to hold the chatbot instance
chatbot_model = "gpt-4o-mini"

SYSTEM_PROMPT = (
    "You are a friendly, thoughtful AI assistant. Be concise, helpful, and empathetic. "
    "If you are unsure about something or lack information, say so clearly. Do not provide clinical, legal, or financial advice."
)

def startChatbot(model_name: str | None = None, temperature: float = 0.7) -> None:
    """Starts up the Chatbot using LangChain."""
    global chatbot, chatbot_model

    if not os.environ.get("OPENAI_API_KEY"):
        raise Exception(
            "OPENAI_API_KEY environment variable is not set. Please set it to your OpenAI API key."
        )

    if chatbot is not None:
        return
    try:
        reasoning = {
            "effort": "medium",  # 'low', 'medium', or 'high'
            "summary": "auto",  # 'detailed', 'auto', or None
        }
        chatbot = ChatOpenAI(model=chatbot_model)
    except Exception as e:
        chatbot = None
        raise Exception(f"Failed to start the chatbot: {str(e)}")




def queryChatbot(query: str):
    """Stream tokens from the chatbot in response to ``query``.

    The caller is responsible for assembling the final message. This
    function only yields token chunks and records metrics.
    """
    # Create a LangChain agent
    try:
        messages = [SystemMessage(content=SYSTEM_PROMPT)]
        messages.append(HumanMessage(content=query))
        for chunk in chatbot.stream(messages):
            if chunk.content:
                yield chunk.content
    except Exception as e:
        raise