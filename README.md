# Service Register

## Dockerfile

```
  docker build -t hareg .
```

## Create Service

```
  docker run -d -P -e SERVICE_DOMAIN="app01:/api" microservice
  docker run -d -P -e SERVICE_DOMAIN="app02:/icon" icon
```

## Run

```
   docker run -d \
      --name register \
      --net=host \
      -v ~/.docker:/tls \
      -e ETCD_HOST=http://localhost:2379 \
      -e DOCKER_HOST=https://localhost:2375 \
      -e DOCKER_CERT_PATH=/tls \
      hareg
```
