# TWHelp API

Graphql API with TribalWars servers data.

## Development

**Required env variables to run this API** (you can set them directly on your system or create .env.development file):

```
DB_USER=your_pgdb_user
DB_NAME=your_pgdb_name
DB_PORT=your_pgdb_port
DB_HOST=your_pgdb_host
DB_PASSWORD=your_pgdb_password
REDIS_HOST=your_redis_host
REDIS_PORT=your_redis_port
LIMIT_WHITELIST=127.0.0.1,::1
LOG_DB_QUERIES=[true|false]
```

### Prerequisites

1. Golang
2. PostgreSQL database
3. Configured [cron](https://github.com/tribalwarshelp/cron)

### Installing

1. Clone this repo.
2. Navigate to the directory where you have cloned this repo.
3. Set the required env variables directly on your system or create .env.development file.
4. go run main.go
