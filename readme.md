# Go Project - Client Server API

## Projeto é um desafio do curso <b>Go Expert</b> da Full Cycle. Segue requisitos:

O client.go deverá realizar uma requisição HTTP no server.go solicitando a cotação do dólar.
 
O server.go deverá consumir a API contendo o câmbio de Dólar e Real no endereço: https://economia.awesomeapi.com.br/json/last/USD-BRL e em seguida deverá retornar no formato JSON o resultado para o cliente.
 
Usando o package "context", o server.go deverá registrar no banco de dados SQLite cada cotação recebida, sendo que o timeout máximo para chamar a API de cotação do dólar deverá ser de 200ms e o timeout máximo para conseguir persistir os dados no banco deverá ser de 10ms.
 
O client.go precisará receber do server.go apenas o valor atual do câmbio (campo "bid" do JSON). Utilizando o package "context", o client.go terá um timeout máximo de 300ms para receber o resultado do server.go.
 
Os 3 contextos deverão retornar erro nos logs caso o tempo de execução seja insuficiente.
 
O client.go terá que salvar a cotação atual em um arquivo "cotacao.txt" no formato: Dólar: {valor}
 
O endpoint necessário gerado pelo server.go para este desafio será: /cotacao e a porta a ser utilizada pelo servidor HTTP será a 8080.



### Este repositório contém os projetos `client` e `server` escritos em Go.

## Projeto Client

Para executar o projeto Client, siga os passos abaixo:

1. Abra o terminal
2. Navegue até a pasta 'client' (`cd client`)
3. Execute o arquivo `client.go` com o comando `go run client.go`

## Projeto Server

Para executar o projeto Server, siga os passos abaixo:

1. Abra o terminal
2. Navegue até a pasta 'server' (`cd server`)
3. Execute o arquivo `server.go` com o comando `go run server.go`

## Observações

Certifique-se de ter o Go instalado em sua máquina antes de tentar executar estes projetos.
