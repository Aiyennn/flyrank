# Assignment 3: Task API with Docker and Postgres

This version of the Task API stores tasks in Postgres instead of memory. Postgres runs in Docker, and the API connects to it using the connection string from `.env`.

## Start the stack

1. Copy `.env.example` to `.env`.
2. Add your local Postgres password to `.env`.
3. Make sure Docker Desktop is running.
4. Start the stack:

```powershell
docker compose up --build
```

The API runs at `http://localhost:8080`.

Postgres uses a named Docker volume called `postgres_data`, so the database data remains after the containers are stopped or recreated.

## Configuration

The database connection string is stored in `DATABASE_URL` inside `.env`.

`.env` is gitignored and should not be committed. `.env.example` is committed as a template that shows the required configuration.

The database table is created by `db/init.sql`. Docker mounts this file into Postgres's `docker-entrypoint-initdb.d/` directory, so it runs automatically when the Postgres volume is initialized for the first time.

## API routes

| Method   | Route         | Description      |
| -------- | ------------- | ---------------- |
| `GET`    | `/tasks`      | List all tasks   |
| `POST`   | `/tasks`      | Create a task    |
| `GET`    | `/tasks/{id}` | Get a task by ID |
| `PUT`    | `/tasks/{id}` | Replace a task   |
| `DELETE` | `/tasks/{id}` | Delete a task    |

## Storage

The service and HTTP handlers were kept unchanged from Assignment 2.

The existing `TaskRepository` interface is still the boundary between the application and the storage layer. For Assignment 3, `PostgresTaskRepository` was added to implement that interface using `database/sql` and the `pgx` driver.

The original `MemoryTaskRepository` is still in the project for comparison.

The repository selected in `main.go` changed from:

```go
// Assignment 2
repository := NewMemoryTaskRepository()
```

to:

```go
// Assignment 3
repository := NewPostgresTaskRepository(db)
```

No changes were needed in `service.go` or `handlers.go`.

## Persistence check

I checked that tasks survive both an application restart and a container restart.

First, I started the stack:

```powershell
docker compose up --build -d
```

Then I created a task:

```powershell
curl.exe -X POST http://localhost:8080/tasks `
  -H "Content-Type: application/json" `
  -d '{"title":"Persistence proof","done":false}'
```

The API returned:

```text
{"id":1,"title":"Persistence proof","done":false}
```

Next, I stopped and removed the containers:

```powershell
docker compose down
```

I did not remove the Docker volume.

I then started the stack again:

```powershell
docker compose up --build -d
```

Finally, I requested the tasks again:

```powershell
curl.exe http://localhost:8080/tasks
```

The task was still present:

```text
[{"id":1,"title":"Persistence proof","done":false}]
```

This confirms that the data is stored in the Postgres volume rather than inside the application or Postgres container itself.

> **Note:** `docker compose down -v` also removes the named volume, which deletes the stored database data. Use `docker compose down` when you want to stop the stack without removing the data.

## Project files

| File                     | Purpose                                                                   |
| ------------------------ | ------------------------------------------------------------------------- |
| `Dockerfile`             | Builds the Go API using a multi-stage Docker build                        |
| `docker-compose.yml`     | Starts the API and Postgres containers and configures the database volume |
| `db/init.sql`            | Creates the `tasks` table when the database volume is initialized         |
| `.env`                   | Local database configuration; gitignored                                  |
| `.env.example`           | Template for the required environment variables                           |
| `repository.go`          | Defines the `TaskRepository` interface                                    |
| `postgres_repository.go` | Postgres implementation of `TaskRepository`                               |
| `memory_repository.go`   | Original in-memory implementation                                         |
| `service.go`             | Task validation and business logic; unchanged from Assignment 2           |
| `handlers.go`            | HTTP handlers; unchanged from Assignment 2                                |
