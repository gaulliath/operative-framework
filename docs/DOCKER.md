# operative-framework
[![Go Report Card](https://goreportcard.com/badge/github.com/graniet/operative-framework)](https://goreportcard.com/report/github.com/graniet/operative-framework) [![GoDoc](https://godoc.org/github.com/graniet/operative-framework?status.svg)](http://godoc.org/github.com/graniet/operative-framework) [![GitHub release](https://img.shields.io/github/release/graniet/operative-framework.svg)](https://github.com/graniet/operative-framework/releases/latest) [![LICENSE](https://img.shields.io/github/license/graniet/operative-framework.svg)](https://github.com/graniet/operative-framework/blob/master/LICENSE)

## Installing

### Running as a Docker Container

#### Pre-requisite
You can run operative-framework in a Docker container to avoid installing Golang locally. To install Docker check out the [official Docker documentation](https://docs.docker.com/engine/getstarted/step_one/#step-1-get-docker).

#### Pull pre-built image
```
docker pull graniet/operative-framework
```

#### Start a container
Once you have docker installed you can run operative-framework:

    $ docker run -ti --rm graniet/operative-framework

You can use Docker volumes to let gitsome access your working directory, your local .gitsomeconfig and .gitconfig:

    $ docker run -ti --rm -v $(pwd):/src/              \
       -v ${HOME}/.opf/.env:/root/.opf/.env \
       -v ${HOME}/.opf/services:/root/.opf/services          \
       graniet/operative-framework

If you are running this command often you will probably want to define an alias:

    $ alias opf="docker run -ti --rm -v $(pwd):/src/              \
                      -v ${HOME}/.opf/.env:/root/.opf/.env  \
                      -v ${HOME}/.opf/services:/root/.opf/services          \
                      graniet/operative-framework"

To build the Docker image from sources:

    $ git clone https://github.com/graniet/operative-framework.git
    $ cd operative-framework
    $ docker build -t opf . (or make docker)

### Running as a Docker-Compose container
```
docker-compose run operative-framework
```
