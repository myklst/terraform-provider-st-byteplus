package byteplus

import (
	"context"

	byteplusCdnClient "github.com/byteplus-sdk/byteplus-sdk-golang/service/cdn"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
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
	client *byteplusCdnClient.CDN
}

type cdnDomainDataSourceModel struct {
	ClientConfig *clientConfig `tfsdk:"client_config"`
	Domain       types.String  `tfsdk:"domain_name"`
	Cname        types.String  `tfsdk:"cname"`
	Status       types.String  `tfsdk:"status"`
}

// Metadata returns the data source type name.
func (d *cdnDomainDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cdn_domain"
}

func (d *cdnDomainDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source provides the CDN Instance of the current Byteplus user.",
		Attributes: map[string]schema.Attribute{
			"domain_name": schema.StringAttribute{
				Description: "Domain name of CDN domain.",
				Required:    true,
			},
			"cname": schema.StringAttribute{
				Description: "Domain CName of CDN domain.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Status of CDN domain.",
				Computed:    true,
			},
		},

		Blocks: map[string]schema.Block{
			"client_config": schema.SingleNestedBlock{
				Description: "Config to override default client created in Provider. " +
					"This block will not be recorded in state file.",
				Attributes: map[string]schema.Attribute{
					"region": schema.StringAttribute{
						Description: "The region of the CDN domains. Default to " +
							"use region configured in the provider.",
						Optional: true,
					},
					"access_key": schema.StringAttribute{
						Description: "The access key that have permissions to list " +
							"CDN domains. Default to use access key configured in " +
							"the provider.",
						Optional: true,
					},
					"secret_key": schema.StringAttribute{
						Description: "The secret key that have permissions to list " +
							"CDN domains. Default to use secret key configured in " +
							"the provider.",
						Optional: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *cdnDomainDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(byteplusClients).cdnClient
}

func (d *cdnDomainDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan, state cdnDomainDataSourceModel
	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ClientConfig == nil {
		plan.ClientConfig = &clientConfig{}
	}

	initClient, clientCredentialsConfig, initClientDiags := initNewClient(d.client.Client, plan.ClientConfig)
	if initClientDiags.HasError() {
		resp.Diagnostics.Append(initClientDiags...)
		return
	}

	if initClient {
		d.client = byteplusCdnClient.NewInstance()
		d.client.Client.SetCredential(*clientCredentialsConfig)
	}

	domainName := plan.Domain.ValueString()

	if domainName == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("domain_name"),
			"Missing CDN Domain Name",
			"Domain name must not be empty",
		)
		return
	}

	ListCdnDomainsRequest := &byteplusCdnClient.ListCdnDomainsRequest{
		Domain: &domainName,
	}

	// Call the API
	response, err := d.client.ListCdnDomains(ListCdnDomainsRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Byteplus CDN Domains",
			err.Error(),
		)
		return
	}

	cdnDomains := response.Result.Data

	for _, cdnDomain := range cdnDomains {
		state.Domain = types.StringValue(cdnDomain.Domain)
		state.Cname = types.StringValue(cdnDomain.Cname)
		state.Status = types.StringValue(cdnDomain.Status)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
