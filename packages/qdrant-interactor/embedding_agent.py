from qdrant_client import QdrantClient, models
import uuid

client = None

def startInteractor(collection_name: str, model="BAAI/bge-small-en") -> bool:
    """
    Starts the Qdrant interactor service for given collection name.
    """
    global client
    if client is not None:
        print("Qdrant interactor is already running.")
        return True
    client = QdrantClient(url="http://localhost:6333")
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

def addDocument(content: str, source: str, collection_name: str, model="BAAI/bge-small-en") -> bool:
    """
    Adds a document to the specified collection.
    """
    global client
    if client is None:
        print("Qdrant interactor is not running. Please start it first.")
        return False
    try:
        id = str(uuid.uuid4())
        payload = {
            "content": content,
            "source": source
        }
        client.upsert(
            collection_name=collection_name,
            points=[
                models.PointStruct(
                    id=id,
                    payload=payload,
                    vector=models.Document(text=content, model=model)
                ),
            ],
        )
        print(f"Document added to collection {collection_name}.")
        return True
    except Exception as e:
        print(f"Failed to add document to collection {collection_name}: {e}")
        return False
