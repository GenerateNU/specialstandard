# Special Standard Backend

## Getting Started

The backend is written in [Golang](https://go.dev/learn/) and handles code dependencies with Go modules managed within go.mod and go.sum.
In order to run the backend, the most straightforward way is to navigate to backend/cmd and execute **go run main.go** or
**go build -o main . && ./main**.

Alternatively, we will use **Docker** to build and run isolated containers to ensure that environment dependencies and runtimes are consistent across our machines.

### Steps to use Docker

Install [Docker Desktop](https://docs.docker.com/get-started/get-docker/) (or just [Docker Engine](https://docs.docker.com/engine/install/)).

In a terminal of your choice:

```bash
cd /specialstandard
docker compose up --build --watch
```

This will compose a cluster consisting of the backend and frontend containers.
Docker Watch is utilized for hot/live reloading to make development easier. The Docker Engine
watches for changes within the backend and frontend, syncs file changes from
the host to the respective container, and then restarts the respective container.  

To end development, in the terminal press ctrl+c / cmd+c OR in another terminal execute:

```bash
docker compose down
```

Open <http://localhost:8080/health> with your browser to see the result. Requests *should* be logged in your terminal.

## Postman

For further development and testing, install [Postman](https://www.postman.com/downloads/), which simplifies making network requests.

## Linting

We use [golangci-lint](https://golangci-lint.run/) to ensure consistent code quality and catch common issues. All code must pass linting before it can be merged into the main branch.

To check for linting issues:

```bash
golangci-lint run
```

To automatically fix linting issues:

```bash
golangci-lint run --fix
```

## Learn More

- [Go Modules](https://faun.pub/understanding-go-mod-and-go-sum-5fd7ec9bcc34) - article about go.mod and go.sum.
- [Tour of Go](https://go.dev/tour/welcome/1) - guided tour of Golang.
- [Go Video](https://youtu.be/8uiZC0l4Ajw?si=YJq6z9nqTN-B-c8c) - fantastic build up of data structures and syntax to write a complicated API.
- [Fiber Framework](https://docs.gofiber.io/) - web framework, similar to Express, SpringBoot, Flask, FastAPI, etc.
- [Docker Engine](https://docs.docker.com/engine/) - documentation for docker building and running docker containers.
- [pgx](https://pkg.go.dev/github.com/jackc/pgx) - driver and toolkit for PostgreSQL which can be used to interact with Supabase.
- [Supabase](https://supabase.com/docs) - database hosting service, with Auth service alongside.  

This tech stack is totally flexible--suggestions are welcome!
