package opnsense

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"strings"
)

type Client interface {
	CreateHostOverride(hostOverride HostOverride) (bool, error)
	CreateAliasOverride(aliasOverride AliasOverride) (bool, error)
	GetHostOverrides() ([]HostOverride, error)
	GetAliasOverrides() ([]AliasOverride, error)
	GetAliasOverridesForHost(host string) ([]AliasOverride, error)
	GetHostOverride(fqdn string) (HostOverride, error)
	GetAliasOverride(fqdn string) (AliasOverride, error)
	DoesHostOverrideExist(fqdn string) (bool, error)
	DeleteHostOverride(fqdn string) (bool, error)
	DeleteAliasOverride(fqdn string) (bool, error)
	SyncAliases(host string, aliases []string, domain string) (bool, error)
}

type apiKeyClient struct {
	apiKey    string
	apiSecret string
	address   string
	client    *resty.Client
}

func NewClient(address, apiKey, apiSecret string) Client {
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	return &apiKeyClient{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		address:   address,
		client:    client,
	}
}

func (c *apiKeyClient) newRequest() *resty.Request {
	request := c.client.R()
	request.SetBasicAuth(c.apiKey, c.apiSecret)
	return request
}

func (c *apiKeyClient) CreateHostOverride(hostOverride HostOverride) (bool, error) {
	endpoint := fmt.Sprintf("%s/api/unbound/settings/addhostoverride", c.address)
	response, err := c.newRequest().
		SetHeader("Content-Type", "application/json").
		SetBody(addHostOverrideContainer{Host: hostOverride}).
		Post(endpoint)
	if err != nil {
		return false, err
	}
	return response.IsSuccess(), nil
}

func (c *apiKeyClient) CreateAliasOverride(aliasOverride AliasOverride) (bool, error) {
	endpoint := fmt.Sprintf("%s/api/unbound/settings/addHostAlias", c.address)
	if aliasOverride.IsHostFQDN() {
		host, err := c.GetHostOverride(aliasOverride.Host)
		if err != nil {
			return false, err
		}
		aliasOverride.Host = host.UUID
	}
	response, err := c.newRequest().
		SetHeader("Content-Type", "application/json").
		SetBody(addHostAliasContainer{Alias: aliasOverride}).
		Post(endpoint)
	if err != nil {
		return false, err
	}
	return response.IsSuccess(), nil
}

func (c *apiKeyClient) GetHostOverrides() ([]HostOverride, error) {
	endpoint := fmt.Sprintf("%s/api/unbound/settings/searchHostOverride/", c.address)
	resp, err := c.newRequest().SetResult(getHostOverridesContainer{}).Get(endpoint)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, errors.New(resp.Status())
	}
	container := resp.Result().(*getHostOverridesContainer)
	return container.Rows, nil
}

func (c *apiKeyClient) GetAliasOverrides() ([]AliasOverride, error) {
	endpoint := fmt.Sprintf("%s/api/unbound/settings/searchHostAlias", c.address)
	resp, err := c.newRequest().SetResult(getHostAliasesContainer{}).Get(endpoint)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, errors.New(resp.Status())
	}
	container := resp.Result().(*getHostAliasesContainer)
	return container.Rows, nil
}

func (c *apiKeyClient) GetAliasOverridesForHost(host string) ([]AliasOverride, error) {
	hostOverride, err := c.GetHostOverride(host)
	if err != nil {
		return nil, err
	}
	aliasOverrides, err := c.GetAliasOverrides()
	if err != nil {
		return nil, err
	}
	var aliasOverridesFromHost []AliasOverride
	for _, alias := range aliasOverrides {
		if alias.Host == hostOverride.GetFQDN() {
			aliasOverridesFromHost = append(aliasOverridesFromHost, alias)
		}
	}
	return aliasOverridesFromHost, nil
}

func (c *apiKeyClient) GetHostOverride(fqdn string) (HostOverride, error) {
	allHostOverrides, err := c.GetHostOverrides()
	if err != nil {
		return HostOverride{}, err
	}
	for _, hostOverride := range allHostOverrides {
		if fqdn == hostOverride.GetFQDN() {
			return hostOverride, nil
		}
	}
	return HostOverride{}, errors.New(fmt.Sprintf("Host override %v does not exist", fqdn))
}

