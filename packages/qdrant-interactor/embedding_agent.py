from qdrant_client import QdrantClient, models
from schemas.QueryResponse import QueryResponse, QueryResult
import os
import uuid

client = None

def startInteractor(collection_name: str, model="BAAI/bge-small-en") -> bool:
    """
    Starts the Qdrant interactor service for given collection name.
    """
    qdrant_host = os.environ.get("QDRANT_HOST", "localhost")
    qdrant_port = os.environ.get("QDRANT_PORT", "6333")
    global client
    if client is not None:
        print("Qdrant interactor is already running!")
        return True
    client = QdrantClient(url=f"http://{qdrant_host}:{qdrant_port}")
    if client.collection_exists(collection_name):
        print(f"Collection {collection_name} exists.")
        return True
    else:
        try:
            client.create_collection(
                collection_name=collection_name,
                vectors_config=models.VectorParams(
                    size=client.get_embedding_size(model),
                    distance=models.Distance.COSINE
                )
            )
            print(f"Collection {collection_name} created.")
            return True
        except Exception as e:
            print(f"Failed to create collection {collection_name}: {e}")
            return False

def addDocument(content: str, title: str, source: str, collection_name: str, model="BAAI/bge-small-en") -> bool:
    """
    Adds a document to the specified collection.
    """
    global client
    if client is None:
        print("Qdrant interactor is not running. Please start it first.")
        return False
    if not client.collection_exists(collection_name):
        print(f"Collection {collection_name} does not exist. Please create it first.")
        client.create_collection(
            collection_name=collection_name,
            vectors_config=models.VectorParams(
                size=client.get_embedding_size(model),
                distance=models.Distance.COSINE
            )
        )
    try:
        id = str(uuid.uuid4())
        full_content = f"{title}\n{content}"
        payload = {
            "content": full_content,
            "title": title,
            "source": source
        }
        client.upsert(
            collection_name=collection_name,
            points=[
                models.PointStruct(
                    id=id,
                    payload=payload,
                    vector=models.Document(text=full_content, model=model)
                ),
            ],
        )
        print(f"Document added to collection {collection_name}.")
        return True
    except Exception as e:
        print(f"Failed to add document to collection {collection_name}: {e}")
        return False

def queryDatabase(query: str, collection_name: str, limit: int = 10, model="BAAI/bge-small-en") -> QueryResponse:
    """
    Queries the specified collection and returns a QueryResponse object with results.
    """
    global client
    if client is None:
        print("Qdrant interactor is not running. Please start it first.")
        return QueryResponse(results=[])
    try:
        results = client.query_points(
            collection_name=collection_name,
            query=models.Document(
                text=query,
                model=model
            ),
            limit=limit
        ).points
        query_results = []
        for result in results:
            if result.score >= 0.7:
                query_results.append(QueryResult(
                    id=str(result.id),
                    content=result.payload.get("content") or "",
                    source=result.payload.get("source") or "",
                    score=result.score or 0.0
                ))
        return QueryResponse(results=query_results)
    except Exception as e:
        print(f"Failed to query collection {collection_name}: {e}")
        return QueryResponse(results=[])