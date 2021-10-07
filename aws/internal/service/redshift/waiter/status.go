package waiter

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/aws/internal/service/redshift/finder"
	"github.com/hashicorp/terraform-provider-aws/aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	tfredshift "github.com/hashicorp/terraform-provider-aws/internal/service/redshift"
	tfredshift "github.com/hashicorp/terraform-provider-aws/internal/service/redshift"
)

func statusCluster(conn *redshift.Redshift, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := tfredshift.FindClusterByID(conn, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.ClusterStatus), nil
	}
}
