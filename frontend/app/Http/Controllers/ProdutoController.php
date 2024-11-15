<?php


namespace App\Http\Controllers;

use Illuminate\Support\Facades\Http;

class ProdutoController extends Controller
{
    public function index()
    {
        // Fazendo uma requisição GET para a API Go
        $response = Http::get('http://127.0.0.1:8080/produtos');

        // Verificando se a requisição foi bem-sucedida
        if ($response->successful()) {
            $produtos = $response->json();
            return view('produtos.index', compact('produtos'));
        } else {
            // Caso ocorra algum erro ao fazer a requisição
            return view('produtos.index', ['produtos' => []]);
        }
    }
}
