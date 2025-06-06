# kubeS3
This project lets you create and manage S3 buckets in AWS using Kubernetes CRDs.
This can upload your pod's data to your desired S3 bucket.
This is build with the help of [Kubebuilder](https://book.kubebuilder.io/) tool.

## Description
kubeS3 is a Kubernetes operator that manages Amazon S3 resources using Custom
Resource Definitions (CRDs). It allows you to declaratively create buckets and
transfer pod data to them directly from your cluster.

## How it works
The project defines two CRDs under the `storage.awsresources.com` API group:

- **S3Bucket** – describes the bucket to create, including its name, region, and
  optional access policies.
- **S3Data** – references an existing bucket and a path inside a pod whose
  contents should be uploaded. A deletion policy controls whether data is removed
  from S3 when the resource is deleted.

Controllers reconcile these resources using the AWS SDK. Finalizers ensure S3
buckets and data are cleaned up when the CRs are removed.

## Possible Use Cases
- Backing up logs from Pods to S3.
- Storing sensitive data from Pods to S3.
- Backing up pod data for disaster recovery.
- Managing application artifacts and dependencies in S3.
- Archiving old data to reduce storage costs.
- Storing configuration files and secrets securely.
- Facilitating data retrieval for analytics and reporting.

## Diagram

[Diagram here](./docs/Images/kubeS3.png)

## Getting Started
To get started with kubeS3, you need to have the following prerequisites installed on your system:
### Prerequisites
- go version v1.22.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.
- Local cluster setup using kind or minikube.
- Access to an AWS account and the aws keys injected on your environment.
