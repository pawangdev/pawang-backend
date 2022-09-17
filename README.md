# Pawang Rest API

## Tech
Pawang Rest API uses a number of open source projects to work properly:
- [Golang](https://go.dev/) - Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.
- [Fiber](https://gofiber.io/) - Fiber is an Express inspired web framework built on top of Fasthttp, the fastest HTTP engine for Go.
- [GORM](https://gorm.io/) - The fantastic ORM library for Golang.
- [Docker](https://www.docker.com/) - Develop faster. Run anywhere.
- [MariaDB](https://mariadb.org/) - MariaDB Server is one of the most popular open source relational databases.

## Installation
Pawang REST API is very easy to install and deploy in a Docker container.
Pawang REST API requires [Docker](https://www.docker.com/) to run.

Install the dependencies and devDependencies and start the server.

```sh
cd pawang-backend
docker-composer --env-file .env up -d --build
```

Verify the deployment by navigating to your server address in
your preferred browser.

```sh
127.0.0.1:1234/api
```
