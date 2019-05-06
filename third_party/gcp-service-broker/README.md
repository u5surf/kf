# GCP Service Broker Installation

## Install
These instructions will allow you to install the GCP Service Broker in your cluster without Helm or Tiller.

```
kubectl apply --recursive --filename ./manifests
```

## About
This directory contains the rendered Service Catalog Helm chart to eliminate the need run Helm and Tiller in your cluster. To regenerate the chart, edit `values/catalog.yml`, then run:

```
helm template \
  --set "broker.service_account_json=_" \
  --name "gcp-service-broker" \
  --values ./values.yaml \
  --output-dir ./manifests \
    ./
```
