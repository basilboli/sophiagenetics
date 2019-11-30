#! /bin/bash

mongoimport --host mongodb --db sg --collection clients --type json --file /mongo-seed/clients.json --jsonArray
mongoimport --host mongodb --db sg --collection clients --type json --file /mongo-seed/usages.json --jsonArray --mode merge --upsertFields salesforceId