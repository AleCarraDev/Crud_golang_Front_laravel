package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:Carra1992@tcp(127.0.0.1:3306)/produtos")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Conexão com o banco de dados estabelecida com sucesso!")

	mux := http.NewServeMux()

	mux.HandleFunc("/produtos", produtosHandler)                      // Rota para listar e adicionar produtos
	mux.HandleFunc("/produtos/importar", importarProdutos)            // Rota para importar produtos
	mux.HandleFunc("/produtos/{id:[0-9]+}", produtoEspecificoHandler) // Rota para produto específico com ID
	mux.HandleFunc("/produtos/excluir/", excluirProduto)              // Rota para excluir produto
	mux.HandleFunc("/produtos/{id}", editarProduto)                   // Rota para editar produto (PUT)

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func configurarCORS(w http.ResponseWriter, r *http.Request) {
	// Permite qualquer origem. Para segurança, pode especificar um domínio no lugar de "*"
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func produtosHandler(w http.ResponseWriter, r *http.Request) {
	// Configura CORS para esta rota
	configurarCORS(w, r)

	if r.Method == http.MethodGet {
		listarProdutos(w, r)
	} else if r.Method == http.MethodPost {
		adicionarProduto(w, r)
	}
}

func listarProdutos(w http.ResponseWriter, r *http.Request) {
	nome := r.URL.Query().Get("nome")
	categoria := r.URL.Query().Get("categoria")
	precoMin := r.URL.Query().Get("preco_min")
	precoMax := r.URL.Query().Get("preco_max")
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	pageInt, _ := strconv.Atoi(page)
	if pageInt <= 0 {
		pageInt = 1
	}
	limitInt, _ := strconv.Atoi(limit)
	if limitInt <= 0 {
		limitInt = 10
	}

	query := "SELECT * FROM produtos WHERE 1=1"
	args := []interface{}{}

	if nome != "" {
		query += " AND nome LIKE ?"
		args = append(args, "%"+nome+"%")
	}
	if categoria != "" {
		query += " AND categoria = ?"
		args = append(args, categoria)
	}
	if precoMin != "" {
		query += " AND preco >= ?"
		args = append(args, precoMin)
	}
	if precoMax != "" {
		query += " AND preco <= ?"
		args = append(args, precoMax)
	}

	query += " LIMIT ? OFFSET ?"
	args = append(args, limitInt, (pageInt-1)*limitInt)

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var produtos []Produto
	for rows.Next() {
		var p Produto
		if err := rows.Scan(&p.ID, &p.Nome, &p.Descricao, &p.Preco, &p.Categoria, &p.DataCriacao, &p.DataAtualizacao); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		produtos = append(produtos, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(produtos)
}

func adicionarProduto(w http.ResponseWriter, r *http.Request) {
	var produto Produto
	if err := json.NewDecoder(r.Body).Decode(&produto); err != nil {
		http.Error(w, "Erro ao ler o corpo da requisição", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO produtos (nome, descricao, preco, categoria) VALUES (?, ?, ?, ?)", produto.Nome, produto.Descricao, produto.Preco, produto.Categoria)
	if err != nil {
		http.Error(w, "Erro ao adicionar produto", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(produto)
}

func editarProduto(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/produtos/"):]

	if strings.TrimSpace(idStr) == "" {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var produto Produto
	if err := json.NewDecoder(r.Body).Decode(&produto); err != nil {
		http.Error(w, "Erro ao ler o corpo da requisição", http.StatusBadRequest)
		return
	}

	produto.DataAtualizacao = fmt.Sprintf("%s", time.Now().Format("2006-01-02 15:04:05"))

	_, err = db.Exec("UPDATE produtos SET nome = ?, descricao = ?, preco = ?, categoria = ?, data_atualizacao = ? WHERE id = ?",
		produto.Nome, produto.Descricao, produto.Preco, produto.Categoria, produto.DataAtualizacao, id)
	if err != nil {
		http.Error(w, "Erro ao editar produto", http.StatusInternalServerError)
		return
	}

	var p Produto
	err = db.QueryRow("SELECT * FROM produtos WHERE id = ?", id).Scan(&p.ID, &p.Nome, &p.Descricao, &p.Preco, &p.Categoria, &p.DataCriacao, &p.DataAtualizacao)
	if err != nil {
		http.Error(w, "Produto não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Retorna status 200 OK
	json.NewEncoder(w).Encode(p)
}

func excluirProduto(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/produtos/excluir/"):]

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM produtos WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Erro ao excluir produto", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"mensagem": "Produto excluído com sucesso"})
}

func produtoEspecificoHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/produtos/"):]

	if strings.TrimSpace(idStr) == "" {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var p Produto
	err = db.QueryRow("SELECT * FROM produtos WHERE id = ?", id).Scan(&p.ID, &p.Nome, &p.Descricao, &p.Preco, &p.Categoria, &p.DataCriacao, &p.DataAtualizacao)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Produto não encontrado", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func importarProdutos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var produtos []Produto
	if err := json.NewDecoder(r.Body).Decode(&produtos); err != nil {
		http.Error(w, "Erro ao ler o corpo da requisição", http.StatusBadRequest)
		return
	}

	for _, produto := range produtos {
		_, err := db.Exec("INSERT INTO produtos (nome, descricao, preco, categoria) VALUES (?, ?, ?, ?)",
			produto.Nome, produto.Descricao, produto.Preco, produto.Categoria)
		if err != nil {
			http.Error(w, "Erro ao importar produtos", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"mensagem": "Produtos importados com sucesso"})
}

type Produto struct {
	ID              int     `json:"id"`
	Nome            string  `json:"nome"`
	Descricao       string  `json:"descricao"`
	Preco           float64 `json:"preco"`
	Categoria       string  `json:"categoria"`
	DataCriacao     string  `json:"data_criacao"`
	DataAtualizacao string  `json:"data_atualizacao"`
}
