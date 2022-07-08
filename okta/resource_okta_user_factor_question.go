package okta

import (
	"context"
	"fmt"
	"net/http"
	"strings"

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
		Importer:      createNestedResourceImporter([]string{"user_id", "id"}),
		Description:   "Resource to manage a question factor for a user",
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of a Okta User",
				ForceNew:    true,
			},
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique key for question",
			},
			"answer": {
				Type:             schema.TypeString,
				Required:         true,
				Sensitive:        true,
				ValidateDiagFunc: stringLenBetween(4, 1000),
				Description:      "User password security answer",
			},
			"text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Display text for question",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User factor status.",
			},
		},
	}
}

func resourceUserFactorQuestionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := validateQuestionKey(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	sq := buildUserFactorQuestion(d)
	_, _, err = getOktaClientFromMetadata(m).UserFactor.EnrollFactor(ctx, d.Get("user_id").(string), sq, nil)

	if err != nil {
		return diag.Errorf("failed to enroll user question factor: %v", err)
	}
	d.SetId(sq.Id)
	return resourceUserFactorQuestionRead(ctx, d, m)
}

func resourceUserFactorQuestionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var uf okta.SecurityQuestionUserFactor
	_, resp, err := getOktaClientFromMetadata(m).UserFactor.GetFactor(ctx, d.Get("user_id").(string), d.Id(), &uf)
	if err := suppressErrorOn404(resp, err); err != nil {
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

func resourceUserFactorQuestionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := validateQuestionKey(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	sq := &okta.SecurityQuestionUserFactor{
		Profile: &okta.SecurityQuestionUserFactorProfile{
			Answer:   d.Get("answer").(string),
			Question: d.Get("key").(string),
		},
	}
	_, err = getSupplementFromMetadata(m).UpdateUserFactor(ctx, d.Get("user_id").(string), d.Id(), sq)
	if err != nil {
		return diag.Errorf("failed to update user question factor: %v", err)
	}
	return resourceUserFactorQuestionRead(ctx, d, m)
}

func resourceUserFactorQuestionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resp, err := getOktaClientFromMetadata(m).UserFactor.DeleteFactor(ctx, d.Get("user_id").(string), d.Id())
	if err != nil {
		// disabled factor can not be removed
		if resp != nil && resp.StatusCode == http.StatusBadRequest {
			return nil
		}
		return diag.Errorf("failed to delete user question factor: %v", err)
	}
	return nil
}

func buildUserFactorQuestion(d *schema.ResourceData) *okta.SecurityQuestionUserFactor {
	return &okta.SecurityQuestionUserFactor{
		FactorType: "question",
		Provider:   "OKTA",
		Profile: &okta.SecurityQuestionUserFactorProfile{
			Answer:   d.Get("answer").(string),
			Question: d.Get("key").(string),
		},
	}
}

func validateQuestionKey(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	sq, _, err := getOktaClientFromMetadata(m).UserFactor.ListSupportedSecurityQuestions(ctx, d.Get("user_id").(string))
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
