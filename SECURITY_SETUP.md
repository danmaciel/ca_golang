# Configuração local e limites de segurança

Este projeto é um **laboratório de estudo**, não um serviço recomendado para exposição pública.

1. Copie [`.env.example`](.env.example) para `.env` **na raiz**, e substitua os placeholders por senhas únicas para MySQL e RabbitMQ. Esse arquivo é lido pelo Docker Compose e ignorado pelo Git.
2. Copie [`cmd/.env.example`](cmd/.env.example) para `cmd/.env`, ajustando `DB_PASSWORD` para coincidir com `MYSQL_ROOT_PASSWORD` da raiz. O aplicativo roda no host e carrega `cmd/.env` opcionalmente; variáveis de ambiente têm precedência.
3. Defina `MIGRATE_DATABASE_URL` no shell antes de `make migrate`. O Makefile não armazena mais uma senha fixa nem imprime a URL intencionalmente.
4. Execute `docker compose config --quiet` antes de subir a infraestrutura. **A validação sintática do Compose não substitui a conferência de valores obrigatórios:** se as variáveis não estiverem configuradas, a lista `environment` poderá omiti-las (e RabbitMQ poderá usar padrões inseguros). Use o `.env` local com valores reais exclusivos, sem publicar esses dados. As portas publicadas de MySQL (3306), AMQP (5672) e painel RabbitMQ (15672) estão restritas a `127.0.0.1` do host. A rede interna dos contêineres continua com comunicação entre serviços.

**Atenção:** MySQL 5.7 está fora de suporte; foi preservado para evitar uma migração incompatível e inadvertida do laboratório. Essa restrição de portas **não equivale** a isolamento ou segurança de produção. A aplicação REST, GraphQL e gRPC não implementa autenticação/autorização, e seu próprio binding de rede deve ser revisado antes de qualquer implantação.

A senha do MySQL é definida na **primeira inicialização** do volume de dados: editar `.env` depois disso não altera automaticamente a senha do banco existente. Não remova volumes para ajustar a senha sem avaliar a perda de dados.

Os valores históricos que eventualmente tenham sido versionados em commits anteriores **não desaparecem** com essas alterações. Se algum valor já tiver sido usado fora de um ambiente descartável, rotacione as credenciais na origem. Reescrita de histórico, se necessária, é uma medida separada e exige avaliação prévia.

Não copie credenciais reais para issues, pull requests, logs, exemplos ou documentação pública.
