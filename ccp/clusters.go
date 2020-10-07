/*Copyright (c) 2019 Cisco and/or its affiliates.

This software is licensed to you under the terms of the Cisco Sample
Code License, Version 1.0 (the "License"). You may obtain a copy of the
License at

               https://developer.cisco.com/docs/licenses

All use of the material herein must be in accordance with the terms of
the License. All rights not expressly granted by the License are
reserved. Unless required by applicable law or agreed to separately in
writing, software distributed under the License is distributed on an "AS
IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
or implied.*/

package ccp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	validator "gopkg.in/validator.v2"
)

// Cluster v3 cluster
type Cluster struct {
	//  Cluster Variable Name in Struct
	//								Go Type			Reference in JSON
	UUID                 *string               `json:"id,omitempty"` //
	Type                 *string               `json:"type,omitempty"  `
	Name                 *string               `json:"name,omitempty"  validate:"nonzero"`
	InfraProviderUUID    *string               `json:"provider,omitempty" `
	Status               *string               `json:"status,omitempty" `
	KubernetesVersion    *string               `json:"kubernetes_version,omitempty" validate:"nonzero"`
	KubeConfig           *string               `json:"kubeconfig,omitempty"`
	IPAllocationMethod   *string               `json:"ip_allocation_method,omitempty" validate:"nonzero"`
	MasterVIP            *string               `json:"master_vip,omitempty"`
	LoadBalancerIPNum    *int64                `json:"load_balancer_num,omitempty"`
	SubnetUUID           *string               `json:"subnet_id,omitempty"`
	NTPPools             *[]string             `json:"ntp_pools,omitempty"`
	NTPServers           *[]string             `json:"ntp_servers,omitempty"`
	RegistriesRootCA     *[]string             `json:"root_ca_registries,omitempty"`
	RegistriesSelfSigned *RegistriesSelfSigned `json:"self_signed_registries,omitempty"`
	RegistriesInsecure   *[]string             `json:"insecure_registries,omitempty"`
	DockerProxyHTTP      *string               `json:"docker_http_proxy,omitempty"`
	DockerProxyHTTPS     *string               `json:"docker_https_proxy,omitempty"`
	DockerBIP            *string               `json:"docker_bip,omitempty"`
	Infra                *Infra                `json:"vsphere_infra,omitempty"  validate:"nonzero"`
	MasterNodePool       *MasterNodePool       `json:"master_group,omitempty"  validate:"nonzero" `
	WorkerNodePool       *[]WorkerNodePool     `json:"node_groups,omitempty"  validate:"nonzero" `
	NetworkPlugin        *NetworkPlugin        `json:"network_plugin_profile,omitempty" validate:"nonzero"`
	IngressAsLB          *bool                 `json:"ingress_as_lb,omitempty"`
	NginxIngressClass    *string               `json:"nginx_ingress_class,omitempty"`
	ETCDEncrypted        *bool                 `json:"etcd_encrypted,omitempty"`
	SkipManagement       *bool                 `json:"skip_management,omitempty"`
	DockerNoProxy        *[]string             `json:"docker_no_proxy,omitempty"`
	RoutableCIDR         *string               `json:"routable_cidr,omitempty"`
	ImagePrefix          *string               `json:"image_prefix,omitempty"`
	ACIProfileUUID       *string               `json:"aci_profile,omitempty"`
	Description          *string               `json:"description,omitempty"`
	AWSIamEnabled        *bool                 `json:"aws_iam_enabled,omitempty"`
}

// WorkerNodePool are the worker nodes - updated for v3
type WorkerNodePool struct {
	Name              *string   `json:"name,omitempty" validate:"nonzero"`      //v3
	Size              *int64    `json:"size,omitempty" validate:"nonzero"`      //v3
	Template          *string   `json:"template,omitempty" validate:"nonzero"`  //v2
	VCPUs             *int64    `json:"vcpus,omitempty" validate:"nonzero"`     //v2
	Memory            *int64    `json:"memory_mb,omitempty" validate:"nonzero"` //v2
	GPUs              *[]string `json:"gpus,omitempty"`                         //v3
	SSHUser           *string   `json:"ssh_user,omitempty"`                     //v3
	SSHKey            *string   `json:"ssh_key,omitempty"`                      //v3
	Nodes             *[]Node   `json:"nodes,omitempty"`                        //v3
	KubernetesVersion *string   `json:"kubernetes_version,omitempty"`           //v3
}

// RegistriesSelfSigned v3
type RegistriesSelfSigned struct {
	Cert *string `json:"selfsignedca,omitempty" `
}

// Infra updated for v3
type Infra struct { // checked for v3
	Datacenter   *string   `json:"datacenter,omitempty"  validate:"nonzero"`
	Datastore    *string   `json:"datastore,omitempty"  validate:"nonzero"`
	Cluster      *string   `json:"cluster,omitempty" validate:"nonzero"`
	Networks     *[]string `json:"networks,omitempty"  validate:"nonzero"`
	ResourcePool *string   `json:"resource_pool,omitempty"`
}

// MasterNodePool updated for v3
type MasterNodePool struct {
	Name              *string   `json:"name,omitempty"`                         // v2
	Size              *int64    `json:"size,omitempty"`                         // v2
	Template          *string   `json:"template,omitempty" validate:"nonzero"`  //v3
	VCPUs             *int64    `json:"vcpus,omitempty" validate:"nonzero"`     //v3
	Memory            *int64    `json:"memory_mb,omitempty" validate:"nonzero"` //v3
	GPUs              *[]string `json:"gpus,omitempty"`                         //v3
	SSHUser           *string   `json:"ssh_user,omitempty"`                     //v3
	SSHKey            *string   `json:"ssh_key,omitempty"`                      //v3
	Nodes             *[]Node   `json:"nodes,omitempty"`                        //v3
	KubernetesVersion *string   `json:"kubernetes_version,omitempty"`           //v3
}

// Node updated for v3
type Node struct {
	// v3 clusters
	Name         *string `json:"name,omitempty"`
	Status       *string `json:"status,omitempty"`
	StatusDetail *string `json:"status_detail,omitempty" `
	StatusReason *string `json:"status_reason,omitempty" `
	PublicIP     *string `json:"public_ip,omitempty"`
	PrivateIP    *string `json:"private_ip,omitempty"`
	Phase        *string `json:"phase,omitempty"`
	//	State        *string `json:"status,omitempty"`

}

// Label updated for v3
type Label struct {
	Key   *string `json:"key,omitempty"`
	Value *string `json:"value,omitempty"`
}

// Deployer updated for v3
type Deployer struct {
	ProxyCMD     *string   `json:"proxy_cmd,omitempty"`
	ProviderType *string   `json:"provider_type,omitempty" validate:"nonzero"`
	Provider     *Provider `json:"provider,omitempty" validate:"nonzero"`
}

// NetworkPlugin now caters for PluginDetails
type NetworkPlugin struct {
	Name    *string               `json:"name,omitempty"`
	Details *NetworkPluginDetails `json:"details,omitempty"`
}

// NetworkPluginDetails updated for v3
type NetworkPluginDetails struct {
	PodCIDR *string `json:"pod_cidr,omitempty"`
}

