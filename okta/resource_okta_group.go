package okta

import (
	"context"
	"net/http"

	"github.com/okta/okta-sdk-golang/v2/okta"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Exists: resourceGroupExists,
		Read:   resourceGroupRead,
		Update: resourceGroupUpdate,
		Delete: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Group name",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Group description",
			},
			"users": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Users associated with the group. This can also be done per user.",
			},
		},
	}
}

func buildGroup(d *schema.ResourceData) *okta.Group {
	return &okta.Group{
		Profile: &okta.GroupProfile{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
		},
	}
}

func resourceGroupCreate(d *schema.ResourceData, m interface{}) error {
	group := buildGroup(d)
	responseGroup, _, err := getOktaClientFromMetadata(m).Group.CreateGroup(context.Background(), *group)
	if err != nil {
		return err
	}

	d.SetId(responseGroup.Id)
	if err := updateGroupUsers(d, m); err != nil {
		return err
	}

	return resourceGroupRead(d, m)
}

func resourceGroupExists(d *schema.ResourceData, m interface{}) (bool, error) {
	g, err := fetchGroup(d, m)

	return err == nil && g != nil, err
}

func resourceGroupRead(d *schema.ResourceData, m interface{}) error {
	g, err := fetchGroup(d, m)

	if g == nil {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	_ = d.Set("name", g.Profile.Name)
	_ = d.Set("description", g.Profile.Description)
	if err := syncGroupUsers(d, m); err != nil {
		return err
	}

	return nil
}

func resourceGroupUpdate(d *schema.ResourceData, m interface{}) error {
	group := buildGroup(d)
	_, _, err := getOktaClientFromMetadata(m).Group.UpdateGroup(context.Background(), d.Id(), *group)
	if err != nil {
		return err
	}

	if err := updateGroupUsers(d, m); err != nil {
		return err
	}

	return resourceGroupRead(d, m)
}

func resourceGroupDelete(d *schema.ResourceData, m interface{}) error {
	_, err := getOktaClientFromMetadata(m).Group.DeleteGroup(context.Background(), d.Id())

	return err
}

func fetchGroup(d *schema.ResourceData, m interface{}) (*okta.Group, error) {
	g, resp, err := getOktaClientFromMetadata(m).Group.GetGroup(context.Background(), d.Id())

	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	return g, err
}

func syncGroupUsers(d *schema.ResourceData, m interface{}) error {
	// Only sync when the user opts in by outlining users in the group config
	if _, exists := d.GetOk("users"); !exists {
		return nil
	}
	userIDList, err := listGroupUserIDs(m, d.Id())
	if err != nil {
		return err
	}

	return d.Set("users", convertStringSetToInterface(userIDList))
}

func updateGroupUsers(d *schema.ResourceData, m interface{}) error {
	// Only sync when the user opts in by outlining users in the group config
	// To remove all users, define an empty set
	arr, exists := d.GetOk("users")
	if !exists {
		return nil
	}

	client := getOktaClientFromMetadata(m)
	existingUserList, _, err := client.Group.ListGroupUsers(context.Background(), d.Id(), nil)
	if err != nil {
		return err
	}

	rawArr := arr.(*schema.Set).List()
	userIDList := make([]string, len(rawArr))

	for i, ifaceID := range rawArr {
		userID := ifaceID.(string)
		userIDList[i] = userID

		if !containsUser(existingUserList, userID) {
			resp, err := client.Group.AddUserToGroup(context.Background(), d.Id(), userID)
			if err != nil {
				return responseErr(resp, err)
			}
		}
	}

	for _, user := range existingUserList {
		if !contains(userIDList, user.Id) {
			err := suppressErrorOn404(client.Group.RemoveUserFromGroup(context.Background(), d.Id(), user.Id))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func containsUser(users []*okta.User, id string) bool {
	for _, user := range users {
		if user.Id == id {
			return true
		}
	}
	return false
}
