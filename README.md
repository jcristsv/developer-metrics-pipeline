# Developer Metrics Pipeline 🚀

Este projeto é uma solução escalável para processamento de métricas de desenvolvedores, utilizando **Go**, **AWS SQS** e **AWS DynamoDB** (simulados via **LocalStack**).

## 🏗️ Arquitetura
O sistema foi desenhado seguindo os princípios de **Clean Architecture** (Organização do código/Separando Responsabilidades) e **S.O.L.I.D.** (	
Boas práticas de código/Código Sustentável), dividido em:
- **Processor Service:** Responsável por consumir eventos brutos, validar integridade (UUID, campos obrigatórios) e enriquecer os dados.
- **Aggregator Service:** Consome os eventos validados, garante a idempotência e persiste as métricas agregadas no DynamoDB.
- **API REST:** Exposta pelo Aggregator para consulta de saúde e sumário de métricas.

## 🛠️ Tecnologias Utilizadas
- **Go (Golang)** 1.21+
- **Docker & Docker Compose**
- **LocalStack** (SQS & DynamoDB)(Fila de mensagens/Banco de dados NoSQL)
- **Gin Gonic** (Framework Web)
- **AWS SDK for Go v2**

## 🚀 Como Rodar
1. Certifique-se de ter o Docker instalado e rodando em sua máquina.
2. Na raiz do projeto, execute o comando principal:
   docker-compose up --build
3. Aguarde a inicialização: O sistema estará pronto quando você visualizar a mensagem Recursos criados com sucesso! nos logs do LocalStack.

## 📊 Endpoints Disponíveis
Health Check: GET http://localhost:8080/health
Developer Summary: GET http://localhost:8080/metrics/{developer_id}/summary

## 🧪 Como Testar (Simulação de Evento )
Para validar o pipeline ponta a ponta, você pode injetar um evento manualmente. Com o sistema rodando, execute o comando abaixo em um novo terminal:
docker exec localstack awslocal sqs send-message --queue-url http://localhost:4566/000000000000/raw-events --message-body "{\"event_id\": \"550e8400-e29b-41d4-a716-446655440000\", \"developer_id\": \"dev-123\", \"metric_type\": \"commits\", \"value\": 10}"

Após o envio, verifique o resultado no navegador acessando o endpoint de summary.

## 🎥 Demonstração
O vídeo com a explicação técnica e demonstração do sistema funcionando pode ser acessado no link abaixo:

## 🛠️🚀 O que faria diferente com mais tempo
Isso demonstra visão de produto e engenharia de longo prazo. Sugestões baseadas em boas práticas:
- Observabilidade: Implementaria métricas (Prometheus/Grafana) e tracing (AWS X-Ray) para monitorar o pipeline.
- Infraestrutura como Código (IaC): Utilizaria Terraform para definir os recursos da AWS em vez de scripts manuais.
- Segurança: Implementaria autenticação/autorização (JWT) na API REST.
