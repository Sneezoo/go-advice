# Go Advice

Small advice microservice with REST-API

## Setup

Install [docker](https://docs.docker.com/install/) and [docker-compose](https://docs.docker.com/compose/install/).
Then you are good to *go*! Just run `docker-compose up --build` and all Docker images will be built and started.

## Usage

If you don't change docker-compose.yml, you can reach the service at *:8080*.
There is one endpoint: http://localhost:8080/advice, which you can reach with *POST* or *GET*

### POST
receives an Object like this in the body:
```json
{
	"Advice": "Berries are sweet",
	"Keywords": ["Berries", "Sweet"],
	"Funny": 0,
	"Serious": 0
}
```

### GET
can receive an optional parameter *term*
If *term* is found in *keywords* or in the *advice* itself, than you get the most *serious* Advice.
Otherwise you get a random one.

Have fun!