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
This will generate a `.env` file within the project folder with all those values to be taken by the `docker-compose` file and build/start the containers.
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

## CI/CD Plan
- CI/CD platform of choice: Jenkins

###Plan:
By using Jenkins pipelines plugins, it could be configured different stages, where there will be executed different set of tests, these stages are:
1. `Build`: as the first step, this stage will run the build of the executable and put it into a docker image, which will have as a tag the very same commit hash (since every commit could be a potential release candidate, we keep them all in `Nexus` repository for docker images)
2. `CI-Test`: once the build passes, the commit hash AKA version of the docker image, will get deployed onto a test environment, with all the dependencies, lets say an ephemeral environment run by a `docker-compose` to perform all the functional/e2e testing by using tools like `JMeter` to test all the endpoints and check DB values.
3. `Stage`: once previous step goes green, this service will be deployed in a production like environment, with a real instance and scalling solutions availables, where would be performed several functional and stress tests to identify possible bugs, caused by the integration with other real dependencies, bottle necks and performance.
4. `Production`: finally when stage instance goes well it gets deployed into production cluster.     

###Some other considerations
- APM for performance metrics like New Relic or SignalFX
- Logs post to a log provider like Kibana
- Business metrics using integration of Prometheus and Grafana
- Zero downtime by using A/B testing deployments strategy or Strangle pattern by start routing traffic gradually to the new versions gradually
- We can use other features like feature flags, Split.io could be one solution among others, to start routing traffic into the new features.

#### Assumptions 
- This will work with an UI
- UI will be in charge of validating the data they send for Country, City and State parameters, since they have to reach the server in a standard format in order to perform request with right names to the weather provider.
- Or else, they will have to send country codes and then backend translate it to the proper names (this could be another approach, but the services was built under the above premise).

##### YOU BETTER GET THE BEERS COLD FOR THE MEETUP TO BE SUCCESSFUL! USE ME TO KNOW THAT! ;) Cheers!