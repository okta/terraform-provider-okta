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

import ()

type PasswordPolicyPasswordSettingsComplexity struct {
	Dictionary        *PasswordDictionary `json:"dictionary,omitempty"`
	ExcludeAttributes []string            `json:"excludeAttributes,omitempty"`
	ExcludeUsername   *bool               `json:"excludeUsername,omitempty"`
	MinLength         int64               `json:"minLength,omitempty"`
	MinLowerCase      int64               `json:"minLowerCase,omitempty"`
	MinNumber         int64               `json:"minNumber,omitempty"`
	MinSymbol         int64               `json:"minSymbol,omitempty"`
	MinUpperCase      int64               `json:"minUpperCase,omitempty"`
}
