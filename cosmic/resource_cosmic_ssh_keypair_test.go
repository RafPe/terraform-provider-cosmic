package cosmic

import (
	"fmt"
	"strings"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/cosmic"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCosmicSSHKeyPair_basic(t *testing.T) {
	var sshkey cosmic.SSHKeyPair

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicSSHKeyPairDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCosmicSSHKeyPair_create,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicSSHKeyPairExists("cosmic_ssh_keypair.foo", &sshkey),
					testAccCheckCosmicSSHKeyPairAttributes(&sshkey),
					testAccCheckCosmicSSHKeyPairCreateAttributes("terraform-test-keypair"),
				),
			},
		},
	})
}

func TestAccCosmicSSHKeyPair_register(t *testing.T) {
	var sshkey cosmic.SSHKeyPair

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicSSHKeyPairDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCosmicSSHKeyPair_register,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicSSHKeyPairExists("cosmic_ssh_keypair.foo", &sshkey),
					testAccCheckCosmicSSHKeyPairAttributes(&sshkey),
					resource.TestCheckResourceAttr(
						"cosmic_ssh_keypair.foo",
						"public_key",
						COSMIC_SSH_PUBLIC_KEY),
				),
			},
		},
	})
}

func testAccCheckCosmicSSHKeyPairExists(n string, sshkey *cosmic.SSHKeyPair) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No key pair ID is set")
		}

		cs := testAccProvider.Meta().(*cosmic.CosmicClient)
		p := cs.SSH.NewListSSHKeyPairsParams()
		p.SetName(rs.Primary.ID)

		list, err := cs.SSH.ListSSHKeyPairs(p)
		if err != nil {
			return err
		}

		if list.Count != 1 || list.SSHKeyPairs[0].Name != rs.Primary.ID {
			return fmt.Errorf("Key pair not found")
		}

		*sshkey = *list.SSHKeyPairs[0]

		return nil
	}
}

func testAccCheckCosmicSSHKeyPairAttributes(
	keypair *cosmic.SSHKeyPair) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		fpLen := len(keypair.Fingerprint)
		if fpLen != 47 {
			return fmt.Errorf("SSH key: Attribute private_key expected length 47, got %d", fpLen)
		}

		return nil
	}
}

func testAccCheckCosmicSSHKeyPairCreateAttributes(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		found := false

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "cosmic_ssh_keypair" {
				continue
			}

			if rs.Primary.ID != name {
				continue
			}

			if !strings.Contains(rs.Primary.Attributes["private_key"], "PRIVATE KEY") {
				return fmt.Errorf(
					"SSH key: Attribute private_key expected 'PRIVATE KEY' to be present, got %s",
					rs.Primary.Attributes["private_key"])
			}

			found = true
			break
		}

		if !found {
			return fmt.Errorf("Could not find key pair %s", name)
		}

		return nil
	}
}

func testAccCheckCosmicSSHKeyPairDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cosmic.CosmicClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cosmic_ssh_keypair" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No key pair ID is set")
		}

		p := cs.SSH.NewListSSHKeyPairsParams()
		p.SetName(rs.Primary.ID)

		list, err := cs.SSH.ListSSHKeyPairs(p)
		if err != nil {
			return err
		}

		for _, keyPair := range list.SSHKeyPairs {
			if keyPair.Name == rs.Primary.ID {
				return fmt.Errorf("Key pair %s still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

var testAccCosmicSSHKeyPair_create = fmt.Sprintf(`
resource "cosmic_ssh_keypair" "foo" {
  name = "terraform-test-keypair"
}`)

var testAccCosmicSSHKeyPair_register = fmt.Sprintf(`
resource "cosmic_ssh_keypair" "foo" {
  name = "terraform-test-keypair"
  public_key = "%s"
}`, COSMIC_SSH_PUBLIC_KEY)
