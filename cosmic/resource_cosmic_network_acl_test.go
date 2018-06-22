package cosmic

import (
	"fmt"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/cosmic"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCosmicNetworkACL_basic(t *testing.T) {
	var acl cosmic.NetworkACLList
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicNetworkACLDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCosmicNetworkACL_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicNetworkACLExists(
						"cosmic_network_acl.foo", &acl),
					testAccCheckCosmicNetworkACLBasicAttributes(&acl),
				),
			},
		},
	})
}

func testAccCheckCosmicNetworkACLExists(
	n string, acl *cosmic.NetworkACLList) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No network ACL ID is set")
		}

		cs := testAccProvider.Meta().(*cosmic.CosmicClient)
		acllist, _, err := cs.NetworkACL.GetNetworkACLListByID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if acllist.Id != rs.Primary.ID {
			return fmt.Errorf("Network ACL not found")
		}

		*acl = *acllist

		return nil
	}
}

func testAccCheckCosmicNetworkACLBasicAttributes(
	acl *cosmic.NetworkACLList) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if acl.Name != "terraform-acl" {
			return fmt.Errorf("Bad name: %s", acl.Name)
		}

		if acl.Description != "terraform-acl-text" {
			return fmt.Errorf("Bad description: %s", acl.Description)
		}

		return nil
	}
}

func testAccCheckCosmicNetworkACLDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cosmic.CosmicClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cosmic_network_acl" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No network ACL ID is set")
		}

		_, _, err := cs.NetworkACL.GetNetworkACLListByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Network ACl list %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

var testAccCosmicNetworkACL_basic = fmt.Sprintf(`
resource "cosmic_vpc" "foobar" {
  name = "terraform-vpc"
  cidr = "%s"
  vpc_offering = "%s"
  zone = "%s"
}

resource "cosmic_network_acl" "foo" {
  name = "terraform-acl"
  description = "terraform-acl-text"
  vpc_id = "${cosmic_vpc.foobar.id}"
}`,
	COSMIC_VPC_CIDR_1,
	COSMIC_VPC_OFFERING,
	COSMIC_ZONE)
