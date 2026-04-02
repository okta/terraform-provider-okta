package idaas

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceUserFactorQuestion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserFactorQuestionCreate,
		ReadContext:   resourceUserFactorQuestionRead,
		UpdateContext: resourceUserFactorQuestionUpdate,
		DeleteContext: resourceUserFactorQuestionDelete,
		Importer:      utils.CreateNestedResourceImporter([]string{"user_id", "id"}),
		Description:   "Creates security question factor for a user. This resource allows you to create and configure security question factor for a user.",
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the user. Resource will be recreated when `user_id` changes.",
				ForceNew:    true,
			},
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Security question unique key. ",
			},
			"answer": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Security question answer. Note here that answer won't be set during the resource import.",
			},
			"text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Display text for security question.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the security question factor.",
			},
		},
	}
}

func resourceUserFactorQuestionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := validateQuestionKey(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	sq := buildUserFactorQuestion(d)
	_, _, err = getOktaClientFromMetadata(meta).UserFactor.EnrollFactor(ctx, d.Get("user_id").(string), sq, nil)
	if err != nil {
		return diag.Errorf("failed to enroll user question factor: %v", err)
	}
	d.SetId(sq.Id)
	return resourceUserFactorQuestionRead(ctx, d, meta)
}

func resourceUserFactorQuestionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var uf sdk.SecurityQuestionUserFactor
	_, resp, err := getOktaClientFromMetadata(meta).UserFactor.GetFactor(ctx, d.Get("user_id").(string), d.Id(), &uf)
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get user question factor: %v", err)
	}

	if uf.Id == "" {
		d.SetId("")
		return nil
	}
	_ = d.Set("status", uf.Status)
	_ = d.Set("key", uf.Profile.Question)
	_ = d.Set("text", uf.Profile.QuestionText)
	return nil
}

func resourceUserFactorQuestionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := validateQuestionKey(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	sq := &sdk.SecurityQuestionUserFactor{
		Profile: &sdk.SecurityQuestionUserFactorProfile{
			Answer:   d.Get("answer").(string),
			Question: d.Get("key").(string),
		},
	}
	_, err = getAPISupplementFromMetadata(meta).UpdateUserFactor(ctx, d.Get("user_id").(string), d.Id(), sq)
	if err != nil {
		return diag.Errorf("failed to update user question factor: %v", err)
	}
	return resourceUserFactorQuestionRead(ctx, d, meta)
}

func resourceUserFactorQuestionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resp, err := getOktaClientFromMetadata(meta).UserFactor.DeleteFactor(ctx, d.Get("user_id").(string), d.Id())
	if err != nil {
		// disabled factor can not be removed
		if resp != nil && resp.StatusCode == http.StatusBadRequest {
			return nil
		}
		return diag.Errorf("failed to delete user question factor: %v", err)
	}
	return nil
}

func buildUserFactorQuestion(d *schema.ResourceData) *sdk.SecurityQuestionUserFactor {
	return &sdk.SecurityQuestionUserFactor{
		FactorType: "question",
		Provider:   "OKTA",
		Profile: &sdk.SecurityQuestionUserFactorProfile{
			Answer:   d.Get("answer").(string),
			Question: d.Get("key").(string),
		},
	}
}

func validateQuestionKey(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	sq, _, err := getOktaClientFromMetadata(meta).UserFactor.ListSupportedSecurityQuestions(ctx, d.Get("user_id").(string))
	if err != nil {
		return fmt.Errorf("failed to list security question keys: %v", err)
	}
	keys := make([]string, len(sq))
	for i := range sq {
		if sq[i].Question == d.Get("key").(string) {
			return nil
		}
		keys[i] = sq[i].Question
	}
	return fmt.Errorf("'%s' is missing from the available questions keys, please use one of [%s]",
		d.Get("key").(string), strings.Join(keys, ","))
}
