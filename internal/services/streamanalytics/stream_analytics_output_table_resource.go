package streamanalytics

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/preview/streamanalytics/mgmt/2020-03-01-preview/streamanalytics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/parse"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)


type OutputTableResource struct {
}

//var _ sdk.ResourceWithUpdate = OutputTableResource{}

var _ sdk.ResourceWithCustomImporter = OutputTableResource{}

type OutputTableResourceModel struct {
	Name					string		`tfschema:"name"`
	StreamAnalyticsJob		string		`tfschema:"stream_analytics_job_name"`
	ResourceGroup			string		`tfschema:"resource_group_name"`
	StorageAccount			string		`tfschema:"storage_account_name"`
	StorageAccountKey		string		`tfschema:"storage_account_key"`
	Table					string		`tfschema:"table"`
	PartitionKey			string		`tfschema:"partition_key"`
	RowKey					string		`tfschema:"row_key"`
	BatchSize				int32		`tfschema:"batch_size"`
}

func (r OutputTableResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type: pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"stream_analytics_job_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": azure.SchemaResourceGroupName(),

		"storage_account_name": {
			Type:         pluginsdk.TypeString,
			Required: 	  true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"storage_account_key": {
			Type:         pluginsdk.TypeString,
			Required: 	  true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"table": {
			Type:         pluginsdk.TypeString,
			Required: 	  true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"partition_key": {
			Type:         pluginsdk.TypeString,
			Required: 	  true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"row_key": {
			Type:         pluginsdk.TypeString,
			Required: 	  true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"batch_size": {
			Type:         pluginsdk.TypeInt,
			Required: 	  true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r OutputTableResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r OutputTableResource) ModelObject() interface{} {
	return &OutputTableResourceModel{}
}

func (r OutputTableResource) ResourceType() string{
	return "azurerm_stream_analytics_output_table"
}

func (r OutputTableResource) IDValidationFunc() pluginsdk.SchemaValidateFunc{
	return validate.OutputID
}

func (r OutputTableResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model OutputTableResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.StreamAnalytics.OutputsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := parse.NewOutputID(subscriptionId, model.ResourceGroup, model.StreamAnalyticsJob, model.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, id.StreamingjobName, id.Name)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			props := streamanalytics.Output{
				Name: utils.String(model.Name),
				OutputProperties: &streamanalytics.OutputProperties{
					Datasource: &streamanalytics.AzureTableOutputDataSource{
						Type: streamanalytics.TypeMicrosoftStorageTable,
						AzureTableOutputDataSourceProperties: &streamanalytics.AzureTableOutputDataSourceProperties{
							AccountName: utils.String(model.StorageAccount),
							AccountKey: utils.String(model.StorageAccountKey),
							Table: utils.String(model.Table),
							PartitionKey: utils.String(model.PartitionKey),
							RowKey: utils.String(model.RowKey),
							BatchSize: utils.Int32(model.BatchSize),
						},
					},
					//Serialization: serialization,
				},
			}

			_, err = client.CreateOrReplace(ctx, props, id.ResourceGroup, id.StreamingjobName, id.Name, "", "")
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r OutputTableResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {

			return nil
		},
	}
}

//func (r OutputTableResource) Update() sdk.ResourceFunc {
//	return sdk.ResourceFunc{
//		Timeout: 5 * time.Minute,
//		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
//
//			return nil
//		},
//	}
//}

func (r OutputTableResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {

			return nil
		},
	}
}

func (r OutputTableResource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		return nil
	}
}