package cosmic

import (
	"fmt"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/cosmic"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCosmicDisk_basic(t *testing.T) {
	var disk cosmic.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicDiskDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCosmicDisk_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicDiskExists(
						"cosmic_disk.foo", &disk),
					testAccCheckCosmicDiskAttributes(&disk),
				),
			},
		},
	})
}

func TestAccCosmicDisk_deviceID(t *testing.T) {
	var disk cosmic.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicDiskDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCosmicDisk_deviceID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicDiskExists(
						"cosmic_disk.foo", &disk),
					testAccCheckCosmicDiskAttributes(&disk),
					resource.TestCheckResourceAttr(
						"cosmic_disk.foo", "device_id", "4"),
				),
			},
		},
	})
}

func TestAccCosmicDisk_update(t *testing.T) {
	var disk cosmic.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicDiskDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCosmicDisk_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicDiskExists(
						"cosmic_disk.foo", &disk),
					testAccCheckCosmicDiskAttributes(&disk),
				),
			},

			resource.TestStep{
				Config: testAccCosmicDisk_resize,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicDiskExists(
						"cosmic_disk.foo", &disk),
					testAccCheckCosmicDiskResized(&disk),
					resource.TestCheckResourceAttr(
						"cosmic_disk.foo", "disk_offering", COSMIC_DISK_OFFERING_2),
				),
			},
		},
	})
}

func testAccCheckCosmicDiskExists(
	n string, disk *cosmic.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No disk ID is set")
		}

		cs := testAccProvider.Meta().(*cosmic.CosmicClient)
		volume, _, err := cs.Volume.GetVolumeByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if volume.Id != rs.Primary.ID {
			return fmt.Errorf("Disk not found")
		}

		*disk = *volume

		return nil
	}
}

func testAccCheckCosmicDiskAttributes(
	disk *cosmic.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if disk.Name != "terraform-disk" {
			return fmt.Errorf("Bad name: %s", disk.Name)
		}

		if disk.Diskofferingname != COSMIC_DISK_OFFERING_1 {
			return fmt.Errorf("Bad disk offering: %s", disk.Diskofferingname)
		}

		return nil
	}
}

func testAccCheckCosmicDiskResized(
	disk *cosmic.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if disk.Diskofferingname != COSMIC_DISK_OFFERING_2 {
			return fmt.Errorf("Bad disk offering: %s", disk.Diskofferingname)
		}

		return nil
	}
}

func testAccCheckCosmicDiskDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cosmic.CosmicClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cosmic_disk" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No disk ID is set")
		}

		_, _, err := cs.Volume.GetVolumeByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Disk %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

var testAccCosmicDisk_basic = fmt.Sprintf(`
resource "cosmic_disk" "foo" {
  name = "terraform-disk"
  attach = false
  disk_offering = "%s"
  zone = "%s"
}`,
	COSMIC_DISK_OFFERING_1,
	COSMIC_ZONE)

var testAccCosmicDisk_deviceID = fmt.Sprintf(`
resource "cosmic_instance" "foobar" {
  name = "terraform-test"
  display_name = "terraform"
  service_offering= "%s"
  network_id = "%s"
  template = "%s"
  zone = "%s"
  expunge = true
}

resource "cosmic_disk" "foo" {
  name = "terraform-disk"
  attach = true
  device_id = 4
  disk_offering = "%s"
  virtual_machine_id = "${cosmic_instance.foobar.id}"
  zone = "${cosmic_instance.foobar.zone}"
}`,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_NETWORK_1,
	COSMIC_TEMPLATE,
	COSMIC_ZONE,
	COSMIC_DISK_OFFERING_1)

var testAccCosmicDisk_update = fmt.Sprintf(`
resource "cosmic_instance" "foobar" {
  name = "terraform-test"
  display_name = "terraform"
  service_offering= "%s"
  network_id = "%s"
  template = "%s"
  zone = "%s"
  expunge = true
}

resource "cosmic_disk" "foo" {
  name = "terraform-disk"
  attach = true
  disk_offering = "%s"
  virtual_machine_id = "${cosmic_instance.foobar.id}"
  zone = "${cosmic_instance.foobar.zone}"
}`,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_NETWORK_1,
	COSMIC_TEMPLATE,
	COSMIC_ZONE,
	COSMIC_DISK_OFFERING_1)

var testAccCosmicDisk_resize = fmt.Sprintf(`
resource "cosmic_instance" "foobar" {
  name = "terraform-test"
  display_name = "terraform"
  service_offering= "%s"
  network_id = "%s"
  template = "%s"
  zone = "%s"
  expunge = true
}

resource "cosmic_disk" "foo" {
  name = "terraform-disk"
  attach = true
  disk_offering = "%s"
	virtual_machine_id = "${cosmic_instance.foobar.id}"
  zone = "${cosmic_instance.foobar.zone}"
}`,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_NETWORK_1,
	COSMIC_TEMPLATE,
	COSMIC_ZONE,
	COSMIC_DISK_OFFERING_2)
