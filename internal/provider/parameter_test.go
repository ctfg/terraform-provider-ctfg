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
				  name = "example"
				  value = "example"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ctfg_parameter.parameter", "name", "example"),
					resource.TestCheckResourceAttr("ctfg_parameter.parameter", "value", "example"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "scaffolding_example.test",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"configurable_attribute", "defaulted"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
