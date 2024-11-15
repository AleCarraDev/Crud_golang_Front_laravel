
# CRUD Golang - API de Produtos

Este projeto implementa uma API RESTful em Golang para gerenciar um sistema de produtos, com integração ao banco de dados MySQL. A API oferece operações de CRUD (Create, Read, Update, Delete) para produtos, além de funcionalidades adicionais como a busca e importação de produtos.

## Funcionalidades

A API oferece os seguintes endpoints:

### 1. **GET /produtos**
   Lista todos os produtos com suporte para paginação.

   **Exemplo de uso**:
   ```
   GET /produtos?page=1&limit=10
   ```

### 2. **GET /produtos/{id}**
   Retorna detalhes de um produto específico.

   **Exemplo de uso**:
   ```
   GET /produtos/1
   ```

### 3. **GET /produtos?nome={nome}&categoria={categoria}&preco_min={preco_min}&preco_max={preco_max}**
   Permite buscar produtos por nome e filtrar por categoria e faixa de preço.

   **Exemplo de uso**:
   ```
   GET /produtos?nome=camisa&categoria=roupas&preco_min=10&preco_max=100
   ```

### 4. **POST /produtos**
   Adiciona um novo produto.

   **Exemplo de uso**:
   ```json
   POST /produtos
   {
     "nome": "Camiseta",
     "descricao": "Camiseta de algodão",
     "preco": 29.99,
     "categoria": "Roupas"
   }
   ```

### 5. **POST /produtos/importar**
   Importa produtos em massa a partir de um arquivo JSON.

   **Exemplo de uso**:
   ```json
   POST /produtos/importar
   [
     {
       "nome": "Tênis",
       "descricao": "Tênis de corrida",
       "preco": 199.99,
       "categoria": "Calçados"
     },
     {
       "nome": "Mochila",
       "descricao": "Mochila de viagem",
       "preco": 99.99,
       "categoria": "Acessórios"
     }
   ]
   ```

### 6. **PUT /produtos/{id}**
   Edita um produto existente.

   **Exemplo de uso**:
   ```json
   PUT /produtos/1
   {
     "nome": "Camiseta Estampada",
     "descricao": "Camiseta de algodão com estampa",
     "preco": 39.99,
     "categoria": "Roupas"
   }
   ```

### 7. **DELETE /produtos/{id}**
   Exclui um produto pelo ID.

   **Exemplo de uso**:
   ```
   DELETE /produtos/1
   ```

## Requisitos

- Go (Golang) v1.16 ou superior
- MySQL v5.7 ou superior

## Como Executar o Projeto

### 1. **Instalação do Banco de Dados (MySQL)**

Certifique-se de ter o MySQL instalado e funcionando. Crie um banco de dados com o nome `produtos`:

```sql
CREATE DATABASE produtos;
```

### 2. **Instalação do Go (Golang)**

Se você não tiver o Go instalado, faça o download em: https://golang.org/dl/

### 3. **Instalação das Dependências**

Clone o repositório e instale as dependências:

```bash
git clone https://github.com/seu-usuario/crud-golang.git
cd crud-golang
go mod tidy
```

### 4. **Configuração do Banco de Dados**

No arquivo `main.go`, ajuste a string de conexão com o MySQL conforme suas credenciais:

```go
db, err = sql.Open("mysql", "root:senha@tcp(127.0.0.1:3306)/produtos")
```

### 5. **Rodando o Servidor**

Execute o servidor Go:

```bash
go run main.go
```

A API estará disponível em `http://localhost:8080`.

## Estrutura do Banco de Dados

A tabela `produtos` deve ter a seguinte estrutura:

```sql
CREATE TABLE produtos (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nome VARCHAR(255) NOT NULL,
  descricao TEXT,
  preco DECIMAL(10, 2) NOT NULL,
  categoria VARCHAR(100) NOT NULL,
  data_criacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  data_atualizacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

## Exemplo de Dados

Aqui estão alguns exemplos de dados que podem ser adicionados:

```json
[
  {
    "nome": "Camiseta",
    "descricao": "Camiseta de algodão",
    "preco": 29.99,
    "categoria": "Roupas"
  },
  {
    "nome": "Tênis",
    "descricao": "Tênis de corrida",
    "preco": 199.99,
    "categoria": "Calçados"
  }
]
```

## Contribuindo

Se você deseja contribuir para o projeto, siga os seguintes passos:

1. Fork este repositório.
2. Crie uma nova branch para suas alterações (`git checkout -b minha-branch`).
3. Realize as alterações e faça commit delas (`git commit -am 'Adicionando nova funcionalidade'`).
4. Faça push para a sua branch (`git push origin minha-branch`).
5. Abra um Pull Request para o repositório principal.

