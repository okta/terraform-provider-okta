data "okta_brands" "test" {
}

data "okta_themes" "test" {
  brand_id = tolist(data.okta_brands.test.brands)[0].id
}

data "okta_theme" "test" {
  brand_id = tolist(data.okta_brands.test.brands)[0].id
  theme_id = tolist(data.okta_themes.test.themes)[0].id
}
