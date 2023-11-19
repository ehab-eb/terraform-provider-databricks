package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/databricks/databricks-sdk-go"
	"github.com/databricks/terraform-provider-databricks/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func matchesIP(s *instanceProfileData, re *regexp.Regexp, field string) bool {
	var x map[string]interface{}
	m, _ := json.Marshal(s)
	_ = json.Unmarshal(m, &x)
	f := fmt.Sprint(x[field])
	return re.Match([]byte(regexp.QuoteMeta(f)))
}

type instanceProfileData struct {
	Name    string `json:"name,omitempty" tf:"computed"`
	Arn     string `json:"arn,omitempty" tf:"computed"`
	RoleArn string `json:"role_arn,omitempty" tf:"computed"`
	IsMeta  bool   `json:"is_meta,omitempty" tf:"computed"`
}

func DataSourceInstanceProfiles() *schema.Resource {
	type instanceProfileFilter struct {
		Name    string `json:"name,omitempty" tf:"required"`
		Pattern string `json:"pattern,omitempty" tf:"required"`
	}
	return common.WorkspaceData(func(ctx context.Context, data *struct {
		InstanceProfiles []instanceProfileData `json:"instance_profiles,omitempty" tf:"computed"`
		Filter           instanceProfileFilter `json:"filter,omitempty"`
	}, w *databricks.WorkspaceClient) error {
		instanceProfiles, err := w.InstanceProfiles.ListAll(ctx)
		if err != nil {
			return err
		}

		var InstanceProfilesData []instanceProfileData
		for _, v := range instanceProfiles {
			arnSlices := strings.Split(v.InstanceProfileArn, "/")
			name := arnSlices[len(arnSlices)-1]
			ipData := instanceProfileData{
				Name:    name,
				Arn:     v.InstanceProfileArn,
				RoleArn: v.IamRoleArn,
				IsMeta:  v.IsMetaInstanceProfile,
			}
			if data.Filter != (instanceProfileFilter{}) {
				re := regexp.MustCompile(data.Filter.Pattern)
				if !matchesIP(&ipData, re, data.Filter.Name) {
					continue
				}
			}
			InstanceProfilesData = append(InstanceProfilesData, ipData)
		}

		data.InstanceProfiles = InstanceProfilesData
		return nil
	})
}
