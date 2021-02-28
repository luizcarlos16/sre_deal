**Arquivo REPLACE.py faz as alterações em massas dentro da pasta “sre_deal”, segue exemplo.**
# python3 replace.py

*Exemplo:*

Qual diretorio quer analisar?
sre-test-1
Qual string quer procurar?
testing/sre-test-1
Qual string quer adicionar?
luizcarlos16/sre_deal


**Comando para buildar image dentro da pasta `sre_deal`, só executar.**
# docker build -t sre_deal .

**Comando para Rodar Aplicação docker dentro da pasta `sre_deal`, só executar**
# docker run -d -p 8080:8080 -p 9090:9090 --name sre_deal sre_deal

**Comando para subir Prometheus, Grafana dentro da pasta `prometheus`, só executar.**
# docker-compose up -d
