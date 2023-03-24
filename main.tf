resource "azurerm_container_registry" "acr" {
  count               = var.resource_count
  name                = "${var.environment}${var.location}${var.name}"
  location            = var.location
  resource_group_name = var.resource_group_name

  admin_enabled = false
  sku           = var.sku
  tags          = var.tags

  dynamic "georeplications" {
    for_each = var.georeplications
    content {
      location                = georeplications.value.location
      tags                    = try(georeplications.value.tags, null)
      zone_redundancy_enabled = try(georeplications.value.zone_redundancy_enabled, null)
    }
  }
}