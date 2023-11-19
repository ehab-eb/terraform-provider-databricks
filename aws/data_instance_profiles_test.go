package aws

import (
	"testing"

	"github.com/databricks/databricks-sdk-go/service/compute"
	"github.com/databricks/terraform-provider-databricks/qa"
)

var testResponse compute.ListInstanceProfilesResponse = compute.ListInstanceProfilesResponse{
	InstanceProfiles: []compute.InstanceProfile{
		{
			IamRoleArn:            "arn:aws:iam::123456789012:role/S3Access",
			InstanceProfileArn:    "arn:aws:iam::123456789012:instance-profile/S3Access",
			IsMetaInstanceProfile: true,
		},
		{
			IamRoleArn:            "arn:aws:iam::123456789098:role/KMSAccess",
			InstanceProfileArn:    "arn:aws:iam::123456789098:instance-profile/KMSAccess",
			IsMetaInstanceProfile: false,
		},
		{
			InstanceProfileArn:    "arn:aws:iam::123456789098:instance-profile/different",
			IamRoleArn:            "arn:aws:iam::123456789098:role/value",
			IsMetaInstanceProfile: false,
		},
		{
			IamRoleArn:            "arn:aws:iam::123456789098:role/Accesses",
			InstanceProfileArn:    "arn:aws:iam::123456789098:instance-profile/Accesses",
			IsMetaInstanceProfile: false,
		},
	},
}

func TestInstanceProfilesData(t *testing.T) {
	qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/instance-profiles/list",
				Response: testResponse,
			},
		},
		Resource:    DataSourceInstanceProfiles(),
		Read:        true,
		NonWritable: true,
		ID:          "_",
	}.ApplyAndExpectData(t, map[string]any{
		"instance_profiles": []interface{}{
			map[string]interface{}{
				"name":     "S3Access",
				"arn":      "arn:aws:iam::123456789012:instance-profile/S3Access",
				"role_arn": "arn:aws:iam::123456789012:role/S3Access",
				"is_meta":  true,
			},
			map[string]interface{}{
				"name":     "KMSAccess",
				"arn":      "arn:aws:iam::123456789098:instance-profile/KMSAccess",
				"role_arn": "arn:aws:iam::123456789098:role/KMSAccess",
				"is_meta":  false,
			},
			map[string]interface{}{
				"name":     "different",
				"arn":      "arn:aws:iam::123456789098:instance-profile/different",
				"role_arn": "arn:aws:iam::123456789098:role/value",
				"is_meta":  false,
			},
			map[string]interface{}{
				"name":     "Accesses",
				"arn":      "arn:aws:iam::123456789098:instance-profile/Accesses",
				"role_arn": "arn:aws:iam::123456789098:role/Accesses",
				"is_meta":  false,
			},
		},
	})
}

func TestInstanceProfilesDataEmpty(t *testing.T) {
	qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/instance-profiles/list",
				Response: compute.ListInstanceProfilesResponse{
					InstanceProfiles: []compute.InstanceProfile{},
				},
			},
		},
		Resource:    DataSourceInstanceProfiles(),
		Read:        true,
		NonWritable: true,
		ID:          "_",
	}.ApplyAndExpectData(t, map[string]any{
		"instance_profiles": []interface{}{},
	})
}

func TestInstanceProfilesDataDuplicate(t *testing.T) {
	qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/instance-profiles/list",
				Response: compute.ListInstanceProfilesResponse{
					InstanceProfiles: []compute.InstanceProfile{
						{
							IamRoleArn:            "arn:aws:iam::123456789012:role/S3Access",
							InstanceProfileArn:    "arn:aws:iam::123456789012:instance-profile/S3Access",
							IsMetaInstanceProfile: true,
						},
						{
							IamRoleArn:            "arn:aws:iam::123456789012:role/S3Access",
							InstanceProfileArn:    "arn:aws:iam::123456789012:instance-profile/S3Access",
							IsMetaInstanceProfile: true,
						},
					},
				},
			},
		},
		Resource:    DataSourceInstanceProfiles(),
		Read:        true,
		NonWritable: true,
		ID:          "_",
	}.ApplyAndExpectData(t, map[string]any{
		"instance_profiles": []interface{}{
			map[string]interface{}{
				"name":     "S3Access",
				"arn":      "arn:aws:iam::123456789012:instance-profile/S3Access",
				"role_arn": "arn:aws:iam::123456789012:role/S3Access",
				"is_meta":  true,
			},
			map[string]interface{}{
				"name":     "S3Access",
				"arn":      "arn:aws:iam::123456789012:instance-profile/S3Access",
				"role_arn": "arn:aws:iam::123456789012:role/S3Access",
				"is_meta":  true,
			},
		},
	})
}

