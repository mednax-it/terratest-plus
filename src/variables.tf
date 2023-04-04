

variable "test_string" {
  type = string
}

variable "test_map" {
  type = map(any)
}

variable "test_array" {
  type = list(any)

}

variable "tenant_id" {
  description = "Set as TF_VAR_tenant_id in scripts/setup_for_deployment.sh"
  type        = string
}
variable "client_id" {
  description = "Set as TF_VAR_client_id in scripts/setup_for_deployment.sh"
  type        = string
}
variable "client_secret" {
  description = "Set as TF_VAR_client_secret in scripts/setup_for_deployment.sh"
  type        = string
}
