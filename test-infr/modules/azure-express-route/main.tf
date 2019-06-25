data "azurerm_resource_group" "main" {
  name = var.resource_group_name
}

# --------------------------------------------------
# Express Route Circuit
# --------------------------------------------------
resource "azurerm_express_route_circuit" "main" {
  name                  = "terraform-acc-express-route"
  location              = data.azurerm_resource_group.main.location
  resource_group_name   = data.azurerm_resource_group.main.name
  service_provider_name = "Equinix"
  peering_location      = var.peering_location[data.azurerm_resource_group.main.location]
  bandwidth_in_mbps     = 100

  sku {
    tier   = "Standard"
    family = "MeteredData"
  }

  tags = {
    Environment = "production"
    Purpose     = "AcceptanceTests"
  }
}