// Provider vsphere provider for v2
type Provider struct {
	VsphereDataCenter         *string              `json:"vsphere_datacenter,omitempty"`
	VsphereDatastore          *string              `json:"vsphere_datastore,omitempty"`
	VsphereSCSIControllerType *string              `json:"vsphere_scsi_controller_type,omitempty"`
	VsphereWorkingDir         *string              `json:"vsphere_working_dir,omitempty"`
	VsphereClientConfigUUID   *string              `json:"vsphere_client_config_uuid,omitempty" validate:"nonzero"`
	ClientConfig              *VsphereClientConfig `json:"client_config,omitempty"`
}

// VsphereClientConfig for provider
type VsphereClientConfig struct {
	IP       *string `json:"ip,omitempty"`
	Port     *int64  `json:"port,omitempty"`
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty"`
}

// AddonsCatalogue List of all of the CCP Add-Ons
// generated from https://mholt.github.io/json-to-go/
// with future versions of CCP this may need to be re-generated with updated JSON catalogue from
// path /v3/clusters/<clusteruuid>/catalog
// this is working for CCP 6.x
type AddonsCatalogue struct {
	CcpMonitor struct {
		DisplayName string `json:"displayName"`
		Name        string `json:"name"`
		Namespace   string `json:"namespace"`
		Description string `json:"description"`
		URL         string `json:"url"`
	} `json:"_ccp-monitor"`
	CcpEfk struct {
		DisplayName string `json:"displayName"`
		Name        string `json:"name"`
		Namespace   string `json:"namespace"`
		Description string `json:"description"`
		URL         string `json:"url"`
	} `json:"_ccp-efk"`
	CcpKubernetesDashboard struct {
		DisplayName   string   `json:"displayName"`
		Name          string   `json:"name"`
		Namespace     string   `json:"namespace"`
		Description   string   `json:"description"`
		URL           string   `json:"url"`
		OverrideFiles []string `json:"overrideFiles"`
	} `json:"_ccp-kubernetes-dashboard"`
	CcpIstioOperator struct {
		DisplayName  string   `json:"displayName"`
		Name         string   `json:"name"`
		Namespace    string   `json:"namespace"`
		Description  string   `json:"description"`
		URL          string   `json:"url"`
		Conflicts    []string `json:"conflicts"`
		Dependencies struct {
			CcpIstio struct {
				DisplayName string `json:"displayName"`
				Name        string `json:"name"`
				Namespace   string `json:"namespace"`
				Description string `json:"description"`
				URL         string `json:"url"`
			} `json:"_ccp-istio"`
		} `json:"dependencies"`
	} `json:"_ccp-istio-operator"`
	CcpHarborOperator struct {
		DisplayName  string   `json:"displayName"`
		Name         string   `json:"name"`
		Namespace    string   `json:"namespace"`
		Description  string   `json:"description"`
		URL          string   `json:"url"`
		Conflicts    []string `json:"conflicts"`
		Dependencies struct {
			CcpHarbor struct {
				DisplayName string `json:"displayName"`
				Name        string `json:"name"`
				Namespace   string `json:"namespace"`
				Description string `json:"description"`
				URL         string `json:"url"`
			} `json:"_ccp-harbor"`
		} `json:"dependencies"`
	} `json:"_ccp-harbor-operator"`
	CcpKubeflow struct {
		Name        string   `json:"name"`
		Namespace   string   `json:"namespace"`
		DisplayName string   `json:"displayName"`
		Description string   `json:"description"`
		URL         string   `json:"url"`
		Conflicts   []string `json:"conflicts"`
		Overrides   string   `json:"overrides"`
	} `json:"_ccp-kubeflow"`
	CcpHxcsi struct {
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
		Description string `json:"description"`
		URL         string `json:"url"`
		Overrides   string `json:"overrides"`
		Namespace   string `json:"namespace,omitempty"`
	} `json:"_ccp-hxcsi"`
}

// ClusterInstalledAddons list of installed AddOn
type ClusterInstalledAddons struct {
	Count    int64      `json:"count"`
	Next     int64      `json:"next"`
	Previous int64      `json:"previous"`
	Results  []struct { // results
		Name        string   `json:"name"`
		Namespace   string   `json:"namespace"`
		DisplayName string   `json:"displayName"`
		Description string   `json:"description"`
		AddonStatus struct { // status
			Name       string `json:"name"`
			HelmStatus string `json:"helmStatus"`
			Status     string `json:"status"`
		} `json:"status"`
	} `json:"results"`
}

// Results - results for AddOns that have been installed or are installing
type Results struct {
	Name          *string   `json:"name,omitempty"`
	Namespace     *string   `json:"namespace,omitempty"`
	Overrides     *string   `json:"overrides,omitempty"`
	OverrideFiles *[]string `json:"overrideFiles,omitempty"`
	URL           *string   `json:"url,omitempty"`
	DisplayName   *string   `json:"displayName,omitempty"`
	AddonStatus   *Status   `json:"status,omitempty"`
}

// Status - status for AddOns that have been installed or are installing
type Status struct {
	Name             *string `json:"name,omitempty"`
	HelmStatus       *string `json:"helmStatus,omitempty"`
	Status           *string `json:"status,omitempty"`
	OverrideHash     *string `json:"overrideHash,omitempty"`
	StatusDetail     *string `json:"statusDetail,omitempty"`
	URLInstalled     *string `json:"urlInstalled,omitempty"`
	VersionInstalled *string `json:"versionInstalled,omitempty"`
	ServiceURL       *string `json:"serviceUrl,omitempty"`
}

// GetClusters function for v3
func (s *Client) GetClusters() ([]Cluster, error) {
	Debug(1, "GetClusters")

	url := s.BaseURL + "/v3/clusters"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	bytes, err := s.doRequest(req)
	if err != nil {
		return nil, err
	}

	Debug(3, "Cluster JSON Payload:\n"+string(bytes))

	// Create an Array of Clusters
	var data []Cluster

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	// Print out list of Clusters and their index
	Debug(2, "Found "+strconv.Itoa(len(data))+" clusters")
	for i, cl := range data {
		Debug(2, "Found cluster "+strconv.Itoa(i)+" named "+*cl.Name+" with UUID "+*cl.UUID)
	}

	return data, nil
}

// GetClusterStatusByName get all clusters, iterate through to find slice matching clusterName
func (s *Client) GetClusterStatusByName(clusterName string) (*string, error) {
	Debug(1, "GetClusterStatusByName")

	clusters, err := s.GetClusters()
	if err != nil {
		return nil, err
	}

	for i, x := range clusters {
		Debug(3, "Iteration "+strconv.Itoa(i)+" Cluster found: "+string(*x.Name)+"\n")
		if string(clusterName) == string(*x.Name) {
			Debug(2, "Found matching cluster "+clusterName+" = "+*x.Name)
			return x.Status, nil
		}
	}
	return nil, errors.New("Cannot find cluster " + clusterName)
}

// GetClusterByName get all clusters, iterate through to find slice matching clusterName
func (s *Client) GetClusterByName(clusterName string) (*Cluster, error) {
	Debug(1, "GetClusterByName")

	clusters, err := s.GetClusters()
	if err != nil {
		return nil, err
	}

	for i, x := range clusters {
		Debug(3, "Iteration "+strconv.Itoa(i)+" Cluster found: "+string(*x.Name)+"\n")
		if string(clusterName) == string(*x.Name) {
			Debug(2, "Found matching cluster "+clusterName+" = "+*x.Name)
			return &x, nil
		}
	}
	return nil, errors.New("Cannot find cluster " + clusterName)
}

