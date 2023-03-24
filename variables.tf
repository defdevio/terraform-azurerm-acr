variable "resource_count" {
  type        = number
  default     = 0
  description = "The number of ACR resources to provision."
}

variable "environment" {
  type        = string
  description = "The name of the deployment environment for the resource."
}

variable "georeplications" {
  type        = map(any)
  default     = {}
  description = "A map containing the Azure `location` of the replicated registry, whether `zone_redundancy_enabled` is enabled for the replica, and any `tags` that are desired."
}

variable "location" {
  type        = string
  description = "The Azure region where the resource will be provisioned."
}

variable "name" {
  type        = string
  description = "The name to provide to the resource."
}

variable "resource_group_name" {
  type        = string
  description = "The Azure resource group where the resources will be provisioned."
}

variable "sku" {
  type        = string
  default     = "Basic"
  description = "The SKU name of the container registry. Possible values are `Basic`, `Standard` and `Premium`."
}

variable "tags" {
  type        = map(string)
  default     = {}
  description = "A map of tags to apply to the resources."
}