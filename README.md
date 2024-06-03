# ARS2024

docker build -t server .
docker run -p 8000:8000 --name alati2024 server
docker start alati2024
docker stop alati2024

#URL to swagger
http://localhost:8080/swagger/index.html#/

#Metrics

#Ukupan broj get zahteva u prethodna 24 sata
sum(increase(http_requests_total[24h]))

#Broj uspešnih zahteva (status kodovi odgovora 2xx, 3xx) za prethodna 24 sata:
sum(increase(http_requests_total{status=~"2..|3.."}[24h]))

#Broj neuspešnih zahteva (status kodovi odgovora 4xx, 5xx) za prethodna 24 sata:
sum(increase(http_requests_total{status=~"4..|5.."}[24h]))

#Prosečno vreme izvršavanja zahteva za sve endpointe
sum(http_response_time_sum) / sum(http_response_time_count)

#Prosečno vreme izvršavanja zahteva za svaki endpoint
sum(http_response_time_sum{job="myapp"}) by (endpoint) / sum(http_response_time_count{job="myapp"}) by (endpoint)

#Broj zahteva u jedinici vremena (minut ili sekund) za svaki endpoint za prethodna 24 sata:
sum by (endpoint) (rate(http_requests_total[1m]))
