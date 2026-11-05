# Configuração local segura

Os arquivos de configuração `.env` não devem ser versionados.

1. Copie o arquivo `.env.example` do exercício correspondente para `.env` na mesma pasta.
2. Preencha localmente senhas e segredos fortes. Não envie valores reais ao GitHub, issues, PRs ou logs.
3. Se algum valor dos arquivos antigos tiver sido usado fora de um ambiente descartável, substitua-o no sistema correspondente.

**Importante:** esta alteração remove os arquivos `.env` da versão atual, mas **não apaga os valores publicados em commits históricos**. Uma reescrita de histórico, se necessária, deve ser planejada separadamente; isso também não substitui a rotação de credenciais.

Os arquivos de exemplo preservam nomes de variáveis para facilitar o estudo, mas não contêm senhas nem segredos utilizáveis.
