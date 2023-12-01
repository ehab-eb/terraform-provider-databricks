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
	v, ok := x[field]
	if !ok {
		panic(fmt.Sprintf("no field `%s` found", field))
	}
	f := fmt.Sprint(v)
	return re.Match([]byte(regexp.QuoteMeta(f)))
}

type instanceProfileData struct {
	Name    string `json:"name"`
	Arn     string `json:"arn"`
	RoleArn string `json:"role_arn"`
	IsMeta  bool   `json:"is_meta"`
}

func DataSourceInstanceProfiles() *schema.Resource {
	type instanceProfileFilter struct {
		Name    string `json:"name"`
		Pattern string `json:"pattern,omitempty"`
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
