package cosmic

import (
	"fmt"
	"log"
	"strings"

	"github.com/MissionCriticalCloud/go-cosmic/cosmic"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCosmicSSHKeyPair() *schema.Resource {
	return &schema.Resource{
		Create: resourceCosmicSSHKeyPairCreate,
		Read:   resourceCosmicSSHKeyPairRead,
		Delete: resourceCosmicSSHKeyPairDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"public_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"private_key": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"fingerprint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCosmicSSHKeyPairCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cosmic.CosmicClient)

	name := d.Get("name").(string)
	publicKey := d.Get("public_key").(string)

	if publicKey != "" {
		// Register supplied key
		p := cs.SSH.NewRegisterSSHKeyPairParams(name, publicKey)

		// If there is a project supplied, we retrieve and set the project id
		if err := setProjectid(p, cs, d); err != nil {
			return err
		}

		_, err := cs.SSH.RegisterSSHKeyPair(p)
		if err != nil {
			return err
		}
	} else {
		// No key supplied, must create one and return the private key
		p := cs.SSH.NewCreateSSHKeyPairParams(name)

		// If there is a project supplied, we retrieve and set the project id
		if err := setProjectid(p, cs, d); err != nil {
			return err
		}

		r, err := cs.SSH.CreateSSHKeyPair(p)
		if err != nil {
			return err
		}
		d.Set("private_key", r.Privatekey)
	}

	log.Printf("[DEBUG] Key pair successfully generated at Cosmic")
	d.SetId(name)

	return resourceCosmicSSHKeyPairRead(d, meta)
}

func resourceCosmicSSHKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cosmic.CosmicClient)

	log.Printf("[DEBUG] looking for key pair with name %s", d.Id())

	p := cs.SSH.NewListSSHKeyPairsParams()
	p.SetName(d.Id())

	// If there is a project supplied, we retrieve and set the project id
	if err := setProjectid(p, cs, d); err != nil {
		return err
	}

	r, err := cs.SSH.ListSSHKeyPairs(p)
	if err != nil {
		return err
	}
	if r.Count == 0 {
		log.Printf("[DEBUG] Key pair %s does not exist", d.Id())
		d.SetId("")
		return nil
	}

	//SSHKeyPair name is unique in a cosmic account so dont need to check for multiple
	d.Set("name", r.SSHKeyPairs[0].Name)
	d.Set("fingerprint", r.SSHKeyPairs[0].Fingerprint)

	return nil
}

func resourceCosmicSSHKeyPairDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cosmic.CosmicClient)

	// Create a new parameter struct
	p := cs.SSH.NewDeleteSSHKeyPairParams(d.Id())

	// If there is a project supplied, we retrieve and set the project id
	if err := setProjectid(p, cs, d); err != nil {
		return err
	}

	// Remove the SSH Keypair
	_, err := cs.SSH.DeleteSSHKeyPair(p)
	if err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"A key pair with name '%s' does not exist for account", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error deleting key pair: %s", err)
	}

	return nil
}
