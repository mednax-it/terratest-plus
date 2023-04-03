provider "azuread" {
}

provider "azurerm" {
  features {}

  client_id       = var.client_id
  subscription_id = "24df1984-169b-47a7-95bf-08a1d9434cb2"
  tenant_id       = var.tenant_id
  client_secret   = var.client_secret
}