func TestInstanceProfilesDataFilterContains(t *testing.T) {
	qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/instance-profiles/list",
				Response: testResponse,
			},
		},
		Resource:    DataSourceInstanceProfiles(),
		Read:        true,
		NonWritable: true,
		ID:          "_",
		HCL: `
		filter {
			name = "name"
			pattern = "Access"
		}`,
	}.ApplyAndExpectData(t, map[string]any{
		"instance_profiles": []interface{}{
			map[string]interface{}{
				"name":     "S3Access",
				"arn":      "arn:aws:iam::123456789012:instance-profile/S3Access",
				"role_arn": "arn:aws:iam::123456789012:role/S3Access",
				"is_meta":  true,
			},
			map[string]interface{}{
				"name":     "KMSAccess",
				"arn":      "arn:aws:iam::123456789098:instance-profile/KMSAccess",
				"role_arn": "arn:aws:iam::123456789098:role/KMSAccess",
				"is_meta":  false,
			},
			map[string]interface{}{
				"name":     "Accesses",
				"arn":      "arn:aws:iam::123456789098:instance-profile/Accesses",
				"role_arn": "arn:aws:iam::123456789098:role/Accesses",
				"is_meta":  false,
			},
		},
	})
}

func TestInstanceProfilesDataFilterEqualsEmpty(t *testing.T) {
	qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/instance-profiles/list",
				Response: testResponse,
			},
		},
		Resource:    DataSourceInstanceProfiles(),
		Read:        true,
		NonWritable: true,
		ID:          "_",
		HCL: `
		filter {
			name = "name"
			pattern = "^Access$"
		}`,
	}.ApplyAndExpectData(t, map[string]any{
		"instance_profiles": []interface{}{},
	})
}

func TestInstanceProfilesDataFilterEquals(t *testing.T) {
	qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/instance-profiles/list",
				Response: testResponse,
			},
		},
		Resource:    DataSourceInstanceProfiles(),
		Read:        true,
		NonWritable: true,
		ID:          "_",
		HCL: `
		filter {
			name = "name"
			pattern = "^KMSAccess$"
		}`,
	}.ApplyAndExpectData(t, map[string]any{
		"instance_profiles": []interface{}{
			map[string]interface{}{
				"name":     "KMSAccess",
				"arn":      "arn:aws:iam::123456789098:instance-profile/KMSAccess",
				"role_arn": "arn:aws:iam::123456789098:role/KMSAccess",
				"is_meta":  false,
			},
		},
	})
}

func TestInstanceProfilesDataFilterCase(t *testing.T) {
	qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/instance-profiles/list",
				Response: testResponse,
			},
		},
		Resource:    DataSourceInstanceProfiles(),
		Read:        true,
		NonWritable: true,
		ID:          "_",
		HCL: `
		filter {
			name = "role_arn"
			pattern = "kms"
		}`,
	}.ApplyAndExpectData(t, map[string]any{
		"instance_profiles": []interface{}{},
	})
}

func TestInstanceProfilesDataFilterEndsWith(t *testing.T) {
	qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/instance-profiles/list",
				Response: testResponse,
			},
		},
		Resource:    DataSourceInstanceProfiles(),
		Read:        true,
		NonWritable: true,
		ID:          "_",
		HCL: `
		filter {
			name = "name"
			pattern = "Access$"
		}`,
	}.ApplyAndExpectData(t, map[string]any{
		"instance_profiles": []interface{}{
			map[string]interface{}{
				"name":     "S3Access",
				"arn":      "arn:aws:iam::123456789012:instance-profile/S3Access",
				"role_arn": "arn:aws:iam::123456789012:role/S3Access",
				"is_meta":  true,
			},
			map[string]interface{}{
				"name":     "KMSAccess",
				"arn":      "arn:aws:iam::123456789098:instance-profile/KMSAccess",
				"role_arn": "arn:aws:iam::123456789098:role/KMSAccess",
				"is_meta":  false,
			},
		},
	})
}

