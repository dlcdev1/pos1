# Auction GoExpert

## Pré-requisitos

- Docker e Docker Compose instalados

## Como rodar o projeto

1. Execultar o docker-compose:
   ```sh
   docker-compose up -d
   
2. Execultar teste unitário:
   ```sh
    go test ./internal/infra/database/auction/ 

3. Criar um novo leilão:
   ```sh
   Criar um novo leilão:
   ```sh
   curl --location 'http://localhost:8080/auction' \ --header 'Content-Type: application/json' \ --data '{"product_name": "Notebook Dell XPS","category": "Eletrônicos","description": "Notebook usado, ótimo estado","condition": 1}'
   
4. Buscar os leilões:
   ```sh
   curl --location 'http://localhost:8080/auction?status=0'