// GetClusterByUUID v3 cluster by UUID
func (s *Client) GetClusterByUUID(clusterUUID string) (*Cluster, error) {
	Debug(1, "GetClusterByUUID")

	url := fmt.Sprintf(s.BaseURL + "/v3/clusters/" + clusterUUID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	bytes, err := s.doRequest(req)
	if err != nil {
		return nil, err
	}
	var data *Cluster

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// ScaleCluster spec for JSON scale
type ScaleCluster struct {
	Name *string `json:"name" validate:"nonzero"`
	Size *int    `json:"size" validate:"nonzero"`
}

// ScaleCluster scales an existing cluster
func (s *Client) ScaleCluster(clusterUUID, workerPoolName string, size int) (*Cluster, error) {

	Debug(1, "Func: ScaleCluster")

	url := s.BaseURL + "/v3/clusters/" + clusterUUID + "/node-pools/" + workerPoolName + "/"
	Debug(2, "PATCH URL: "+url)

	cluserScale := ScaleCluster{
		Name: String(workerPoolName),
		Size: Int(size),
	}
	j, err := json.Marshal(cluserScale)
	if err != nil {
		return nil, err
	}

	Debug(3, "Sending JSON patch: "+string(j))

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(j))
	if err != nil {
		Debug(1, "http.NewRequest PATCH error "+err.Error())
		return nil, err
	}

	bytes, err := s.doRequest(req)
	if err != nil {
		Debug(1, "http.doRequest error "+err.Error())
		Debug(1, string(bytes))
		return nil, err
	}

	var data Cluster
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// ConvertJSONToCluster convers JSON
func (s *Client) ConvertJSONToCluster(jsonFile string) (*Cluster, error) {
	Debug(1, "Entered ConvertJSONToCluster")

	// Debug(2, "Cluster Struct for cluster named "+string(*cluster.Name))
	jsonBody, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return nil, err
	}

	var newCluster Cluster
	err = json.Unmarshal([]byte(jsonBody), &newCluster)

	return &newCluster, nil
}

// AddClusterOld creates a new cluster without much error checking
func (s *Client) AddClusterOld(cluster *Cluster) (*Cluster, error) {
	Debug(1, "Entered AddCluster for "+string(*cluster.Name))

	Debug(2, "Start validating Cluster struct")
	errs := validator.Validate(cluster)
	if errs != nil {
		Debug(1, "Errors validating Cluster struct with validator.Validate(): "+string(errs.Error()))
		return nil, errs
	}
	Debug(3, "No Errors validating Cluster struct")

	url := s.BaseURL + "/v3/clusters/"

	j, err := json.Marshal(cluster)
	if err != nil {
		Debug(1, "Errors marshaling with json.Marshal(): "+string(err.Error()))
		return nil, err
	}
	Debug(3, "No errors Marshaling JSON")

	Debug(2, "About to POST to url "+url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	if err != nil {
		Debug(1, "Errors POSTing with http.NewRequest: "+string(err.Error()))
		return nil, err
	}

	bytes, err := s.doRequest(req)
	if err != nil {
		Debug(1, "Errors POSTing with s.doRequest: "+string(err.Error()))
		Debug(1, "POST response: "+string(bytes))
		return nil, err
	}
	Debug(3, "POST response: "+string(bytes))

	var data Cluster

	// err = json.Unmarshal(bytes, &data)
	Debug(2, "Unmarshaling response")
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		Debug(1, "Errors unmarshaling with json.Unmarshal: "+string(err.Error()))
		return nil, err
	}
	Debug(2, "Unmarshaled response successfully")

	Debug(2, "CCP API responded with JSON payload for cluster named "+*data.Name+" with UUID "+*data.UUID)
	if *data.UUID == "" {
		Debug(1, "CCP API created cluster named "+*data.Name+" with UUID "+*data.UUID)
	}
	return &data, nil
}

// AddCluster creates a new cluster with error checking (Conor Murphy updates)
func (s *Client) AddCluster(cluster *Cluster) (*Cluster, error) {

	Debug(1, "Entered AddCluster for "+string(*cluster.Name))

	Debug(2, "Start validating Cluster struct")
	errs := validator.Validate(cluster)
	if errs != nil {
		Debug(1, "Errors validating Cluster struct with validator.Validate(): "+string(errs.Error()))
		return nil, errs
	}
	Debug(3, "No Errors validating Cluster struct")

	// https://stackoverflow.com/questions/44320960/omitempty-doesnt-omit-interface-nil-values-in-json
	// *cluster.MasterNodePool.Nodes returns &[] and since this is not nil, omitempty, won't omit it when we marshal. Instead it includes nodes: null
	// which CCP doesn't like. Therefore we need to check if it's empty and then set it to nil.

	if cluster.MasterNodePool.Nodes != nil { // this check is for the CCP clientlibrary / ccpctl only
		if len(*cluster.MasterNodePool.Nodes) == 0 { // this check is for the TF library creating a zero sized array
			cluster.MasterNodePool.Nodes = nil
		}
	}
	// Same as above but there can be multiple pools of worker nodes, therefore we need to iterate through, check if the nodes are empty,
	// and if so set them to nil
	var attr WorkerNodePool

	if cluster.WorkerNodePool != nil { // check if this exists, if it does then process it
		for i := 0; i < len(*cluster.WorkerNodePool); i++ {
			attr = (*cluster.WorkerNodePool)[i]

			if attr.Nodes != nil {
				if len(*attr.Nodes) == 0 {
					attr.Nodes = nil
				}
				(*cluster.WorkerNodePool)[i] = attr
			}
		}
	}

	if cluster.NTPPools != nil { // check if this exists, if it does then process it
		if len(*cluster.NTPPools) == 0 {
			cluster.NTPPools = nil
		}
	}

	if cluster.NTPServers != nil {
		if len(*cluster.NTPServers) == 0 {
			cluster.NTPServers = nil
		}
	}

	if cluster.DockerNoProxy != nil {
		if len(*cluster.DockerNoProxy) == 0 {
			cluster.DockerNoProxy = nil
		}
	}

	if cluster.RegistriesRootCA != nil {
		if len(*cluster.RegistriesRootCA) == 0 {
			cluster.RegistriesRootCA = nil
		}
	}

	if cluster.RegistriesInsecure != nil {
		if len(*cluster.RegistriesInsecure) == 0 {
			cluster.RegistriesInsecure = nil
		}
	}

	url := s.BaseURL + "/v3/clusters/"

	j, err := json.Marshal(&cluster)

	if err != nil {
		Debug(1, "Errors marshaling with json.Marshal(): "+string(err.Error()))
		return nil, err
	}
	Debug(3, "No errors Marshaling JSON")

	Debug(2, "About to POST to url "+url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	if err != nil {
		Debug(1, "Errors POSTing with http.NewRequest: "+string(err.Error()))
		return nil, err
	}

	bytes, err := s.doRequest(req)
	if err != nil {
		Debug(1, "Errors POSTing with s.doRequest: "+string(err.Error()))
		Debug(1, "POST response: "+string(bytes))
		return nil, err
	}
	Debug(3, "POST response: "+string(bytes))

	var data Cluster

	// err = json.Unmarshal(bytes, &data)
	Debug(2, "Unmarshaling response")
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		Debug(1, "Errors unmarshaling with json.Unmarshal: "+string(err.Error()))
		return nil, err
	}
	Debug(2, "Unmarshaled response successfully")

	Debug(2, "CCP API responded with JSON payload for cluster named "+*data.Name+" with UUID "+*data.UUID)
	if *data.UUID == "" {
		Debug(1, "CCP API created cluster named "+*data.Name+" with UUID "+*data.UUID)
	}

	return &data, nil
}

// AddClusterSynchronous creates a new cluster but waits until the cluster is created before returning
func (s *Client) AddClusterSynchronous(cluster *Cluster) (*Cluster, error) {

	errs := validator.Validate(cluster)
	if errs != nil {
		return nil, errs
	}

	// https://stackoverflow.com/questions/44320960/omitempty-doesnt-omit-interface-nil-values-in-json
	// *cluster.MasterNodePool.Nodes returns &[] and since this is not nil, omitempty, won't omit it when we marshal. Instead it includes nodes: null
	// which CCP doesn't like. Therefore we need to check if it's empty and then set it to nil.

	if cluster.MasterNodePool.Nodes != nil { // this check is for the CCP clientlibrary / ccpctl only
		if len(*cluster.MasterNodePool.Nodes) == 0 { // this check is for the TF library creating a zero sized array
			cluster.MasterNodePool.Nodes = nil
		}
	}
	// Same as above but there can be multiple pools of worker nodes, therefore we need to iterate through, check if the nodes are empty,
	// and if so set them to nil
	var attr WorkerNodePool

	if cluster.WorkerNodePool != nil { // check if this exists, if it does then process it
		for i := 0; i < len(*cluster.WorkerNodePool); i++ {
			attr = (*cluster.WorkerNodePool)[i]

			if attr.Nodes != nil {
				if len(*attr.Nodes) == 0 {
					attr.Nodes = nil
				}
				(*cluster.WorkerNodePool)[i] = attr
			}
		}
	}

	if cluster.NTPPools != nil { // check if this exists, if it does then process it
		if len(*cluster.NTPPools) == 0 {
			cluster.NTPPools = nil
		}
	}

	if cluster.NTPServers != nil {
		if len(*cluster.NTPServers) == 0 {
			cluster.NTPServers = nil
		}
	}

	if cluster.DockerNoProxy != nil {
		if len(*cluster.DockerNoProxy) == 0 {
			cluster.DockerNoProxy = nil
		}
	}

	if cluster.RegistriesRootCA != nil {
		if len(*cluster.RegistriesRootCA) == 0 {
			cluster.RegistriesRootCA = nil
		}
	}

	if cluster.RegistriesInsecure != nil {
		if len(*cluster.RegistriesInsecure) == 0 {
			cluster.RegistriesInsecure = nil
		}
	}

	// Need this because even though we don't have the networks setting configured for contiv-aci, if we leave it
	// out then CCP will complain it's empty. We also can't marshall an empty slice as it converts to null which doesn't work
	// see here for details https://apoorvam.github.io/blog/2017/golang-json-marshal-slice-as-empty-array-not-null
	// Also, it seems CCP complains when load_balancer_ip_num is missing even though it's not used for the aci cni

	if cluster.NetworkPlugin.Name != nil {
		if *cluster.NetworkPlugin.Name == "contiv-aci" {
			*cluster.Infra.Networks = make([]string, 0)
			*cluster.LoadBalancerIPNum = 1
		}
	}
	url := s.BaseURL + "/v3/clusters/"

	j, err := json.Marshal(&cluster)

	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}

	bytes, err := s.doRequest(req)
	if err != nil {
		return nil, err
	}
	var data Cluster

	// err = json.Unmarshal(bytes, &data)
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	time.Sleep(10 * time.Second)

	status, err := s.GetClusterStatusByName(*cluster.Name)

	if err != nil {
		return nil, err
	}

	for *status == "CREATING" {

		status, err = s.GetClusterStatusByName(*cluster.Name)

		if err != nil {
			return nil, err
		}

		time.Sleep(10 * time.Second)

	}

	return &data, nil

}