func TestInstanceProfilesDataFilterName(t *testing.T) {
	qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/instance-profiles/list",
				Response: testResponse,
			},
		},
		Resource:    DataSourceInstanceProfiles(),
		Read:        true,
		NonWritable: true,
		ID:          "_",
		HCL: `
		filter {
			name = "name"
			pattern = "Access$"
		}`,
	}.ApplyAndExpectData(t, map[string]any{
		"instance_profiles": []interface{}{
			map[string]interface{}{
				"name":     "S3Access",
				"arn":      "arn:aws:iam::123456789012:instance-profile/S3Access",
				"role_arn": "arn:aws:iam::123456789012:role/S3Access",
				"is_meta":  true,
			},
			map[string]interface{}{
				"name":     "KMSAccess",
				"arn":      "arn:aws:iam::123456789098:instance-profile/KMSAccess",
				"role_arn": "arn:aws:iam::123456789098:role/KMSAccess",
				"is_meta":  false,
			},
		},
	})
}
func TestInstanceProfilesDataFilterArn(t *testing.T) {
	qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/instance-profiles/list",
				Response: testResponse,
			},
		},
		Resource:    DataSourceInstanceProfiles(),
		Read:        true,
		NonWritable: true,
		ID:          "_",
		HCL: `
		filter {
			name = "arn"
			pattern = "arn:aws:iam::123456789098:instance-profile/different"
		}`,
	}.ApplyAndExpectData(t, map[string]any{
		"instance_profiles": []interface{}{
			map[string]interface{}{
				"name":     "different",
				"arn":      "arn:aws:iam::123456789098:instance-profile/different",
				"role_arn": "arn:aws:iam::123456789098:role/value",
				"is_meta":  false,
			},
		},
	})
}

func TestInstanceProfilesDataFilterRoleArn(t *testing.T) {
	qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/instance-profiles/list",
				Response: testResponse,
			},
		},
		Resource:    DataSourceInstanceProfiles(),
		Read:        true,
		NonWritable: true,
		ID:          "_",
		HCL: `
		filter {
			name = "role_arn"
			pattern = "KMS"
		}`,
	}.ApplyAndExpectData(t, map[string]any{
		"instance_profiles": []interface{}{
			map[string]interface{}{
				"name":     "KMSAccess",
				"arn":      "arn:aws:iam::123456789098:instance-profile/KMSAccess",
				"role_arn": "arn:aws:iam::123456789098:role/KMSAccess",
				"is_meta":  false,
			},
		},
	})
}

func TestInstanceProfilesDataFilterIsMeta(t *testing.T) {
	qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/instance-profiles/list",
				Response: testResponse,
			},
		},
		Resource:    DataSourceInstanceProfiles(),
		Read:        true,
		NonWritable: true,
		ID:          "_",
		HCL: `
		filter {
			name = "is_meta"
			pattern = "true"
		}`,
	}.ApplyAndExpectData(t, map[string]any{
		"instance_profiles": []interface{}{
			map[string]interface{}{
				"name":     "S3Access",
				"arn":      "arn:aws:iam::123456789012:instance-profile/S3Access",
				"role_arn": "arn:aws:iam::123456789012:role/S3Access",
				"is_meta":  true,
			},
		},
	})
}

func TestInstanceProfilesDataFilterBadName(t *testing.T) {
	qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/instance-profiles/list",
				Response: testResponse,
			},
		},
		Resource:    DataSourceInstanceProfiles(),
		Read:        true,
		NonWritable: true,
		ID:          "_",
		HCL: `
		filter {
			name = "does_not_exist"
			pattern = "value"
		}`,
	}.ApplyAndExpectData(t, map[string]any{})
}

func TestInstanceProfilesDataFilterNameOnly(t *testing.T) {
	qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/instance-profiles/list",
				Response: testResponse,
			},
		},
		Resource:    DataSourceInstanceProfiles(),
		Read:        true,
		NonWritable: true,
		ID:          "_",
		HCL: `
		filter {
			pattern = "does_not_exist"
		}`,
	}.ApplyAndExpectData(t, map[string]any{})
}
