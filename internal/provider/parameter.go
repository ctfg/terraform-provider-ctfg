package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type Option struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Value       types.String `tfsdk:"value"`
	Icon        types.String `tfsdk:"icon"`
}

var _ resource.Resource = &ParameterResource{}

func NewParameterResource() resource.Resource {
	return &ParameterResource{}
}

type ParameterResource struct {
}

type Parameter struct {
	Value       types.String `tfsdk:"value"`
	Name        types.String `tfsdk:"name"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
	Icon        types.String `tfsdk:"icon"`
	Options     types.List   `tfsdk:"options"`
	Optional    types.Bool   `tfsdk:"optional"`
}

func (p *ParameterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_parameter"
}

func (p *ParameterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to configure editable options for workspaces.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the parameter",
				Required:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "Value of the parameter",
				Required:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name of the parameter",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the parameter",
				Optional:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of the parameter",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("string"),
				Validators: []validator.String{
					stringvalidator.OneOf("number", "string", "bool", "json"),
				},
			},
			"icon": schema.StringAttribute{
				MarkdownDescription: "Icon of the parameter, see in https://tabler-icons.io/ \n\n e.g. `database`",
				Optional:            true,
			},
			"options": schema.ListNestedAttribute{
				MarkdownDescription: "Options of the parameter",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the option",
							Required:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Description of the option",
							Optional:            true,
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "Value of the option",
							Required:            true,
						},
						"icon": schema.StringAttribute{
							MarkdownDescription: "Icon of the option, see in https://tabler-icons.io/ \n\n e.g. `database`",
							Optional:            true,
						},
					},
				},
			},
			"optional": schema.BoolAttribute{
				MarkdownDescription: "Optional of the parameter",
				Optional:            true,
			},
		},
	}
}

func (p *ParameterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
}

func (p *ParameterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var parameter Parameter

	resp.Diagnostics.Append(req.Plan.Get(ctx, &parameter)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(p.verifyParameter(ctx, &parameter)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &parameter)...)
}

func (p *ParameterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (p *ParameterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var parameter Parameter

	resp.Diagnostics.Append(req.Plan.Get(ctx, &parameter)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(p.verifyParameter(ctx, &parameter)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &parameter)...)
}

func (p *ParameterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var parameter Parameter

	resp.Diagnostics.Append(req.State.Get(ctx, &parameter)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &parameter)...)
}

func (p *ParameterResource) verifyParameter(ctx context.Context, parameter *Parameter) (diagnostics diag.Diagnostics) {
	value := parameter.Value.ValueString()

	diagnostics.Append(valueIsType(parameter.Type.ValueString(), value))
	if diagnostics.HasError() {
		return
	}

	options := make([]types.Object, 0, len(parameter.Options.Elements()))

	diagnostics.Append(parameter.Options.ElementsAs(ctx, &options, false)...)

	if len(options) > 0 {
		names := map[string]interface{}{}
		valid := false

		for _, optionObj := range options {
			var option Option
			diagnostics.Append(optionObj.As(ctx, &option, basetypes.ObjectAsOptions{})...)
			_, exists := names[option.Name.ValueString()]
			if exists {
				diagnostics.Append(diag.NewErrorDiagnostic(fmt.Errorf("multiple options cannot have the same name %q", option.Name.ValueString()).Error(), ""))
			}
			diagnostics.Append(valueIsType(parameter.Type.ValueString(), option.Value.ValueString()))
			names[option.Name.ValueString()] = nil

			if option.Value.ValueString() == parameter.Value.ValueString() {
				valid = true
			}
		}

		if !valid {
			diagnostics.Append(diag.NewErrorDiagnostic(fmt.Errorf("value %q is not a valid option", parameter.Value.ValueString()).Error(), ""))
			return
		}
	}

	if !parameter.Optional.ValueBool() && parameter.Value.ValueString() == "" {
		diagnostics.Append(diag.NewErrorDiagnostic("Parameter is optional and not set", ""))
	}

	return
}

func valueIsType(typ, value string) diag.Diagnostic {
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
