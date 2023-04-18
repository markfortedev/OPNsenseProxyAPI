package opnsense

import (
	"crypto/tls"
	"github.com/go-resty/resty/v2"
	"os"
	"reflect"
	"testing"
)

type apiClientFields struct {
	apiKey    string
	apiSecret string
	address   string
	client    *resty.Client
}

func generateAPIClientFields() apiClientFields {
	apiKey := os.Getenv("API_KEY")
	apiSecret := os.Getenv("API_SECRET")
	address := os.Getenv("OPNSENSE_ADDRESS")
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	return apiClientFields{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		address:   address,
		client:    client,
	}
}

func Test_apiKeyClient_CreateAliasOverride(t *testing.T) {

	fieldConst := generateAPIClientFields()

	type args struct {
		aliasOverride AliasOverride
	}

	tests := []struct {
		name    string
		fields  apiClientFields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "Create Alias Override with UUID",
			fields:  fieldConst,
			args:    struct{ aliasOverride AliasOverride }{aliasOverride: NewAliasOverride("testOPNsenseProxyAPI01", "example.com", "f2a5edee-1b46-4a08-9041-4f51e02932f5")},
			wantErr: false,
			want:    true,
		},
		{
			name:    "Create Alias Override with Hostname",
			fields:  fieldConst,
			args:    struct{ aliasOverride AliasOverride }{aliasOverride: NewAliasOverride("testOPNsenseProxyAPI02", "example.com", "test.testdomain.com")},
			wantErr: false,
			want:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &apiKeyClient{
				apiKey:    tt.fields.apiKey,
				apiSecret: tt.fields.apiSecret,
				address:   tt.fields.address,
				client:    tt.fields.client,
			}
			got, err := c.CreateAliasOverride(tt.args.aliasOverride)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAliasOverride() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreateAliasOverride() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_apiKeyClient_CreateHostOverride(t *testing.T) {
	fieldConst := generateAPIClientFields()
	type args struct {
		hostOverride HostOverride
	}
	tests := []struct {
		name    string
		fields  apiClientFields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "Create Host Override",
			fields:  fieldConst,
			args:    struct{ hostOverride HostOverride }{hostOverride: NewHostOverride("opnsenseProxyAPIOverrideTest01", "example.com", "10.0.2.0")},
			wantErr: false,
			want:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &apiKeyClient{
				apiKey:    tt.fields.apiKey,
				apiSecret: tt.fields.apiSecret,
				address:   tt.fields.address,
				client:    tt.fields.client,
			}
			got, err := c.CreateHostOverride(tt.args.hostOverride)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateHostOverride() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreateHostOverride() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_apiKeyClient_GetAliasOverrides(t *testing.T) {
	fieldConsts := generateAPIClientFields()
	tests := []struct {
		name    string
		fields  apiClientFields
		wantErr bool
	}{
		{
			name:    "Get Alias Overrides",
			fields:  fieldConsts,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &apiKeyClient{
				apiKey:    tt.fields.apiKey,
				apiSecret: tt.fields.apiSecret,
				address:   tt.fields.address,
				client:    tt.fields.client,
			}
			_, err := c.GetAliasOverrides()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAliasOverrides() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_apiKeyClient_GetHostOverrides(t *testing.T) {
	fieldConsts := generateAPIClientFields()
	tests := []struct {
		name    string
		fields  apiClientFields
		wantErr bool
	}{
		{
			name:    "Get Host Overrides",
			fields:  fieldConsts,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &apiKeyClient{
				apiKey:    tt.fields.apiKey,
				apiSecret: tt.fields.apiSecret,
				address:   tt.fields.address,
				client:    tt.fields.client,
			}
			_, err := c.GetHostOverrides()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHostOverrides() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_apiKeyClient_DeleteAliasOverride(t *testing.T) {
	fieldConsts := generateAPIClientFields()
	type args struct {
		fqdn string
	}
	tests := []struct {
		name    string
		fields  apiClientFields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "Delete alias that exists",
			fields:  fieldConsts,
			args:    struct{ fqdn string }{fqdn: "testDelete.testdomain.com"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Delete alias that does not exists",
			fields:  fieldConsts,
			args:    struct{ fqdn string }{fqdn: "testDeleteThatDoesNotExist.testdomain.com"},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &apiKeyClient{
				apiKey:    tt.fields.apiKey,
				apiSecret: tt.fields.apiSecret,
				address:   tt.fields.address,
				client:    tt.fields.client,
			}
			got, err := c.DeleteAliasOverride(tt.args.fqdn)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteAliasOverride() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DeleteAliasOverride() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_apiKeyClient_DeleteHostOverride(t *testing.T) {
	fieldConsts := generateAPIClientFields()
	type args struct {
		fqdn string
	}
	tests := []struct {
		name    string
		fields  apiClientFields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "Test Delete Host that exists",
			fields:  fieldConsts,
			args:    struct{ fqdn string }{fqdn: "testDeleteHost.testdomain.com"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Test Delete Host that does not exists",
			fields:  fieldConsts,
			args:    struct{ fqdn string }{fqdn: "doesnotExits.tld"},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &apiKeyClient{
				apiKey:    tt.fields.apiKey,
				apiSecret: tt.fields.apiSecret,
				address:   tt.fields.address,
				client:    tt.fields.client,
			}
			got, err := c.DeleteHostOverride(tt.args.fqdn)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteHostOverride() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DeleteHostOverride() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_apiKeyClient_getAliasesToCreateAndDelete(t *testing.T) {
	fieldConsts := generateAPIClientFields()
	type args struct {
		currentAliases  []string
		existingAliases []AliasOverride
	}
	tests := []struct {
		name         string
		fields       apiClientFields
		args         args
		wantToCreate []string
		wantToDelete []string
	}{
		{
			name:   "Test Sync",
			fields: fieldConsts,
			args: struct {
				currentAliases  []string
				existingAliases []AliasOverride
			}{
				currentAliases: []string{
					"1.tld",
					"2.tld",
					"3.tld",
					"new.tld",
				},
				existingAliases: []AliasOverride{
					{
						Hostname: "1",
						Domain:   "tld",
					},
					{
						Hostname: "2",
						Domain:   "tld",
					},
					{
						Hostname: "3",
						Domain:   "tld",
					},
					{
						Hostname: "delete",
						Domain:   "tld",
					},
				},
			},
			wantToCreate: []string{
				"new.tld",
			},
			wantToDelete: []string{
				"delete.tld",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &apiKeyClient{
				apiKey:    tt.fields.apiKey,
				apiSecret: tt.fields.apiSecret,
				address:   tt.fields.address,
				client:    tt.fields.client,
			}
			got, got1 := c.getAliasesToCreateAndDelete(tt.args.currentAliases, tt.args.existingAliases)
			if !reflect.DeepEqual(got, tt.wantToCreate) {
				t.Errorf("getAliasesToCreateAndDelete() got = %v, want %v", got, tt.wantToCreate)
			}
			if !reflect.DeepEqual(got1, tt.wantToDelete) {
				t.Errorf("getAliasesToCreateAndDelete() got1 = %v, want %v", got1, tt.wantToDelete)
			}
		})
	}
}
