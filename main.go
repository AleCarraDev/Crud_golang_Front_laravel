package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

	http.HandleFunc("/produtos", listarProdutos)
	http.HandleFunc("/produtos/", getProduto)

	log.Fatal(http.ListenAndServe(":8080", nil))
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

func listarProdutos(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM produtos")
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

func getProduto(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/produtos/"):]
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
