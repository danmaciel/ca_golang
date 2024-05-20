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

Depois ainda na raiz da aplicacao, entre no diretorio cmd e execute o seguinte comando:

```
go run main.go wire_gen.go
```

### Para testar via Rest
Acesse o diretorio "api" na raiz da aplicacao e execute os arquivos .http, create_order.http para criar e list_orders.http para listar

### Para testar via GraphQl
Com a aplicacao rodando, entre abra no browser http://localhost:8080/ e execute:

```
query queryListOrders{
  orders {
    id, Price, Tax, FinalPrice
  }
}
```


### Para testar gRPC

Instale o evans de acordo com o repositorio <link>https://github.com/ktr0731/evans</link> e com a aplicacao rodando execute em um terminal:

```
evans -r rpl
```

Depois

```
package pb
```

Depois

```
service OrderService
```

E para retornar a lista de Order

```
call GetListOrder
```

Se for necessario criar uma Order

```
call CreateOrder
```

