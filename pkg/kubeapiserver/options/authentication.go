/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package options

import (
	"github.com/spf13/pflag"

	"time"

	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/kubernetes/pkg/kubeapiserver/authenticator"
)

type BuiltInAuthenticationOptions struct {
	PasswordFile    *PasswordFileAuthenticationOptions
	ServiceAccounts *ServiceAccountAuthenticationOptions

	TokenSuccessCacheTTL time.Duration
	TokenFailureCacheTTL time.Duration
}

type PasswordFileAuthenticationOptions struct {
	BasicAuthFile string
}

type ServiceAccountAuthenticationOptions struct {
	KeyFiles     []string
	Lookup       bool
	Issuer       string
	APIAudiences []string
}

func NewBuiltInAuthenticationOptions() *BuiltInAuthenticationOptions {
	return &BuiltInAuthenticationOptions{
		TokenSuccessCacheTTL: 10 * time.Second,
		TokenFailureCacheTTL: 0 * time.Second,
	}
}

func (s *BuiltInAuthenticationOptions) WithAll() *BuiltInAuthenticationOptions {
	return s.
		WithPasswordFile().
		WithServiceAccounts()
}

func (s *BuiltInAuthenticationOptions) WithPasswordFile() *BuiltInAuthenticationOptions {
	s.PasswordFile = &PasswordFileAuthenticationOptions{}
	return s
}

func (s *BuiltInAuthenticationOptions) WithServiceAccounts() *BuiltInAuthenticationOptions {
	s.ServiceAccounts = &ServiceAccountAuthenticationOptions{Lookup: true}
	return s
}

// Validate checks invalid config combination
func (s *BuiltInAuthenticationOptions) Validate() []error {
	allErrors := []error{}
	return allErrors
}

func (s *BuiltInAuthenticationOptions) AddFlags(fs *pflag.FlagSet) {
	if s.PasswordFile != nil {
		fs.StringVar(&s.PasswordFile.BasicAuthFile, "basic-auth-file", s.PasswordFile.BasicAuthFile, ""+
			"If set, the file that will be used to admit requests to the secure port of the API server "+
			"via http basic authentication.")
	}

	if s.ServiceAccounts != nil {
		fs.StringArrayVar(&s.ServiceAccounts.KeyFiles, "service-account-key-file", s.ServiceAccounts.KeyFiles, ""+
			"File containing PEM-encoded x509 RSA or ECDSA private or public keys, used to verify "+
			"ServiceAccount tokens. The specified file can contain multiple keys, and the flag can "+
			"be specified multiple times with different files. If unspecified, "+
			"--tls-private-key-file is used. Must be specified when "+
			"--service-account-signing-key is provided")

		fs.BoolVar(&s.ServiceAccounts.Lookup, "service-account-lookup", s.ServiceAccounts.Lookup,
			"If true, validate ServiceAccount tokens exist in etcd as part of authentication.")

		fs.StringVar(&s.ServiceAccounts.Issuer, "service-account-issuer", s.ServiceAccounts.Issuer, ""+
			"Identifier of the service account token issuer. The issuer will assert this identifier "+
			"in \"iss\" claim of issued tokens. This value is a string or URI.")

		fs.StringSliceVar(&s.ServiceAccounts.APIAudiences, "service-account-api-audiences", s.ServiceAccounts.APIAudiences, ""+
			"Identifiers of the API. The service account token authenticator will validate that "+
			"tokens used against the API are bound to at least one of these audiences.")
	}
}

func (s *BuiltInAuthenticationOptions) ToAuthenticationConfig() authenticator.AuthenticatorConfig {
	ret := authenticator.AuthenticatorConfig{
		TokenSuccessCacheTTL: s.TokenSuccessCacheTTL,
		TokenFailureCacheTTL: s.TokenFailureCacheTTL,
	}

	if s.PasswordFile != nil {
		ret.BasicAuthFile = s.PasswordFile.BasicAuthFile
	}

	if s.ServiceAccounts != nil {
		ret.ServiceAccountKeyFiles = s.ServiceAccounts.KeyFiles
		ret.ServiceAccountLookup = s.ServiceAccounts.Lookup
	}

	return ret
}

func (o *BuiltInAuthenticationOptions) ApplyTo(c *genericapiserver.Config) error {
	if o == nil {
		return nil
	}

	c.Authentication.SupportsBasicAuth = o.PasswordFile != nil && len(o.PasswordFile.BasicAuthFile) > 0
	return nil
}

// ApplyAuthorization will conditionally modify the authentication options based on the authorization options
func (o *BuiltInAuthenticationOptions) ApplyAuthorization(authorization *BuiltInAuthorizationOptions) {
	if o == nil || authorization == nil {
		return
	}
}
