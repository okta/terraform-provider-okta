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

package okta

import (
	"net/http"
	"time"

	"github.com/okta/okta-sdk-golang/okta/cache"
)

type Config struct {
	BackoffEnabled bool          `yaml:"withBackoff" envconfig:"OKTA_BACK_OFF_ENABLED"`
	MaxRetries     int32         `yaml:"maxRetries" envconfig:"OKTA_MAX_RETRIES"`
	MinWait        time.Duration `yaml:"minWait"`
	MaxWait        time.Duration `yaml:"maxWait"`
	Okta           struct {
		Client struct {
			Cache struct {
				Enabled    bool  `yaml:"enabled" envconfig:"OKTA_CLIENT_CACHE_ENABLED"`
				DefaultTtl int32 `yaml:"defaultTtl" envconfig:"OKTA_CLIENT_CACHE_DEFAULT_TTL"`
				DefaultTti int32 `yaml:"defaultTti" envconfig:"OKTA_CLIENT_CACHE_DEFAULT_TTI"`
			} `yaml:"cache"`
			Proxy struct {
				Port     int32  `yaml:"port" envconfig:"OKTA_CLIENT_PROXY_PORT"`
				Host     string `yaml:"host" envconfig:"OKTA_CLIENT_PROXY_HOST"`
				Username string `yaml:"username" envconfig:"OKTA_CLIENT_PROXY_USERNAME"`
				Password string `yaml:"password" envconfig:"OKTA_CLIENT_PROXY_PASSWORD"`
			} `yaml:"proxy"`
			ConnectionTimeout int32  `yaml:"connectionTimeout" envconfig:"OKTA_CLIENT_CONNECTION_TIMEOUT"`
			OrgUrl            string `yaml:"orgUrl" envconfig:"OKTA_CLIENT_ORGURL"`
			Token             string `yaml:"token" envconfig:"OKTA_CLIENT_TOKEN"`
		} `yaml:"client"`
		Testing struct {
			DisableHttpsCheck bool `yaml:"disableHttpsCheck" envconfig:"OKTA_TESTING_DISABLE_HTTPS_CHECK"`
		} `yaml:"testing"`
	} `yaml:"okta"`
	UserAgentExtra string
	HttpClient     http.Client
	CacheManager   cache.Cache
}

type ConfigSetter func(*Config)

func WithCache(cache bool) ConfigSetter {
	return func(c *Config) {
		c.Okta.Client.Cache.Enabled = cache
	}
}

func WithCacheManager(cacheManager cache.Cache) ConfigSetter {
	return func(c *Config) {
		c.CacheManager = cacheManager
	}
}

func WithCacheTtl(i int32) ConfigSetter {
	return func(c *Config) {
		c.Okta.Client.Cache.DefaultTtl = i
	}
}

func WithCacheTti(i int32) ConfigSetter {
	return func(c *Config) {
		c.Okta.Client.Cache.DefaultTti = i
	}
}

func WithConnectionTimeout(i int32) ConfigSetter {
	return func(c *Config) {
		c.Okta.Client.ConnectionTimeout = i
	}
}

func WithProxyPort(i int32) ConfigSetter {
	return func(c *Config) {
		c.Okta.Client.Proxy.Port = i
	}
}

func WithProxyHost(host string) ConfigSetter {
	return func(c *Config) {
		c.Okta.Client.Proxy.Host = host
	}
}

func WithProxyUsername(username string) ConfigSetter {
	return func(c *Config) {
		c.Okta.Client.Proxy.Username = username
	}
}

func WithProxyPassword(pass string) ConfigSetter {
	return func(c *Config) {
		c.Okta.Client.Proxy.Password = pass
	}
}

func WithOrgUrl(url string) ConfigSetter {
	return func(c *Config) {
		c.Okta.Client.OrgUrl = url
	}
}

func WithToken(token string) ConfigSetter {
	return func(c *Config) {
		c.Okta.Client.Token = token
	}
}

func WithUserAgentExtra(userAgent string) ConfigSetter {
	return func(c *Config) {
		c.UserAgentExtra = userAgent
	}
}

func WithHttpClient(httpClient http.Client) ConfigSetter {
	return func(c *Config) {
		c.HttpClient = httpClient
	}
}

func WithTestingDisableHttpsCheck(httpsCheck bool) ConfigSetter {
	return func(c *Config) {
		c.Okta.Testing.DisableHttpsCheck = httpsCheck
	}
}

func WithBackoff(backoff bool) ConfigSetter {
	return func(c *Config) {
		c.BackoffEnabled = backoff
	}
}

func WithMinWait(wait time.Duration) ConfigSetter {
	return func(c *Config) {
		c.MinWait = wait
	}
}

func WithMaxWait(wait time.Duration) ConfigSetter {
	return func(c *Config) {
		c.MaxWait = wait
	}
}

func WithRetries(retries int32) ConfigSetter {
	return func(c *Config) {
		c.MaxRetries = retries
	}
}
