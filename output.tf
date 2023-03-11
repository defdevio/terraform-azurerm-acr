output "resource_ids" {
  value       = azurerm_container_registry.acr[*].id
  description = "The resource ids for the Azure Container Registries."
}