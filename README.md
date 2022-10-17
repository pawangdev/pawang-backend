# Pawang Rest API

## Tech
Pawang Rest API uses a number of open source projects to work properly:
- [NodeJS](https://nodejs.org/) - Node.js® is an open-source, cross-platform JavaScript runtime environment.
- [ExpressJS](https://expressjs.com/) - Node.js® is an open-source, cross-platform JavaScript runtime environment.
- [Prisma ORM](https://www.prisma.io/) - Next-generation Node.js and TypeScript ORM.
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
127.0.0.1:5000/api
```
