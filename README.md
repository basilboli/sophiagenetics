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
returns index page

---/hc
returns sample health check page

---/version
returns page containing current running application commit hash and build time

---/clients 
returns page containing clients data
```

API module contains integrations tests which are run before every build. 

- `frontend`

UI is implemented using reactjs, d3.

To [keep it simple simple](https://fr.wikipedia.org/wiki/Principe_KISS) and due to time limitations no ready-to-go framework for charting was used.

Current implementation is using directly d3 apis to graph two histograms for actual and predicted usage.

x axis corresponds to time period: 12 months numerated from 1 to 12. 

y axis corresponds to usage data.

As far as usage data are heterogeneous we need to scale the y axis to have a consistent view independently of the clients data.

Possible improvents: 

- Currently due to time limitations we load data all at once. This can be improved by lazy loading only the data we need for visualization. 

- User experience can be improved consistently. Current interface is pretty basic due to time limitations. 

- Periods where actual usage differs consistenly from predicted one we can put markers to easily identify the differences

- Have more relevant and complex visualizations to easily compare the differences between actual and predicted usage. 

Instead of showing histogram we can use different visualizations:

 1. stacked barcharts [example](https://river.datawrapper.de/_/bR6ZS)

 2. calendar heatmap [example](https://reaviz.io/?path=/docs/demos-heatmap-calendar--year-calendar)
 
 3. show percentage increase / decrease instead of charts [example](https://www.pinterest.fr/pin/309552174372406889)

- `database`

For this given task MongoDB database was chosen for its flexible schema approach. 

To import data into mongodb we run dedicated service with docker-compose (see `mongo-seed`)
which runs `mongo-seed/import.sh` script.

To go faster we have manually extracted two collections from input data `sophia-test.json`: 

`cat sophia-test.json | jq '.domainY' > clients.json`

`cat sophia-test.json | jq '.data' > usages.json`

The above mentioned script `import.sh` is using [mongoimport](https://docs.mongodb.com/manual/reference/program/mongoimport/) flag [--upsertFields](https://docs.mongodb.com/manual/reference/program/mongoimport/#cmdoption-mongoimport-upsertfields) which lets you merge two collections by providing merge fields. 

After import in database we will have only one collection `clients` where each given record will contain client info as well predicted and actual usage data :

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
