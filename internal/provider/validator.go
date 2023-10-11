package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type parameterValueValidate struct {
}

func parameterValueValidator() validator.String {
	return &parameterValueValidate{}
}

func (v *parameterValueValidate) Description(_ context.Context) string {
	return "value must match type and options"
}

func (v *parameterValueValidate) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *parameterValueValidate) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	var parameter Parameter

	resp.Diagnostics.Append(req.Config.Get(ctx, &parameter)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "parameterValueValidate.ValidateString", map[string]interface{}{
		"parameter": parameter,
	})

	value := req.ConfigValue

	if !parameter.Optional.ValueBool() && value.IsNull() {
		resp.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			req.Path,
			"parameter is not optional",
			value.String(),
		))
		return
	}

	resp.Diagnostics.Append(valueIsType(parameter.Type, value))
	if resp.Diagnostics.HasError() {
		return
	}

	options, diags := parameter.GetOptions(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(options) > 0 {
		valid := false
		for _, option := range options {
			resp.Diagnostics.Append(valueIsType(parameter.Type, option.Value))
			if resp.Diagnostics.HasError() {
				return
			}
			if option.Value.Equal(option.Value) {
				valid = true
				break
			}
		}
		if !valid {
			resp.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
				req.Path,
				"value is not a valid option",
				value.String(),
			))
			return
		}
	}
}

func valueIsType(t, v basetypes.StringValue) diag.Diagnostic {
	if v.IsUnknown() || t.IsUnknown() {
		return nil
	}
	typ := t.ValueString()
	value := v.ValueString()
	switch typ {
	case "number":
		_, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return diag.NewErrorDiagnostic(fmt.Sprintf("%q is not a number", value), "")
		}
	case "bool":
		_, err := strconv.ParseBool(value)
		if err != nil {
			return diag.NewErrorDiagnostic(fmt.Sprintf("%q is not a bool", value), "")
		}
	case "json":
		var items []string
		err := json.Unmarshal([]byte(value), &items)
		if err != nil {
			return diag.NewErrorDiagnostic(fmt.Sprintf("%q is not a json", value), "")
		}
	case "string":
		// Anything is a string!
	default:
		return diag.NewErrorDiagnostic(fmt.Sprintf("%q is not a valid type", typ), "")
	}
	return nil
}
