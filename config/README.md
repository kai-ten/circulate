# Configure the ArangoDB instance

This container creates all of the necessary collections in the ArangoDB instance that will be used to build out 
the graph database as data is collected from the data sources.

If a collection already exists, this lambda will gracefully handle an attempt at creating a collection for a second time. 
This means that nodes and edges can be added over time as more data is gathered and as needs change.

<br />

## Requirements to run locally:

- Docker (e.g. [Rancher Desktop](https://rancherdesktop.io/))
- [ArangoDB Docker](https://hub.docker.com/_/arangodb )
- Golang (see go.mod for current version we are on)

<br />

## Running locally:

1. Start an ArangoDB instance on your machine: <br />
    ```
    docker run \
        -e ARANGO_ROOT_PASSWORD=password \
        -p 8529:8529 \
        -v /arangodb:/var/lib/arangodb3 \
        -v /arangodb/logs/:/var/log/arangodb3 \
        --restart always \
        --name arangodb \
        arangodb:latest \
        --log.level debug
    ```

1. Navigate to localhost:8529 in your browser to check out your new database
    - Username: root
    - Password: password
    - ** Never use this in production...

1. Build the Configurator container: <br />
    ```
    docker build . -t adb_configurator:latest
    ```


1. Run the docker container locally: <br />
    ```
    docker run \
        -e DB_ENDPOINT="tcp://172.17.0.2:8529" \
        -e ROOT_PASSWORD="password" \
        -e CA_CERT="" \
        -p 9000:8080 adb_configurator:latest /main
    ```


1. Trigger the lambda from another terminal: <br />
    ```
    curl -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{}'
    ```
