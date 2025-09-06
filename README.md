# API de Autenticação com Go

API de autenticação e autorização desenvolvida em Go usando Gin Framework, MongoDB e Redis.

## Funcionalidades

- Registro de usuários
- Login com token JWT
- Refresh de token
- Controle de acesso baseado em roles e permissões
- Documentação Swagger

## Requisitos

- Go 1.18+
- MongoDB
- Redis

## Como executar

1. Clone o repositório
2. Configure as variáveis de ambiente (veja abaixo)
3. Execute `go run cmd/api/main.go`

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