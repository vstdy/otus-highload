# OTUS-Highload

OTUS Highload Architect course project.  

## REST API

By default, server starts at `8080` HTTP port with the following endpoints:

- `POST /login` — login user;
- `POST /user/register` — register user;
- `GET /user/get/{id}` — get user data;

For details check out [***http-client.http***](./http-client.http) file


## Code
### Packages used

App architecture and configuration:

- [viper](https://github.com/spf13/viper) - app configuration;
- [cobra](https://github.com/spf13/cobra) - CLI;
- [zerolog](https://github.com/rs/zerolog) - logger;

Networking:

- [go-chi](https://github.com/go-chi/chi) - HTTP router;

SQL database interface provider:

- [pgx](https://github.com/jackc/pgx) - Go driver and toolkit for PostgreSQL;
- [scany](https://github.com/georgysavva/scany) - scanner toolkit;

## CLI

All CLI commands have the following flags:
- `--log_level`: (optional) logging level (default: `info`);
- `--config`: (optional) path to configuration file (default: `./config.yaml`);
- `--timeout`: (optional) request timeout (default: `5s`);
- `-d --database_dsn`: (optional) database source name (default: `postgres://user:password@localhost:5432/project?sslmode=disable`);

Root only command flags:
- `-a --server_address`: (optional) server address (default: `0.0.0.0:8080`);

If config file not specified, defaults are used. Defaults can be overwritten using ENV variables.

### Migrations

    project migrate up --config ./my-confs/config-1.yaml

Command migrates DB to the latest version

    project migrate down --config ./my-confs/config-1.yaml

Command rolls back a single migration from the current version

## How to run
### Docker

    docker-compose -f build/docker-compose.yml up
