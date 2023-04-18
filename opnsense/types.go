package opnsense

import (
	"fmt"
	"strings"
)

type HostOverride struct {
	UUID        string `json:"uuid"`
	Enabled     string `json:"enabled"`
	Hostname    string `json:"hostname"`
	Domain      string `json:"domain"`
	Server      string `json:"server"`
	Type        string `json:"rr"`
	Description string `json:"description"`
}

func NewHostOverride(hostname, domain, server string) HostOverride {
	override := HostOverride{
		Enabled:  "1",
		Hostname: hostname,
		Domain:   domain,
		Server:   server,
		Type:     "A",
	}
	override.Description = fmt.Sprintf("%s Automatically created by OPNsenseProxyAPI", override.GetFQDN())
	return override
}

func (hostOverride HostOverride) GetFQDN() string {
	return fmt.Sprintf("%s.%s", hostOverride.Hostname, hostOverride.Domain)
}

type AliasOverride struct {
	UUID        string `json:"uuid"`
	Enabled     string `json:"enabled"`
	Host        string `json:"host"`
	Hostname    string `json:"hostname"`
	Domain      string `json:"domain"`
	Description string `json:"description"`
}

func NewAliasOverride(hostname, domain, host string) AliasOverride {
	override := AliasOverride{
		Enabled:  "1",
		Host:     host,
		Hostname: hostname,
		Domain:   domain,
	}
	override.Description = fmt.Sprintf("%s Automatically created by OPNsenseProxyAPI", override.GetFQDN())
	return override
}

func (aliasOverride AliasOverride) IsHostFQDN() bool {
	return strings.Contains(aliasOverride.Host, ".")
}

func (aliasOverride AliasOverride) GetFQDN() string {
	return fmt.Sprintf("%s.%s", aliasOverride.Hostname, aliasOverride.Domain)
}

type addHostOverrideContainer struct {
	Host HostOverride `json:"host"`
}

type getHostOverridesContainer struct {
	Rows     []HostOverride `json:"rows"`
	RowCount int            `json:"rowCount"`
	Total    int            `json:"total"`
	Current  int            `json:"current"`
}

type getHostAliasesContainer struct {
	Rows     []AliasOverride `json:"rows"`
	RowCount int             `json:"rowCount"`
	Total    int             `json:"total"`
	Current  int             `json:"current"`
}

type addHostAliasContainer struct {
	Alias AliasOverride `json:"alias"`
}

type deleteResponse struct {
	Result string `json:"result"`
}

func (r *deleteResponse) Succeeded() bool {
	return r.Result == "deleted"
}
