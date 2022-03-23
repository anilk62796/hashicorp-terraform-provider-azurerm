package firewall_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type FirewallResource struct{}

func TestAccFirewall_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall", "test")
	r := FirewallResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_configuration.0.name").HasValue("configuration"),
				check.That(data.ResourceName).Key("ip_configuration.0.private_ip_address").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewall_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall", "test")
	r := FirewallResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewall_enableDNS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall", "test")
	r := FirewallResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.enableDNS(data, "1.1.1.1", "8.8.8.8"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.enableDNS(data, "1.1.1.1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewall_withManagementIp(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall", "test")
	r := FirewallResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withManagementIp(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_configuration.0.name").HasValue("configuration"),
				check.That(data.ResourceName).Key("ip_configuration.0.private_ip_address").Exists(),
				check.That(data.ResourceName).Key("management_ip_configuration.0.name").HasValue("management_configuration"),
				check.That(data.ResourceName).Key("management_ip_configuration.0.public_ip_address_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewall_withMultiplePublicIPs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall", "test")
	r := FirewallResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multiplePublicIps(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_configuration.0.name").HasValue("configuration"),
				check.That(data.ResourceName).Key("ip_configuration.0.private_ip_address").Exists(),
				check.That(data.ResourceName).Key("ip_configuration.1.name").HasValue("configuration_2"),
				check.That(data.ResourceName).Key("ip_configuration.1.public_ip_address_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewall_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall", "test")
	r := FirewallResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_firewall"),
		},
	})
}

func TestAccFirewall_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall", "test")
	r := FirewallResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("Production"),
				check.That(data.ResourceName).Key("tags.cost_center").HasValue("MSFT"),
			),
		},
		{
			Config: r.withUpdatedTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("staging"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewall_withZones(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall", "test")
	r := FirewallResource{}
	zones := []string{"1"}
	zonesUpdate := []string{"1", "3"}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withZones(data, zones),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zones.#").HasValue("1"),
				check.That(data.ResourceName).Key("zones.0").HasValue("1"),
			),
		},
		{
			Config: r.withZones(data, zonesUpdate),
			Check: acceptance.ComposeTestCheckFunc(

				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zones.#").HasValue("2"),
				check.That(data.ResourceName).Key("zones.0").HasValue("1"),
				check.That(data.ResourceName).Key("zones.1").HasValue("3"),
			),
		},
	})
}