// DeleteCluster deletes a cluster
func (s *Client) DeleteCluster(clusterUUID string) error {
	Debug(1, "Entered DeleteCluster for UUID "+clusterUUID)

	if clusterUUID == "" {
		return errors.New("Cluster UUID to delete is required")
	}

	url := s.BaseURL + "/v3/clusters/" + clusterUUID + "/"

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	_, err = s.doRequest(req)
	if err != nil {
		return err
	}

	Debug(2, "Request sent to API with success response")
	return nil
}

// SetDebug sets the debug level
func (s *Client) SetDebug(debug int) {
	debuglvl = debug
	Debug(1, "Debug level set to "+string(debuglvl))
}

// GetKubeVerFromImage splits the image name and gets the kube ver
func GetKubeVerFromImage(value string) string {
	// https://www.dotnetperls.com/between-before-after-go
	// Get substring between two strings.
	a := "image-"
	b := "-ubuntu18"
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}

// AddClusterBasic add a v3 cluster the easy way
func (s *Client) AddClusterBasic(cluster *Cluster) (*Cluster, error) {
	Debug(1, "Entered AddClusterBasic for cluster "+string(*cluster.Name))
	/*

		This function was added in order to provide users a better experience with adding clusters. The list of required
		fields has been shortend with all defaults and computed values such as UUIDs to be automatically configured on behalf of the user.

		The following fields and values will be configured. The remainder to be specified by the user

		ProviderClientConfigUUID
		KubernetesVersion - default will be set to 1.10.1
		Type - default will be set to 1
		Deployer
			ProviderType will be set to "vsphere"
			Provider
				VsphereDataCenter - already specified as part of Cluster struct so will use this same value
				VsphereClientConfigUUID
				VsphereDatastore - already specified as part of Cluster struct so will use this same value
				VsphereWorkingDir - default will be set to /VsphereDataCenter/vm
		NetworkPlugin
			Name - default will be set to contiv-vpp
			Status - default will be set to ""
			Details - default will be set to "{\"pod_cidr\":\"192.168.0.0/16\"}"
		WorkerNodePool
			VCPUs - default will be set to 2
			Memory - default will be set to 16384
		MasterNodePool
			VCPUs - default will be set to 2
			Memory - default will be set to 8192

	*/

	var data Cluster

	// The following will configured the defaults for the cluster as specified above as well as check that the minimum
	// fields are provided

	if nonzero(cluster.Name) {
		return nil, errors.New("Cluster.Name is missing")
	}
	if nonzero(cluster.Infra.Datacenter) {
		return nil, errors.New("Cluster.Infra.Datacenter is missing")
	}

	// if nonzero(cluster.ResourcePool) {
	// 	return nil, errors.New("Cluster.ResourcePool is missing")
	// }
	if nonzero(cluster.MasterNodePool.SSHUser) {
		return nil, errors.New("cluster.MasterNodePool.SSHUser is missing")
	}
	if nonzero(cluster.MasterNodePool.SSHKey) {
		return nil, errors.New("cluster.MasterNodePool.SSHKey is missing")
	}

	// loop over array of WorkerNodePool
	for _, v := range *cluster.WorkerNodePool {

		if nonzero(v.SSHUser) {
			return nil, errors.New("v.SSHUser is missing")
		}

		if nonzero(v.SSHKey) {
			return nil, errors.New("v.SSHKey is missing")
		}
		if nonzero(v.Size) {
			return nil, errors.New("v.Size is missing")
		}
		if nonzero(v.Template) {
			return nil, errors.New("v.Template is missing")
		}
	}

	if nonzero(cluster.MasterNodePool.Size) {
		return nil, errors.New("cluster.MasterNodePool.Size is missing")
	}

	if nonzero(cluster.MasterNodePool.Template) {
		return nil, errors.New("cluster.MasterNodePool.Template is missing")
	}

	// check that cluster.MasterNodePool.Template and cluster.WorkerNodePool.Template are the same

	// Retrieve the provider client config UUID rather than have the user need to provide this themselves.
	// This is also built for a single provider client config and as of CCP 1.5 this wll be Vsphere
	providerClientConfigs, err := s.GetInfraProviderByName("vsphere")
	if err != nil {
		return nil, err
	}

	networkPlugin := NetworkPlugin{
		Name: String("calico"),
		// Details: String("{\"pod_cidr\":\"192.168.0.0/16\"}"),
		Details: &NetworkPluginDetails{
			PodCIDR: String("192.168.0.0/16"),
		},
	}

	workerNodePool := WorkerNodePool{
		Size:     Int64(1),
		VCPUs:    Int64(2),
		Memory:   Int64(32768),
		Template: String(*cluster.MasterNodePool.Template), // use same template as master
	}

	masterNodePool := MasterNodePool{
		Size:     Int64(1),
		VCPUs:    Int64(2),
		Memory:   Int64(16384),
		Template: String(*cluster.MasterNodePool.Template),
	}

	// Since it returns a list we will use the UUID from the first element
	cluster.InfraProviderUUID = String(*providerClientConfigs.UUID)
	// cluster.KubernetesVersion = String("1.16.3") // todo: fetch this somehow
	// below should work, but issues with *string / string types
	// kubever := GetKubeVerFromImage(*cluster.MasterNodePool.Template)
	// cluster.KubernetesVersion = &kubever
	cluster.KubernetesVersion = String(GetKubeVerFromImage(*cluster.MasterNodePool.Template))

	// cluster.Type = Int64(1)
	cluster.NetworkPlugin = &networkPlugin
	// cluster.Deployer = &deployer

	//	cluster.WorkerNodePool = &workerNodePool
	cluster.WorkerNodePool = &[]WorkerNodePool{workerNodePool}

	cluster.MasterNodePool = &masterNodePool

	url := s.BaseURL + "/v3/clusters"

	j, err := json.Marshal(cluster)

	if err != nil {

		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}

	bytes, err := s.doRequest(req)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &data)

	if err != nil {
		return nil, err
	}

	cluster = &data

	return cluster, nil
}

