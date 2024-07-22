import sys
from sentence_transformers import SentenceTransformer

def get_embedding(sentence):
    model = SentenceTransformer('all-MiniLM-L6-v2')
    embedding = model.encode([sentence])
    return embedding[0].tolist()

if __name__ == "__main__":
    sentence = sys.argv[1]
    embedding = get_embedding(sentence)
    print(embedding)
