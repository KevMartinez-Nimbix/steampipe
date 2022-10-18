package gcp

import (
	"context"
	"strings"

	//"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableGcpStorageFile(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_storage_file",
		Description: "GCP Filestore",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getGcpStorageFile,
		},
		List: &plugin.ListConfig{
			Hydrate:           listGcpStorageFile,
			ShouldIgnoreError: isIgnorableError([]string{"403"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "createTime",
				Description: "The time when the instance was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "description",
				Description: "The description of the instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "etag",
				Description: "HTTP 1.1 Entity tag for the bucket.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "fileshares",
				Description: "File system shares on the instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "kmsKeyName",
				Description: "A Cloud KMS key that will be used to encrypt objects inserted into this bucket, if no encryption method is specified.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Encryption.DefaultKmsKeyName"),
			},
			{
				Name:        "labels",
				Description: "Labels that apply to this Filestore.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "name",
				Description: "The name of the Filestore instance",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "networks",
				Description: "VPC Networks to which the instance is connected",
				Type:        proto.ColumnType_JSON,
			},
			{
				//SatifiesPzs
				Name:        "satisfiesPzs",
				Description: "Output only, Reserved for future use",
				Type:        proto.ColumnType_BOOL,
			},
			{
				//state
				Name:        "state",
				Description: "The instance state",
				Type:        proto.ColumnType_STRING,
			},
			{
				//statusmsg
				Name:        "statusMessage",
				Description: "Output only, additional information about the instance",
				Type:        proto.ColumnType_STRING,
			},
			{
				//suspensionreasons
				Name:        "suspensionReasons",
				Description: "Output only, field indicates all the reasons the instance is suspended.",
				Type:        proto.ColumnType_JSON,
			},
			{
				//tier
				Name:        "tier",
				Description: "The service tier of the instance",
				Type:        proto.ColumnType_STRING,
			},
		},
	}
}

// Functions
func listGcpStorageFile(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	//Get project details
	getProjectCached := plugin.HydrateFunc(getProject).WithCache()
	projectId, err := getProjectCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)
	projectPath := strings.Split(project, "/")[0]

	//Create Service Connection
	service, err := FileStoreService(ctx, d)
	if err != nil {
		return nil, err
	}

	//stringPath := "projects/nimbix-cloud/locations/-"
	stringPath := "projects/" + projectPath + "/locations/-"
	resp, err := service.Projects.Locations.Instances.List(stringPath).Do()
	if err == nil {
		for _, instance := range resp.Instances {
			d.StreamListItem(ctx, instance)
		}
	} else {
		return nil, err
	}
	return nil, err
}

func getGcpStorageFile(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// project := projectName
	name := d.KeyColumnQuals["name"].GetStringValue()

	// Create Service Connection
	service, err := FileStoreService(ctx, d)
	if err != nil {
		return nil, err
	}

	req, err := service.Projects.Locations.Instances.Get(name).Do()
	if err != nil {
		plugin.Logger(ctx).Trace("getGcpStorageFile", "Error", err)
		return nil, err
	}

	return req, nil
}
