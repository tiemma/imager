# Imager

Get list of images present with repositories

## Use cases
 - Auditing images
 - Verifying tags are correct 
 - Remote container scans

## How to use

### CLI
The `imager` cli provides a command line utility to test on.

> NOTE: You would need to verify access to the repositories, it would skip erroneous instances.

```shell script
go run cmd/main.go --input <comma separated list of repositories>
```

### Docker Image
You can likewise run the script by using the docker image

```shell script
docker run imager /imager --input <comma separated list of repositories>
```

### Kubernetes
You can also run the script as a job to extract the information 

```shell script
kubectl apply -f imager.yaml
```
