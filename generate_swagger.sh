#!/bin/bash

echo "Gerando documentação Swagger..."
cd "$(dirname "$0")"
swag init -g cmd/api/main.go -o ./docs
echo "Documentação Swagger gerada com sucesso!"
