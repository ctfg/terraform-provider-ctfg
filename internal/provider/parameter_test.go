package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccParameter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: `
				resource "ctfg_parameter" "parameter" {
				  id = "example"
				  value = "example"
				  type = "string"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ctfg_parameter.parameter", "id", "example"),
					resource.TestCheckResourceAttr("ctfg_parameter.parameter", "value", "example"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
