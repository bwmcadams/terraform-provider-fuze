package main

import (
	"github.com/bwmcadams/terraform-provider-fuze/fuze"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: fuze.Provider,
	})
}
