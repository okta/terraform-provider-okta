/*
* Copyright 2018 - Present Okta, Inc.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

// AUTO-GENERATED!  DO NOT EDIT FILE DIRECTLY

package okta

import (
	"context"
	"fmt"
	"io/ioutil"
	"os/user"

	"github.com/okta/okta-sdk-golang/v2/okta/cache"

	"github.com/go-yaml/yaml"
	"github.com/kelseyhightower/envconfig"
)

const Version = "2.0.0"

type Client struct {
	config *config

	requestExecutor *RequestExecutor

	resource resource

	Application         *ApplicationResource
	AuthorizationServer *AuthorizationServerResource
	EventHook           *EventHookResource
	Feature             *FeatureResource
	Group               *GroupResource
	IdentityProvider    *IdentityProviderResource
	InlineHook          *InlineHookResource
	LogEvent            *LogEventResource
	LinkedObject        *LinkedObjectResource
	UserType            *UserTypeResource
	Policy              *PolicyResource
	Session             *SessionResource
	SmsTemplate         *SmsTemplateResource
	TrustedOrigin       *TrustedOriginResource
	User                *UserResource
	UserFactor          *UserFactorResource
}

type resource struct {
	client *Client
}

type clientContextKey struct {
}

func NewClient(ctx context.Context, conf ...ConfigSetter) (context.Context, *Client, error) {
	config := &config{}

	setConfigDefaults(config)
	config = readConfigFromSystem(*config)
	config = readConfigFromApplication(*config)
	config = readConfigFromEnvironment(*config)

	for _, confSetter := range conf {
		confSetter(config)
	}

	var oktaCache cache.Cache
	if !config.Okta.Client.Cache.Enabled {
		oktaCache = cache.NewNoOpCache()
	} else {
		if config.CacheManager == nil {
			oktaCache = cache.NewGoCache(config.Okta.Client.Cache.DefaultTtl,
				config.Okta.Client.Cache.DefaultTti)
		} else {
			oktaCache = config.CacheManager
		}
	}

	config.CacheManager = oktaCache

	config, err := validateConfig(config)
	if err != nil {
		panic(err)
	}

	c := &Client{}
	c.config = config
	c.requestExecutor = NewRequestExecutor(&config.HttpClient, oktaCache, config)

	c.resource.client = c

	c.Application = (*ApplicationResource)(&c.resource)
	c.AuthorizationServer = (*AuthorizationServerResource)(&c.resource)
	c.EventHook = (*EventHookResource)(&c.resource)
	c.Feature = (*FeatureResource)(&c.resource)
	c.Group = (*GroupResource)(&c.resource)
	c.IdentityProvider = (*IdentityProviderResource)(&c.resource)
	c.InlineHook = (*InlineHookResource)(&c.resource)
	c.LogEvent = (*LogEventResource)(&c.resource)
	c.LinkedObject = (*LinkedObjectResource)(&c.resource)
	c.UserType = (*UserTypeResource)(&c.resource)
	c.Policy = (*PolicyResource)(&c.resource)
	c.Session = (*SessionResource)(&c.resource)
	c.SmsTemplate = (*SmsTemplateResource)(&c.resource)
	c.TrustedOrigin = (*TrustedOriginResource)(&c.resource)
	c.User = (*UserResource)(&c.resource)
	c.UserFactor = (*UserFactorResource)(&c.resource)

	contextReturn := context.WithValue(ctx, clientContextKey{}, c)

	return contextReturn, c, nil
}

func ClientFromContext(ctx context.Context) (*Client, bool) {
	u, ok := ctx.Value(clientContextKey{}).(*Client)
	return u, ok
}

func (c *Client) GetConfig() *config {
	return c.config
}

func (c *Client) GetRequestExecutor() *RequestExecutor {
	return c.requestExecutor
}

func setConfigDefaults(c *config) {
	var conf []ConfigSetter

	conf = append(conf,
		WithConnectionTimeout(30),
		WithCache(true),
		WithCacheTtl(300),
		WithCacheTti(300),
		WithUserAgentExtra(""),
		WithTestingDisableHttpsCheck(false),
		WithRequestTimeout(0),
		WithRateLimitMaxRetries(2),
		WithAuthorizationMode("SSWS"))

	for _, confSetter := range conf {
		confSetter(c)
	}
}

func readConfigFromFile(location string, c config) (*config, error) {
	yamlConfig, err := ioutil.ReadFile(location)

	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlConfig, &c)
	if err != nil {
		return nil, err
	}

	return &c, err
}

func readConfigFromSystem(c config) *config {
	currUser, err := user.Current()
	if err != nil {
		return &c
	}
	if currUser.HomeDir == "" {
		return &c
	}

	conf, err := readConfigFromFile(currUser.HomeDir+"/.okta/okta.yaml", c)

	if err != nil {
		return &c
	}

	return conf
}

func readConfigFromApplication(c config) *config {
	conf, err := readConfigFromFile(".okta.yaml", c)

	if err != nil {
		return &c
	}

	return conf
}

func readConfigFromEnvironment(c config) *config {
	err := envconfig.Process("okta", &c)
	if err != nil {
		fmt.Println("error parsing")
		return &c
	}
	return &c
}
