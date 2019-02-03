# go-martini-test-app
Simple Golang (Martini framework) based test application:
- create/update/delete post
- output posts list
- [mongo-go driver](https://github.com/mongodb/mongo-go-driver)
- docker
- [mongo-express](https://github.com/mongo-express/mongo-express) (adminer for mongo)

## Install Docker 

https://docs.docker.com/install/

## Install docker-compose 

https://docs.docker.com/compose/install/

## Install docker-hostmanager

https://github.com/iamluc/docker-hostmanager

Run manager

```
$ docker run -d --name docker-hostmanager --restart=always -v /var/run/docker.sock:/var/run/docker.sock -v /etc/hosts:/hosts iamluc/docker-hostmanager
```

## Env

create .env file from .env.dist and set correct env var values

## Start with Docker

```
$ cd /project/path
$ docker compose up -d --build
```


## Application URLs:

- Mongo Adminer: http://localhost:8081/
- Application: http://localhost:3000/

