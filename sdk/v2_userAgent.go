package sdk

import "runtime"

type UserAgent struct {
	goVersion string

	osName string

	osVersion string

	config *config
}

func NewUserAgent(config *config) UserAgent {
	ua := UserAgent{}
	ua.config = config
	ua.goVersion = runtime.Version()
	ua.osName = runtime.GOOS
	ua.osVersion = runtime.GOARCH

	return ua
}

func (ua UserAgent) String() string {
	userAgentString := "local-v2-sdk/0.0.0 "
	userAgentString += "golang/" + ua.goVersion + " "
	userAgentString += ua.osName + "/" + ua.osVersion

	if ua.config.UserAgentExtra != "" {
		userAgentString += " " + ua.config.UserAgentExtra
	}

	return userAgentString
}
