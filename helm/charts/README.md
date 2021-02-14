# zeebo

Installs [zeebo](https://github.com/neverping/zeebo-app) using an Ingress controller and two services that uses [gCRP](https://grpc.io/) protocol.

## TL;DR;

```bash
$ helm install zeebo
```

## Introduction

This chart installs the [zeebo](https://github.com/neverping/zeebo-app) app which contains an Ingress controller on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager. It also creates an Golang app that receives the connection from the ingress controller and communicates with an internal app in Python.

## Prerequisites
  - Kubernetes 1.19.0 or above.
  - An ingress controller installed. Any supported version by the Kubernetes cluster you are running.
  - Helm 3.3.x or above.

## Installing the Chart

To install the chart with the release name called `my-release`:

```bash
$ helm install my-release --wait zeebo
```

Optionally you can select any values files from the `targets` directory. For example, for running on Minikube enviroment, you can run the command below:

```bash
$ helm install my-release --wait -f ../targets/minikube.yaml zeebo
```

## Uninstalling the Chart

To uninstall/delete the `my-release` chart deployed in your cluster:

```bash
$ helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Minikube notes

Please be aware to install the ingress plugin using this command below, otherwise the IP won't be exposed.

```bash
$ minikube addons enable ingress
```
