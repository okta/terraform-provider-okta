package okta

func linksValue(links interface{}, keys ...string) string {
	if links == nil {
		return ""
	}
	sl, ok := links.([]interface{})
	if ok {
		if len(sl) == 0 {
			links = map[string]interface{}{}
		} else {
			links = sl[0]
		}
	}
	if len(keys) == 0 {
		v, ok := links.(string)
		if !ok {
			return ""
		}
		return v
	}
	l, ok := links.(map[string]interface{})
	if !ok {
		return ""
	}
	if len(keys) == 1 {
		return linksValue(l[keys[0]])
	}
	return linksValue(l[keys[0]], keys[1:]...)
}
