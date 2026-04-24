package idaas

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/crc32"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

func dataSourceGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGroupsRead,
		Schema: map[string]*schema.Schema{
			"q": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Searches the name property of groups for matching value",
				ConflictsWith: []string{"search"},
			},
			"search": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Searches for groups with a supported filtering expression for all attributes except for '_embedded', '_links', and 'objectClass'",
				ConflictsWith: []string{"type", "q"},
			},
			"type": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Type of the group. When specified in the terraform resource, will act as a filter when searching for the groups",
				ConflictsWith: []string{"search"},
			},
			"limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "The maximum number of groups returned by the Okta API, between 1 and 10000.",
				Default:      utils.DefaultPaginationLimit,
				ValidateFunc: validation.IntBetween(1, 10000),
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Group ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Group name.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Group type, either 'APP_GROUP' or 'OKTA_GROUP'.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Group description.",
						},
						"source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the application the Group is sourced/imported from (only present for groups of type APP_GROUP).",
						},
						"custom_profile_attributes": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Raw JSON containing all custom profile attributes. Likely only useful on groups of type",
						},
					},
				},
			},
		},
		Description: "Get a list of groups from Okta.",
	}
}

func dataSourceGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getOktaV5ClientFromMetadata(meta)

	apiRequest := client.GroupAPI.ListGroups(ctx)
	apiRequest = apiRequest.Limit(int32(utils.DefaultPaginationLimit))
	qp := &query.Params{Limit: utils.DefaultPaginationLimit}

	if groupType, ok := d.GetOk("type"); ok {
		filter := fmt.Sprintf("type eq \"%s\"", groupType.(string))
		apiRequest = apiRequest.Filter(filter)
		qp.Filter = filter
	}

	if limit, ok := d.GetOk("limit"); ok {
		// override default page_size if user specified a custom limit value
		apiRequest = apiRequest.Limit(int32(limit.(int)))
	}

	if q, ok := d.GetOk("q"); ok {
		qp.Limit = 10000 // keeping this here to avoid potentially changing datasource ID generation behavior
		apiRequest = apiRequest.Q(q.(string))
		qp.Q = q.(string)
	}

	if search, ok := d.GetOk("search"); ok {
		apiRequest = apiRequest.Search(search.(string))
		qp.Search = search.(string)
	}

	okta_groups, resp, err := apiRequest.Execute()
	if err != nil {
		d.SetId("")
		return diag.Errorf("failed to list groups: %v", err)
	}

	// handle pagination
	for {
		if !resp.HasNextPage() {
			break
		}
		var moreGroups []v5okta.Group
		var err error
		resp, err = resp.Next(&moreGroups)
		if err != nil {
			return diag.Errorf("failed to get next page of groups: %v", err)
		}
		okta_groups = append(okta_groups, moreGroups...)

	}

	// generate a unique ID for the data source based on the query parameters
	dataSourceId := fmt.Sprintf("%d", crc32.ChecksumIEEE([]byte(qp.String())))
	d.SetId(dataSourceId)

	// convert the groups to a list of maps
	arr := make([]map[string]interface{}, len(okta_groups))
	for i := range okta_groups {
		arr[i] = map[string]interface{}{}

		arr[i]["id"] = okta_groups[i].Id
		arr[i]["name"] = okta_groups[i].Profile.Name
		arr[i]["type"] = okta_groups[i].Type
		arr[i]["description"] = okta_groups[i].Profile.Description

		additionalProperties := okta_groups[i].AdditionalProperties
		src, ok := additionalProperties["source"].(map[string]any)
		if ok && src["id"] != nil {
			arr[i]["source"] = src["id"].(string)
		}

		delete(additionalProperties, "source")

		// Use Profile.AdditionalProperties for custom profile attributes
		customProfileMap := okta_groups[i].Profile.AdditionalProperties
		customProfile, err := json.Marshal(customProfileMap)
		if err != nil {
			return diag.Errorf("failed to read custom profile attributes from group ID: %s", *okta_groups[i].Id)
		}
		arr[i]["custom_profile_attributes"] = string(customProfile)
	}
	_ = d.Set("groups", arr)
	return nil
}
