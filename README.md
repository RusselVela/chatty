[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=RusselVela_chatty&metric=ncloc)](https://sonarcloud.io/summary/new_code?id=RusselVela_chatty)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=RusselVela_chatty&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=RusselVela_chatty)
[![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=RusselVela_chatty&metric=duplicated_lines_density)](https://sonarcloud.io/summary/new_code?id=RusselVela_chatty)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=RusselVela_chatty&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=RusselVela_chatty)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=RusselVela_chatty&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=RusselVela_chatty)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=RusselVela_chatty&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=RusselVela_chatty)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=RusselVela_chatty&metric=bugs)](https://sonarcloud.io/summary/new_code?id=RusselVela_chatty)
# chatty
Simple real-time messages microservice 

## Features
* Signup. Create a new account using username and password
* Login. Login to the service via JWT
* Direct Messages. Send messages to another user instantly
* Channels. Join a channel to share with others
* Message queue. Messages sent to an offline account are stored and sent on the next session
* User list. See who's online
* Channels list. Look at the public available channels to join

## Setup
The service is ready to be deployed following these steps on a terminal at the root folder:
* `make build-linux` Creates a linux binary to be used by the docker image 
* `make image Builds` a new image based on the binary
* `make push TAG=X.X.X` pushes the image to the given repository with the tag defined
* `make deploy-nginx` Deploys the nginx ingress engine to expose the service to the internet. This command must be run only once. It can be omitted on subsequent deploys
* `make deploy` Deploys the service in a Kubernetes instance using the image created on the previous step 

## Config
By default, the service will expose at localhost. This value can be changed in [values.yaml](./deploy/charts/chatty/values.yaml) under `ingress.hostname`

The docker repository where the Docker image is pushed can be changed at the [Makefile](./Makefile) at `CONTAINER_REGISTRY` 

## OAS
The OAS definition is located at [api folder](./api/chatty-service-api.yaml)