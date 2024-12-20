package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/byteplus-sdk/byteplus-sdk-golang/base"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &byteplusProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &byteplusProvider{
			version: version,
		}
	}
}

// byteplusProvider is the provider implementation.
type byteplusProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *byteplusProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "Byteplus"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *byteplusProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Byteplus provider is used to interact with the many resources supported by Byteplus. " +
			"The provider needs to be configured with the proper credentials before it can be used.",
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Description: "Region for Byteplus API.",
				Optional:    true,
			},
			"access_key": schema.StringAttribute{
				Description: "Access Key for Byteplus API.",
				Optional:    true,
			},
			"secret_key": schema.StringAttribute{
				Description: "Secret key for Byteplus API.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

// ByteplusProviderModel maps provider schema data to a Go type.
type byteplusProviderModel struct {
	Region    types.String `tfsdk:"region"`
	AccessKey types.String `tfsdk:"access_key"`
	SecretKey types.String `tfsdk:"secret_key"`
}

func (p *byteplusProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config byteplusProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Region.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("region"),
			"Unknown Byteplus region",
			"The provider cannot create the Byteplus API client as there is an unknown configuration value for the"+
				"Byteplus API region. Set the value statically in the configuration, or use the Byteplus_REGION environment variable.",
		)
	}

	if config.AccessKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("access_key"),
			"Unknown Byteplus access key",
			"The provider cannot create the Byteplus API client as there is an unknown configuration value for the"+
				"Byteplus API access key. Set the value statically in the configuration, or use the Byteplus_ACCESS_KEY environment variable.",
		)
	}

	if config.SecretKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("secret_key"),
			"Unknown Byteplus secret key",
			"The provider cannot create the Byteplus API client as there is an unknown configuration value for the"+
				"Byteplus secret key. Set the value statically in the configuration, or use the Byteplus_SECRET_KEY environment variable.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	var region, accessKey, secretKey string
	if !config.Region.IsNull() {
		region = config.Region.ValueString()
	} else {
		region = os.Getenv("Byteplus_REGION")
	}

	if !config.AccessKey.IsNull() {
		accessKey = config.AccessKey.ValueString()
	} else {
		accessKey = os.Getenv("Byteplus_ACCESS_KEY")
	}

	if !config.SecretKey.IsNull() {
		secretKey = config.SecretKey.ValueString()
	} else {
		secretKey = os.Getenv("Byteplus_SECRET_KEY")
	}

	// If any of the expected configuration are missing, return
	// errors with provider-specific guidance.
	if region == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("region"),
			"Missing Byteplus API region",
			"The provider cannot create the Byteplus API client as there is a "+
				"missing or empty value for the Byteplus API region. Set the "+
				"region value in the configuration or use the Byteplus_REGION "+
				"environment variable. If either is already set, ensure the value "+
				"is not empty.",
		)
	}

	if accessKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("access_key"),
			"Missing Byteplus API access key",
			"The provider cannot create the Byteplus API client as there is a "+
				"missing or empty value for the Byteplus API access key. Set the "+
				"access key value in the configuration or use the Byteplus_ACCESS_KEY "+
				"environment variable. If either is already set, ensure the value "+
				"is not empty.",
		)
	}

	if secretKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("secret_key"),
			"Missing Byteplus secret key",
			"The provider cannot create the Byteplus API client as there is a "+
				"missing or empty value for the Byteplus API Secret Key. Set the "+
				"secret key value in the configuration or use the Byteplus_SECRET_KEY "+
				"environment variable. If either is already set, ensure the value "+
				"is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new Byteplus client using the configuration values
	client := base.NewClient(&base.ServiceInfo{
		Credentials: base.Credentials{
			AccessKeyID:     accessKey,
			SecretAccessKey: secretKey,
		},
	}, nil)

	// Validate that the client was created successfully (depending on its implementation)
	if client == nil {
		resp.Diagnostics.AddError(
			"Unable to Create Byteplus API Client",
			"An unexpected error occurred when creating the Byteplus API client. "+
				"If the error is not clear, please contact the provider developers.",
		)
		return
	}

	// Make the Byteplus client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *byteplusProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCdnDomainDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *byteplusProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
