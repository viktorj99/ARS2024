# ARS2024

docker build -t server .
docker run -p 8000:8000 --name alati2024 server
docker start alati2024
docker stop alati2024

#URL to swagger
http://localhost:8080/swagger/index.html#/

#Metrics

#Ukupan broj get zahteva
#sum(http_method{method="GET"})

#Prosecan broj get zahteva u roku od 5 minuta
##rate(http_requests_total{method="GET"}[5m])
