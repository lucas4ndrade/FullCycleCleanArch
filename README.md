## Clean Architecture challenge

Repositório para guardar o desafio de clean architecture do curso GO Expert da Fullcycle

## Dependências
- Ter o `docker` e o `docker-compose` instalado na máquina
- Ter um motor para rodar `Makefile`s na máquina

## Execução

### Preparando o ambiente
A raiz do projeto contém um arquivo `.env` com as configurações do mesmo.

Se necessário, deve-se configurar este arquivo com as configurações da sua máquina local.

```.env
# MYSQL DB
DB_DRIVER={driver de banco utilizado}
DB_HOST={host do banco}
DB_PORT={porta do banco}
DB_USER={usuário do banco}
DB_PASSWORD={senha do banco}
DB_NAME={nome do banco de dados que será utilizado}

# AMQP
AMQP_URL={url aonde está rodando o rabbitmq}

# SERVER PORTS
WEB_SERVER_PORT={porta desejada para subir a API HTTP REST}
GRPC_SERVER_PORT={porta desejada para subir o servidor gRPC}
GRAPHQL_SERVER_PORT={porta desejada para subir o playground GRAPHQL}
```

### Rodando recursos de infra na máquina local
O projeto tem como dependência o RabbitMQ e o MySQL.

Para executá-los, basta rodar:
```shell
make run-infra
```
O script make utilizará o arquivo `docker-compose.yaml` para rodar os recursos de infra nas portas padrão.

### Preparando o banco de dados
Após subir os recursos de infra, convém preparar o banco para executar as ações do projeto.

Para isso, basta utilizar o comando `migration-up` do `Makefile`
```shell
make migration-up
```

O comando tem alguns parâmetros que podem ser utilizados para configurar qual banco será utilizado
- DB_HOST -> (default: `localhost`)
- DB_PORT -> (default: `3306`)
- DB_USER -> (default: `root`)
- DB_PASS -> (default: `root`)
- MIGRATION_UP_FILE -> (default: `./cmd/migrations/migration_up.sql`)

### Rodando o projeto
Por fim, após rodar e configurar as dependências, basta executar
```shell
make run
```
Que o projeto será executado.

O projeto é composto por três servidores, um HTTP, um gRPC, e um playground GraphQL.
As portas aonde estes servidores estarão disponíveis estão configuradas no arquivo `.env` do projeto.
(Vide [Preparando o ambiente](#preparando-o-ambiente) para mais informações)
