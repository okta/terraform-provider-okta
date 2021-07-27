package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func resourceUserFactorQuestion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserFactorQuestionCreate,
		ReadContext:   resourceUserFactorQuestionRead,
		UpdateContext: resourceUserFactorQuestionUpdate,
		DeleteContext: resourceUserFactorQuestionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				_ = d.Set("user_id", d.Id())
				d.SetId(d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},
		Description: "Resource to manage a set of Factors for a specific user,",
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of a Okta User",
				ForceNew:    true,
			},
			"security_question_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User Password Security Question",
			},
			"security_answer": {
				Type:             schema.TypeString,
				Required:         true,
				Sensitive:        true,
				ValidateDiagFunc: stringLenBetween(4, 1000),
				Description:      "User Password Security Answer",
			},
		},
	}
}

func resourceUserFactorQuestionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	userId := d.Get("user_id").(string)
	factorProfile := okta.NewSecurityQuestionUserFactorProfile()
	factorProfile.Question = d.Get("security_question_key").(string)
	factorProfile.Answer = d.Get("security_answer").(string)
	factor := okta.NewSecurityQuestionUserFactor()
	factor.Profile = factorProfile
	responseFactor, _, err := getOktaClientFromMetadata(m).UserFactor.EnrollFactor(ctx, userId, factor, nil)
	d.SetId(responseFactor.(*okta.UserFactor).Id)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceUserFactorQuestionRead(ctx, d, m)
}

func resourceUserFactorQuestionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var uf *okta.SecurityQuestionUserFactor
	_, _, err := getOktaClientFromMetadata(m).UserFactor.GetFactor(ctx, d.Get("user_id").(string), d.Id(), uf)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceUserFactorQuestionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := resourceUserFactorQuestionDelete(ctx, d, m)
	if err != nil {
		return err
	}
	err = resourceUserFactorQuestionCreate(ctx, d, m)
	if err != nil {
		return err
	}
	return resourceUserFactorQuestionRead(ctx, d, m)
}

func resourceUserFactorQuestionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, err := getOktaClientFromMetadata(m).UserFactor.DeleteFactor(ctx, d.Get("user_id").(string), d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
