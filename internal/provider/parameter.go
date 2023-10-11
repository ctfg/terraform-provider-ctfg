package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type Option struct {
	DisplayName types.String `tfsdk:"display_name"`
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
	ID          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
	Icon        types.String `tfsdk:"icon"`
	Options     types.List   `tfsdk:"options"`
	Optional    types.Bool   `tfsdk:"optional"`
}

func (p *Parameter) GetOptions(ctx context.Context) ([]Option, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	optionObjects := make([]types.Object, 0, len(p.Options.Elements()))
	options := make([]Option, 0, len(p.Options.Elements()))

	diagnostics.Append(p.Options.ElementsAs(ctx, &optionObjects, false)...)

	if len(optionObjects) > 0 {
		valid := false

		for idx, optionObj := range optionObjects {
			diagnostics.Append(optionObj.As(ctx, &options[idx], basetypes.ObjectAsOptions{})...)
		}

		if !valid {
			diagnostics.Append(diag.NewErrorDiagnostic(fmt.Errorf("value %q is not a valid option", p.Value.ValueString()).Error(), ""))
			return nil, diagnostics
		}
	}

	return options, diagnostics
}

func (p *ParameterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_parameter"
}

func (p *ParameterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to configure editable options for workspaces.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the parameter",
				Required:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "Value of the parameter",
				Required:            true,
				Validators: []validator.String{
					parameterValueValidator(),
				},
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
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("number", "string", "bool", "json"),
				},
			},
			"icon": schema.StringAttribute{
				MarkdownDescription: "Icon of the parameter, see in https://tabler-icons.io/  \n   e.g. `database`",
				Optional:            true,
			},
			"options": schema.ListNestedAttribute{
				MarkdownDescription: "Options of the parameter",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"display_name": schema.StringAttribute{
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
							MarkdownDescription: "Icon of the option, see in https://tabler-icons.io/  \n   e.g. `database`",
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
