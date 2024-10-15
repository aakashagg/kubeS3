# kubeS3
This project lets you create and manage S3 buckets in AWS using Kubernetes CRDs.
This can upload your pod's data to your desired S3 bucket.
This is build with tht help of [Kubebuilder](https://book.kubebuilder.io/) tool.

## Description
kubeS3 is a Kubernetes controller that allows you to create and manage S3 buckets in AWS using Kubernetes CRDs. It provides a simple way to create and manage S3 buckets in AWS using Kubernetes CRDs.

## Possible Use Cases
- We can use this to back up logs from the pods to S3.
- We can use this to store any other sensitive data from the pods to S3.

## Features
You can automate the management of S3 buckets in AWS using Kubernetes CRDs. The following are the features of kubeS3:

## Getting Started
To get started with kubeS3, you need to have the following prerequisites installed on your system:

### Prerequisites
- go version v1.22.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.
- Local cluster setup using kind or minikube.
- Access to an AWS account and the aws keys injected on your environment.

### To Deploy on the cluster
**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/kubes3:tag
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/kubes3:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Project Distribution

Following are the steps to build the installer and distribute this project to users.

1. Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=<some-registry>/kubes3:tag
```

NOTE: The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without
its dependencies.

2. Using the installer

Users can just run kubectl apply -f <URL for YAML BUNDLE> to install the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/kubes3/<tag or branch>/dist/install.yaml
```