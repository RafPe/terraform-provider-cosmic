package cosmic

import (
	"fmt"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/cosmic"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCosmicStaticNAT_basic(t *testing.T) {
	var ipaddr cosmic.PublicIpAddress

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicStaticNATDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCosmicStaticNAT_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicStaticNATExists(
						"cosmic_static_nat.foo", &ipaddr),
					testAccCheckCosmicStaticNATAttributes(&ipaddr),
				),
			},
		},
	})
}

func testAccCheckCosmicStaticNATExists(
	n string, ipaddr *cosmic.PublicIpAddress) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No static NAT ID is set")
		}

		cs := testAccProvider.Meta().(*cosmic.CosmicClient)
		ip, _, err := cs.PublicIPAddress.GetPublicIpAddressByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if ip.Id != rs.Primary.ID {
			return fmt.Errorf("Static NAT not found")
		}

		if !ip.Isstaticnat {
			return fmt.Errorf("Static NAT not enabled")
		}

		*ipaddr = *ip

		return nil
	}
}

func testAccCheckCosmicStaticNATAttributes(
	ipaddr *cosmic.PublicIpAddress) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if ipaddr.Associatednetworkid != COSMIC_NETWORK_1 {
			return fmt.Errorf("Bad network ID: %s", ipaddr.Associatednetworkid)
		}

		return nil
	}
}

func testAccCheckCosmicStaticNATDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cosmic.CosmicClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cosmic_static_nat" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No static NAT ID is set")
		}

		ip, _, err := cs.PublicIPAddress.GetPublicIpAddressByID(rs.Primary.ID)
		if err == nil && ip.Isstaticnat {
			return fmt.Errorf("Static NAT %s still enabled", rs.Primary.ID)
		}
	}

	return nil
}

var testAccCosmicStaticNAT_basic = fmt.Sprintf(`
resource "cosmic_instance" "foobar" {
  name = "terraform-test"
  display_name = "terraform-test"
  service_offering= "%s"
  network_id = "%s"
  template = "%s"
  zone = "%s"
  user_data = "foobar\nfoo\nbar"
  expunge = true
}

resource "cosmic_ipaddress" "foo" {
  network_id = "${cosmic_instance.foobar.network_id}"
}

resource "cosmic_static_nat" "foo" {
	ip_address_id = "${cosmic_ipaddress.foo.id}"
  virtual_machine_id = "${cosmic_instance.foobar.id}"
}`,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_NETWORK_1,
	COSMIC_TEMPLATE,
	COSMIC_ZONE,
)
