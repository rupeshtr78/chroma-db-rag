# chroma-db

This repository includes a Go project for working with Chroma DB and ollama embedding models. It can set up and run a vector database using Chroma DB, handle user queries, and interact with an embedding model. It also includes a Dockerfile for running the project locally and a docker compose file to run the project in a container.
## Sample Results
```txt
<|user|> what is mirostat_tau?</s>:-
Based on the provided content, I can answer your query.

**Query Result:** Mirostat_tau Controls the balance between coherence and diversity of the output. A lower value will result in more focused and coherent text. (Default: 5.0)

**Document Content:**

mirostat_tau Controls the balance between coherence and diversity of the output. A lower value will result in more focused and coherent text. (Default: 5.0)
float
mirostat_tau 5.0

**Additional Information on this Topic:**

Here are three main points related to Mirostat_tau:

1. **Coherence vs Diversity:** Mirostat_tau controls the balance between coherence and diversity of the output, which means it determines how focused or creative the generated text will be.
2. **Lower Values Mean More Focus:** A lower value for mirostat_tau results in more focused and coherent text, while a higher value allows for more diverse and potentially less coherent output.
3. **Default Value:** The default value for Mirostat_tau is 5.0, which means that if no specific value is provided, the model will generate text with a balance between coherence and diversity.

Please note that these points are based solely on the provided content and do not go beyond it.%    
```
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

- **docker-compose.yaml**: Docker Compose configuration file for setting up the Chroma DB service.

## Configuration

Adjust configuration values in `internal/constants/constants.go` to fit your needs. This includes settings like:

Chroma DB URL, Tenant name, Database & Namespace.
Ollama model type and URL.

### Functionality

#### Running VectorDB

Start the VectorDb with the following command:

```sh
docker compose up
```

#### Chat with Ollama

Execute chat-related operations:

```sh
go run ./cmd/main.go
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



