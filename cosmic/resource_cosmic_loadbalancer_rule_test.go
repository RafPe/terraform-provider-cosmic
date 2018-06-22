package cosmic

import (
	"fmt"
	"strings"
	"testing"

	"github.com/MissionCriticalCloud/go-cosmic/cosmic"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCosmicLoadBalancerRule_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicLoadBalancerRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCosmicLoadBalancerRule_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicLoadBalancerRuleExist("cosmic_loadbalancer_rule.foo", nil),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "name", "terraform-lb"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "algorithm", "roundrobin"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "public_port", "80"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "private_port", "80"),
				),
			},
		},
	})
}

func TestAccCosmicLoadBalancerRule_update(t *testing.T) {
	var id string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicLoadBalancerRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCosmicLoadBalancerRule_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicLoadBalancerRuleExist("cosmic_loadbalancer_rule.foo", &id),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "name", "terraform-lb"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "algorithm", "roundrobin"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "public_port", "80"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "private_port", "80"),
				),
			},

			resource.TestStep{
				Config: testAccCosmicLoadBalancerRule_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicLoadBalancerRuleExist("cosmic_loadbalancer_rule.foo", &id),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "name", "terraform-lb-update"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "algorithm", "leastconn"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "public_port", "80"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "private_port", "80"),
				),
			},
		},
	})
}

func TestAccCosmicLoadBalancerRule_forceNew(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicLoadBalancerRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCosmicLoadBalancerRule_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicLoadBalancerRuleExist("cosmic_loadbalancer_rule.foo", nil),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "name", "terraform-lb"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "algorithm", "roundrobin"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "public_port", "80"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "private_port", "80"),
				),
			},

			resource.TestStep{
				Config: testAccCosmicLoadBalancerRule_forcenew,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicLoadBalancerRuleExist("cosmic_loadbalancer_rule.foo", nil),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "name", "terraform-lb-update"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "algorithm", "leastconn"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "public_port", "443"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "private_port", "443"),
				),
			},
		},
	})
}

func TestAccCosmicLoadBalancerRule_vpc(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicLoadBalancerRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCosmicLoadBalancerRule_vpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicLoadBalancerRuleExist("cosmic_loadbalancer_rule.foo", nil),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "name", "terraform-lb"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "algorithm", "roundrobin"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "public_port", "80"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "private_port", "80"),
				),
			},
		},
	})
}

func TestAccCosmicLoadBalancerRule_vpcUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosmicLoadBalancerRuleDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCosmicLoadBalancerRule_vpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicLoadBalancerRuleExist("cosmic_loadbalancer_rule.foo", nil),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "name", "terraform-lb"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "algorithm", "roundrobin"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "public_port", "80"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "private_port", "80"),
				),
			},

			resource.TestStep{
				Config: testAccCosmicLoadBalancerRule_vpc_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosmicLoadBalancerRuleExist("cosmic_loadbalancer_rule.foo", nil),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "name", "terraform-lb-update"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "algorithm", "leastconn"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "public_port", "443"),
					resource.TestCheckResourceAttr(
						"cosmic_loadbalancer_rule.foo", "private_port", "443"),
				),
			},
		},
	})
}

func testAccCheckCosmicLoadBalancerRuleExist(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No loadbalancer rule ID is set")
		}

		if id != nil {
			if *id != "" && *id != rs.Primary.ID {
				return fmt.Errorf("Resource ID has changed!")
			}

			*id = rs.Primary.ID
		}

		cs := testAccProvider.Meta().(*cosmic.CosmicClient)
		_, count, err := cs.LoadBalancer.GetLoadBalancerRuleByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if count == 0 {
			return fmt.Errorf("Loadbalancer rule %s not found", n)
		}

		return nil
	}
}

