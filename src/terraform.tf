provider "azuread" {
}

provider "azurerm" {
  features {}

  client_id       = var.client_id
  subscription_id = "24df1984-169b-47a7-95bf-08a1d9434cb2"
  tenant_id       = var.tenant_id
  client_secret   = var.client_secret
}

resource "azurerm_resource_group" "example" {
  name     = "TEMP-DeleteMe-Testing"
  location = "eastus"
}

resource "azurerm_storage_account" "example" {
  name                     = "sdbideletemetemp"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "Testing"
    info        = "Temp-Testing-Script"
  }
}

resource "azurerm_storage_account" "exampleCount" {
  count = 1

  name                     = "sdbideletemetemp${count.index}"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "Testing"
    info        = "Temp-Testing-Script"
  }
}

resource "azurerm_storage_account" "exampleforeach" {
  for_each = tomap({
    "key" = "value"
  })

  name                     = "sdbideleteme${each.key}"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "Testing"
    info        = "Temp-Testing-Script"
  }
}
