
# Projeto Neoway V1

## Descrição do Projeto

Este projeto é uma aplicação escrita em Go, com o objetivo de processar arquivos de dados e armazená-los em um banco de dados PostgreSQL. A aplicação lida com grandes volumes de dados, processando-os em lotes e inserindo-os em tabelas relacionadas a clientes, lojas e transações. Além disso, o sistema valida e atualiza o status de CPFs e CNPJs conforme as regras de negócio.

## Estrutura do Projeto

A estrutura do projeto foi organizada da seguinte maneira:

```
Neoway V1/
├── cmd/
│   └── main.go               # Ponto de entrada da aplicação
├── internal/
│   ├── batch/                    # Lógica de processamento em lote
│   ├── db/                       # Conexão com banco de dados e manipulação de tabelas
│   ├── fileprocessor/            # Processamento de arquivos de entrada
│   └── validation/               # Validação de CPFs e CNPJs
├── assets/                       # Arquivos de dados para processamento
│   └── base_teste.txt
├── go.mod                        # Arquivo de gerenciamento de dependências Go
└── README.md                     # Documentação do projeto
```

### Descrição dos Diretórios
- **cmd/app**: Contém o arquivo `main.go`, que é o ponto de entrada do projeto, onde o serviço é inicializado.
- **internal/batch**: Contém a lógica para inserção de dados em lotes no banco de dados.
- **internal/db**: Lida com a conexão ao banco de dados PostgreSQL e a criação de tabelas, além de funções para atualizar os status de CPFs e CNPJs.
- **internal/fileprocessor**: Lida com o processamento de arquivos de entrada, convertendo-os em um formato que possa ser inserido no banco de dados.
- **internal/validation**: Contém funções para validar CPFs e CNPJs, além de formatação adequada dos mesmos.
- **assets**: Contém arquivos de dados de teste, que podem ser utilizados no processamento de entrada.

## Tecnologias Utilizadas

- **Go**: Linguagem principal utilizada no desenvolvimento do projeto.
- **PostgreSQL**: Banco de dados relacional utilizado para armazenar os dados processados.
- **pgx**: Biblioteca Go utilizada para conectar e interagir com o banco de dados PostgreSQL.
- **godotenv**: Biblioteca Go utilizada para carregar variáveis de ambiente do arquivo `.env`.

Aqui está uma versão mais resumida da **Estrutura do Banco de Dados**:

## Estrutura do Banco de Dados

O banco de dados utiliza PostgreSQL para armazenar informações sobre clientes, lojas e transações. A seguir, a descrição das tabelas:

### Tabela `customers`

Armazena informações dos clientes, como CPF, status de privacidade, e CNPJs das lojas mais frequentes e recentemente visitadas.

```sql
CREATE TABLE IF NOT EXISTS customers (
    cpf TEXT PRIMARY KEY,                        
    private BOOLEAN,                             
    incomplete BOOLEAN,                          
    status_cpf TEXT,                             
    most_frequent_store_cnpj TEXT,               
    last_store_cnpj TEXT,                        
    status_cnpj_last_store TEXT,                 
    status_cnpj_frequent_store TEXT              
);
```

### Tabela `stores`

Armazena os CNPJs das lojas e os status de validação dos CNPJs vinculados aos clientes.

```sql
CREATE TABLE IF NOT EXISTS stores (
    id SERIAL PRIMARY KEY,                       
    cnpj TEXT UNIQUE,                            
    status_cnpj_last_store TEXT,                 
    status_cnpj_frequent_store TEXT              
);
```

### Tabela `transactions`

Armazena as transações realizadas pelos clientes, vinculando-as ao CPF de cada cliente.

```sql
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,                       
    cpf TEXT REFERENCES customers(cpf),          
    last_purchase_date DATE,                     
    average_ticket NUMERIC(10, 2),               
    last_ticket NUMERIC(10, 2)                   
);
```

### Relacionamentos

- **customers** e **transactions**: As transações estão associadas aos clientes por meio do CPF.
- **customers** e **stores**: Clientes têm CNPJs de lojas visitadas, vinculados logicamente à tabela `stores`.

## Pré-requisitos

Para rodar este projeto, é necessário ter instalado:
- Go (versão 1.16 ou superior)
- PostgreSQL (versão 12 ou superior)
- Git (para clonar o repositório)

## Configuração do Projeto

### 1. Clonar o Repositório

```bash
git clone https://github.com/monalizaloren/Neoway.git
cd Neoway
```

### 2. Instalar Dependências

Certifique-se de que está no diretório raiz do projeto e rode o seguinte comando para instalar as dependências:

```bash
go mod tidy
```

### 3. Configuração do Banco de Dados

- Instale e configure o PostgreSQL na sua máquina.
- Crie um banco de dados no PostgreSQL, por exemplo, `neoway_db`.
  
No PostgreSQL, crie um banco de dados:

```sql
CREATE DATABASE neoway_db;
```

### 4. Configuração de Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto, contendo a seguinte variável de ambiente:

```
DATABASE_URL=postgres://usuario:senha@localhost:5432/neoway_db
```

Substitua `usuario`, `senha`, e `neoway_db` pelos seus próprios valores.

### 5. Inicialização do Banco de Dados

A aplicação automaticamente criará as tabelas necessárias no banco de dados ao iniciar. Certifique-se de que as tabelas `customers`, `stores` e `transactions` sejam criadas corretamente.

### 6. Executando o Projeto

Para rodar a aplicação, execute o seguinte comando:

```bash
go run ./cmd 
```

A aplicação irá:
1. Configurar a conexão com o banco de dados.
2. Verificar se as tabelas necessárias existem e criá-las se necessário.
3. Processar o arquivo de entrada (`assets/base_teste.txt`) e inserir os dados no banco de dados.
4. Atualizar os status dos CPFs e CNPJs com base nas regras de validação.

## Notas Adicionais

- **Gerenciamento de Lotes**: O processamento de dados é feito em lotes de 50.000 registros, para garantir performance e evitar problemas de memória.
