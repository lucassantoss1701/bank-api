<div align="center">
  <br>
  <h1>BANK API</h1>
  <br>
</div>

<p align="center">
  <a href="#-sobre-o-projeto">Sobre</a> •
  <a href="#-funcionalidades">Funcionalidades</a> •
  <a href="#-como-executar-o-projeto">Como executar</a> •
  <a href="#-como-executar-os-testes">Como testar</a> •
  <a href="#-rotas">Rotas</a>
</p>

## 💻 Sobre o projeto

O projeto tem como objetivo disponibilizar recursos para operações comuns que ocorrem dentro de um banco.

---

## 📢 Funcionalidades

- [x] Criar uma conta.
- [x] Visualizar contas existentes (id, nome, saldo e data de criação).
- [x] Pesquisar o saldo de uma conta específica.
- [x] Realizar transferências entre diferentes contas.
- [x] Visualizar transferências realizadas do usuário.
- [x] Fazer login de um usuário.

---

## 🚀 Como executar o projeto

Antes de começar, você vai precisar ter instalado em sua máquina o [DOCKER](https://docs.docker.com/engine/install/) e [DOCKER-COMPOSE](https://docs.docker.com/compose/install/).

#### 🎲 Adquirindo o repositório do projeto

```bash
# Clone este repositório
$ git clone https://github.com/lucassantoss1701/bank-api.git
```

#### 🎲 Executando a aplicação

```bash
# Rode seguinte comando no terminal (root)
$ docker compose up --build -d
# ou
make run
```

<p>✅ Pronto, a api estará rodando no host: (http://localhost:8000/)</p>

---

## 🚀 Como executar os testes

Para rodar os testes temos também duas formas para conseguir rodar.

#### 🎲 Executando os testes

```bash

# Rode seguinte comando no terminal
$ go test ./...
# ou
make tests
```

## 🕹 Rotas


### POST - /login

Realiza o login de um usuário.

curl
```bash
curl --location --request POST 'http://localhost:8000/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "cpf": "73249636096",
    "secret": "supersecret"
}'
```

resposta

```bash
{
    "token": "xpto"
}
```

### GET - /accounts?limit=10&offset=0

Retorna as contas que foram criadas no banco.

curl

```bash
curl --location --request GET 'http://localhost:8000/accounts?limit=10&offset=0' \
--header 'Authorization: Bearer token'
```

resposta

```bash
[
    {
        "id": "0b8b418c-da4a-4856-8b6a-eec63d6c7a6d",
        "name": "jaque",
        "balance": 200000,
        "created_at": "2023-08-13T19:02:19Z"
    },
    {
        "id": "640f2bea-4f97-4842-b514-0cc0b23a41f5",
        "name": "lucas",
        "balance": 200000,
        "created_at": "2023-08-13T18:51:18Z"
    }
]

```


### POST - /accounts

Cria uma nova conta.

curl

```bash
curl --location --request POST 'http://localhost:8000/accounts' \
--header 'Content-Type: application/json' \
--data-raw '{
  "name": "Lucas",
  "cpf": "73249636096",
  "secret": "supersecret",
  "balance": 200000
}
'
```

resposta

```bash
{
    "id": "0b8b418c-da4a-4856-8b6a-eec63d6c7a6d",
    "name": "jaque",
    "balance": 200000,
    "created_at": "2023-08-13T19:02:19Z"
}
```

### GET - /accounts/{id}/balance

Busca o saldo de uma conta específica.

curl 

```bash
curl --location --request GET 'http://localhost:8000/accounts/0b8b418c-da4a-4856-8b6a-eec63d6c7a6d/balance' \
--header 'Authorization: Bearer token'
```

resposta 
```bash
{
    "balance": 200000
}
```


### POST - /transfers

Realiza uma transfêrencia entre a conta logada e a conta informada no request body(conta logada é identificada atráves do token).

curl 

```bash
curl --location --request POST 'http://localhost:8000/transfers' \
--header 'Authorization: Bearer token' \
--header 'Content-Type: application/json' \
--data-raw '{
    "destination_account":{
        "id": "0b8b418c-da4a-4856-8b6a-eec63d6c7a6d"
    },
    "amount": 5000
}'
```

resposta

```bash
{
    "id": "2cb151d1-b28c-44a4-90c7-3ba18ec47c9c",
    "amount": 5000,
    "origin_account": {
        "id": "640f2bea-4f97-4842-b514-0cc0b23a41f5",
        "name": "lucas"
    },
    "destination_account": {
        "id": "0b8b418c-da4a-4856-8b6a-eec63d6c7a6d",
        "name": "jaque"
    },
    "created_at": "2023-08-13T19:59:31Z"
}
```

### GET - /transfers?limit=10&offset=0

Busca as transferências realizadas pelo usuário logado(conta logada é identificada atráves do token).

curl

```bash
curl --location --request GET 'http://localhost:8000/transfers?limit=2&offset=0' \
--header 'Authorization: Bearer token'
```

resposta

```bash
[
    {
        "id": "2cb151d1-b28c-44a4-90c7-3ba18ec47c9c",
        "destination_account": {
            "id": "0b8b418c-da4a-4856-8b6a-eec63d6c7a6d",
            "name": "jaque"
        },
        "amount": 5000,
        "created_at": "2023-08-13T19:59:32Z"
    }
]
```