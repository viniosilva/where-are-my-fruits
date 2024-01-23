# Where are my fruits?

Gerenciamento de frutas em baldes conforme [desafio](./docs/desafio_backend_planne.pdf)

`[Desenho da arquitetura da base de dados aqui]`

## Requisitos

Para executar a API é necessário instalar as seguintes ferramentas:

- [go](https://tip.golang.org/doc/go1.20)
- [mockgen](https://github.com/golang/mock)
- [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- [swaggo](https://github.com/swaggo/swag)
- [docker](https://www.docker.com/)
- [docker-compose](https://docs.docker.com/compose/)

## Construído com

- [gin](https://gin-gonic.com)
- [gorm](https://gorm.io)
- [mysql](https://www.mysql.com)
- [swagger](https://swagger.io)

## Instalação

```bash
$ go get
```

## Configuração

- Copie o arquivo `.env.example` para `.env`

### Migração da base de dados

Depois de rodar `docker-compose up -d`, é necessário esperar alguns segundos para rodar o `make migrate` para a criação do schema do banco de dados MySQL.

```bash
$ docker-compose up -d
$ make migrate
```

## Rodando

```bash
$ make run
```

É possível acessar a documentação local das rotas no [swagger](http:localhost:8080/api/swagger/index.html)

---

## Testes

```bash
# Testes unitários
$ make test

# Cobertura dos testes unitários
$ make test/cov

# Testes E2E
$ make test/e2e
```

`[Estrutura de diretórios aqui]`