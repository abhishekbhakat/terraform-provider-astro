package models

import (
	"context"

	"github.com/astronomer/astronomer-terraform-provider/internal/clients/platform"
	"github.com/astronomer/astronomer-terraform-provider/internal/provider/schemas"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ClustersDataSource describes the data source data model.
type ClustersDataSource struct {
	Clusters      types.List   `tfsdk:"clusters"`
	CloudProvider types.String `tfsdk:"cloud_provider"` // query parameter
	Names         types.List   `tfsdk:"names"`          // query parameter
}

func (data *ClustersDataSource) ReadFromResponse(
	ctx context.Context,
	clusters []platform.Cluster,
) diag.Diagnostics {
	values := make([]attr.Value, len(clusters))
	for i, deployment := range clusters {
		var singleClusterData Cluster
		diags := singleClusterData.ReadFromResponse(ctx, &deployment)
		if diags.HasError() {
			return diags
		}

		objectValue, diags := types.ObjectValueFrom(ctx, schemas.ClustersElementAttributeTypes(), singleClusterData)
		if diags.HasError() {
			return diags
		}
		values[i] = objectValue
	}
	var diags diag.Diagnostics
	data.Clusters, diags = types.ListValue(types.ObjectType{AttrTypes: schemas.ClustersElementAttributeTypes()}, values)
	if diags.HasError() {
		return diags
	}

	return nil
}