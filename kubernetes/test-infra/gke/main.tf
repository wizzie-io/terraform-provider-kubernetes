provider "google" {
  // Provider settings to be provided via ENV variables
  project = "${var.gcp_project}"
  region  = "${var.gcp_region}"
}

data "google_compute_zones" "available" {}

resource "random_id" "cluster_name" {
  byte_length = 10
}

resource "random_id" "username" {
  byte_length = 14
}

resource "random_id" "password" {
  byte_length = 16
}

variable gcp_project {}
variable gcp_region {}

variable gcp_network {
  default = "default"
}

variable gcp_subnetwork {
  default = "default"
}

# See https://cloud.google.com/container-engine/supported-versions
variable "kubernetes_version" {
  description = <<EOF
The GKE Kubernetes version.
EXAMPLE:
  '1.9'
  '1.9.6-gke.1'.

See https://cloud.google.com/container-engine/supported-versions
EOF
}

resource "google_container_cluster" "primary" {
  name               = "tf-acc-test-${random_id.cluster_name.hex}"
  zone               = "${data.google_compute_zones.available.names[0]}"
  initial_node_count = 2
  node_version       = "${var.kubernetes_version}"
  min_master_version = "${var.kubernetes_version}"

  network    = "${var.gcp_network}"
  subnetwork = "${var.gcp_subnetwork}"

  additional_zones = [
    "${data.google_compute_zones.available.names[1]}",
  ]

  master_auth {
    username = "${random_id.username.hex}"
    password = "${random_id.password.hex}"
  }

  node_config {
    machine_type = "n1-standard-2"

    oauth_scopes = [
      "https://www.googleapis.com/auth/compute",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]
  }
}

output "google_zone" {
  value = "${data.google_compute_zones.available.names[0]}"
}

output "endpoint" {
  value = "${google_container_cluster.primary.endpoint}"
}

output "username" {
  value = "${google_container_cluster.primary.master_auth.0.username}"
}

output "password" {
  value = "${google_container_cluster.primary.master_auth.0.password}"
}

output "client_certificate_b64" {
  value = "${google_container_cluster.primary.master_auth.0.client_certificate}"
}

output "client_key_b64" {
  value = "${google_container_cluster.primary.master_auth.0.client_key}"
}

output "cluster_ca_certificate_b64" {
  value = "${google_container_cluster.primary.master_auth.0.cluster_ca_certificate}"
}

output "node_version" {
  value = "${google_container_cluster.primary.node_version}"
}

data template_file kube_config {
  vars {
    endpoint    = "${google_container_cluster.primary.endpoint}"
    certificate = "${google_container_cluster.primary.master_auth.0.cluster_ca_certificate}"
    password    = "${google_container_cluster.primary.master_auth.0.password}"
    username    = "${google_container_cluster.primary.master_auth.0.username}"
    client_cert = "${google_container_cluster.primary.master_auth.0.client_certificate}"
    client_key  = "${google_container_cluster.primary.master_auth.0.client_key}"
  }

  template = <<EOF
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: $${certificate}
    server: https://$${endpoint}
  name: acctest
contexts:
- context:
    cluster: acctest
    user: acctest
  name: acctest
current-context: acctest
kind: Config
preferences: 
  colors: true
users:
- name: acctest
  user:
    password: $${password}
    username: $${username}
    client-certificate-data: $${client_cert}
    client-key-data: $${client_key}
EOF
}

resource local_file kube_ca {
  content  = "${base64decode(google_container_cluster.primary.master_auth.0.cluster_ca_certificate)}"
  filename = "kube.ca"
}

resource local_file kube_cert {
  content  = "${base64decode(google_container_cluster.primary.master_auth.0.client_certificate)}"
  filename = "client.cert"
}

resource local_file kube_key {
  content  = "${base64decode(google_container_cluster.primary.master_auth.0.client_key)}"
  filename = "client.key"
}

resource local_file kube_config {
  content  = "${data.template_file.kube_config.rendered}"
  filename = ".kube_config"
}
