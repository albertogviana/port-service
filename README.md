# Port Domain Service

It is responsible for keeping the port information up to date. The user will run the import command and pass a ports.json file, and the system will handle the creation or update of each record.

# Automation

Every day we need to perform many tasks in a project. So I created a Makefile that automates some of those tasks. Windows, Linux, and MacOsx support the make command, making it flexible to automate some of our tasks.

## Linter

To evaluate whether the project adheres to the golang standard, we can use [Golangci](https://golangci-lint.run/). To run golangci you will need the command below:

```
make lint
```

## Tests

To execute the tests, you can run the command below:

```
make test
```

It is also possible to check the code coverage of the project by running the command below:

```
make coverage
```

## Building a docker container

To build a docker container, you just need to run the command below:

```
make docker-build
```

# Running the application

So after the project is built, you can run the application by running the command below:

```
docker-compose up -d
```

The docker-compose command will start a MySQL database and the Port Domain Service application.

The API runs on the `http://localhost:8080` address and has only two endpoints.

## Endpoints

### Create or Update Port

Endpoint: `/v1/port`

HTTP Method: `POST`

Sample Payload:

```json
{
  "name": "Ajman",
  "city": "Ajman",
  "country": "United Arab Emirates",
  "alias": [],
  "regions": [],
  "coordinates": [
    55.5136433,
    25.4052165
  ],
  "province": "Ajman",
  "timezone": "Asia/Dubai",
  "unlocs": [
    "AEAJM"
  ],
  "code": "52000"
}
```

### Get a Port

Endpoint: `/v1/port/:unloc` - where the `:unloc` value should be changed for a valid one.

HTTP Method: `GET`

Sample Response

```json
{
  "name": "Ajman",
  "city": "Ajman",
  "country": "United Arab Emirates",
  "alias": [],
  "regions": [],
  "coordinates": [
    55.5136433,
    25.4052165
  ],
  "province": "Ajman",
  "timezone": "Asia/Dubai",
  "unlocs": [
    "AEAJM"
  ],
  "code": "52000"
}
```

## Importing Data

It is easy to import a JSON file to the system.

The command below will run the application and will show the help menu.

```
docker exec -it port-service_port_1 /port-service --help
```

For the application to find and process the file, you must place it inside the `import` directory. Then, the application will read the file inside the `import` directory. But, first, it is required to pass the file's filename as a parameter. And this way, the app will start the import process.

```
docker exec -it port-service_port_1 /port-service import --file my-file.json
```
