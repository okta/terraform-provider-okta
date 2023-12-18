# an Okta Group's memberships can be imported via the Okta group ID.
terraform import okta_group_memberships.test &#60;group id&#62;

# optional parameter track all users will also import all user id currently assigned to the group
terraform import okta_group_memberships.test &#60;group id&#62;/&#60;true&#62;