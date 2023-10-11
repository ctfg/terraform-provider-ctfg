//go:build tools

package tools

import (
	// Documentation generation
	_ "github.com/Khan/genqlient/generate"
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)
