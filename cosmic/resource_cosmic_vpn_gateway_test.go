package cosmic

import (
	"fmt"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/cosmic"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCosmicVPNGateway_basic(t *testing.T) {
	var vpnGateway cosmic.VpnGateway

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicVPNGatewayDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCosmicVPNGateway_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicVPNGatewayExists(
						"cosmic_vpn_gateway.foo", &vpnGateway),
				),
			},
		},
	})
}

func testAccCheckCosmicVPNGatewayExists(
	n string, vpnGateway *cosmic.VpnGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPN Gateway ID is set")
		}

		cs := testAccProvider.Meta().(*cosmic.CosmicClient)
		v, _, err := cs.VPN.GetVpnGatewayByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if v.Id != rs.Primary.ID {
			return fmt.Errorf("VPN Gateway not found")
		}

		*vpnGateway = *v

		return nil
	}
}

func testAccCheckCosmicVPNGatewayDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cosmic.CosmicClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cosmic_vpn_gateway" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPN Gateway ID is set")
		}

		_, _, err := cs.VPN.GetVpnGatewayByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("VPN Gateway %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

var testAccCosmicVPNGateway_basic = fmt.Sprintf(`
resource "cosmic_vpc" "foo" {
  name = "terraform-vpc"
  display_text = "terraform-vpc-text"
  cidr = "%s"
  vpc_offering = "%s"
  zone = "%s"
}

resource "cosmic_vpn_gateway" "foo" {
  vpc_id = "${cosmic_vpc.foo.id}"
}`,
	COSMIC_VPC_CIDR_1,
	COSMIC_VPC_OFFERING,
	COSMIC_ZONE)
