package cosmic

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"errors"
	"github.com/go-ini/ini"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_url": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("COSMIC_API_URL", nil),
				ConflictsWith: []string{"config", "profile"},
			},

			"api_key": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("COSMIC_API_KEY", nil),
				ConflictsWith: []string{"config", "profile"},
			},

			"secret_key": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("COSMIC_SECRET_KEY", nil),
				ConflictsWith: []string{"config", "profile"},
			},

			"http_get_only": {
				Type:        schema.TypeBool,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("COSMIC_HTTP_GET_ONLY", false),
			},

			"timeout": {
				Type:        schema.TypeInt,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("COSMIC_TIMEOUT", 900),
			},

			"config": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"api_url", "api_key", "secret_key"},
			},

			"profile": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"api_url", "api_key", "secret_key"},
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"cosmic_affinity_group":       resourceCosmicAffinityGroup(),
			"cosmic_disk":                 resourceCosmicDisk(),
			"cosmic_egress_firewall":      resourceCosmicEgressFirewall(),
			"cosmic_firewall":             resourceCosmicFirewall(),
			"cosmic_instance":             resourceCosmicInstance(),
			"cosmic_ipaddress":            resourceCosmicIPAddress(),
			"cosmic_loadbalancer_rule":    resourceCosmicLoadBalancerRule(),
			"cosmic_network":              resourceCosmicNetwork(),
			"cosmic_network_acl":          resourceCosmicNetworkACL(),
			"cosmic_network_acl_rule":     resourceCosmicNetworkACLRule(),
			"cosmic_nic":                  resourceCosmicNIC(),
			"cosmic_port_forward":         resourceCosmicPortForward(),
			"cosmic_private_gateway":      resourceCosmicPrivateGateway(),
			"cosmic_secondary_ipaddress":  resourceCosmicSecondaryIPAddress(),
			"cosmic_ssh_keypair":          resourceCosmicSSHKeyPair(),
			"cosmic_static_nat":           resourceCosmicStaticNAT(),
			"cosmic_static_route":         resourceCosmicStaticRoute(),
			"cosmic_template":             resourceCosmicTemplate(),
			"cosmic_vpc":                  resourceCosmicVPC(),
			"cosmic_vpn_connection":       resourceCosmicVPNConnection(),
			"cosmic_vpn_customer_gateway": resourceCosmicVPNCustomerGateway(),
			"cosmic_vpn_gateway":          resourceCosmicVPNGateway(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	apiURL, apiURLOK := d.GetOk("api_url")
	apiKey, apiKeyOK := d.GetOk("api_key")
	secretKey, secretKeyOK := d.GetOk("secret_key")
	config, configOK := d.GetOk("config")
	profile, profileOK := d.GetOk("profile")

	switch {
	case apiURLOK, apiKeyOK, secretKeyOK:
		if !(apiURLOK && apiKeyOK && secretKeyOK) {
			return nil, errors.New("'api_url', 'api_key' and 'secret_key' should all have values")
		}
	case configOK, profileOK:
		if !(configOK && profileOK) {
			return nil, errors.New("'config' and 'profile' should both have a value")
		}
	default:
		return nil, errors.New(
			"either 'api_url', 'api_key' and 'secret_key' or 'config' and 'profile' should have values")
	}

	if configOK && profileOK {
		cfg, err := ini.Load(config.(string))
		if err != nil {
			return nil, err
		}

		section, err := cfg.GetSection(profile.(string))
		if err != nil {
			return nil, err
		}

		apiURL = section.Key("url").String()
		apiKey = section.Key("apikey").String()
		secretKey = section.Key("secretkey").String()
	}

	cfg := Config{
		APIURL:      apiURL.(string),
		APIKey:      apiKey.(string),
		SecretKey:   secretKey.(string),
		HTTPGETOnly: d.Get("http_get_only").(bool),
		Timeout:     int64(d.Get("timeout").(int)),
	}

	return cfg.NewClient()
}
