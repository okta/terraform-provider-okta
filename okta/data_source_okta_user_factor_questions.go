package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUserSecurityQuestions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserSecurityQuestionsQuestionsRead,
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of a Okta User",
			},
			"questions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"text": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceUserSecurityQuestionsQuestionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sq, _, err := getOktaClientFromMetadata(m).UserFactor.ListSupportedSecurityQuestions(ctx, d.Get("user_id").(string))
	if err != nil {
		return diag.Errorf("failed to list security questions for '%s' user: %v", d.Get("user_id").(string), err)
	}
	arr := make([]map[string]interface{}, len(sq))
	for i := range sq {
		arr[i] = map[string]interface{}{
			"key":  sq[i].Question,
			"text": sq[i].QuestionText,
		}
	}
	d.SetId(d.Get("user_id").(string))
	_ = d.Set("questions", arr)
	return nil
}
