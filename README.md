RPSLS Service
========

## Overview

This project offers an API for a simple, yet extensible, Rock-Paper-Scissors game. 
The service was written in Go (1.16) and uses a Neo4j database, as well as a Redis cache. 

## Server and dependencies configuration

The service sets up both the server and the DB connection with configuration from environment variables. Godotenv is used 
to manage the environment variables and handles different profiles / environments with preconfigured settings.

* RPSLS_ENV: the environment in which the server will run. Initially only two envs are configured: development (default) 
  and production. New environments with default settings can be easily added by creating a .env.**newenvironment** file.
* SERVER_ADDR: the port on which the server should be started.
* DB_DATABASE: the name of the database to use. Defaults to **rpsls**.
* DB_URI: the URI for the connection to the database. Includes scheme, host and port, e.g. **neo4j://localhost:7687**
* DB_USERNAME and DB_PASSWORD.
* REDIS_ADDR: the address of the Redis server, e.g. localhost:6379
* REDIS_PASSWORD
* REDIS_DB: the Redis database to be used.
* RANDOM_NUMBER_SERVER: the URL of the external random number server to be used.

## Running the server

1. Install and configure Go.
1. Access to a Neo4j server is required. If you do not wish to install a server locally, you can simply run a container 
   using the provided docker-compose file: `docker-compose up -d`
1. Redis: item 2 also applies.
1. Get all the dependencies by running the command `go get`
1. Make sure Wire has injected everything we need by running the command `go generate`.
1. Build the project: `go build -o rpsls ./cmd/rpsls/main.go`
1. Start the server: `rpsls`

## Running the server as a container

The server can also be run as a Docker container. Simply follow steps 1 and 2 from the *Running the server* section, and 
proceed with these:

1. Build the image: `docker build -t rpsls:dev .`
1. Run the image: `docker run -it --rm --name my-running-app -p 3000:3000 --env DB_URI="neo4j://host.docker.internal:7687" rpsls:dev`

After running the commands above, the server will be listening for requests on port 3000. The hostname `host.docker.internal` 
should be used on Windows or Mac hosts. For Linux hosts running older Docker versions (< 20.04), this may not work. Please 
refer to https://stackoverflow.com/a/62431165

## Database migration

The project uses go-migrate to handle database migrations. If, for some reason, it fails to populate the database with the 
initial choices, the script for creating them is located at db/migrations/000001_create_choice_nodes.up.cypher 

## Extending the game

Future versions of this project will include an endpoint to add more choices and relationships between them and the 
existing ones. Until then, it is possible to add more choices directly to the database. Here's a sample cypher query:

// all existing choices must be matched and related to the new one
MATCH (rock:Choice) WHERE rock.name = "Rock"
MATCH (paper:Choice) WHERE paper.name = "Paper"
MATCH (scissors:Choice) WHERE scissors.name = "Scissors"
MATCH (lizard:Choice) WHERE lizard.name = "Lizard"
MATCH (spock:Choice) WHERE spock.name = "Spock"
// this part creates the new choice
CREATE (k:Choice { name: "Kitten" }),
// and here we associate it with the existing ones
(k)-[:BEATS {with: "scratches"}]->(paper),
(k)-[:BEATS {with: "eats"}]->(lizard),
(k)-[:BEATS {with: "mesmerizes"}]->(spock),
(rock)-[:BEATS {with: "crushes"}]->(k),
(scissors)-[:BEATS {with: "cuts"}]->(k);

## Scoreboard

The game includes two endpoints to get the scoreboard with the 10 most recent results and to clear it:

* `GET /scoreboard`
  returns and array of the same object received when playing a round, representing the 10 most recent results, in 
  descending order (most recent first)
* `DELETE /scoreboard`
