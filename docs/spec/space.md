# Space

Cloud Foundry Spaces are mapped to Kf Spaces. In Kf, the Space is the only unit
of configuration and organization. Each space maintains its own configuration,
quota, and roles to give developers.

Because of this consolidation, clusters can be shared more effectively between
teams with varying requirements and can even host multiple environments.

## Actions

This table gives a comparison between the `cf` and `kf` command line tools for
various actions you may want to complete.

| Action | `cf` Command | `kf` Command |
|--------|-------------|-----------|
| **Basic** | | |
| List spaces | `cf spaces` | `kf spaces` |
| Show space info | `cf space SPACE` | `kf space SPACE` |
| Create a space | `cf create-space SPACE` | `kf create-space SPACE` |
| Delete a space | `cf delete-space SPACE` | `kf delete-space SPACE` |
| Rename a space | `cf rename-space` | Not supported by Kubernetes |
| **SSH** | | |
| Allow SSH access for the space | `cf allow-space-ssh` | Not yet supported |
| Disallow SSH access for the space | `cf disallow-space-ssh` | Not yet supported |
| Report whether SSH is allowed in a space | `cf space-ssh-allowed` | Not yet supported |
| **Environment Variables** | | |
| Retrieve the running environment variable group | `cf running-environment-variable-group` | `kf space SPACE` |
| Retrieve the staging environment variable group | `cf staging-environment-variable-group` | `kf space SPACE` |
| Update the Running environment variable group | `cf set-running-environment-variable-group` | `kf configure-space (un)set-env SPACE` |
| Update the staging environment variable group | `cf set-staging-environment-variable-group` | `kf configure-space (un)set-buildpack-env SPACE` |
| **Feature Flags** | | |
| Retrieve list of feature flags with status | `cf feature-flags` | `kf space SPACE` |
| Retrieve an individual feature flag with status | `cf feature-flag FLAG` | Not supported |
| Allow use of a feature | `cf enable-feature-flag FLAG` | `kf configure-space ...` |
| Prevent use of a feature | `cf disable-feature-flag FLAG` | `kf configure-space ...` |
| **Quotas** | | |
| List resource quotas | `cf (space-)quotas` | Not supported by Kubernetes |
| Show quota info | `cf (space-)quota` | `kf space SPACE` |
| Show quota info | `cf (space-)quota` | `kf space SPACE` |
| Define a new space resource quota | `create-space-quota` | `kf configure-space quota` |
| Update an existing space quota | `update-(space-)quota` | `kf configure-space quota` |
| Delete a quota | `delete-(space-)quota` | `kf configure-space quota` |
| **Domains** | | |
| Share a private domain | `cf share-private-domain` | Not yet supported |
| Unshare a private domain | `cf unshare-private-domain` | Not yet supported |

## Implementation

Spaces are implemented by the `kf.dev` spaces CRD. Each space is responsible for
creating a Kubernetes namespace, RBAC roles, and quotas.

### CRD Structure

An annotated version is below:

```.yaml
apiVersion: kf.dev/v1alpha1
kind: Space
metadata:
  labels:
    # The managed-by label is automatically added by kf.
    app.kubernetes.io/managed-by: kf
  # The name of this space.
  name: space-name
spec:
  # The buildpackBuild section holds configuration for application builds
  # with buildpacks.
  buildpackBuild:
    # The image to use as the builder, this is the default and contains
    # buildpacks for Kf's natively supported langauges.
    builderImage: gcr.io/kf-releases/buildpack-builder:latest

    # The container registry that built apps will be stored in.
    containerRegistry: gcr.io/kf-source

    # Environment variables injected into the build.
    env:
    - name: JAVA_VERSION
      value: "11"
  # The execution section contains settings for running applications.
  execution:
    # Environment variables that will be injected into all running apps.
    env:
    - name: ENVIRONMENT
      value: production
  resourceLimits: {}

  # Security contains configuration properties to set the space's RBAC roles.
  security:
    # should the developer RBAC role get access to Kubernetes logs?
    # this is required for
    enableDeveloperLogsAccess: true
```
