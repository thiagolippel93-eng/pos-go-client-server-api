▶️ Como Executar
1. Clonar o projeto
git clone https://github.com/seu-usuario/desafio-go-cotacao.git
cd desafio-go-cotacao

2. Inicializar módulo Go
go mod init desafio-go
go mod tidy

3. Executar o servidor
go run server.go

O servidor estará disponível em:

http://localhost:8080/cotacao
4. Executar o cliente

Em outro terminal:

go run client.go