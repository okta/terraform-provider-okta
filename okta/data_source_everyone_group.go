package okta

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

// data source to retrieve information on the Everyone Group

func dataSourceEveryoneGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEveryoneGroupRead,

		Schema: map[string]*schema.Schema{},
	}
}

func dataSourceEveryoneGroupRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Data Source Everyone Group Read")
	client := m.(*Config).oktaClient

	groups, _, err := client.Groups.ListGroups("q=Everyone")
	if err != nil {
		return fmt.Errorf("[ERROR] ListGroups Everyone query error: %v", err)
	}
	if groups != nil {
		if len(groups.Groups) > 1 {
			return fmt.Errorf("[ERROR] Group Everyone query resulted in more than one group.")
		}
		d.SetId(groups.Groups[0].ID)
	} else {
		return fmt.Errorf("[ERROR] Group Everyone query resulted in no groups.")
	}
	return nil
}
