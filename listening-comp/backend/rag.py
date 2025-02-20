import chromadb
# setup Chroma in-memory, for easy prototyping. Can add persistence easily!
client = chromadb.Client()

# Create collection. get_collection, get_or_create_collection, delete_collection also available!
collection = client.create_collection("jlptn5-listening-comphresion")

# Add docs to the collection. Can also update and delete. Row-based API coming soon!
# Sample documents
doc1 = "This is a sample document about JLPT N5 listening comprehension. It covers basic Japanese conversations."
doc2 = "JLPT N5 tests basic understanding of Japanese through simple dialogues and short passages."

collection.add(
    documents=[doc1,doc2],
    metadatas=[
        {"source": "doc1.txt"},
        {"source": "doc2.txt"} 
    ],
    ids=["doc1","doc2"],
)

# Query/search 2 most similar results. You can also .get by id
results = collection.query(
    query_texts=["This is a query document"],
    n_results=2,
    # where={"metadata_field": "is_equal_to_this"}, # optional filter
    # where_document={"$contains":"search_string"}  # optional filter
)
print(results)