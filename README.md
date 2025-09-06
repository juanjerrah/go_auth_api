# API de Autenticação com Go

API de autenticação e autorização desenvolvida em Go usando Gin Framework, MongoDB e Redis.

## Funcionalidades

- Registro de usuários
- Login com token JWT
- Refresh de token
- Controle de acesso baseado em roles e permissões
- Documentação Swagger

## Requisitos

### Execução Local
- Go 1.18+
- MongoDB
- Redis

### Execução com Docker
- Docker
- Docker Compose

## Como executar

### Execução Local

1. Clone o repositório
2. Configure as variáveis de ambiente (veja abaixo)
3. Execute `go run cmd/api/main.go`

### Execução com Docker

1. Clone o repositório
2. Crie um arquivo `.env` na raiz do projeto com as variáveis de ambiente necessárias (veja exemplo abaixo)
3. Execute `docker compose up -d`

```bash
# Exemplo de arquivo .env para Docker
SERVER_PORT=8080
JWT_SECRET=your-secret-key
TOKEN_EXPIRES_IN=3600
MONGODB_USER=admin
MONGODB_PASSWORD=password
MONGODB_DATABASE=Users
REDIS_PASSWORD=password
REDIS_DB=0
```

4. A aplicação estará disponível em `http://localhost:8080`
5. MongoDB estará acessível em `localhost:27018` e Redis em `localhost:6380`
6. Para parar os containers, execute `docker-compose down`

## Variáveis de ambiente

```
SERVER_PORT=8080
JWT_SECRET=your-secret-key
TOKEN_EXPIRES_IN=3600
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=Users
REDIS_URI=127.0.0.1:6379
REDIS_PASSWORD=
REDIS_DB=0
```

## Documentação Swagger

A API inclui documentação Swagger para facilitar o entendimento e teste dos endpoints.

### Acesso à documentação

Após iniciar a aplicação, acesse:

```
http://localhost:8080/swagger/index.html
```

### Regenerar documentação

Se você fizer mudanças nos endpoints ou modelos, regenere a documentação usando:

```
./generate_swagger.sh
```

### Autenticação no Swagger UI

1. Faça login com POST `/api/auth/login` para obter um token
2. Clique no botão "Authorize" no topo da página do Swagger
3. Insira o token no formato `Bearer seu_token_aqui`
4. Agora você pode acessar os endpoints protegidos

## Estrutura do projeto

- `cmd/api`: Ponto de entrada da aplicação
- `internal/config`: Configurações da aplicação
- `internal/delivery/http`: Handlers HTTP e configuração de rotas
- `internal/domain`: Regras de negócio e entidades
- `internal/infrastructure`: Implementação de repositórios (MongoDB, Redis)
- `internal/utils`: Utilitários compartilhados
- `pkg/middleware`: Middlewares como autenticação
- `pkg/types`: Tipos compartilhados

## Arquitetura Técnica

```
                                   ┌──────────────────────┐
                                   │                      │
                                   │     Cliente HTTP     │
                                   │   (Navegador/App)    │
                                   │                      │
                                   └──────────┬───────────┘
                                              │
                                              │ HTTP Requests
                                              ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                             │
│                              Go Auth API                                    │
│                                                                             │
│  ┌───────────────┐      ┌───────────────┐       ┌──────────────────┐        │
│  │               │      │               │       │                  │        │
│  │  HTTP Routes  │─────▶│    Handlers   │──────▶│     Services     │        │
│  │   (Gin)      │      │               │       │                  │        │
│  │               │      │               │       │                  │        │
│  └───────────────┘      └───────┬───────┘       └────────┬─────────┘        │
│                                 │                        │                   │
│                         Auth Middleware                  │                   │
│                                 │                        │                   │
│                                 ▼                        ▼                   │
│                         ┌───────────────┐       ┌──────────────────┐        │
│                         │               │       │                  │        │
│                         │  JWT Manager  │       │  Repositories    │        │
│                         │               │       │                  │        │
│                         └───────────────┘       └────────┬─────────┘        │
│                                                          │                   │
└─────────────────────────────────────────────────────────┼───────────────────┘
                                                          │
                          ┌─────────────────┬─────────────┴──────────────┐
                          │                 │                            │
                          ▼                 ▼                            ▼
             ┌─────────────────────┐ ┌─────────────────┐      ┌──────────────────┐
             │                     │ │                 │      │                  │
             │      MongoDB        │ │      Redis      │      │    Swagger       │
             │  (Persistência)     │ │  (Cache/Tokens) │      │  Documentação    │
             │                     │ │                 │      │                  │
             └─────────────────────┘ └─────────────────┘      └──────────────────┘
```

## Fluxo de Autenticação

1. **Registro de Usuário**:
   - O cliente envia dados de registro para `/api/auth/register`
   - O sistema valida os dados, criptografa a senha e armazena no MongoDB
   - Retorna confirmação de registro

2. **Login**:
   - O cliente envia credenciais para `/api/auth/login`
   - O sistema valida as credenciais contra os registros no MongoDB
   - Se válido, gera token JWT de acesso e refresh token
   - O refresh token é armazenado no Redis com TTL
   - Retorna ambos os tokens para o cliente

3. **Autorização**:
   - O cliente inclui o token JWT no header de autorização das requisições
   - O middleware de autenticação valida o token
   - Verifica permissões de acesso baseadas em roles do usuário
   - Permite ou nega o acesso ao recurso solicitado

4. **Refresh de Token**:
   - Quando o token de acesso expira, o cliente envia o refresh token para `/api/auth/refresh`
   - O sistema valida o refresh token contra o Redis
   - Se válido, gera um novo token JWT de acesso
   - Atualiza ou mantém o refresh token no Redis
   - Retorna o novo token de acesso para o cliente

5. **Logout**:
   - O cliente envia o refresh token para `/api/auth/logout`
   - O sistema remove o refresh token do Redis
   - Confirma o logout bem-sucedido