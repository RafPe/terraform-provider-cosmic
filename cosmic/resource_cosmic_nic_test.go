package cosmic

import (
	"fmt"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/cosmic"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCosmicNIC_basic(t *testing.T) {
	var nic cosmic.Nic

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicNICDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCosmicNIC_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicNICExists(
						"cosmic_instance.foobar", "cosmic_nic.foo", &nic),
					testAccCheckCosmicNICAttributes(&nic),
				),
			},
		},
	})
}

func TestAccCosmicNIC_update(t *testing.T) {
	var nic cosmic.Nic

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicNICDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCosmicNIC_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicNICExists(
						"cosmic_instance.foobar", "cosmic_nic.foo", &nic),
					testAccCheckCosmicNICAttributes(&nic),
				),
			},

			resource.TestStep{
				Config: testAccCosmicNIC_ipaddress,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicNICExists(
						"cosmic_instance.foobar", "cosmic_nic.foo", &nic),
					testAccCheckCosmicNICIPAddress(&nic),
					resource.TestCheckResourceAttr(
						"cosmic_nic.foo", "ip_address", COSMIC_2ND_NIC_IPADDRESS),
				),
			},
		},
	})
}

func testAccCheckCosmicNICExists(
	v, n string, nic *cosmic.Nic) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rsv, ok := s.RootModule().Resources[v]
		if !ok {
			return fmt.Errorf("Not found: %s", v)
		}

		if rsv.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		rsn, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rsn.Primary.ID == "" {
			return fmt.Errorf("No NIC ID is set")
		}

		cs := testAccProvider.Meta().(*cosmic.CosmicClient)
		vm, _, err := cs.VirtualMachine.GetVirtualMachineByID(rsv.Primary.ID)

		if err != nil {
			return err
		}

		for _, n := range vm.Nic {
			if n.Id == rsn.Primary.ID {
				*nic = n
				return nil
			}
		}

		return fmt.Errorf("NIC not found")
	}
}

func testAccCheckCosmicNICAttributes(
	nic *cosmic.Nic) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if nic.Networkid != COSMIC_2ND_NIC_NETWORK {
			return fmt.Errorf("Bad network ID: %s", nic.Networkid)
		}

		return nil
	}
}

func testAccCheckCosmicNICIPAddress(
	nic *cosmic.Nic) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if nic.Networkid != COSMIC_2ND_NIC_NETWORK {
			return fmt.Errorf("Bad network ID: %s", nic.Networkname)
		}

		if nic.Ipaddress != COSMIC_2ND_NIC_IPADDRESS {
			return fmt.Errorf("Bad IP address: %s", nic.Ipaddress)
		}

		return nil
	}
}

func testAccCheckCosmicNICDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cosmic.CosmicClient)

	// Deleting the instance automatically deletes any additional NICs
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cosmic_instance" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		_, _, err := cs.VirtualMachine.GetVirtualMachineByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Virtual Machine %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

var testAccCosmicNIC_basic = fmt.Sprintf(`
resource "cosmic_instance" "foobar" {
  name = "terraform-test"
  display_name = "terraform"
  service_offering= "%s"
  network_id = "%s"
  template = "%s"
  zone = "%s"
  expunge = true
}

resource "cosmic_nic" "foo" {
  network_id = "%s"
  virtual_machine_id = "${cosmic_instance.foobar.id}"
}`,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_NETWORK_1,
	COSMIC_TEMPLATE,
	COSMIC_ZONE,
	COSMIC_2ND_NIC_NETWORK)

var testAccCosmicNIC_ipaddress = fmt.Sprintf(`
resource "cosmic_instance" "foobar" {
  name = "terraform-test"
  display_name = "terraform"
  service_offering= "%s"
  network_id = "%s"
  template = "%s"
  zone = "%s"
  expunge = true
}

resource "cosmic_nic" "foo" {
  network_id = "%s"
  ip_address = "%s"
  virtual_machine_id = "${cosmic_instance.foobar.id}"
}`,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_NETWORK_1,
	COSMIC_TEMPLATE,
	COSMIC_ZONE,
	COSMIC_2ND_NIC_NETWORK,
	COSMIC_2ND_NIC_IPADDRESS)