func (c *apiKeyClient) DoesHostOverrideExist(fqdn string) (bool, error) {
	allHostOverrides, err := c.GetHostOverrides()
	if err != nil {
		return false, err
	}
	for _, hostOverride := range allHostOverrides {
		if fqdn == hostOverride.GetFQDN() {
			return true, nil
		}
	}
	return false, nil
}

func (c *apiKeyClient) GetAliasOverride(fqdn string) (AliasOverride, error) {
	allAliasOverrides, err := c.GetAliasOverrides()
	if err != nil {
		return AliasOverride{}, err
	}
	for _, aliasOverride := range allAliasOverrides {
		if fqdn == aliasOverride.GetFQDN() {
			return aliasOverride, nil
		}
	}
	return AliasOverride{}, errors.New(fmt.Sprintf("Alias override %v does not exist", fqdn))
}

func (c *apiKeyClient) DeleteHostOverride(fqdn string) (bool, error) {
	hostOverride, err := c.GetHostOverride(fqdn)
	if err != nil {
		return false, err
	}
	endpoint := fmt.Sprintf("%s/api/unbound/settings/delHostOverride/%s", c.address, hostOverride.UUID)
	return c.performDelete(fqdn, endpoint)
}

func (c *apiKeyClient) DeleteAliasOverride(fqdn string) (bool, error) {
	aliasOverride, err := c.GetAliasOverride(fqdn)
	if err != nil {
		return false, err
	}
	endpoint := fmt.Sprintf("%s/api/unbound/settings/delHostAlias/%s", c.address, aliasOverride.UUID)
	return c.performDelete(fqdn, endpoint)
}

func (c *apiKeyClient) performDelete(fqdn string, endpoint string) (bool, error) {
	resp, err := c.newRequest().SetResult(deleteResponse{}).Post(endpoint)
	if err != nil {
		return false, err
	}
	result := resp.Result().(*deleteResponse)
	if !result.Succeeded() {
		return false, errors.New(fmt.Sprintf("%v does not exist", fqdn))
	}
	return true, nil
}

func (c *apiKeyClient) SyncAliases(host string, currentAliases []string, domain string) (bool, error) {
	existingAliases, err := c.GetAliasOverridesForHost(host)
	if err != nil {
		return false, err
	}
	aliasesToCreate, aliasesToDelete := c.getAliasesToCreateAndDelete(currentAliases, existingAliases)
	if len(aliasesToDelete) > 0 {
		log.Infof("Deleting %v aliases for %v: [%v]", len(aliasesToDelete), host, strings.Join(aliasesToDelete, ", "))
	}
	if len(aliasesToCreate) > 0 {
		log.Infof("Creating %v aliases for %v: [%v]", len(aliasesToCreate), host, strings.Join(aliasesToCreate, ", "))
	}
	for _, aliasToCreate := range aliasesToCreate {
		hostname := strings.Replace(aliasToCreate, fmt.Sprintf(".%v", domain), "", -1)
		aliasOverride := NewAliasOverride(hostname, domain, host)
		_, err = c.CreateAliasOverride(aliasOverride)
		if err != nil {
			return false, err
		}
	}
	for _, aliasToDelete := range aliasesToDelete {
		_, err = c.DeleteAliasOverride(aliasToDelete)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (c *apiKeyClient) getAliasesToCreateAndDelete(currentAliases []string, existingAliases []AliasOverride) ([]string, []string) {
	var aliasesToCreate []string
	var aliasesToDelete []string

	// get current aliases to create
	for _, currentAliasFQDN := range currentAliases {
		doesCurrentAliasExist := false
		for _, existingAlias := range existingAliases {
			if currentAliasFQDN == existingAlias.GetFQDN() {
				doesCurrentAliasExist = true
			}
		}
		if !doesCurrentAliasExist {
			aliasesToCreate = append(aliasesToCreate, currentAliasFQDN)
		}
	}
	// get current aliases to delete
	for _, existingAlias := range existingAliases {
		shouldAliasExist := false
		for _, currentAliasFQDN := range currentAliases {
			if currentAliasFQDN == existingAlias.GetFQDN() {
				shouldAliasExist = true
			}
		}
		if !shouldAliasExist {
			aliasesToDelete = append(aliasesToDelete, existingAlias.GetFQDN())
		}
	}
	return aliasesToCreate, aliasesToDelete
}
