package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/MissionCriticalCloud/terraform-provider-cosmic/cosmic"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: cosmic.Provider,
	})
}
