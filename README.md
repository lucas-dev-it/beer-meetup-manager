# Santander - Meetup and Beers

## Stack Components
- Meetup Manager: core service with 2 main endpoints for weather and beer packs count based on the provided forecast, and one endpoint for authentication.
- Posgres DB: db to persist meetup's and user's metadata
- Redis: as cache layer, to keep available there forecast metadata per `country-city`
- Traefik: as reverse proxy and service discovery, so instances could be spawn up and make it available through a single endpoint (in this case it is just a round robin load balancing strategy), not exposing the actual service and port to public.
- Swagger UI: just to render the API docs  

## How to run the project
There is a `Makefile` available to perform the setup for you. Just it is necessary to run a command to create the whole environment (which already includes test data for you to play around). By running:
```
$ make docker 
```
it will prompt you to fill up this values:
```
Postgres User: <your desired DB user>
Postgres Pass: <password for postgres DB>
Postgres DB Name: <whatever name you want for this DB>
WeatherBit API Key: <your weather API key>
JWT Token Signing Key: <signature value to sign the JWT tokens>
```
(in order for this project to work subscribe to the free [Weather Bit API](https://www.weatherbit.io) in order to get you own API key to use this meetup service).

After you have completed that step just run:
```
$ make run-all 
```
for having the whole stack up and running, and then you can start to use the API. 

For further details please refer to the [Swagger docs](http://localhost:8181)

To tear down everything, you can just use (bear in mind this cleans up the test data, so you have to start over if you wish to keep using these services):
```
$ make docker-down 
```

For the sake of simplify the tests, there was included the postman collection so just need to import it and start playing around it.

### Test Data
When the service is started up for the first time there is a script that runs to include meetup and user's tests data. By using the meetup's IDs to call the endpoints you will notice:
- Meetup ID -> `1`: valid location so you must receive forecast data from the provider
- Meetup ID -> `2`: invalid location, so you will receive a `406` error since the weather bit provider will respond with no content
- Meetup ID -> `3`: valid location, but meetup start date is next year, so there is no forecast available for that date, so you will get a `204` response code

Users with IDs from `1` to `3` are `ADMIN`s the others from `4` to `10` are just `USER`s

 