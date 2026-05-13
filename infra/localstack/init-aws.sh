#!/bin/bash
echo "Iniciando criação de recursos..."

# Esperar o LocalStack estar pronto
sleep 10

# Usar awslocal que já vem pré-instalado no container
awslocal sqs create-queue --queue-name raw-events
awslocal sqs create-queue --queue-name processed-events

awslocal dynamodb create-table \
    --table-name events \
    --attribute-definitions AttributeName=event_id,AttributeType=S \
    --key-schema AttributeName=event_id,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST

awslocal dynamodb create-table \
    --table-name developer_summary \
    --attribute-definitions AttributeName=developer_id,AttributeType=S \
    --key-schema AttributeName=developer_id,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST

echo "Recursos criados com sucesso!"
