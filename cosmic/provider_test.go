package cosmic

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"cosmic": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testSetValueOnResourceData(t *testing.T) {
	d := schema.ResourceData{}
	d.Set("id", "name")

	setValueOrID(&d, "id", "name", "54711781-274e-41b2-83c0-17194d0108f7")

	if d.Get("id").(string) != "name" {
		t.Fatal("err: 'id' does not match 'name'")
	}
}

func testSetIDOnResourceData(t *testing.T) {
	d := schema.ResourceData{}
	d.Set("id", "54711781-274e-41b2-83c0-17194d0108f7")

	setValueOrID(&d, "id", "name", "54711781-274e-41b2-83c0-17194d0108f7")

	if d.Get("id").(string) != "54711781-274e-41b2-83c0-17194d0108f7" {
		t.Fatal("err: 'id' doest not match '54711781-274e-41b2-83c0-17194d0108f7'")
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("COSMIC_API_URL"); v == "" {
		t.Fatal("COSMIC_API_URL must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_API_KEY"); v == "" {
		t.Fatal("COSMIC_API_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_SECRET_KEY"); v == "" {
		t.Fatal("COSMIC_SECRET_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_2ND_NIC_IPADDRESS"); v == "" {
		t.Fatal("COSMIC_2ND_NIC_IPADDRESS must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_2ND_NIC_NETWORK"); v == "" {
		t.Fatal("COSMIC_2ND_NIC_NETWORK must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_DISK_OFFERING_1"); v == "" {
		t.Fatal("COSMIC_DISK_OFFERING_1 must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_DISK_OFFERING_2"); v == "" {
		t.Fatal("COSMIC_DISK_OFFERING_2 must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_HYPERVISOR"); v == "" {
		t.Fatal("COSMIC_HYPERVISOR must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_SERVICE_OFFERING_1"); v == "" {
		t.Fatal("COSMIC_SERVICE_OFFERING_1 must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_SERVICE_OFFERING_2"); v == "" {
		t.Fatal("COSMIC_SERVICE_OFFERING_2 must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_NETWORK_1"); v == "" {
		t.Fatal("COSMIC_NETWORK_1 must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_NETWORK_1_IPADDRESS1"); v == "" {
		t.Fatal("COSMIC_NETWORK_1_IPADDRESS1 must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_NETWORK_1_IPADDRESS2"); v == "" {
		t.Fatal("COSMIC_NETWORK_1_IPADDRESS2 must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_NETWORK_2"); v == "" {
		t.Fatal("COSMIC_NETWORK_2 must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_NETWORK_2_CIDR"); v == "" {
		t.Fatal("COSMIC_NETWORK_2_CIDR must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_NETWORK_2_GATEWAY"); v == "" {
		t.Fatal("COSMIC_NETWORK_2_GATEWAY must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_NETWORK_2_OFFERING"); v == "" {
		t.Fatal("COSMIC_NETWORK_2_OFFERING must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_NETWORK_2_IPADDRESS"); v == "" {
		t.Fatal("COSMIC_NETWORK_2_IPADDRESS must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_VPC_CIDR_1"); v == "" {
		t.Fatal("COSMIC_VPC_CIDR_1 must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_VPC_CIDR_2"); v == "" {
		t.Fatal("COSMIC_VPC_CIDR_2 must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_VPC_OFFERING"); v == "" {
		t.Fatal("COSMIC_VPC_OFFERING must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_VPC_NETWORK_CIDR"); v == "" {
		t.Fatal("COSMIC_VPC_NETWORK_CIDR must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_VPC_NETWORK_GATEWAY"); v == "" {
		t.Fatal("COSMIC_VPC_NETWORK_GATEWAY must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_VPC_NETWORK_OFFERING"); v == "" {
		t.Fatal("COSMIC_VPC_NETWORK_OFFERING must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_VPC_NETWORK_IPADDRESS"); v == "" {
		t.Fatal("COSMIC_VPC_NETWORK_IPADDRESS must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_PUBLIC_IPADDRESS"); v == "" {
		t.Fatal("COSMIC_PUBLIC_IPADDRESS must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_SSH_PUBLIC_KEY"); v == "" {
		t.Fatal("COSMIC_SSH_PUBLIC_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_TEMPLATE"); v == "" {
		t.Fatal("COSMIC_TEMPLATE must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_TEMPLATE_FORMAT"); v == "" {
		t.Fatal("COSMIC_TEMPLATE_FORMAT must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_TEMPLATE_URL"); v == "" {
		t.Fatal("COSMIC_TEMPLATE_URL must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_TEMPLATE_OS_TYPE"); v == "" {
		t.Fatal("COSMIC_TEMPLATE_OS_TYPE must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_PROJECT_NAME"); v == "" {
		t.Fatal("COSMIC_PROJECT_NAME must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_PROJECT_NETWORK"); v == "" {
		t.Fatal("COSMIC_PROJECT_NETWORK must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_ZONE"); v == "" {
		t.Fatal("COSMIC_ZONE must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_PRIVNW_CIDR"); v == "" {
		t.Fatal("COSMIC_PRIVNW_CIDR must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_PRIVNW_OFFERING"); v == "" {
		t.Fatal("COSMIC_PRIVNW_OFFERING must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_PRIVGW_IPADDRESS"); v == "" {
		t.Fatal("COSMIC_PRIVGW_IPADDRESS must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_STATIC_ROUTE_CIDR"); v == "" {
		t.Fatal("COSMIC_STATIC_ROUTE_CIDR must be set for acceptance tests")
	}
	if v := os.Getenv("COSMIC_STATIC_ROUTE_NEXTHOP"); v == "" {
		t.Fatal("COSMIC_STATIC_ROUTE_NEXTHOP must be set for acceptance tests")
	}
}

// Name of a valid disk offering
var COSMIC_DISK_OFFERING_1 = os.Getenv("COSMIC_DISK_OFFERING_1")

// Name of a disk offering that COSMIC_DISK_OFFERING_1 can resize to
var COSMIC_DISK_OFFERING_2 = os.Getenv("COSMIC_DISK_OFFERING_2")

// Name of a valid service offering
var COSMIC_SERVICE_OFFERING_1 = os.Getenv("COSMIC_SERVICE_OFFERING_1")

// Name of a service offering that COSMIC_SERVICE_OFFERING_1 can resize to
var COSMIC_SERVICE_OFFERING_2 = os.Getenv("COSMIC_SERVICE_OFFERING_2")

// Name of a network that already exists
var COSMIC_NETWORK_1 = os.Getenv("COSMIC_NETWORK_1")

// A valid IP address in COSMIC_NETWORK_1
var COSMIC_NETWORK_1_IPADDRESS1 = os.Getenv("COSMIC_NETWORK_1_IPADDRESS1")

// A valid IP address in COSMIC_NETWORK_1
var COSMIC_NETWORK_1_IPADDRESS2 = os.Getenv("COSMIC_NETWORK_1_IPADDRESS2")

// Name for a network that will be created
var COSMIC_NETWORK_2 = os.Getenv("COSMIC_NETWORK_2")

// Any range
var COSMIC_NETWORK_2_CIDR = os.Getenv("COSMIC_NETWORK_2_CIDR")

// An IP address in COSMIC_NETWORK_2_CIDR to be used as the gateway
var COSMIC_NETWORK_2_GATEWAY = os.Getenv("COSMIC_NETWORK_2_GATEWAY")

// Name of an available network offering with specifyvlan=false
var COSMIC_NETWORK_2_OFFERING = os.Getenv("COSMIC_NETWORK_2_OFFERING")

// An IP address in COSMIC_NETWORK_2_CIDR
var COSMIC_NETWORK_2_IPADDRESS = os.Getenv("COSMIC_NETWORK_2_IPADDRESS")

// A network that already exists and isn't COSMIC_NETWORK_1
var COSMIC_2ND_NIC_NETWORK = os.Getenv("COSMIC_2ND_NIC_NETWORK")

// An IP address in COSMIC_2ND_NIC_NETWORK
var COSMIC_2ND_NIC_IPADDRESS = os.Getenv("COSMIC_2ND_NIC_IPADDRESS")

// Any range
var COSMIC_VPC_CIDR_1 = os.Getenv("COSMIC_VPC_CIDR_1")

// Any range that doesn't overlap to COSMIC_VPC_CIDR_1, will be VPNed
var COSMIC_VPC_CIDR_2 = os.Getenv("COSMIC_VPC_CIDR_2")

// An available VPC offering
var COSMIC_VPC_OFFERING = os.Getenv("COSMIC_VPC_OFFERING")

// A sub-range of COSMIC_VPC_CIDR_1 with same starting point
var COSMIC_VPC_NETWORK_CIDR = os.Getenv("COSMIC_VPC_NETWORK_CIDR")

// A sub-range of COSMIC_VPC_CIDR_1 with same starting point
var COSMIC_VPC_NETWORK_GATEWAY = os.Getenv("COSMIC_VPC_NETWORK_GATEWAY")

// Name of an available network offering with forvpc=true
var COSMIC_VPC_NETWORK_OFFERING = os.Getenv("COSMIC_VPC_NETWORK_OFFERING")

// An IP address in the range of the COSMIC_VPC_NETWORK_CIDR
var COSMIC_VPC_NETWORK_IPADDRESS = os.Getenv("COSMIC_VPC_NETWORK_IPADDRESS")

// Path to a public IP that exists for COSMIC_NETWORK_1
var COSMIC_PUBLIC_IPADDRESS = os.Getenv("COSMIC_PUBLIC_IPADDRESS")

// Path to a public key on local disk
var COSMIC_SSH_PUBLIC_KEY = os.Getenv("COSMIC_SSH_PUBLIC_KEY")

// Name of a template that exists already for building VMs
var COSMIC_TEMPLATE = os.Getenv("COSMIC_TEMPLATE")

// Details of a template that will be added
var COSMIC_TEMPLATE_FORMAT = os.Getenv("COSMIC_TEMPLATE_FORMAT")
var COSMIC_HYPERVISOR = os.Getenv("COSMIC_HYPERVISOR")
var COSMIC_TEMPLATE_URL = os.Getenv("COSMIC_TEMPLATE_URL")
var COSMIC_TEMPLATE_OS_TYPE = os.Getenv("COSMIC_TEMPLATE_OS_TYPE")

// Name of a project that exists already
var COSMIC_PROJECT_NAME = os.Getenv("COSMIC_PROJECT_NAME")

// Name of a network that exists already in COSMIC_PROJECT_NAME
var COSMIC_PROJECT_NETWORK = os.Getenv("COSMIC_PROJECT_NETWORK")

// Name of a zone that exists already
var COSMIC_ZONE = os.Getenv("COSMIC_ZONE")

// Details of the private gateway that will be created
var COSMIC_PRIVNW_CIDR = os.Getenv("COSMIC_PRIVNW_CIDR")
var COSMIC_PRIVNW_OFFERING = os.Getenv("COSMIC_PRIVNW_OFFERING")
var COSMIC_PRIVGW_IPADDRESS = os.Getenv("COSMIC_PRIVGW_IPADDRESS")

// Details of the static route that will be added to private gateway testing this.
// nexthop should be in COSMIC_PRIVNW_CIDR
var COSMIC_STATIC_ROUTE_CIDR = os.Getenv("COSMIC_STATIC_ROUTE_CIDR")
var COSMIC_STATIC_ROUTE_NEXTHOP = os.Getenv("COSMIC_STATIC_ROUTE_NEXTHOP")
