# ARS2024

docker build -t server .
docker run -p 8000:8000 --name alati2024 server
docker start alati2024
docker stop alati2024