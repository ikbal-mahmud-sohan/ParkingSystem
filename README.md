# Parking Mangement System

The Parking Mnagement System is designed as a gRPC gateway-based solution intended to establish a robust software architecture suitable for production environments. It can also serve as a microservice.

## Reqiurments to run

- Go v1.19+
- Docker v20+
- Docker Compose v2.23.2+

## Getting started
### Steps to run Parking Mangement System with go run .
- Open env and update config file with your db password name .
- create db as like config dbname
- open terminal 

```
$ cd api
$ go run .

```  

### Steps to run Parking Mangement System with docker

Using docker
Open terminal and follow this command
```
$ go mod tidy
$ cd deploy/compose
$ docker compose up --build -d
```
## Note: Keep Postman collection in 'postman-collection' directory