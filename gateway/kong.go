package gateway

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type KongGateway struct {
	GatewayURL string
	config     *MicroserviceConfig
	client     *http.Client
}

type MicroserviceConfig struct {
	MicroserviceName string
	MicroservicePort int
	VirtualHost      string
	Hosts            []string
	Weight           int
	ServicesMaxSlots int
	UpstreamURL      string
}

func NewKongGateway(adminUrl string, client *http.Client, config *MicroserviceConfig) *KongGateway {
	return &KongGateway{
		GatewayURL: adminUrl,
		config:     config,
		client:     client,
	}
}

func NewKongGatewayFromConfigFile(adminUrl string, client *http.Client, configFile string) (*KongGateway, error) {
	var config MicroserviceConfig
	cnf, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(cnf, &config)
	if err != nil {
		return nil, err
	}

	return NewKongGateway(adminUrl, client, &config), nil
}

func (kong *KongGateway) SelfRegister() error {
	err := kong.createOrUpdateUpstream(kong.config.VirtualHost, kong.config.ServicesMaxSlots)
	if err != nil {
		return err
	}

	apiConf := NewAPIConf()
	// TODO: map here from config
	apiConf.Name = kong.config.MicroserviceName
	apiConf.Hosts = kong.config.Hosts
	apiConf.UpstreamURL = kong.config.UpstreamURL

	_, err = kong.createOrUpdateAPI(apiConf)
	if err != nil {
		return err
	}

	_, err = kong.addSelfAsTarget(kong.config.VirtualHost, kong.config.MicroservicePort, kong.config.Weight)
	return err
}

func (kong *KongGateway) Unregister() error {
	_, err := kong.addSelfAsTarget(kong.config.VirtualHost, kong.config.MicroservicePort, 0)
	return err
}

type upstream struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	OrderList []int  `json:"orderlist,omitempty"`
	Slots     int    `json:"slots,omitempty"`
	CreatedAt int    `json:"created_at,omitempty"`
}

type upstreamTarget struct {
	ID         string `json:"id,omitempty"`
	Target     string `json:"target,omitempty"`
	Weight     int    `json:"weight,omitempty"`
	UpstreamID string `json:"upstream_id,omitempty"`
	CreatedAt  int    `json:"created_at,omitempty"`
}

type API struct {
	ID                     string   `json:"id,omitempty"`
	CreatedAt              int      `json:"created_at,omitempty"`
	Hosts                  []string `json:"hosts,omitempty"`
	URIs                   []string `json:"-"`
	Methods                []string `json:"-"`
	HTTPIfTerminated       bool     `json:"http_if_terminated,omitempty"`
	HTTPSOnly              bool     `json:"https_only,omitempty"`
	Name                   string   `json:"name,omitempty"`
	PreserveHost           bool     `json:"preserve_host,omitempty"`
	Retries                int      `json:"retries,omitempty"`
	StripURI               bool     `json:"strip_uri,omitempty"`
	UpstreamConnectTimeout int      `json:"upstream_connect_timeout,omitempty"`
	UpstreamReadTimeout    int      `json:"upstream_read_timeout,omitempty"`
	UpstreamSendTimeout    int      `json:"upstream_send_timeout,omitempty"`
	UpstreamURL            string   `json:"upstream_url,omitempty"`
}

func NewAPIConf() *API {
	api := API{
		Hosts:                  []string{},
		URIs:                   []string{},
		Methods:                []string{},
		HTTPIfTerminated:       true,
		HTTPSOnly:              false,
		PreserveHost:           false,
		Retries:                5,
		StripURI:               true,
		UpstreamConnectTimeout: 60000,
		UpstreamReadTimeout:    60000,
		UpstreamSendTimeout:    60000,
	}
	return &api
}

func (api *API) AddHost(host string) {
	api.Hosts = append(api.Hosts, host)
}

func getServiceIP() (string, error) {
	return "", nil
}

func (kong *KongGateway) getKongURL(path string) string {
	return fmt.Sprintf("%s/%s", kong.GatewayURL, path)
}

