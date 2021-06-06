# tribalwarshelp.com API

A GraphQL API designed for developers who want to create something meaningful for the game [Tribal Wars](https://tribalwars.net).

## Limits

It is possible to fetch in one GraphQL query:

1. 1000 daily player/tribe stats records
2. 200 ennoblements
3. 30 versions
4. 200 players/tribes
5. 1000 villages
6. 100 player/tribe history records
7. 100 servers
8. 60 server stats records
9. 100 tribe changes

## Sample queries

1. Fetch all bonus villages with 10% more population

```graphql
query {
    villages(server: "en115", filter: { playerID: [0], bonus: 4 }) {
        total
        items {
            id
            name
            x
            y
            points
        }
    }
}

```

2. Top 30 players without a tribe, ordered by points.

```graphql
query {
    players(
        server: "pl148"
        filter: { tribeID: [0] }
        sort: ["points DESC"]
        limit: 30
    ) {
        total
        items {
            id
            name
            rank
            points
            totalVillages
            rankAtt
            rankDef
            rankTotal
            rankSup
        }
    }
}

```

3. Search a player by a nickname fragment.

```graphql
query {
    players(server: "pl148", filter: { nameIEQ: "%pablo%" }) {
        total
        items {
            id
            name
        }
    }
}
```

More examples [here](https://github.com/tribalwarshelp/scripts).

## Map service

You can generate a server map with this API. The current endpoint is http(s)://youraddress/map/server (replace "server"
with the server you're interested in, for example, pl151).

### Available query params:

| Param                                            | Default                                                    |
| ------------------------------------------------ | ---------------------------------------------------------- |
| showBarbarians                                   | false                                                      |
| largerMarkers                                    | false                                                      |
| markersOnly                                      | false                                                      |
| centerX                                          | 500                                                        |
| centerY                                          | 500                                                        |
| scale                                            | 1 (max 5)                                                  |
| showGrid                                         | false                                                      |
| showContinentNumbers                             | false                                                      |
| backgroundColor                                  | #000                                                       |
| gridLineColor                                    | #fff                                                       |
| continentNumberColor                             | #fff                                                       |
| tribe(this param you can define multiple times)  | format: tribeid,hexcolor (for example, tribe=631,#0000ff)   |
| player(this param you can define multiple times) | format: playerid,hexcolor (for example, player=631,#0000ff) |

### Example

**pl151**
```
https://api.tribalwarshelp.com/map/pl151?showBarbarian=true&tribe=124,%230000ff&tribe=631,%230000ff&tribe=1675,%230000ff&onlyMarkers=false&scale=1&showGrid=true&showContinentNumbers=true
```
![Map](https://api.tribalwarshelp.com/map/pl151?showBarbarian=true&tribe=124,%230000ff&tribe=631,%230000ff&tribe=1675,%230000ff&onlyMarkers=false&scale=1&showGrid=true&showContinentNumbers=true)

## Development

### Prerequisites

1. Golang
2. PostgreSQL database
3. Configured [cron](https://github.com/tribalwarshelp/cron)

### Installation

**Required ENV variables:**

```
DB_USER=your_pgdb_user
DB_NAME=your_pgdb_name
DB_PORT=your_pgdb_port
DB_HOST=your_pgdb_host
DB_PASSWORD=your_pgdb_password
LIMIT_WHITELIST=127.0.0.1,::1
LOG_DB_QUERIES=true
DISABLE_ACCESS_LOG=false
```

1. Clone this repo.

```
git clone git@github.com:tribalwarshelp/api.git
```

2. Navigate to the directory where you have cloned this repo.
3. Set the required env variables directly in your system or create .env.local file.
4. Run the app.

```
go run main.go
```

## License

Distributed under the MIT License. See ``LICENSE`` for more information.

## Contact

Dawid Wysokiński - [contact@dwysokinski.me](mailto:contact@dwysokinski.me)
