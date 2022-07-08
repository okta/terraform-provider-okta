package okta

import (
	"context"
	"fmt"
	"net/http"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLinkValue() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLinkValueCreate,
		ReadContext:   resourceLinkValueRead,
		UpdateContext: resourceLinkValueUpdate,
		DeleteContext: resourceLinkValueDelete,
		Importer:      createNestedResourceImporter([]string{"primary_name", "primary_user_id"}),
		Schema: map[string]*schema.Schema{
			"primary_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the 'primary' relationship being assigned.",
			},
			"primary_user_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "User ID to be assigned to 'primary' for the 'associated' user in the specified relationship.",
			},
			"associated_user_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Set of User IDs or login values of the users to be assigned the 'associated' relationship.",
			},
		},
	}
}

func resourceLinkValueCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	lo, _, err := client.LinkedObject.GetLinkedObjectDefinition(ctx, d.Get("primary_name").(string))
	if err != nil {
		return diag.Errorf("failed to get linked object by primary name: %v", err)
	}
	if lo.Primary.Name != d.Get("primary_name").(string) {
		return diag.Errorf("primary name should be provided instead of associated one")
	}
	puID := d.Get("primary_user_id").(string)
	d.SetId(fmt.Sprintf("%s/%s", lo.Primary.Name, puID))
	associatedUsers := convertInterfaceToStringSetNullable(d.Get("associated_user_ids"))
	for _, associatedUser := range associatedUsers {
		_, err := client.User.SetLinkedObjectForUser(ctx, associatedUser, lo.Primary.Name, puID)
		if err != nil {
			return diag.Errorf("failed to set linked object value for primary name: "+
				"associatedUser: %s, primaryName: %s, primaryUser: %s, err: %v", associatedUser, lo.Primary.Name, puID, err)
		}
	}
	return nil
}

func resourceLinkValueRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	lo, resp, err := client.LinkedObject.GetLinkedObjectDefinition(ctx, d.Get("primary_name").(string))
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get linked object by primary name: %v", err)
	}
	if lo == nil {
		d.SetId("")
		return nil
	}
	if lo.Primary.Name != d.Get("primary_name").(string) {
		return diag.Errorf("primary name should be provided instead of associated one")
	}
	puID := d.Get("primary_user_id").(string)
	los, resp, err := client.User.GetLinkedObjectsForUser(ctx, puID, lo.Associated.Name, nil)
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	} else if err != nil {
		return diag.Errorf("failed to get associated linked object values: %v", err)
	}
	ids := make([]string, len(los))
	for i := range los {
		ids[i] = path.Base(linksValue(los[i].Links, "self", "href"))
	}
	_ = d.Set("associated_user_ids", convertStringSliceToSet(ids))
	return nil
}

func resourceLinkValueUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	oldUsers, newUsers := d.GetChange("associated_user_ids")
	oldSet := oldUsers.(*schema.Set)
	newSet := newUsers.(*schema.Set)
	usersToAdd := convertInterfaceArrToStringArr(newSet.Difference(oldSet).List())
	usersToRemove := convertInterfaceArrToStringArr(oldSet.Difference(newSet).List())
	for _, u := range usersToAdd {
		_, err := getOktaClientFromMetadata(m).User.SetLinkedObjectForUser(ctx, u, d.Get("primary_name").(string), d.Get("primary_user_id").(string))
		if err != nil {
			return diag.Errorf("failed to set relationship: associatedUser: %s, primaryName: %s, primaryUser: %s, "+
				"err: %v", u, d.Get("primary_name").(string), d.Get("primary_user_id").(string), err)
		}
	}
	for _, u := range usersToRemove {
		_, err := getOktaClientFromMetadata(m).User.RemoveLinkedObjectForUser(ctx, u, d.Get("primary_name").(string))
		if err != nil {
			return diag.Errorf("failed to remove relationship: associatedUser: %s, primaryName: %s, err: %v", u, d.Get("primary_name"), err)
		}
	}
	return nil
}

func resourceLinkValueDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	associatedUsers := convertInterfaceToStringSetNullable(d.Get("associated_user_ids"))
	for _, u := range associatedUsers {
		resp, err := getOktaClientFromMetadata(m).User.RemoveLinkedObjectForUser(ctx, u, d.Get("primary_name").(string))
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			continue
		}
		if err != nil {
			return diag.Errorf("failed to remove relationship: associatedUser: %s, primaryName: %s, err: %v", u, d.Get("primary_name"), err)
		}
	}
	return nil
}
