package cosmic

import (
	"fmt"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/cosmic"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCosmicVPC_basic(t *testing.T) {
	var vpc cosmic.VPC

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicVPCDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCosmicVPC_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicVPCExists(
						"cosmic_vpc.foo", &vpc),
					testAccCheckCosmicVPCAttributes(&vpc),
					resource.TestCheckResourceAttr(
						"cosmic_vpc.foo", "vpc_offering", COSMIC_VPC_OFFERING),
				),
			},
		},
	})
}

func testAccCheckCosmicVPCExists(
	n string, vpc *cosmic.VPC) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPC ID is set")
		}

		cs := testAccProvider.Meta().(*cosmic.CosmicClient)
		v, _, err := cs.VPC.GetVPCByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if v.Id != rs.Primary.ID {
			return fmt.Errorf("VPC not found")
		}

		*vpc = *v

		return nil
	}
}

func testAccCheckCosmicVPCAttributes(
	vpc *cosmic.VPC) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if vpc.Name != "terraform-vpc" {
			return fmt.Errorf("Bad name: %s", vpc.Name)
		}

		if vpc.Displaytext != "terraform-vpc-text" {
			return fmt.Errorf("Bad display text: %s", vpc.Displaytext)
		}

		if vpc.Cidr != COSMIC_VPC_CIDR_1 {
			return fmt.Errorf("Bad VPC CIDR: %s", vpc.Cidr)
		}

		if vpc.Networkdomain != "terraform-domain" {
			return fmt.Errorf("Bad network domain: %s", vpc.Networkdomain)
		}

		return nil
	}
}

func testAccCheckCosmicVPCDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cosmic.CosmicClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cosmic_vpc" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPC ID is set")
		}

		_, _, err := cs.VPC.GetVPCByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("VPC %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

var testAccCosmicVPC_basic = fmt.Sprintf(`
resource "cosmic_vpc" "foo" {
  name = "terraform-vpc"
  display_text = "terraform-vpc-text"
  cidr = "%s"
  vpc_offering = "%s"
  network_domain = "terraform-domain"
  zone = "%s"
}`,
	COSMIC_VPC_CIDR_1,
	COSMIC_VPC_OFFERING,
	COSMIC_ZONE)
