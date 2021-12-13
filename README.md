# productManager_proxy
A REST service to act as the GRPC proxy to ProductManager

#Product Manager Proxy

## Overview

Product Manager Proxy is a backend REST service to act as the GRPC proxy to ProductManager

The application is in GoLang and uses REST to communicate with Web and gRPC to communicate with Product Manager.
It is deployed to Cloud Run which manages containerized applications.


## Running Locally

### Prerequisites

- Access to the repository in GitHub
- Docker for Mac
- Cloud Code for GoLand

### Testing

This service uses [Grinkgo](https://onsi.github.io/ginkgo/) as the framework and [Gomega](https://onsi.github.io/gomega/) its matcher

Cheat sheet
```bash
##Run Tests
ginkgo -r

##Creating a new test suite
cd path/to/package
ginkgo bootstrap

## Creating a test spec
ginkgo generate <file_name>

```

