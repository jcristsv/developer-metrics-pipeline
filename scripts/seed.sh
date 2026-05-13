#!/bin/bash

echo "Enviando mensagens para a fila raw-events..."

# Mensagem Válida 1
aws --endpoint-url=http://localhost:4566 sqs send-message --queue-url http://localhost:4566/000000000000/raw-events --message-body '{
  "event_id": "'$(powershell -c "[guid]::NewGuid( ).ToString()")'",
  "developer_id": "dev-123",
  "metric_type": "commits",
  "value": 10,
  "repository": "org/repo-1",
  "timestamp": "'$(date -u +"%Y-%m-%dT%H:%M:%SZ")'"
}'

# Mensagem Válida 2
aws --endpoint-url=http://localhost:4566 sqs send-message --queue-url http://localhost:4566/000000000000/raw-events --message-body '{
  "event_id": "'$(powershell -c "[guid]::NewGuid( ).ToString()")'",
  "developer_id": "dev-123",
  "metric_type": "pull_requests",
  "value": 2,
  "repository": "org/repo-1",
  "timestamp": "'$(date -u +"%Y-%m-%dT%H:%M:%SZ")'"
}'

# Mensagem Inválida (Sem UUID) para testar DLQ
aws --endpoint-url=http://localhost:4566 sqs send-message --queue-url http://localhost:4566/000000000000/raw-events --message-body '{
  "event_id": "invalido",
  "developer_id": "dev-123",
  "value": 5
}'

echo "Mensagens enviadas!"

