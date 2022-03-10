# :elephant: elephantsql-broker

[![CircleCI](https://circleci.com/gh/JamesClonk/elephantsql-broker.svg?style=svg)](https://circleci.com/gh/JamesClonk/elephantsql-broker)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue)](https://github.com/JamesClonk/elephantsql-broker/blob/master/LICENSE)
[![Platform](https://img.shields.io/badge/platform-Cloud%20Foundry-lightgrey)](https://developer.swisscom.com/)

> #### PostgreSQL as a Service
> Perfectly configured and optimized PostgreSQL databases ready in 2 minutes.

**elephantsql-broker** is an [ElephantSQL](https://www.elephantsql.com/) [service broker](https://www.openservicebrokerapi.org/) for [Cloud Foundry](https://www.cloudfoundry.org/) and [Kubernetes](https://kubernetes.io/)

## Usage

#### Deploy service broker to Cloud Foundry

1. create an API key for your ElephantSQL [account](https://customer.elephantsql.com/apikeys)
2. pick a Cloud Foundry provider.
   I'd suggest the [Swisscom AppCloud](https://developer.swisscom.com/)
3. push the app, providing the API key and a username/password to secure the service broker with
4. register the service broker in your space (`--space-scoped`)
5. check `cf marketplace` to see your new available service plans

![create service broker](https://raw.githubusercontent.com/JamesClonk/elephantsql-broker/recordings/recordings/setup-min.gif "create service broker")

As an alternative to deploying the service broker to CF, you can also [run it in a Docker container](docker.md).

#### Provision new databases

1. create a new service instance (`cf cs`)
2. bind the service instance to your app (`cf bs`), or create a service key (`cf csk`)
3. inspect the service binding/key, have a look at the credentials (`cf env`/`cf sk`)
4. use the given credentials to connect to your new Postgres database
5. enjoy!

![provision service](https://raw.githubusercontent.com/JamesClonk/elephantsql-broker/recordings/recordings/provisioning-min.gif "provision service")

### Default Region

By default the service broker will provision new elephantsql database instances in the configured region `BROKER_API_DEFAULT_REGION` (see `manifest.yml`) or if none configured at all it will use `azure-arm::westeurope` as default value.
When issuing service provisioning requests to the service broker it is also possible to provide the region as an additional parameter.
###### Example:
```
$ cf create-service elephantsql hippo my-db -c '{"region": "amazon-web-services::eu-west-3"}'
```
