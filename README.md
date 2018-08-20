tockx programming challenge.

## Objective

For this challenge I was instructed to create an API that could interact with the size reviews of a shoe. I refer to these size reviews as truetosize. The API allows a client to POST truetosize values for a ShoeId and GET the average truetosize for a ShoeId.

## To Run
#### create user-defined bridge network

    sudo docker network create --driver bridge shoes-net

#### run api and db docker container

    sudo docker run -dit -p 8080:8080 --name shoes --network shoes-net penutty/shoes:latest
    sudo docker run -dit --name shoesdb --network shoes-net penutty/shoesdb:latest

## Postgres Database
A Postgres database is being used to store shoes data. I've written the initialization script below. The actual initialization script may be found at [github.com/tjper/shoesdb].(https://github.com/tjper/shoesdb.git)

``` SQL
CREATE USER shoes;
CREATE DATABASE shoes;
\connect shoes

-- CREATE TABLES
CREATE TABLE shoes (
        id SERIAL PRIMARY KEY,
        name VARCHAR (64) NOT NULL
);

CREATE TABLE truetosize (
        shoes_id integer NOT NULL,
        truetosize smallint NOT NULL
);

CREATE INDEX shoes_id_idx ON truetosize (shoes_id);

-- CREATE FOREIGN KEYS
ALTER TABLE truetosize
        ADD CONSTRAINT truetosize_shoes_fk
        FOREIGN KEY(shoes_id)
        REFERENCES shoes (id) MATCH FULL;

-- GRANT USER PRIVILEGES
GRANT SELECT, INSERT ON TABLE shoes, truetosize TO shoes;

-- Populate shoes table
INSERT INTO shoes (name) VALUES
        ('Jordan 4 Retro Raptors'),
        ('adidas NMD Hu Pharrell Solar Pack Red'),
        ('Vans Authetic Slim Red'),
        ('adidas Yeezy Boost 350 V2 Butter'),
        ('LeBron 6 Bred'),
        ('Air Force 1 Low Travis Scott Sail'),
        ('ADIDAS YEEZY POWERPHASE')
```
This database holds 7 shoes and have been inserted in the script above. Their corresponding Ids are 1, 2, 3, 4, 5, 6,  and 7. 

## API /shoes/truetosize
### POST /shoes/truetosize 
**Content-Type:** "application/json"
**Data Params:** 
 - ShoeId int 
 - TrueToSize int [1, 5]
 
**Success Response**
 - Code: 201
 
**Error Response**
- Code: 400

**Sample Call**
`curl -d '{"ShoeId":1, "TrueToSize":2}' -H "Content-Type: application/json" -X POST http://localhost:8080/shoes/truetosize`

### GET /shoes/truetosize
**URL Params**
- shoeId int

**Success Response**
- Code: 200

**Error Response**
- Code: 400

**Sample Call**

    curl http://localhost:8080/shoes/truetosize?shoeId=1

## Thoughts
### Go
The Go API logic is fairly simple for this challenge, but I re-structured my project a couple times trying to find a good combination of simplicity and the ability to expand in the future. I know this is a just a challenge, but I like to think of application design this way. I Design with future changes in mind. I also invested a fair amount of time researching different approaches to testing. My unit tests utilize interfaces and table driven testing.

### Postgres
I normalized the small database I created. Rather than having shoe names repeated in each row within the truetosize table, I created the shoes table which pairs an id and name of a shoe.

### Docker 
I would say most of my time during this challenge was invested into Docker. While I've been familiar with docker for some time, I wasn't aware of the methods to reduce image size for Go programs. I was able to get my API down to 4.88MB and I know have a much better grasp of Docker's overall feature set utility.