// InstallAddonIstioOp Installs the Istio Operator
func (s *Client) InstallAddonIstioOp(clusterUUID string) error {
	Debug(1, "Entered InstallAddonIstio for UUID "+clusterUUID)

	if clusterUUID == "" {
		return errors.New("Cluster UUID is required")
	}

	url := s.BaseURL + "/v3/clusters/" + clusterUUID + "/addons/"

	jsonBody := []byte(`
	{
        "displayName": "Istio Operator",
        "name": "ccp-istio-operator",
        "namespace": "ccp",
        "description": "Istio Operator",
        "url": "/opt/ccp/charts/ccp-istio-operator.tgz",
        "conflicts": [
            "ccp-kubeflow",
            "ccp-harbor-operator"
        ],
        "dependencies": {
            "_ccp-istio": {
                "displayName": "Istio",
                "name": "ccp-istio-cr",
                "namespace": "ccp",
                "description": "Istio (REQUIRES ISTIO OPERATOR)",
                "url": "/opt/ccp/charts/ccp-istio-cr.tgz"
            }
		}
	}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	resp, err := s.doRequest(req)
	if err != nil {
		return err
	}
	Debug(3, "Response is:")
	Debug(3, string(resp))

	Debug(2, "Request sent to API with success response")
	return nil
}

// InstallAddonIstioInstance Installs the Istio Instance (install the Operator first)
func (s *Client) InstallAddonIstioInstance(clusterUUID string) error {
	Debug(1, "Entered InstallAddonIstioInstance for UUID "+clusterUUID)

	if clusterUUID == "" {
		return errors.New("Cluster UUID is required")
	}

	url := s.BaseURL + "/v3/clusters/" + clusterUUID + "/addons/"

	jsonBody := []byte(`
	{
		"displayName": "Istio",
		"name": "ccp-istio-cr",
		"namespace": "ccp",
		"description": "Istio (REQUIRES ISTIO OPERATOR)",
		"url": "/opt/ccp/charts/ccp-istio-cr.tgz"
	}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	resp, err := s.doRequest(req)
	if err != nil {
		return err
	}
	Debug(3, "Response is:")
	Debug(3, string(resp))

	Debug(2, "Request sent to API with success response")
	return nil
}

// InstallAddonIstio install both
func (s *Client) InstallAddonIstio(clusterUUID string) error {
	err := s.InstallAddonIstioOp(clusterUUID)
	if err != nil {
		Debug(1, "Failed to add Add-On Istio Operator: "+string(err.Error()))
		return err
	}
	time.Sleep(2 * time.Second) // wait 2 seconds before sending the next request
	err = s.InstallAddonIstioInstance(clusterUUID)
	if err != nil {
		Debug(1, "Failed to add Add-On Istio Instance: "+string(err.Error()))
		return err
	}
	return nil
}

// InstallAddonDashboard Installs the Istio Instance (install the Operator first)
func (s *Client) InstallAddonDashboard(clusterUUID string) error {
	Debug(1, "Entered InstallAddonDashboard for UUID "+clusterUUID)

	if clusterUUID == "" {
		return errors.New("Cluster UUID is required")
	}

	url := s.BaseURL + "/v3/clusters/" + clusterUUID + "/addons/"

	jsonBody := []byte(`
	{
		"displayName": "Dashboard",
		"name": "kubernetes-dashboard",
		"namespace": "ccp",
		"description": "Dashboard",
		"url": "/opt/ccp/charts/kubernetes-dashboard.tgz",
		"overrideFiles": [
			"/opt/ccp/charts/kubernetes-dashboard.yaml"
		]
	}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	resp, err := s.doRequest(req)
	if err != nil {
		return err
	}
	Debug(3, "Response is:")
	Debug(3, string(resp))

	Debug(2, "Request sent to API with success response")
	return nil
}

// InstallAddonMonitoring Installs the Istio Instance (install the Operator first)
func (s *Client) InstallAddonMonitoring(clusterUUID string) error {
	Debug(1, "Entered InstallAddonMonitoring for UUID "+clusterUUID)

	if clusterUUID == "" {
		return errors.New("Cluster UUID is required")
	}

	url := s.BaseURL + "/v3/clusters/" + clusterUUID + "/addons/"

	jsonBody := []byte(`
	{
		"displayName": "Monitoring",
		"name": "ccp-monitor",
		"namespace": "ccp",
		"description": "Monitoring",
		"url": "/opt/ccp/charts/ccp-monitor.tgz"
	}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	resp, err := s.doRequest(req)
	if err != nil {
		return err
	}
	Debug(3, "Response is:")
	Debug(3, string(resp))

	Debug(2, "Request sent to API with success response")
	return nil
}

// InstallAddonLogging Installs the Istio Instance (install the Operator first)
func (s *Client) InstallAddonLogging(clusterUUID string) error {
	Debug(1, "Entered InstallAddonLogging for UUID "+clusterUUID)

	if clusterUUID == "" {
		return errors.New("Cluster UUID is required")
	}

	url := s.BaseURL + "/v3/clusters/" + clusterUUID + "/addons/"

	jsonBody := []byte(`
	{
		"displayName": "Logging",
		"name": "ccp-efk",
		"namespace": "ccp",
		"description": "Logging",
		"url": "/opt/ccp/charts/ccp-efk.tgz"
	}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	resp, err := s.doRequest(req)
	if err != nil {
		return err
	}
	Debug(3, "Response is:")
	Debug(3, string(resp))

	Debug(2, "Request sent to API with success response")
	return nil
}

// InstallAddonHarborOp Installs the Istio Instance (install the Operator first)
func (s *Client) InstallAddonHarborOp(clusterUUID string) error {
	Debug(1, "Entered InstallAddonHarborOp for UUID "+clusterUUID)

	if clusterUUID == "" {
		return errors.New("Cluster UUID is required")
	}

	url := s.BaseURL + "/v3/clusters/" + clusterUUID + "/addons/"

	jsonBody := []byte(`
	{
        "displayName": "Harbor Operator",
        "name": "ccp-harbor-operator",
        "namespace": "ccp",
        "description": "Harbor Operator",
        "url": "/opt/ccp/charts/ccp-harbor-operator.tgz",
        "conflicts": [
            "ccp-istio-operator"
		]
	}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	resp, err := s.doRequest(req)
	if err != nil {
		return err
	}
	Debug(3, "Response is:")
	Debug(3, string(resp))

	Debug(2, "Request sent to API with success response")
	return nil
}

// InstallAddonHarborInstance Installs the Istio Instance (install the Operator first)
func (s *Client) InstallAddonHarborInstance(clusterUUID string) error {
	Debug(1, "Entered InstallAddonHarborInstance for UUID "+clusterUUID)

	if clusterUUID == "" {
		return errors.New("Cluster UUID is required")
	}

	url := s.BaseURL + "/v3/clusters/" + clusterUUID + "/addons/"

	jsonBody := []byte(`
	{
		"displayName": "Harbor",
		"name": "ccp-harbor-cr",
		"namespace": "ccp",
		"description": "Harbor registry",
		"url": "/opt/ccp/charts/ccp-harbor-cr.tgz"
	}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	resp, err := s.doRequest(req)
	if err != nil {
		return err
	}
	Debug(3, "Response is:")
	Debug(3, string(resp))

	Debug(2, "Request sent to API with success response")
	return nil
}

// InstallAddonHarbor install both
func (s *Client) InstallAddonHarbor(clusterUUID string) error {
	err := s.InstallAddonHarborOp(clusterUUID)
	if err != nil {
		Debug(1, "Failed to add Add-On Istio Operator: "+string(err.Error()))
		return err
	}
	time.Sleep(2 * time.Second) // wait 2 seconds before sending the next request
	err = s.InstallAddonHarborInstance(clusterUUID)
	if err != nil {
		Debug(1, "Failed to add Add-On Istio Instance: "+string(err.Error()))
		return err
	}
	return nil
}

// InstallAddon installs addon and waits for it to finish
func (s *Client) InstallAddon(clusterUUID string, addonName string) error {

	var jsonBody []byte

	if clusterUUID == "" {
		return errors.New("Cluster UUID is required")
	}

	switch addonName {

	case "kubernetes-dashboard", "dashboard":
		jsonBody = []byte(`
		{
			"displayName": "Dashboard",
			"name": "kubernetes-dashboard",
			"namespace": "ccp",
			"description": "Dashboard",
			"url": "/opt/ccp/charts/kubernetes-dashboard.tgz",
			"overrideFiles": [
				"/opt/ccp/charts/kubernetes-dashboard.yaml"
			]
		}`)

		err = s.InstallAddonAndWaitUntilInstalled(clusterUUID, addonName, jsonBody)

		if err != nil {
			return err
		}

	case "ccp-efk", "logging", "efk":
		jsonBody = []byte(`
		{
			"displayName": "Logging",
			"name": "ccp-efk",
			"namespace": "ccp",
			"description": "Logging",
			"url": "/opt/ccp/charts/ccp-efk.tgz"
		}`)

		err = s.InstallAddonAndWaitUntilInstalled(clusterUUID, addonName, jsonBody)

		if err != nil {
			return err
		}

	case "ccp-monitor", "monitor", "grafana":
		jsonBody = []byte(`
		{
			"displayName": "Monitoring",
			"name": "ccp-monitor",
			"namespace": "ccp",
			"description": "Monitoring",
			"url": "/opt/ccp/charts/ccp-monitor.tgz"
		}`)

		err = s.InstallAddonAndWaitUntilInstalled(clusterUUID, addonName, jsonBody)

		if err != nil {
			return err
		}

	case "istio":
		jsonBody = []byte(`
		{
			"displayName": "Istio Operator",
			"name": "ccp-istio-operator",
			"namespace": "ccp",
			"description": "Istio Operator",
			"url": "/opt/ccp/charts/ccp-istio-operator.tgz",
			"conflicts": [
				"ccp-kubeflow",
				"ccp-harbor-operator"
			],
			"dependencies": {
				"_ccp-istio": {
					"displayName": "Istio",
					"name": "ccp-istio-cr",
					"namespace": "ccp",
					"description": "Istio (REQUIRES ISTIO OPERATOR)",
					"url": "/opt/ccp/charts/ccp-istio-cr.tgz"
				}
			}
		}`)

		err = s.InstallAddonAndWaitUntilInstalled(clusterUUID, "ccp-istio-operator", jsonBody)

		if err != nil {
			return err
		}

		err = s.InstallAddonIstioInstance(clusterUUID)
		if err != nil {
			return err
		}

	case "harbor":
		jsonBody = []byte(`
		{
			"displayName": "Harbor Operator",
			"name": "ccp-harbor-operator",
			"namespace": "ccp",
			"description": "Harbor Operator",
			"url": "/opt/ccp/charts/ccp-harbor-operator.tgz",
			"conflicts": [
				"ccp-istio-operator"
			]
		}`)

		err = s.InstallAddonAndWaitUntilInstalled(clusterUUID, "ccp-harbor-operator", jsonBody)

		if err != nil {
			return err
		}

		err = s.InstallAddonHarborInstance(clusterUUID)
		if err != nil {
			return err
		}

	case "ccp-kubeflow", "kubeflow":

		jsonBody, err = s.GetKubeflowAddonConfig(clusterUUID)

		if err != nil {
			return err
		}

		err = s.InstallAddonAndWaitUntilInstalled(clusterUUID, addonName, jsonBody)

		if err != nil {
			return err
		}

	default:
		return errors.New("Incorrect addon name provided. Options are 'kubernetes-dashboard', 'ccp-efk', 'ccp-monitor', 'ccp-istio', 'harbor', 'ccp-kubeflow'")

	}

	return nil
}

// InstallAddonAndWaitUntilInstalled install addon and wait until it's completed
func (s *Client) InstallAddonAndWaitUntilInstalled(clusterUUID string, addonName string, jsonBody []byte) error {

	url := s.BaseURL + "/v3/clusters/" + clusterUUID + "/addons/"

	addonInstalled, err := s.IsAddonInstalled(clusterUUID, addonName)

	if err != nil {
		return err
	}

	if !*addonInstalled {
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		if err != nil {
			return err
		}

		_, err = s.doRequest(req)
		if err != nil {
			return err
		}

	}

	addonInstalled, err = s.IsAddonInstalled(clusterUUID, addonName)

	for !*addonInstalled {

		addonInstalled, err = s.IsAddonInstalled(clusterUUID, addonName)

		if err != nil {
			return err
		}

		time.Sleep(5 * time.Second)

	}

	return nil
}

// DeleteAddonLogging deletes the addon
func (s *Client) DeleteAddonLogging(clusterUUID string) error {
	Debug(1, "Entered DeleteAddonLogging for UUID "+clusterUUID)

	if clusterUUID == "" {
		return errors.New("Cluster UUID to delete is required")
	}

	url := s.BaseURL + "/v3/clusters/" + clusterUUID + "/addons/ccp-efk/"
	Debug(2, "Sending HTTP delte to "+url)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	_, err = s.doRequest(req)
	if err != nil {
		return err
	}

	Debug(2, "Request sent to API with success response")
	return nil
}

// DeleteAddonMonitor deletes the addon
func (s *Client) DeleteAddonMonitor(clusterUUID string) error {
	Debug(1, "Entered DeleteAddonMonitor for UUID "+clusterUUID)

	if clusterUUID == "" {
		return errors.New("Cluster UUID to delete is required")
	}

	url := s.BaseURL + "/v3/clusters/" + clusterUUID + "/addons/ccp-monitor/"
	Debug(2, "Sending HTTP delte to "+url)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	_, err = s.doRequest(req)
	if err != nil {
		return err
	}

	Debug(2, "Request sent to API with success response")
	return nil
}

// DeleteAddonIstioInstance deletes the addon
func (s *Client) DeleteAddonIstioInstance(clusterUUID string) error {
	Debug(1, "Entered DeleteAddonIstioInstance for UUID "+clusterUUID)

	if clusterUUID == "" {
		return errors.New("Cluster UUID to delete is required")
	}

	addonInstalled, err := s.IsAddonInstalled(clusterUUID, "ccp-istio-operator")

	if err != nil {
		return err
	}

	if *addonInstalled {

		url := s.BaseURL + "/v3/clusters/" + clusterUUID + "/addons/ccp-istio-cr/"
		Debug(2, "Sending HTTP delte to "+url)

		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			return err
		}
		_, err = s.doRequest(req)
		if err != nil {
			return err
		}
	}

	Debug(2, "Request sent to API with success response")
	return nil
}

// DeleteAddonIstioOp deletes the addon
func (s *Client) DeleteAddonIstioOp(clusterUUID string) error {
	Debug(1, "Entered DeleteAddonIstioOp for UUID "+clusterUUID)

	if clusterUUID == "" {
		return errors.New("Cluster UUID to delete is required")
	}

	url := s.BaseURL + "/v3/clusters/" + clusterUUID + "/addons/ccp-istio-operator/"
	Debug(2, "Sending HTTP delte to "+url)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	_, err = s.doRequest(req)
	if err != nil {
		return err
	}

	Debug(2, "Request sent to API with success response")
	return nil
}

// DeleteAddonDashboard deletes the addon
func (s *Client) DeleteAddonDashboard(clusterUUID string) error {
	Debug(1, "Entered DeleteAddonDashboard for UUID "+clusterUUID)

	if clusterUUID == "" {
		return errors.New("Cluster UUID to delete is required")
	}

	url := s.BaseURL + "/v3/clusters/" + clusterUUID + "/addons/kubernetes-dashboard/"
	Debug(2, "Sending HTTP delte to "+url)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	_, err = s.doRequest(req)
	if err != nil {
		return err
	}

	Debug(2, "Request sent to API with success response")
	return nil
}

// DeleteAddonIstio install both
func (s *Client) DeleteAddonIstio(clusterUUID string) error {
	err := s.DeleteAddonIstioInstance(clusterUUID)
	if err != nil {
		Debug(1, "Failed to delete Add-On Istio Instance: "+string(err.Error()))
		return err
	}
	time.Sleep(2 * time.Second) // wait 2 seconds before sending the next request
	err = s.DeleteAddonIstioOp(clusterUUID)
	if err != nil {
		Debug(1, "Failed to delete Add-On Istio Operator: "+string(err.Error()))
		return err
	}
	return nil
}

// DeleteAddonHarborInstance deletes the addon
func (s *Client) DeleteAddonHarborInstance(clusterUUID string) error {
	Debug(1, "Entered DeleteAddonHarborInstance for UUID "+clusterUUID)

	if clusterUUID == "" {
		return errors.New("Cluster UUID to delete is required")
	}

	addonInstalled, err := s.IsAddonInstalled(clusterUUID, "ccp-harbor-operator")

	if err != nil {
		return err
	}

	if *addonInstalled {

		url := s.BaseURL + "/v3/clusters/" + clusterUUID + "/addons/ccp-harbor-cr/"
		Debug(2, "Sending HTTP delte to "+url)

		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			return err
		}
		_, err = s.doRequest(req)
		if err != nil {
			return err
		}

	}

	Debug(2, "Request sent to API with success response")
	return nil
}

// DeleteAddonHarborOp deletes the addon
func (s *Client) DeleteAddonHarborOp(clusterUUID string) error {
	Debug(1, "Entered DeleteAddonHarborOp for UUID "+clusterUUID)

	if clusterUUID == "" {
		return errors.New("Cluster UUID to delete is required")
	}

	url := s.BaseURL + "/v3/clusters/" + clusterUUID + "/addons/ccp-harbor-operator/"
	Debug(2, "Sending HTTP delte to "+url)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	_, err = s.doRequest(req)
	if err != nil {
		return err
	}

	Debug(2, "Request sent to API with success response")
	return nil
}

// DeleteAddonHarbor delete both
func (s *Client) DeleteAddonHarbor(clusterUUID string) error {
	err := s.DeleteAddonHarborInstance(clusterUUID)
	if err != nil {
		Debug(1, "Failed to delete Add-On Harbor Instance: "+string(err.Error()))
		return err
	}
	time.Sleep(2 * time.Second) // wait 2 seconds before sending the next request
	err = s.DeleteAddonHarborOp(clusterUUID)
	if err != nil {
		Debug(1, "Failed to delete Add-On Harbor Operator: "+string(err.Error()))
		return err
	}
	return nil
}

// GetAddonsCatalogue returns a list of Addons
func (s *Client) GetAddonsCatalogue(clusterUUID string) (*AddonsCatalogue, error) {
	// https://mholt.github.io/json-to-go/
	Debug(3, "GetAddonsCatalogue for cluster "+clusterUUID)

	url := s.BaseURL + "/v3/clusters/" + clusterUUID + "/catalog"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	bytes, err := s.doRequest(req)
	if err != nil {
		return nil, err
	}
	Debug(3, string(bytes))
	var data *AddonsCatalogue

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetClusterInstalledAddons returns a list of Addons
func (s *Client) GetClusterInstalledAddons(clusterUUID string) (*ClusterInstalledAddons, error) {
	Debug(3, "GetClusterInstalledAddons for cluster "+clusterUUID)

	url := s.BaseURL + "/v3/clusters/" + clusterUUID + "/addons/"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	bytes, err := s.doRequest(req)
	if err != nil {
		return nil, err
	}
	Debug(3, string(bytes))
	var data *ClusterInstalledAddons

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// IsAddonInstalled check if addon is installed (bool)
func (s *Client) IsAddonInstalled(clusterUUID string, addonName string) (*bool, error) {

	clusterAddons, err := s.GetClusterInstalledAddons(clusterUUID)

	if err != nil {
		return nil, err
	}

	for _, addon := range clusterAddons.Results {
		if addonName == addon.Name {
			if addon.AddonStatus.Status == "INSTALLED" {
				return Bool(true), nil
			} else {
				return Bool(false), nil
			}
		}
	}

	return Bool(false), nil

}

// InstallAddonHXCSI Installs the Istio Operator
func (s *Client) InstallAddonHXCSI(clusterUUID string) error {
	Debug(1, "Entered InstallAddonHXCSI for UUID "+clusterUUID)

	if clusterUUID == "" {
		return errors.New("Cluster UUID is required")
	}

	url := s.BaseURL + "/v3/clusters/" + clusterUUID + "/addons/"

	Debug(2, "Getting Add-Ons catalog for UUID "+clusterUUID)
	addons, err := s.GetAddonsCatalogue(clusterUUID)
	if err != nil {
		Debug(2, err.Error())
		return err
	}

	// fix missing Namespace from CcpHxcsi struct
	addons.CcpHxcsi.Namespace = "ccp"
	// now prepare the JSON body
	jsonBody, err := json.Marshal(addons.CcpHxcsi)
	if err != nil {
		return err
	}

	Debug(2, "POSTing the JSON body:")
	Debug(2, string(jsonBody))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	resp, err := s.doRequest(req)
	if err != nil {
		return err
	}
	Debug(3, "Response is:")
	Debug(3, string(resp))

	Debug(2, "Request sent to API with success response")
	return nil
}

// DeleteAddonHXCSI deletes the addon
func (s *Client) DeleteAddonHXCSI(clusterUUID string) error {
	Debug(1, "Entered DeleteAddonHXCSI for UUID "+clusterUUID)

	if clusterUUID == "" {
		return errors.New("Cluster UUID to delete is required")
	}

	url := s.BaseURL + "/v3/clusters/" + clusterUUID + "/addons/ccp-hxcsi/"
	Debug(2, "Sending HTTP delete to "+url)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	resp, err := s.doRequest(req)
	if err != nil {
		return err
	}
	Debug(2, "Response:")
	Debug(2, string(resp))

	Debug(2, "Request sent to API with success response")
	return nil
}

// InstallAddonKubeflow Installs the Kubeflow addons
func (s *Client) InstallAddonKubeflow(clusterUUID string) error {
	Debug(1, "Entered InstallAddonKubeflow for UUID "+clusterUUID)

	if clusterUUID == "" {
		return errors.New("Cluster UUID is required")
	}

	url := s.BaseURL + "/v3/clusters/" + clusterUUID + "/addons/"

	Debug(2, "Getting Add-Ons catalog for UUID "+clusterUUID)
	addons, err := s.GetAddonsCatalogue(clusterUUID)
	if err != nil {
		Debug(2, err.Error())
		return err
	}
	jsonBody, err := json.Marshal(addons.CcpKubeflow)
	if err != nil {
		return err
	}

	Debug(2, string(jsonBody))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	resp, err := s.doRequest(req)
	if err != nil {
		return err
	}
	Debug(3, "Response is:")
	Debug(3, string(resp))

	Debug(2, "Request sent to API with success response")
	return nil
}

// GetKubeflowAddonConfig for kubeflow
func (s *Client) GetKubeflowAddonConfig(clusterUUID string) ([]byte, error) {
	Debug(1, "Entered InstallAddonKubeflow for UUID "+clusterUUID)

	if clusterUUID == "" {
		return nil, errors.New("Cluster UUID is required")
	}

	Debug(2, "Getting Add-Ons catalog for UUID "+clusterUUID)
	addons, err := s.GetAddonsCatalogue(clusterUUID)
	if err != nil {
		Debug(2, err.Error())
		return nil, err
	}
	jsonBody, err := json.Marshal(addons.CcpKubeflow)
	if err != nil {
		return nil, err
	}

	return jsonBody, nil
}

// DeleteAddonKubeflow deletes the addon
func (s *Client) DeleteAddonKubeflow(clusterUUID string) error {
	Debug(1, "Entered DeleteAddonKubeflow for UUID "+clusterUUID)

	if clusterUUID == "" {
		return errors.New("Cluster UUID to delete is required")
	}

	url := s.BaseURL + "/v3/clusters/" + clusterUUID + "/addons/ccp-kubeflow/"
	Debug(2, "Sending HTTP delete to "+url)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	_, err = s.doRequest(req)
	if err != nil {
		return err
	}

	Debug(2, "Request sent to API with success response")
	return nil
}

// PatchCluster does the things
func (s *Client) PatchCluster(cluster *Cluster, clusterUUID string) (*Cluster, error) {

	var data Cluster

	//clusterUUID := *cluster.UUID

	url := fmt.Sprintf(s.BaseURL + "/v3/clusters/" + clusterUUID + "/")

	j, err := json.Marshal(cluster)

	if err != nil {

		return nil, err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}

	bytes, err := s.doRequest(req)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &data)

	if err != nil {
		return nil, err
	}

	cluster = &data

	return cluster, nil
}

// DeleteAddon delete an Addon
func (s *Client) DeleteAddon(clusterUUID string, addonName string) error {

	if clusterUUID == "" {
		return errors.New("Cluster UUID is required")
	}

	switch addonName {

	case "kubernetes-dashboard", "ccp-efk", "ccp-monitor", "ccp-kubeflow":
		err = s.DeleteAddonAndConfirm(clusterUUID, addonName)

		if err != nil {
			return err
		}

	case "istio":

		err = s.DeleteAddonAndConfirm(clusterUUID, "ccp-istio-operator")

		if err != nil {
			return err
		}

		err = s.DeleteAddonIstioInstance(clusterUUID)

		if err != nil {
			return err
		}

	case "harbor":

		err = s.DeleteAddonAndConfirm(clusterUUID, "ccp-harbor-operator")

		if err != nil {
			return err
		}

		err = s.DeleteAddonHarborInstance(clusterUUID)

		if err != nil {
			return err
		}

	default:
		return errors.New("Incorrect addon name provided. Options are 'kubernetes-dashboard', 'ccp-efk', 'ccp-monitor', 'cp-istio','harbor', 'ccp-kubeflow'")

	}

	return nil
}

// DeleteAddonAndConfirm delete and confirm
func (s *Client) DeleteAddonAndConfirm(clusterUUID string, addonName string) error {

	url := s.BaseURL + "/v3/clusters/" + clusterUUID + "/addons/" + addonName + "/"

	addonInstalled, err := s.IsAddonInstalled(clusterUUID, addonName)

	if err != nil {
		return err
	}

	if *addonInstalled {

		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			return err
		}
		_, err = s.doRequest(req)
		if err != nil {
			return err
		}

	}

	return nil
}
