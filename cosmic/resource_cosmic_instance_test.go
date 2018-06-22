package cosmic

import (
	"fmt"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/cosmic"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCosmicInstance_basic(t *testing.T) {
	var instance cosmic.VirtualMachine

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCosmicInstance_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicInstanceExists(
						"cosmic_instance.foobar", &instance),
					testAccCheckCosmicInstanceAttributes(&instance),
					resource.TestCheckResourceAttr(
						"cosmic_instance.foobar", "user_data", "0cf3dcdc356ec8369494cb3991985ecd5296cdd5"),
				),
			},
		},
	})
}

func TestAccCosmicInstance_update(t *testing.T) {
	var instance cosmic.VirtualMachine

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCosmicInstance_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicInstanceExists(
						"cosmic_instance.foobar", &instance),
					testAccCheckCosmicInstanceAttributes(&instance),
					resource.TestCheckResourceAttr(
						"cosmic_instance.foobar", "user_data", "0cf3dcdc356ec8369494cb3991985ecd5296cdd5"),
				),
			},

			resource.TestStep{
				Config: testAccCosmicInstance_renameAndResize,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicInstanceExists(
						"cosmic_instance.foobar", &instance),
					testAccCheckCosmicInstanceRenamedAndResized(&instance),
					resource.TestCheckResourceAttr(
						"cosmic_instance.foobar", "name", "terraform-updated"),
					resource.TestCheckResourceAttr(
						"cosmic_instance.foobar", "display_name", "terraform-updated"),
					resource.TestCheckResourceAttr(
						"cosmic_instance.foobar", "service_offering", COSMIC_SERVICE_OFFERING_2),
				),
			},
		},
	})
}

func TestAccCosmicInstance_fixedIP(t *testing.T) {
	var instance cosmic.VirtualMachine

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCosmicInstance_fixedIP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicInstanceExists(
						"cosmic_instance.foobar", &instance),
					resource.TestCheckResourceAttr(
						"cosmic_instance.foobar", "ip_address", COSMIC_NETWORK_1_IPADDRESS1),
				),
			},
		},
	})
}

func TestAccCosmicInstance_keyPair(t *testing.T) {
	var instance cosmic.VirtualMachine

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCosmicInstance_keyPair,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicInstanceExists(
						"cosmic_instance.foobar", &instance),
					resource.TestCheckResourceAttr(
						"cosmic_instance.foobar", "keypair", "terraform-test-keypair"),
				),
			},
		},
	})
}

func TestAccCosmicInstance_project(t *testing.T) {
	var instance cosmic.VirtualMachine

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCosmicInstance_project,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicInstanceExists(
						"cosmic_instance.foobar", &instance),
					resource.TestCheckResourceAttr(
						"cosmic_instance.foobar", "project", COSMIC_PROJECT_NAME),
				),
			},
		},
	})
}

func testAccCheckCosmicInstanceExists(
	n string, instance *cosmic.VirtualMachine) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		cs := testAccProvider.Meta().(*cosmic.CosmicClient)
		vm, _, err := cs.VirtualMachine.GetVirtualMachineByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if vm.Id != rs.Primary.ID {
			return fmt.Errorf("Instance not found")
		}

		*instance = *vm

		return nil
	}
}

func testAccCheckCosmicInstanceAttributes(
	instance *cosmic.VirtualMachine) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if instance.Name != "terraform-test" {
			return fmt.Errorf("Bad name: %s", instance.Name)
		}

		if instance.Displayname != "terraform-test" {
			return fmt.Errorf("Bad display name: %s", instance.Displayname)
		}

		if instance.Serviceofferingname != COSMIC_SERVICE_OFFERING_1 {
			return fmt.Errorf("Bad service offering: %s", instance.Serviceofferingname)
		}

		if instance.Templatename != COSMIC_TEMPLATE {
			return fmt.Errorf("Bad template: %s", instance.Templatename)
		}

		if instance.Nic[0].Networkid != COSMIC_NETWORK_1 {
			return fmt.Errorf("Bad network ID: %s", instance.Nic[0].Networkid)
		}

		return nil
	}
}

func testAccCheckCosmicInstanceRenamedAndResized(
	instance *cosmic.VirtualMachine) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if instance.Name != "terraform-updated" {
			return fmt.Errorf("Bad name: %s", instance.Name)
		}

		if instance.Displayname != "terraform-updated" {
			return fmt.Errorf("Bad display name: %s", instance.Displayname)
		}

		if instance.Serviceofferingname != COSMIC_SERVICE_OFFERING_2 {
			return fmt.Errorf("Bad service offering: %s", instance.Serviceofferingname)
		}

		return nil
	}
}

func testAccCheckCosmicInstanceDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cosmic.CosmicClient)

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

var testAccCosmicInstance_basic = fmt.Sprintf(`
resource "cosmic_instance" "foobar" {
  name = "terraform-test"
  display_name = "terraform-test"
  service_offering= "%s"
  network_id = "%s"
  template = "%s"
  zone = "%s"
  user_data = "foobar\nfoo\nbar"
  expunge = true
}`,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_NETWORK_1,
	COSMIC_TEMPLATE,
	COSMIC_ZONE)

var testAccCosmicInstance_renameAndResize = fmt.Sprintf(`
resource "cosmic_instance" "foobar" {
  name = "terraform-updated"
  display_name = "terraform-updated"
  service_offering= "%s"
  network_id = "%s"
  template = "%s"
  zone = "%s"
  user_data = "foobar\nfoo\nbar"
  expunge = true
}`,
	COSMIC_SERVICE_OFFERING_2,
	COSMIC_NETWORK_1,
	COSMIC_TEMPLATE,
	COSMIC_ZONE)

var testAccCosmicInstance_fixedIP = fmt.Sprintf(`
resource "cosmic_instance" "foobar" {
  name = "terraform-test"
  display_name = "terraform-test"
  service_offering= "%s"
  network_id = "%s"
  ip_address = "%s"
  template = "%s"
  zone = "%s"
  expunge = true
}`,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_NETWORK_1,
	COSMIC_NETWORK_1_IPADDRESS1,
	COSMIC_TEMPLATE,
	COSMIC_ZONE)

var testAccCosmicInstance_keyPair = fmt.Sprintf(`
resource "cosmic_ssh_keypair" "foo" {
  name = "terraform-test-keypair"
}

resource "cosmic_instance" "foobar" {
  name = "terraform-test"
  display_name = "terraform-test"
  service_offering= "%s"
  network_id = "%s"
  ip_address = "%s"
  template = "%s"
  zone = "%s"
	keypair = "${cosmic_ssh_keypair.foo.name}"
  expunge = true
}`,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_NETWORK_1,
	COSMIC_NETWORK_1_IPADDRESS1,
	COSMIC_TEMPLATE,
	COSMIC_ZONE)

var testAccCosmicInstance_project = fmt.Sprintf(`
resource "cosmic_instance" "foobar" {
  name = "terraform-test"
  display_name = "terraform-test"
  service_offering= "%s"
	network_id = "%s"
  template = "%s"
	project = "%s"
  zone = "%s"
  expunge = true
}`,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_PROJECT_NETWORK,
	COSMIC_TEMPLATE,
	COSMIC_PROJECT_NAME,
	COSMIC_ZONE)
