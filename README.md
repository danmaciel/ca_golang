# Clean Architecture com Golang


## Como usar?

Entre na raiz do projeto e execute o seguinte comando para criar e subir o banco de dados Mysql e o servidor do rabbitMq:

```
docker compose up -d
```

Instalar o golang-migrate, tutorial de instalacao disponivel em <link>https://github.com/golang-migrate/migrate</link>

Depois de instalado o migrate, rode o seguinte comando na pasta raiz da aplicacao:

```
make migrate
```

Depoism ainda na raiz da aplicacao, entre no diretorio cmd e execute o seguinte comando:

```
go run main.go wire_gen.go
```