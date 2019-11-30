# SophiA Genetics Home Assignment 

see assets/task.pdf


## Demo version 

Please take a look at deployed [demo](http://206.189.114.13/clients).

## How to build / run

Everything is packaged with docker / docker-compose.

To build project from scratch: 

`make build`

To get running solution locally please run the following command: 

`docker-compose up` or `make run`

Note: to build and run solution solution you should have [docker](https://docs.docker.com/install/) and [docker compose](https://docs.docker.com/compose/install/) installed.

## Architecture explanation

Projects consists of the following modules:

- `backend`

API is implemented with [Golang](https://go.dev/). 

```endpoints
/
index

---/hc
health check

---/version
version 

---/clients 
clients data
```

API module contains integrations tests which are run before every build. 

- `frontend`

UI is implemented using reactjs, d3.

No ready-to-go frameworks for charting are used to stay simple and due to time limitations.

For simplicity current implementation is using directly d3 apis to graph two histograms for actual and predicted usage.

As far as data are heterogeneous we need to programatically scale the y axis to have a consistent view independently of the clients data.

Possible improvents: 

- Currently due to time limitations we load data all at once. This can be improved by lazy loading only data we need for visualization. 

- Have more relevant and complex visualizations to easily compare the differences between actual and predicted usage. 

- User experience can be improved consistently. Current interface is pretty basic due to time limitations. 

Some ideas:

 1. stacked barcharts [example](https://river.datawrapper.de/_/bR6ZS)

 2. calendar heatmap [example](https://reaviz.io/?path=/docs/demos-heatmap-calendar--year-calendar)
 
 3. show percentage increase / decrease instead of charts [example](https://www.pinterest.fr/pin/309552174372406889)

- `database`
For this given task as database MongoDB was chosen for its flexible schema approach. 

To import data into mongodb we run dedicated service with docker-compose (see `mongo-seed`)
which runs `mongo-seed/import.sh`

To go faster I have manually extracted two collections from input data (sophia-test.json)

`cat sophia-test.json | jq '.domainY' > clients.json`

`cat sophia-test.json | jq '.data' > usages.json`

The above mentioned script `import.sh` is using cool flag of [mongoimport](https://docs.mongodb.com/manual/reference/program/mongoimport/) [--upsertFields](https://docs.mongodb.com/manual/reference/program/mongoimport/#cmdoption-mongoimport-upsertfields) which lets you merge two collections by providing merge fields. 

Finally for simplicity reason in database we will have only one collection `clients` where each given record have the following format:

```
{
  "salesforceId": 22699,
  "country": "Morocco",
  "owner": "Paul Piovene",
  "manager": "Aristide Mullane",
  "predictedUsage": [
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    1,
    0,
    0
  ],
  "actualUsage": [
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    0,
    3,
    0,
    0
  ]
}
```
Given record is used directly by page comparing predicted / real usage for the client [example](http://206.189.114.13/g/22699).

- `assets` 

task definition, test data 
