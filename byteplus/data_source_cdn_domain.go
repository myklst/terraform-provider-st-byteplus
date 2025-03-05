package byteplus

import (
	"context"
	"fmt"
	"time"

	byteplusCdnClient "github.com/byteplus-sdk/byteplus-sdk-golang/service/cdn"
	"github.com/cenkalti/backoff/v4"

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

func NewCdnDomainDataSource() datasource.DataSource {
	return &cdnDomainDataSource{}
}

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
			"Missing CDN Domain Name!",
			"Domain name must not be empty.",
		)
		return
	}

	pageNum := int64(1)
	pageSize := int64(100)

	// Create the request
	ListCdnDomainsRequest := &byteplusCdnClient.ListCdnDomainsRequest{
		Domain:   byteplusCdnClient.GetStrPtr(fmt.Sprintf("^%s$", domainName)),
		PageNum:  &pageNum,
		PageSize: &pageSize,
	}

	var response *byteplusCdnClient.ListCdnDomainsResponse
	var err error

	describeCdnDomain := func() (err error) {
		// Call the API
		// Paging handling not needed, because it will always only output 1 CDN domain.
		response, err = d.client.ListCdnDomains(ListCdnDomainsRequest)
		if err != nil {
			if byteErr, ok := err.(byteplusCdnClient.CDNError); ok {
				errCode := byteErr.Code
				if isPermanentCommonError(errCode) || isPermanentCdnError(errCode) {
					return backoff.Permanent(fmt.Errorf("err:\n%s", byteErr))
				}

				return fmt.Errorf("err:\n%s", errCode)
			}
		}

		return
	}

	// Retry with backoff
	reconnectBackoff := backoff.NewExponentialBackOff()
	reconnectBackoff.MaxElapsedTime = 30 * time.Second
	err = backoff.Retry(describeCdnDomain, reconnectBackoff)
	if err != nil {
		resp.Diagnostics.AddError(
			"[API ERROR] Failed to Describe CDN Domain.",
			err.Error(),
		)
		return
	}

	switch {
	case len(response.Result.Data) == 1:
		cdnDomain := response.Result.Data[0]
		state.Domain = types.StringValue(cdnDomain.Domain)
		state.Cname = types.StringValue(cdnDomain.Cname)
		state.Status = types.StringValue(cdnDomain.Status)
	case len(response.Result.Data) > 1:
		// We will only expect 1 result as there will be no repeating CDN domain
		// listed with 'domain name' filter.
		resp.Diagnostics.AddError(
			"[API ERROR] Multiple CDN Domain Found.",
			"Multiple CDN domain found with domain name: "+domainName+".",
		)
	default:
		// If not found, return null to avoid error in data source.
		state.Domain = types.StringNull()
		state.Cname = types.StringNull()
		state.Status = types.StringNull()
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
