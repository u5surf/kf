# Installing Kf

## Pre-requisites

This guide is intended to provide you with all the commands you'll
need to install `kf` in a single place. It assumes you have a the 
ability to run root containers in a cluster with at least 12 vCPUs 
and 45G of memory and a minumum of three nodes.

You will also need a Docker-compatible image registry. We use gcr.io in this doc. 

## Configure your local environment

1. Export your GCP project ID (and its corresponding gcr.io URI):

  ```
  export GCP_PROJECT=<REPLACE>
  export KF_REGISTRY=gcr.io/$GCP_PROJECT
  ```

1. Grant your account `cluster-admin` (necessary to install Service Catalog in a later step):

  ```
  kubectl create clusterrolebinding cluster-admin-binding \
      --clusterrole=cluster-admin \
      --user=$(gcloud config get-value core/account)
  ```

## Install required services and components

> The following instructions assume you have cloned this repository and are at its root.

1. Install Istio in your cluster:

  ```
  kubectl apply --filename https://github.com/knative/serving/releases/download/v0.5.0/istio-crds.yaml
  kubectl apply --filename https://github.com/knative/serving/releases/download/v0.5.0/istio.yaml
  kubectl label namespace default istio-injection=enabled
  ```

2. Install Knative Serve and Build in your cluster:

  ```
  kubectl apply --selector knative.dev/crd-install=true \
    --filename https://github.com/knative/serving/releases/download/v0.5.0/serving.yaml \
    --filename https://github.com/knative/build/releases/download/v0.5.0/build.yaml \
    --filename https://github.com/knative/eventing/releases/download/v0.5.0/release.yaml \
    --filename https://github.com/knative/eventing-sources/releases/download/v0.5.0/eventing-sources.yaml \
    --filename https://github.com/knative/serving/releases/download/v0.5.0/monitoring.yaml \
    --filename https://raw.githubusercontent.com/knative/serving/v0.5.0/third_party/config/build/clusterrole.yaml 

  kubectl apply --filename https://github.com/knative/serving/releases/download/v0.5.2/serving.yaml \
    --filename https://github.com/knative/build/releases/download/v0.5.0/build.yaml \
    --filename https://github.com/knative/eventing/releases/download/v0.5.0/release.yaml \
    --filename https://github.com/knative/eventing-sources/releases/download/v0.5.0/eventing-sources.yaml \
    --filename https://github.com/knative/serving/releases/download/v0.5.2/monitoring.yaml \
    --filename https://raw.githubusercontent.com/knative/serving/v0.5.2/third_party/config/build/clusterrole.yaml

  ```

3. Install Service Catalog in your cluster:

  ```
    kubectl apply -R -f third_party/service-catalog/manifests/catalog/templates
  ```

4. Confirm Service Catalog is visible in the marketplace:

  ```
  kf marketplace
  ```

## Upload buildpacks
Buildpacks are provided by the operator and can be uploaded using `kf`.
Sample buidpacks are included in this repo and can be installed with:

```
cd buildpack-samples
kf upload-buildpacks --container-registry $KF_REGISTRY
cd -
```

## Push your first app
At this point you are ready to deploy your first app using `kf`. Run the following command 
to push your first app. 

```
kf push helloworld --container-registry $KF_REGISTRY
```

## Install a service broker
Once you have the service catalog you'll want to install a service
broker. You can use helm to install the gcp-service-broker from
the third_party directory. 

Configure GCP service account & APIs
```
gcloud iam service-accounts create \
    gcp-service-broker

gcloud iam service-accounts keys \
    create /tmp/key.json --iam-account \
    gcp-service-broker@$GCP_PROJECT.iam.gserviceaccount.com

gcloud projects \
    add-iam-policy-binding \
    $GCP_PROJECT --member \
    serviceAccount:gcp-service-broker@$GCP_PROJECT.iam.gserviceaccount.com \
    --role "roles/owner"

gcloud services enable \
    cloudresourcemanager.googleapis.com \
    iam.googleapis.com
```

Once you have your key.json, copy this in the values.yaml file
in `/third_party/gcp-service-broker` file. 
