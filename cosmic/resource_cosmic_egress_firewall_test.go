package cosmic

import (
	"fmt"
	"strings"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/cosmic"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCosmicEgressFirewall_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicEgressFirewallDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCosmicEgressFirewall_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicEgressFirewallRulesExist("cosmic_egress_firewall.foo"),
					resource.TestCheckResourceAttr(
						"cosmic_egress_firewall.foo", "rule.#", "1"),
					resource.TestCheckResourceAttr(
						"cosmic_egress_firewall.foo",
						"rule.813975931.cidr_list.1700444180",
						"10.10.10.10/32"),
					resource.TestCheckResourceAttr(
						"cosmic_egress_firewall.foo", "rule.813975931.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"cosmic_egress_firewall.foo", "rule.813975931.ports.32925333", "8080"),
				),
			},
		},
	})
}

func TestAccCosmicEgressFirewall_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicEgressFirewallDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCosmicEgressFirewall_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicEgressFirewallRulesExist("cosmic_egress_firewall.foo"),
					resource.TestCheckResourceAttr(
						"cosmic_egress_firewall.foo", "rule.#", "1"),
					resource.TestCheckResourceAttr(
						"cosmic_egress_firewall.foo",
						"rule.813975931.cidr_list.1700444180",
						"10.10.10.10/32"),
					resource.TestCheckResourceAttr(
						"cosmic_egress_firewall.foo", "rule.813975931.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"cosmic_egress_firewall.foo", "rule.813975931.ports.32925333", "8080"),
				),
			},

			resource.TestStep{
				Config: testAccCosmicEgressFirewall_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicEgressFirewallRulesExist("cosmic_egress_firewall.foo"),
					resource.TestCheckResourceAttr(
						"cosmic_egress_firewall.foo", "rule.#", "2"),
					resource.TestCheckResourceAttr(
						"cosmic_egress_firewall.foo",
						"rule.233813027.cidr_list.3722895217",
						"10.10.10.11/32"),
					resource.TestCheckResourceAttr(
						"cosmic_egress_firewall.foo",
						"rule.3600316628.cidr_list.1700444180",
						"10.10.10.10/32"),
					resource.TestCheckResourceAttr(
						"cosmic_egress_firewall.foo", "rule.233813027.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"cosmic_egress_firewall.foo", "rule.3600316628.ports.32925333", "8080"),
					resource.TestCheckResourceAttr(
						"cosmic_egress_firewall.foo",
						"rule.3600316628.cidr_list.3722895217",
						"10.10.10.11/32"),
					resource.TestCheckResourceAttr(
						"cosmic_egress_firewall.foo", "rule.3600316628.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"cosmic_egress_firewall.foo", "rule.233813027.ports.1889509032", "80"),
				),
			},
		},
	})
}

func testAccCheckCosmicEgressFirewallRulesExist(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No firewall ID is set")
		}

		for k, id := range rs.Primary.Attributes {
			if !strings.Contains(k, ".uuids.") || strings.HasSuffix(k, ".uuids.%") {
				continue
			}

			cs := testAccProvider.Meta().(*cosmic.CosmicClient)
			_, count, err := cs.Firewall.GetEgressFirewallRuleByID(id)

			if err != nil {
				return err
			}

			if count == 0 {
				return fmt.Errorf("Firewall rule for %s not found", k)
			}
		}

		return nil
	}
}

func testAccCheckCosmicEgressFirewallDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cosmic.CosmicClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cosmic_egress_firewall" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		for k, id := range rs.Primary.Attributes {
			if !strings.Contains(k, ".uuids.") || strings.HasSuffix(k, ".uuids.%") {
				continue
			}

			_, _, err := cs.Firewall.GetEgressFirewallRuleByID(id)
			if err == nil {
				return fmt.Errorf("Egress rule %s still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

var testAccCosmicEgressFirewall_basic = fmt.Sprintf(`
resource "cosmic_network" "foo" {
  name = "terraform-network"
  cidr = "10.10.10.0/24"
  network_offering = "%s"
  zone = "%s"
  tags = {
    terraform-tag = "true"
  }
}

resource "cosmic_egress_firewall" "foo" {
  network_id = "${cosmic_network.foo.id}"

  rule {
    cidr_list = ["10.10.10.10/32"]
    protocol = "tcp"
    ports = ["8080"]
  }
}`,
	COSMIC_NETWORK_2_OFFERING,
	COSMIC_ZONE)

var testAccCosmicEgressFirewall_update = fmt.Sprintf(`
resource "cosmic_network" "foo" {
  name = "terraform-network"
  cidr = "10.10.10.0/24"
  network_offering = "%s"
  zone = "%s"
  tags = {
    terraform-tag = "true"
  }
}

resource "cosmic_egress_firewall" "foo" {
  network_id = "${cosmic_network.foo.id}"

  rule {
    cidr_list = ["10.10.10.10/32", "10.10.10.11/32"]
    protocol = "tcp"
    ports = ["8080"]
  }

  rule {
    cidr_list = ["10.10.10.11/32"]
    protocol = "tcp"
    ports = ["80", "1000-2000"]
  }
}`,
	COSMIC_NETWORK_2_OFFERING,
	COSMIC_ZONE)
