provider "azuread" {
}

provider "azurerm" {
  features {}

  client_id       = var.client_id
  subscription_id = var.subscription_id # TODO: replace it just with the compliance_subscription?
  tenant_id       = var.tenant_id
  client_secret   = var.client_secret
}
