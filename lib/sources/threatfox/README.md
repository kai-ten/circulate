# ThreatFox Sync

This container retrieves IOCs and the limited Threat Intelligence from the [ThreatFox API](https://threatfox.abuse.ch/api/#recent-iocs). 
If hosting in the cloud, this lambda is expected to run on a cadence defined by a Cloudwatch Event Trigger. 

<br />

## Requirements to run locally:

- Complete the configuration in ${PROJECT_ROOT}/config, you will then be able to build this module

<br />

## Running locally:

1. Ensure you have ArangoDB running and you have configured your environment as mentioned in the Requirements

1. Build the ThreatFox docker container: <br />
    ```
    docker build . -t threatfox:latest
    ```


1. Run the docker container locally: <br />
    ```
    docker run \
        -e DB_ENDPOINT="tcp://172.17.0.2:8529" \
        -e ROOT_PASSWORD="password" \
        -e CA_CERT="" \
        -p 9001:8080 threatfox:latest /main
    ```


1. Trigger the lambda from another terminal: <br />
    ```
    curl -XPOST "http://localhost:9001/2015-03-31/functions/function/invocations" -d '{}'
    ```

<br />

## Roadmap:
- Creating an API endpoint with API GW
- Connecting API GW to React frontend