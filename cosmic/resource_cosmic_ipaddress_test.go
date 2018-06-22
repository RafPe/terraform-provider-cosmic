package cosmic

import (
	"fmt"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/cosmic"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCosmicIPAddress_basic(t *testing.T) {
	var ipaddr cosmic.PublicIpAddress

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicIPAddressDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCosmicIPAddress_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicIPAddressExists(
						"cosmic_ipaddress.foo", &ipaddr),
					testAccCheckCosmicIPAddressAttributes(&ipaddr),
				),
			},
		},
	})
}

func TestAccCosmicIPAddress_vpc(t *testing.T) {
	var ipaddr cosmic.PublicIpAddress

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicIPAddressDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCosmicIPAddress_vpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicIPAddressExists(
						"cosmic_ipaddress.foo", &ipaddr),
				),
			},
		},
	})
}

func testAccCheckCosmicIPAddressExists(
	n string, ipaddr *cosmic.PublicIpAddress) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No IP address ID is set")
		}

		cs := testAccProvider.Meta().(*cosmic.CosmicClient)
		pip, _, err := cs.PublicIPAddress.GetPublicIpAddressByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if pip.Id != rs.Primary.ID {
			return fmt.Errorf("IP address not found")
		}

		*ipaddr = *pip

		return nil
	}
}

func testAccCheckCosmicIPAddressAttributes(
	ipaddr *cosmic.PublicIpAddress) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if ipaddr.Associatednetworkid != COSMIC_NETWORK_1 {
			return fmt.Errorf("Bad network ID: %s", ipaddr.Associatednetworkid)
		}

		return nil
	}
}

func testAccCheckCosmicIPAddressDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cosmic.CosmicClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cosmic_ipaddress" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No IP address ID is set")
		}

		ip, _, err := cs.PublicIPAddress.GetPublicIpAddressByID(rs.Primary.ID)
		if err == nil && ip.Associatednetworkid != "" {
			return fmt.Errorf("Public IP %s still associated", rs.Primary.ID)
		}
	}

	return nil
}

var testAccCosmicIPAddress_basic = fmt.Sprintf(`
resource "cosmic_ipaddress" "foo" {
  network_id = "%s"
}`, COSMIC_NETWORK_1)

var testAccCosmicIPAddress_vpc = fmt.Sprintf(`
resource "cosmic_vpc" "foobar" {
  name = "terraform-vpc"
  cidr = "%s"
  vpc_offering = "%s"
  zone = "%s"
}

resource "cosmic_ipaddress" "foo" {
  vpc_id = "${cosmic_vpc.foobar.id}"
}`,
	COSMIC_VPC_CIDR_1,
	COSMIC_VPC_OFFERING,
	COSMIC_ZONE)
