# Running the broker in a Docker container

Instead of running the service broker as a CF app, as described in the main [README](README.md), you can run it elsewhere, for example in a Docker container. It's not a better approach, but it does show that you can run service brokers anywhere and still connect to and register them as service brokers in CF.

If you want to try this out, there's a [Dockerfile](Dockerfile) in this repository that you can use. In many ways it mirrors the [manifest.yml](manifest.yml) configuration for CF; in particular, it has `ENV` instructions for each of the variables such as the username and password for the broker, and the API key.

## Steps

Here are the steps to getting the service broker up and running in a Docker container, and available to your CF space.

Assumptions:

* you have a Docker engine installed locally
* you have [ngrok](https://ngrok.com/) installed so that you can make the service available beyond your local machine
* you have cloned this repository and are in the root of the clone
* you have the CF command line tool `cf` and also `curl` locally

### Build the image

First, build the Docker image:

```bash
docker build -t elephantsql-broker .
```

The output should look similar to this (redacted for brevity):

```
[+] Building 33.6s (12/12) FINISHED
 => [internal] load build definition from Dockerfile                                                 0.0s
 => => transferring dockerfile: 581B                                                                 0.0s
 => [internal] load .dockerignore                                                                    0.0s
 => => transferring context: 2B                                                                      0.0s
 => [internal] load metadata for docker.io/library/golang:1.17                                       2.4s
 => [auth] library/golang:pull token for registry-1.docker.io                                        0.0s
 => [internal] load build context                                                                    0.3s
 => => transferring context: 14.39MB                                                                 0.3s
 => [1/6] FROM docker.io/library/golang:1.17@sha256:ca709802394f4744685c1ecc083965656c3633799a005e  27.9s
 => => resolve docker.io/library/golang:1.17@sha256:ca709802394f4744685c1ecc083965656c3633799a005e9  0.0s
 => => sha256:4ff1945c672b08a1791df62afaaf8aff14d3047155365f9c3646902937f7ffe6 5.15MB / 5.15MB       0.8s
 [...]
 => => extracting sha256:debee5da592d2d2e96de757c5aff0cf0723782cf20a78682c7c85aa53064b702            0.0s
 => [2/6] WORKDIR /usr/src/app                                                                       0.1s
 => [3/6] COPY go.mod go.sum ./                                                                      0.0s
 => [4/6] RUN go mod download && go mod verify                                                       1.1s
 => [5/6] COPY . .                                                                                   0.1s
 => [6/6] RUN go build -v -o /usr/local/bin/ ./...                                                   1.6s
 => exporting to image                                                                               0.1s
 => => exporting layers                                                                              0.1s
 => => writing image sha256:359e76a8be6eed2f4bb773f009e0396aed8918c9dfb2dd880804bf1575c32c72         0.0s
 => => naming to docker.io/library/elephantsql-broker                                                0.0s
```

### Create a container from the image

Now you have the image, create a container from it, passing the API key and your choice of username and password for the broker, similar to how you would do it if pushing the broker app to CF (as described in the [Deploy service broker to Cloud Foundry](#deploy-service-broker-to-cloud-foundry) of the main README):

```bash
docker run \
  --rm \
  --detach \
  --publish 8080:8080 \
  --env BROKER_USERNAME=brokerusername \
  --env BROKER_PASSWORD=brokerpassword \
  --env BROKER_API_KEY=your-api-key \
  elephantsql-broker
```

The output should be an identifier for the container created, like this:

```
32b9a69d38609aace8e6c790be1957891ccbc727e80921473b94fa28c8e276dc
```

The options passed to the `docker run` invocation are as follows:

* `--rm` remove container when execution ends
* `--detach` run container in background, emitting the ID before returning the prompt
* `--publish 8080:8080` pass through the broker port to the container host
* `--env ...` set the values for the environment variables that the broker uses

### Make the service broker available beyond your local machine

You can check the service broker is running like this:

```bash
curl http://localhost:8080/health
```

The output should look like this:
```
{
  "status": "ok"
}
```

Use `ngrok` to make the service available beyond your local machine:

```
ngrok http 8080
```

This should result in a monitor display that looks something like this:

```
ngrok by @inconshreveable

Session Status                online
Session Expires               1 hour, 59 minutes
Version                       2.3.40
Region                        United States (us)
Web Interface                 http://127.0.0.1:4040
Forwarding                    http://0c01-86-150-217-67.ngrok.io -> http://localhost:8080
Forwarding                    https://0c01-86-150-217-67.ngrok.io -> http://localhost:8080

Connections                   ttl     opn     rt1     rt5     p50     p90
                              0       0       0.00    0.00    0.00    0.00
```

The `https` scheme based forwarding URL is the one you can use to connect from anywhere. You're now ready to register the service broker in your CF space as described in the main [README](README.md).

Here's an example, based on the information in these steps:

```bash
cf create-service-broker \
  elephantsql \
  brokerusername \
  brokerpassword \
  https://0c01-86-150-217-67.ngrok.io \
  --space-scoped
```

This should result in something similar to this:

```
Creating service broker elephantsql in org acbb5e7etrial / space dev as dj.adams@example.com...
OK
```
