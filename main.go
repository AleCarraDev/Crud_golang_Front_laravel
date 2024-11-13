package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"

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

    // Definindo rota básica para teste
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "API de Produtos - Conectada ao banco de dados!")
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}
