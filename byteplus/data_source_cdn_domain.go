package byteplus

import (
	"context"
	"fmt"

	byteplus "github.com/byteplus-sdk/byteplus-sdk-golang/base"
	byteplusCdn "github.com/byteplus-sdk/byteplus-sdk-golang/service/cdn"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &cdnDomainDataSource{}
	_ datasource.DataSourceWithConfigure = &cdnDomainDataSource{}
)

// NewCoffeesDataSource is a helper function to simplify the provider implementation.
func NewCdnDomainDataSource() datasource.DataSource {
	return &cdnDomainDataSource{}
}

// coffeesDataSource is the data source implementation.
type cdnDomainDataSource struct {
	client *byteplus.Client
}

// Metadata returns the data source type name.
func (d *cdnDomainDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cdn_domains"
}

// Read refreshes the Terraform state with the latest data.
func (d *cdnDomainDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state cdnDomainDataSourceModel

	// Call the API to get the list of CDN domains
	apiResp, err := byteplusCdn.DefaultInstance.ListCdnDomains(&byteplusCdn.ListCdnDomainsRequest{})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to retrieve CDN domains",
			err.Error(),
		)
		return
	}

	// Check if the API response is valid
	if apiResp == nil || apiResp.Result.Data == nil {
		resp.Diagnostics.AddError(
			"Invalid API response",
			"Received empty or nil response from the CDN API",
		)
		return
	}

	// Map API response to the state model
	for _, domain := range apiResp.Result.Data {
		domainState := domainsModel{
			Domain: types.StringValue(domain.Domain),
			CName:  types.StringValue(domain.Cname),
		}
		state.Domains = append(state.Domains, domainState)
	}

	// Set the state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *cdnDomainDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*byteplus.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *byteplus.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Schema defines the schema for the data source.
func (d *cdnDomainDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"domains": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"domain": schema.StringAttribute{
							Computed: true,
						},
						"cname": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// coffeesDataSourceModel maps the data source schema data.
type cdnDomainDataSourceModel struct {
	Domains []domainsModel `tfsdk:"domains"`
}

// coffeesModel maps coffees schema data.
type domainsModel struct {
	Domain types.String `tfsdk:"domain"`
	CName  types.String `tfsdk:"cname"`
}