func testAccCheckCosmicLoadBalancerRuleDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cosmic.CosmicClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cosmic_loadbalancer_rule" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Loadbalancer rule ID is set")
		}

		for k, id := range rs.Primary.Attributes {
			if !strings.Contains(k, "uuid") {
				continue
			}

			_, _, err := cs.LoadBalancer.GetLoadBalancerRuleByID(id)
			if err == nil {
				return fmt.Errorf("Loadbalancer rule %s still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

var testAccCosmicLoadBalancerRule_basic = fmt.Sprintf(`
resource "cosmic_instance" "foobar1" {
  name = "terraform-server1"
  display_name = "terraform"
  service_offering= "%s"
  network_id = "%s"
  template = "%s"
  zone = "%s"
  expunge = true
}

resource "cosmic_loadbalancer_rule" "foo" {
  name = "terraform-lb"
  ip_address_id = "%s"
  algorithm = "roundrobin"
  public_port = 80
  private_port = 80
  member_ids = ["${cosmic_instance.foobar1.id}"]
}
`,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_NETWORK_1,
	COSMIC_TEMPLATE,
	COSMIC_ZONE,
	COSMIC_PUBLIC_IPADDRESS)

var testAccCosmicLoadBalancerRule_update = fmt.Sprintf(`
resource "cosmic_instance" "foobar1" {
  name = "terraform-server1"
  display_name = "terraform"
  service_offering= "%s"
  network_id = "%s"
  template = "%s"
  zone = "%s"
  expunge = true
}

resource "cosmic_loadbalancer_rule" "foo" {
  name = "terraform-lb-update"
  ip_address_id = "%s"
  algorithm = "leastconn"
  public_port = 80
  private_port = 80
  member_ids = ["${cosmic_instance.foobar1.id}"]
}
`,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_NETWORK_1,
	COSMIC_TEMPLATE,
	COSMIC_ZONE,
	COSMIC_PUBLIC_IPADDRESS)

var testAccCosmicLoadBalancerRule_forcenew = fmt.Sprintf(`
resource "cosmic_instance" "foobar1" {
  name = "terraform-server1"
  display_name = "terraform"
  service_offering= "%s"
  network_id = "%s"
  template = "%s"
  zone = "%s"
  expunge = true
}

resource "cosmic_loadbalancer_rule" "foo" {
  name = "terraform-lb-update"
  ip_address_id = "%s"
  algorithm = "leastconn"
  public_port = 443
  private_port = 443
  member_ids = ["${cosmic_instance.foobar1.id}"]
}
`,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_NETWORK_1,
	COSMIC_TEMPLATE,
	COSMIC_ZONE,
	COSMIC_PUBLIC_IPADDRESS)

var testAccCosmicLoadBalancerRule_vpc = fmt.Sprintf(`
resource "cosmic_vpc" "foobar" {
	name = "terraform-vpc"
	cidr = "%s"
	vpc_offering = "%s"
	zone = "%s"
}

resource "cosmic_network" "foo" {
  name = "terraform-network"
  cidr = "%s"
  gateway = "%s"
  network_offering = "%s"
  vpc_id = "${cosmic_vpc.foobar.id}"
  zone = "${cosmic_vpc.foobar.zone}"
}

resource "cosmic_ipaddress" "foo" {
  vpc_id = "${cosmic_vpc.foobar.id}"
}

resource "cosmic_instance" "foobar1" {
  name = "terraform-server1"
  display_name = "terraform"
  service_offering= "%s"
  network_id = "${cosmic_network.foo.id}"
  template = "%s"
  zone = "${cosmic_network.foo.zone}"
  expunge = true
}

resource "cosmic_loadbalancer_rule" "foo" {
  name = "terraform-lb"
  ip_address_id = "${cosmic_ipaddress.foo.id}"
  algorithm = "roundrobin"
  network_id = "${cosmic_network.foo.id}"
  public_port = 80
  private_port = 80
  member_ids = ["${cosmic_instance.foobar1.id}"]
}`,
	COSMIC_VPC_CIDR_1,
	COSMIC_VPC_OFFERING,
	COSMIC_ZONE,
	COSMIC_VPC_NETWORK_CIDR,
	COSMIC_VPC_NETWORK_GATEWAY,
	COSMIC_VPC_NETWORK_OFFERING,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_TEMPLATE)

var testAccCosmicLoadBalancerRule_vpc_update = fmt.Sprintf(`
resource "cosmic_vpc" "foobar" {
  name = "terraform-vpc"
  cidr = "%s"
  vpc_offering = "%s"
  zone = "%s"
}

resource "cosmic_network" "foo" {
  name = "terraform-network"
  cidr = "%s"
  gateway = "%s"
  network_offering = "%s"
  vpc_id = "${cosmic_vpc.foobar.id}"
  zone = "${cosmic_vpc.foobar.zone}"
}

resource "cosmic_ipaddress" "foo" {
  vpc_id = "${cosmic_vpc.foobar.id}"
}

resource "cosmic_instance" "foobar1" {
  name = "terraform-server1"
  display_name = "terraform"
  service_offering= "%s"
  network_id = "${cosmic_network.foo.id}"
  template = "%s"
  zone = "${cosmic_network.foo.zone}"
  expunge = true
}

resource "cosmic_instance" "foobar2" {
  name = "terraform-server2"
  display_name = "terraform"
  service_offering= "%s"
  network_id = "${cosmic_network.foo.id}"
  template = "%s"
  zone = "${cosmic_network.foo.zone}"
  expunge = true
}

resource "cosmic_loadbalancer_rule" "foo" {
  name = "terraform-lb-update"
  ip_address_id = "${cosmic_ipaddress.foo.id}"
  algorithm = "leastconn"
  network_id = "${cosmic_network.foo.id}"
  public_port = 443
  private_port = 443
  member_ids = ["${cosmic_instance.foobar1.id}", "${cosmic_instance.foobar2.id}"]
}`,
	COSMIC_VPC_CIDR_1,
	COSMIC_VPC_OFFERING,
	COSMIC_ZONE,
	COSMIC_VPC_NETWORK_CIDR,
	COSMIC_VPC_NETWORK_GATEWAY,
	COSMIC_VPC_NETWORK_OFFERING,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_TEMPLATE,
	COSMIC_SERVICE_OFFERING_1,
	COSMIC_TEMPLATE)
