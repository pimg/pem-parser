# PEM parser

Submit a PEM file and see what it contains.

## Run locally

1. Clone the repository 
2. Issue `make run` to run the application locally.
3. Go to http://localhost:8080

## Run locally via Docker Compose

1. clone the repository
2. build the docker container: `docker build . -t pem-parser`
3. Add entry to hosts file `127.0.0.1 pem-parser.local`
4. run via Docker compose: `docker compose up`
5. go to https://pem-parser.local
6. trust the self-signed certificate

