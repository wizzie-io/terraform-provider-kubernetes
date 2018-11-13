---
layout: "kubernetes"
page_title: "Kubernetes: kubernetes_secret"
sidebar_current: "docs-kubernetes-data-source-secret"
description: |-
  The resource provides mechanisms to inject containers with sensitive information while keeping containers agnostic of Kubernetes.
---

# kubernetes_secret

The secret provides mechanisms to inject containers with sensitive information, such as passwords, while keeping containers agnostic of Kubernetes.
Secrets can be used to store sensitive information either as individual properties or coarse-grained entries like entire files or JSON blobs.
The data source is able to read a secret such as a token of a service account so it can be used e.g. to give a service account token to an external service.

~> Read more about security properties and risks involved with using Kubernetes secrets: https://kubernetes.io/docs/user-guide/secrets/#security-properties

~> **Note:** All arguments including the secret data will be stored in the raw state as plain-text. [Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
data "kubernetes_secret" "example" {
  metadata {
    name = "basic-auth"
  }
}
```

## Example Usage (Extract service-account token)

```hcl
resource "kubernetes_service_account" "example" {
  metadata {
    name      = "example"
  }
}

data "kubernetes_secret" "example" {
	metadata {
    name      = "${kubernetes_service_account.example.default_secret_name}"
	}
}

output "example_token" {
  value = "kubernetes_secret.example.data.token"
}
```

## Argument Reference

The following arguments are supported:

* `metadata` - (Required) Standard service's metadata. More info: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#metadata

## Attributes

* `data` - A map of the secret data.
* `metadata` - Standard secret's metadata. More info: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#metadata
* `type` - The secret type. Defaults to `Opaque`. More info: https://github.com/kubernetes/community/blob/master/contributors/design-proposals/auth/secrets.md#proposed-design

## Nested Blocks

### `metadata`

#### Arguments

* `annotations` - (Optional) An unstructured key value map stored with the secret that may be used to store arbitrary metadata. More info: http://kubernetes.io/docs/user-guide/annotations
* `generate_name` - (Optional) Prefix, used by the server, to generate a unique name ONLY IF the `name` field has not been provided. This value will also be combined with a unique suffix. Read more: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#idempotency
* `labels` - (Optional) Map of string keys and values that can be used to organize and categorize (scope and select) the secret. May match selectors of replication controllers and services. More info: http://kubernetes.io/docs/user-guide/labels
* `name` - (Optional) Name of the secret, must be unique. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/identifiers#names
* `namespace` - (Optional) Namespace defines the space within which name of the secret must be unique.

#### Attributes

* `generation` - A sequence number representing a specific generation of the desired state.
* `resource_version` - An opaque value that represents the internal version of this secret that can be used by clients to determine when secret has changed. Read more: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#concurrency-control-and-consistency
* `self_link` - A URL representing this secret.
* `uid` - The unique in time and space value for this secret. More info: http://kubernetes.io/docs/user-guide/identifiers#uids

## Import

Secret can be imported using its namespace and name, e.g.

```
$ terraform import kubernetes_secret.example default/my-secret
```
