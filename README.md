# chroma-db

This repository includes a Go project for working with Chroma DB and embedding models. It can set up and run a vector database using Chroma DB, handle user queries, and interact with an embedding model. It also includes a Dockerfile for running the project locally and a docker compose file to run the project in a container.
## Getting Started

### Prerequisites

- Go (>=1.22.0)
- Docker
- Docker Compose

### Installation

1. **Clone the Repository**

```sh
git clone https://github.com/yourusername/chroma-db.git
cd chroma-db
```

2. **Install Go Packages**
3. **Build the Go Project**

```sh
go build -o chroma-db cmd/main.go
```

4. **Set Up Docker Containers**

Ensure Docker and Docker Compose are installed. Use the `docker-compose.yaml` to set up the Chroma DB service.

```sh
docker-compose up -d
```

### Running the Project

```sh
./chroma-db
```

## Project Structure

- **cmd/**:
  - **main.go**: Entry point for running the Chroma DB.
  - **chat/**:
    - **ollama_chat.go**: Contains the logic for interacting with the Ollama chat model.

- **internal/constants/**:
  - **constants.go**: Houses all the necessary constants used across the project.

- **test/**:
  - **chroma-ollama.go**: Contains sample queries to interact with Chroma DB.

- **chromadb/**: Directory intended for Chroma DB related files and volumes.

- **docker-compose.yaml**: Docker Compose configuration file for setting up the Chroma DB service.

## Scripts

- **Makefile**: For managing the generation, build, and installation of Python bindings.
- **setup.py**: Setup script for Python bindings.

### Functionality

#### Running VectorDB

Start the VectorDb with the following command:

```sh
db.RunVectorDb(ctx)
```

#### Chat with Ollama

Execute chat-related operations:

```sh
chat.ChatOllama(ctx)
```

#### Sample Query (test/chroma-ollama.go)

```go
func SampleQuery() []exampleCase {
    type filter = map[string]any
    // ... example cases
    return exampleCases
}
```

## Configuration

Default configuration values are provided in `internal/constants/constants.go` and can be adjusted as per your needs. Some of these include:

- `ChromaUrl`, `TenantName`, `Database`, `Namespace`
- `OllamaModel` and `OllamaUrl`

### License

This project is licensed under the BSD 3-Clause License - see the [LICENSE](./LICENSE) file for details.

## Acknowledgments

- [Chroma DB](https://github.com/chroma-db)
- [Ollama](https://ollama-ai.com)

For any issues or contributions, please open an issue or submit a pull request on GitHub.