func (kong *KongGateway) getUpstreamObj(name string) (*upstream, error) {
	var upstreamObj upstream
	resp, err := kong.client.Get(kong.getKongURL(fmt.Sprintf("upstreams/%s", name)))
	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(resp.Body).Decode(&upstreamObj); err != nil {
		return nil, err
	}
	return &upstreamObj, nil
}

func (kong *KongGateway) createUpstreamObj(name string, slots int) (*upstream, error) {
	var upstreamObj upstream
	form := url.Values{}

	form.Add("name", name)
	form.Add("slots", fmt.Sprintf("%d", slots))

	resp, err := kong.client.Post(kong.getKongURL("upstreams/"), "multipart/form-data", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&upstreamObj); err != nil {
		return nil, err
	}
	return &upstreamObj, nil
}

func (kong *KongGateway) createOrUpdateUpstream(name string, slots int) error {
	up, err := kong.getUpstreamObj(name)
	if err != nil {
		return err
	}
	if up == nil {
		if _, err = kong.createUpstreamObj(name, slots); err != nil {
			return err
		}
	}
	return nil
}

func (kong *KongGateway) createKongAPI(apiConf *API) (*API, error) {
	var result API
	form := url.Values{}

	if apiConf.Name != "" {
		form.Add("name", apiConf.Name)
	}

	if apiConf.Hosts != nil {
		form.Add("hosts", strings.Join(apiConf.Hosts, ","))
	}

	if apiConf.URIs != nil {
		form.Add("uris", strings.Join(apiConf.URIs, ","))
	}

	if apiConf.Methods != nil {
		form.Add("methods", strings.Join(apiConf.Methods, ","))
	}

	if apiConf.UpstreamURL != "" {
		form.Add("upstream_url", apiConf.UpstreamURL)
	}

	form.Add("retries", fmt.Sprintf("%d", apiConf.Retries))
	form.Add("upstream_connect_timeout", fmt.Sprintf("%d", apiConf.UpstreamConnectTimeout))
	form.Add("upstream_send_timeout", fmt.Sprintf("%d", apiConf.UpstreamSendTimeout))
	form.Add("upstream_read_timeout", fmt.Sprintf("%d", apiConf.UpstreamReadTimeout))

	form.Add("strip_uri", fmt.Sprintf("%t", apiConf.StripURI))
	form.Add("preserve_host", fmt.Sprintf("%t", apiConf.PreserveHost))
	form.Add("https_only", fmt.Sprintf("%t", apiConf.HTTPSOnly))
	form.Add("http_if_terminated", fmt.Sprintf("%t", apiConf.HTTPIfTerminated))

	resp, err := kong.client.Post(kong.getKongURL("apis/"), "", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	result.URIs = apiConf.URIs
	result.Methods = apiConf.Methods

	return &result, nil
}

func (kong *KongGateway) getAPI(name string) (*API, error) {
	resp, err := kong.client.Get(kong.getKongURL(fmt.Sprintf("apis/%s", name)))
	if err != nil {
		return nil, err
	}
	var api API
	if err := json.NewDecoder(resp.Body).Decode(&api); err != nil {
		return nil, err
	}
	return &api, nil
}

func (kong *KongGateway) createOrUpdateAPI(apiConf *API) (*API, error) {
	api, err := kong.getAPI(apiConf.Name)
	if err != nil {
		return nil, err
	}
	if api == nil {
		api, err = kong.createKongAPI(apiConf)
		if err != nil {
			return nil, err
		}
	}
	return api, nil
}

func (kong *KongGateway) addSelfAsTarget(upstream string, port int, weight int) (*upstreamTarget, error) {
	var target upstreamTarget

	ip, err := getServiceIP()
	if err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Add("target", fmt.Sprintf("%s:%d", ip, port))
	form.Add("weight", fmt.Sprintf("%d", weight))

	resp, err := kong.client.Post(kong.getKongURL(fmt.Sprintf("upstreams/%s/targets", upstream)), "multipart/form-data", strings.NewReader(form.Encode()))

	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(resp.Body).Decode(&target); err != nil {
		return nil, err
	}

	return &target, nil
}
