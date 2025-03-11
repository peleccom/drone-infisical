A secret extension to Drone CI secrets client for infisical. _Please note this project requires Drone server version 1.4 or higher._

## Installation

Create a shared secret:

```console
$ openssl rand -hex 16
bea26a2221fd8090ea38720fc445eca6
```

Download and run the plugin:

```console
$ docker run -d \
  --publish=3000:3000 \
  --env=DRONE_INFISICAL_DEBUG=true \
  --env=DRONE_INFISICAL_DRONE_SECRET=bea26a2221fd8090ea38720fc445eca6 \
  --env DRONE_INFISICAL_URL=https://infistical.example.com \
  --env DRONE_INFISICAL_CLIENT_ID=8253cb18-28d5-4bb0-b674-5f4f2964d076 \
  --env DRONE_INFISICAL_CLIENT_SECRET=0d36de35f1e3b1438c316c137ca25bfc215db44bd7b5958b21f3fbfc51bb3d9a \
  --env DRONE_INFISICAL_PROJECT_ID=2ea3211d-0fc8-4db5-a302-1652e24a5b83 \
  --restart=always \
  --name=secrets foo/bar
```

Update your runner configuration to include the plugin address and the shared secret.

```text
DRONE_SECRET_PLUGIN_ENDPOINT=http://1.2.3.4:3000
DRONE_SECRET_PLUGIN_TOKEN=bea26a2221fd8090ea38720fc445eca6
