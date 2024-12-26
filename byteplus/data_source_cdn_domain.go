package byteplus

import (
	"context"

	byteplus "github.com/byteplus-sdk/byteplus-sdk-golang/base"
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

type cdnDomainDataSourceModel struct {
	ClientConfig *clientConfig `tfsdk:"client_config"`
	DomainName   types.String  `tfsdk:"domain_name"`
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

func (d *cdnDomainDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(byteplusClients).baseClient
}

func (d *cdnDomainDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
}
