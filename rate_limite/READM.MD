# Rate Limiter em Go

## Introdução

Este **Rate Limiter** é uma solução para controlar o volume de requisições recebidas por um serviço HTTP, protegendo-o contra abusos ou picos excessivos de tráfego. Ele limita as requisições baseando-se em:

- **Endereço IP** do cliente, para restrição geral;
- **Token de acesso** passado no header `API_KEY`, permitindo regras diferenciadas para cada token.

O rate limiter é implementado como um **middleware** para o framework Gin, fácil de integrar e configurar.

---

## Como Funciona

1. A cada requisição, o middleware identifica o cliente pelo token recebido no header `API_KEY`, ou, se inexistente, pelo endereço IP.

2. O limiter mantém um contador de requisições para cada chave (token ou IP) no armazenamento configurado (por padrão, Redis).

3. O limite de requisições por segundo (RPS) é configurável:

    - Tokens têm prioridade na limitação. Se o token existe e tem limite configurado, aplica-se este limite.

    - Caso contrário, o limite aplicado é o limite por IP.

4. Sempre que a contagem de requisições ultrapassa o limite, o cliente é bloqueado e o middleware retorna HTTP **429 Too Many Requests** até que o bloqueio expire.

5. O bloqueio dura o tempo configurado (expiração). Após o tempo, a contagem é resetada e o cliente pode fazer requisições novamente.

---

## Componentes Principais

- **Middleware:** intercepta as requisições e chama o limiter para checar a permissão.

- **Limiter:** componente que implementa a lógica de contagem e regra de bloqueio.

- **Storage (Strategy):** interface para abstrair o armazenamento dos contadores (Redis é implementação padrão, mas pode ser outro).

---

## Configurações

As configurações são realizadas via **variáveis de ambiente** no arquivo `.env` ou definidas no ambiente do container/sistema.

| Nome da Variável    | Descrição                                             | Exemplo           | Valor Padrão    |
|---------------------|-------------------------------------------------------|-------------------|-----------------|
| `LIMIT_IP_RPS`      | Limite de requisições por segundo por IP (int)        | `5`               | `5`             |
| `LIMIT_TOKEN_RPS`   | Limite padrão de requisições por segundo por token    | `10`              | `10`            |
| `BLOCK_DURATION_MIN`| Duração do bloqueio após ultrapassar o limite (min)   | `5`               | `5`             |
| `TOKEN_LIMITS`      | Limites específicos para tokens (formato: token:limite, separados por vírgula) | `abc123:100,def456:50` | Vazio (usa default) |
| `REDIS_ADDR`        | Endereço do servidor Redis                              | `localhost:6379`  | `localhost:6379`|
| `REDIS_PASSWORD`    | Senha do Redis, se aplicada                            | (vazio)           | Vazio           |
| `REDIS_DB`          | Número da base de dados Redis                          | `0`               | `0`             |

---

## Como Configurar e Usar

### 1. Configurar `.env`

Crie o arquivo `.env` na raiz do projeto, por exemplo:

```env
LIMIT_IP_RPS=5
LIMIT_TOKEN_RPS=10
BLOCK_DURATION_MIN=5
TOKEN_LIMITS=abc123:100,def456:50

REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
```

### 2. Rodar o Redis

```bash
  docker run -d --name redis -p 6379:6379 redis:7
```

### 3. Iniciar o servidor

```bash
  go run main.go
```

### 4. Ou Executar o docker-compose

```bash
  docker-compose up -d
```

### 5. Testar a API
- Requisições sem token são limitadas pelo IP:


```bash
  curl -i http://localhost:8080/
```

- Requisições com token:

```bash
  curl -H "API_KEY: abc123" http://localhost:8080
```

### 6. Resposta HTTP

````
{
    "ip": "172.0.0.1",
    "message": "you have reached the maximum number of requests or actions allowed within a certain time frame",
    "status code": 429
}
````

### 7. Executar os Testes Unitários

```bash
  go test ./... -v
```

### 8. Para o servidor

```bash
  ctrl + c
```

### 9. Desenvolvido por: [David Correia](https://github.com/dlcdev1)