func TestAccFirewall_withoutZone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall", "test")
	r := FirewallResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withoutZone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewall_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall", "test")
	r := FirewallResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccFirewall_withFirewallPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall", "test")
	r := FirewallResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withFirewallPolicy(data, "pol-01"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withFirewallPolicy(data, "pol-02"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewall_inVirtualHub(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall", "test")
	r := FirewallResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.inVirtualHub(data, 1),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("virtual_hub.0.public_ip_addresses.#").HasValue("1"),
				check.That(data.ResourceName).Key("virtual_hub.0.private_ip_address").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.inVirtualHub(data, 2),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("virtual_hub.0.public_ip_addresses.#").HasValue("2"),
				check.That(data.ResourceName).Key("virtual_hub.0.private_ip_address").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.inVirtualHub(data, 1),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("virtual_hub.0.public_ip_addresses.#").HasValue("1"),
				check.That(data.ResourceName).Key("virtual_hub.0.private_ip_address").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFirewall_privateRanges(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_firewall", "test")
	r := FirewallResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_configuration.0.name").HasValue("configuration"),
				check.That(data.ResourceName).Key("ip_configuration.0.private_ip_address").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.privateRanges(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_configuration.0.name").HasValue("configuration"),
				check.That(data.ResourceName).Key("ip_configuration.0.private_ip_address").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_configuration.0.name").HasValue("configuration"),
				check.That(data.ResourceName).Key("ip_configuration.0.private_ip_address").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func (FirewallResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FirewallID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Firewall.AzureFirewallsClient.Get(ctx, id.ResourceGroup, id.AzureFirewallName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Azure Firewall %s : %v", *id, err)
	}

	return utils.Bool(resp.AzureFirewallPropertiesFormat != nil), nil
}

func (FirewallResource) Destroy(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FirewallID(state.ID)
	if err != nil {
		return nil, err
	}

	future, err := clients.Firewall.AzureFirewallsClient.Delete(ctx, id.ResourceGroup, id.AzureFirewallName)
	if err != nil {
		return nil, fmt.Errorf("deleting Azure Firewall %q: %+v", id.AzureFirewallName, err)
	}

	if err = future.WaitForCompletionRef(ctx, clients.Firewall.AzureFirewallsClient.Client); err != nil {
		return nil, fmt.Errorf("waiting for Deletion on azureFirewallsClient: %+v", err)
	}

	return utils.Bool(true), nil
}

func (FirewallResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fw-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_firewall" "test" {
  name                = "acctestfirewall%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "AZFW_VNet"
  sku_tier            = "Standard"

  ip_configuration {
    name                 = "configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }
  threat_intel_mode = "Deny"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (FirewallResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fw-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_firewall" "test" {
  name                = "acctestfirewall%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "AZFW_VNet"
  sku_tier            = "Standard"

  ip_configuration {
    name                 = "configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }
  threat_intel_mode = "Deny"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (FirewallResource) enableDNS(data acceptance.TestData, dnsServers ...string) string {
	servers := make([]string, len(dnsServers))
	for idx, server := range dnsServers {
		servers[idx] = fmt.Sprintf(`"%s"`, server)
	}

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_firewall" "test" {
  name                = "acctestfirewall%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "AZFW_VNet"
  sku_tier            = "Standard"

  ip_configuration {
    name                 = "configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }
  threat_intel_mode = "Deny"
  dns_servers       = [%s]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, strings.Join(servers, ","))
}

func (FirewallResource) withManagementIp(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fw-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_subnet" "test_mgmt" {
  name                 = "AzureFirewallManagementSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "test_mgmt" {
  name                = "acctestmgmtpip%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_firewall" "test" {
  name                = "acctestfirewall%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "AZFW_VNet"
  sku_tier            = "Standard"

  ip_configuration {
    name                 = "configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }

  management_ip_configuration {
    name                 = "management_configuration"
    subnet_id            = azurerm_subnet.test_mgmt.id
    public_ip_address_id = azurerm_public_ip.test_mgmt.id
  }

  threat_intel_mode = "Alert"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (FirewallResource) multiplePublicIps(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fw-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_public_ip" "test_2" {
  name                = "acctestpip2%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_firewall" "test" {
  name                = "acctestfirewall%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "AZFW_VNet"
  sku_tier            = "Standard"

  ip_configuration {
    name                 = "configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }

  ip_configuration {
    name                 = "configuration_2"
    public_ip_address_id = azurerm_public_ip.test_2.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (FirewallResource) requiresImport(data acceptance.TestData) string {
	template := FirewallResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_firewall" "import" {
  name                = azurerm_firewall.test.name
  location            = azurerm_firewall.test.location
  resource_group_name = azurerm_firewall.test.resource_group_name
  sku_name            = azurerm_firewall.test.sku_name
  sku_tier            = azurerm_firewall.test.sku_tier

  ip_configuration {
    name                 = "configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }
  threat_intel_mode = azurerm_firewall.test.threat_intel_mode
}
`, template)
}

func (FirewallResource) withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fw-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_firewall" "test" {
  name                = "acctestfirewall%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "AZFW_VNet"
  sku_tier            = "Standard"

  ip_configuration {
    name                 = "configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (FirewallResource) withUpdatedTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fw-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_firewall" "test" {
  name                = "acctestfirewall%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "AZFW_VNet"
  sku_tier            = "Standard"

  ip_configuration {
    name                 = "configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (FirewallResource) withZones(data acceptance.TestData, zones []string) string {
	zoneString := strings.Join(zones, ",")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fw-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = [%s]
}

resource "azurerm_firewall" "test" {
  name                = "acctestfirewall%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "AZFW_VNet"
  sku_tier            = "Standard"

  ip_configuration {
    name                 = "configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }

  zones = [%s]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, zoneString, data.RandomInteger, zoneString)
}

func (FirewallResource) withoutZone(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fw-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_firewall" "test" {
  name                = "acctestfirewall%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "AZFW_VNet"
  sku_tier            = "Standard"

  ip_configuration {
    name                 = "configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }

  zones = []
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (FirewallResource) withFirewallPolicy(data acceptance.TestData, policyName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fw-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_firewall_policy" "test" {
  name                = "acctestfirewall-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_firewall" "test" {
  name                = "acctestfirewall%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "AZFW_VNet"
  sku_tier            = "Standard"

  ip_configuration {
    name                 = "configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }

  firewall_policy_id = azurerm_firewall_policy.test.id

  lifecycle {
    create_before_destroy = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, policyName, data.RandomInteger)
}

func (FirewallResource) inVirtualHub(data acceptance.TestData, pipCount int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fw-%[1]d"
  location = "%s"
}

resource "azurerm_firewall_policy" "test" {
  name                = "acctest-firewallpolicy-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctest-virtualwan-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_hub" "test" {
  name                = "acctest-virtualhub-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.1.0/24"
}

resource "azurerm_firewall" "test" {
  name                = "acctest-firewall-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "AZFW_Hub"
  sku_tier            = "Standard"

  virtual_hub {
    virtual_hub_id  = azurerm_virtual_hub.test.id
    public_ip_count = %[3]d
  }

  firewall_policy_id = azurerm_firewall_policy.test.id
}
`, data.RandomInteger, data.Locations.Primary, pipCount)
}

func (FirewallResource) privateRanges(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-fw-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_firewall" "test" {
  name                = "acctestfirewall%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "AZFW_VNet"
  sku_tier            = "Standard"

  ip_configuration {
    name                 = "configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }
  private_ip_ranges = ["IANAPrivateRanges", "255.255.0.0/16"]
  threat_intel_mode = "Deny"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
