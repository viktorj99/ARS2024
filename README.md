# ARS2024

docker build -t server .
docker run -p 8000:8000 --name alati2024 server
docker start alati2024
docker stop alati2024

URL to swagger
Api Documentation
http://localhost:8080/swagger/index.html#/


#Metrics
http://localhost:9090/graph - Prometheus UI

1. Ukupan broj zahteva za prethodna 24 sata
ars_requests_total

2. Broj uspešnih zahteva (status kodova odgovora 2xx, 3xx) za prethodna 24 sata
ars_successful_requests_total

3. Broj neuspešnih zahteva (status kodova odgovora 4xx, 5xx) za prethodna 24 sata
ars_failed_requests_total

4. Prosečno vreme izvršavanja zahteva za svaki endpoint
sum(rate(ars_request_duration_seconds_sum[24h])) / sum(rate(ars_request_duration_seconds_count[24h]))

5. Broj zahteva u jedinici vremena (minut) za svaki endpoint za prethodna 24 sata
sum(rate(ars_requests_total{endpoint="/configGroups/{name}/{version}"}[24h])) / 24 / 60

