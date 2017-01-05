package fuze

import (
	"encoding/json"
	"strconv"
  "errors"

	fuze "github.com/coreos/container-linux-config-transpiler/config"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceFuzeConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFuzeConfigRead,

		Schema: map[string]*schema.Schema{
			"content": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"pretty_print": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"rendered": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "rendered ignition configuration",
			},
		},
	}
}

func dataSourceFuzeConfigRead(d *schema.ResourceData, meta interface{}) error {
	rendered, err := renderFuzeConfig(d)
	if err != nil {
		return err
	}

	d.Set("rendered", rendered)
	d.SetId(strconv.Itoa(hashcode.String(rendered)))
	return nil
}

func renderFuzeConfig(d *schema.ResourceData) (string, error) {
	pretty := d.Get("pretty_print").(bool)
	config := d.Get("content").(string)

	ignition1, pR := fuze.Parse([]byte(config))
	if len(pR.Entries) > 0 {
    return "", errors.New("Initial Fuze to Ignition 1.0 conversion failed:\n" + pR.String())
	}

  // Convert to the 2.0 ignition format
  ignition, cR := fuze.ConvertAs2_0_0(ignition1)
	if len(cR.Entries) > 0 {
    return "", errors.New("Conversion of Ignition 1.0 JSON to Ignition 2.0 JSON failed:\n" + cR.String())
	}

	if pretty {
		ignitionJSON, pErr := json.MarshalIndent(&ignition, "", "  ")
		return string(ignitionJSON), pErr
	}


	ignitionJSON, mErr := json.Marshal(&ignition)
	return string(ignitionJSON), mErr
}

