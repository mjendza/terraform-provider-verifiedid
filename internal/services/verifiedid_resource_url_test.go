package services_test

import (
	"testing"

	"github.com/mjendza/terraform-provider-verifiedid/internal/services"
)

func TestIsContractURL(t *testing.T) {
	cases := []struct {
		name string
		url  string
		want bool
	}{
		{"contracts collection", "verifiableCredentials/authorities/00000000-0000-0000-0000-000000000000/contracts", true},
		{"contracts collection with leading slash", "/verifiableCredentials/authorities/abc/contracts", true},
		{"contracts collection with trailing slash", "verifiableCredentials/authorities/abc/contracts/", true},
		{"contracts collection with api version prefix", "/v1.0/verifiableCredentials/authorities/abc/contracts", true},
		{"contracts collection with beta prefix", "/beta/verifiableCredentials/authorities/abc/contracts", true},
		{"contracts collection mixed case", "VerifiableCredentials/Authorities/abc/Contracts", true},
		{"contracts collection with query string", "verifiableCredentials/authorities/abc/contracts?$select=id", true},

		{"authorities collection", "verifiableCredentials/authorities", false},
		{"authority item", "verifiableCredentials/authorities/abc", false},
		{"action under authority", "verifiableCredentials/authorities/abc/didInfo/signingKeys", false},
		{"applications graph", "applications", false},
		{"users graph", "users", false},
		{"empty", "", false},
		{"unrelated path containing contracts", "fooContracts/bar", false},
		{"missing authority id", "verifiableCredentials/authorities//contracts", false},
		{"ref relationship", "groups/abc/members/$ref", false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := services.IsContractURL(tc.url)
			if got != tc.want {
				t.Fatalf("IsContractURL(%q) = %v, want %v", tc.url, got, tc.want)
			}
		})
	}
}
