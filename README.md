# Stockx programming challenge.

## Objective

For this challenge I was instructed to create an API that could interact with the size reviews of a shoe. I refer to these size reviews as truetosize. The API allows a client to POST truetosize values for a ShoeId and GET the average truetosize for a ShoeId.

## Postgres Database
A Postgres database is being used to store shoes data. I've written the initialization script below. The actual initialization script may be found at github.com/tjper/shoesdb

## API
The shoes API has one endpoint /shoes/truetosize
### POST /shoes/truetosize 
content-type: "application/json"
{"ShoeId":1, "TrueToSize": 1}

ShoeId is constrained by a foreign-key in the database.
TrueToSize must be >= 1 AND <= 5.

Returns status code CREATED (201) on success.

### GET /shoes/truetosize?shoeId=1
Returns status code OK (200) on success and...
{"ShoeId":1, "TrueToSizeAvg": 1.5}
