<!DOCTYPE html>
<html lang="pt-br">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Lista de Produtos</title>
</head>
<body>
    <h1>Produtos Disponíveis</h1>
    <ul>
        @foreach($produtos as $produto)
            <li>
                <strong>Nome:</strong> {{ $produto['nome'] }} <br>
                <strong>Preço:</strong> R$ {{ number_format($produto['preco'], 2, ',', '.') }} <br>
                <strong>Categoria:</strong> {{ $produto['categoria'] }} <br>
                <strong>Descrição:</strong> {{ $produto['descricao'] }} <br>
            </li>
        @endforeach
    </ul>
</body>
</html>
