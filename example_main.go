package main

import (
	"fmt"

	// fork this github repo in to your ~/git/src dir
	// go get github.com/rob-moss/ccp-clientlibrary-go
	"github.com/rob-moss/ccp-clientlibrary-go/ccp"
)

var cpUser = "admin"    // user for CCP Control Plane
var cpPass = "password" // Password for CCP Control Plane
var cpURL = "https://10.100.60.40"

// var cpUser = os.GetEnv("CCPPUSER") // user for CCP Control Plane
// var cpPass = os.GetEnv("CCPPPASS") // Password for CCP Control Plane
// var cpURL = os.GetEnv("CCPURL")

// var jar, err = cookiejar.New(nil)

func main() {

	fmt.Println("* Entered main")

	client := ccp.NewClient(cpUser, cpPass, cpURL)

	err := client.Login(client)
	if err != nil {
		fmt.Println(err)
	}

	// clusters, err := client.GetClusters()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("* clusters = " + strconv.Itoa(len(clusters)))

	// ----
	// fmt.Println("* Get first cluster name")
	// fmt.Println(string(*clusters[0].Name))
	// clustername := string(*clusters[0].Name)

	// GetClusterByName gets all clusters, searches for matching cluster name, returns *Cluster struct
	// clustername := string("romoss-testcp01-tenant01")
	// cluster, err := client.GetClusterByName(clustername)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	//
	// if cluster == nil {
	// 	fmt.Println(err)
	// 	return
	// } else {
	// 	fmt.Printf("* Got cluster UUID %s\n", *cluster.UUID)
	// }

	// // GetClusterByUUID gets cluster by UUID, returns *Cluster struct
	// cluster, err = client.GetClusterByUUID(*cluster.UUID)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	//
	// if cluster == nil {
	// 	fmt.Println(err)
	// 	return
	// } else {
	// 	fmt.Printf("* Got cluster %s by UUID\n", *cluster.Name)
	// }

	// ---- GetInfraProviders
	// providerClientConfigs, err := client.GetInfraProviders()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// // Print out the providerClientConfig details
	// fmt.Println("* Provider Config name: " + *providerClientConfigs[0].Name + " hostname: " + *providerClientConfigs[0].Address + " UUID: " + *providerClientConfigs[0].UUID)

	// ---- GetInfraProviderByName
	// infraProvider, err := client.GetInfraProviderByName("vsphere")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// // Print out the providerClientConfig details
	// fmt.Println("* Provider Config name: " + *infraProvider.Name + " hostname: " + *infraProvider.Address + " UUID: " + *infraProvider.UUID)

	// // Get network provider and Subnet
	// networkProviderSubnets, err := client.GetNetworkProviderSubnets()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// // Print out the providerClientConfig details
	// fmt.Println("* Network Provider  name: " + *networkProviderSubnets[0].Name + " CIDR: " + *networkProviderSubnets[0].CIDR + " UUID: " + *networkProviderSubnets[0].UUID)

	// --- GetNetworkProviderSubnetByName
	// Get network provider by name and return single entry
	// networkProviderSubnet, err := client.GetNetworkProviderSubnetByName("default-network-subnet")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// // Print out the providerClientConfig details
	// fmt.Println("* Network Provider  name: " + *networkProviderSubnet.Name + " CIDR: " + *networkProviderSubnet.CIDR + " UUID: " + *networkProviderSubnet.UUID)

	// --- scale a cluster
	// clusterUUID := string(*cluster.UUID)
	// clusterWorkerPoolName := string("node-group")
	// clusterSize := int64(3)
	// scaleCluster, err := client.ScaleCluster(clusterUUID, clusterWorkerPoolName, clusterSize)
	// if err != nil {
	// 	fmt.Println(err)
	// 	// Print out the Println of bytes
	// 	// to debug: uncomment below. Prints JSON payload
	// 	//fmt.Println(string(scaleCluster))
	// } else {
	// 	fmt.Println("Name: " + *scaleCluster.Name + " Size: " + string(*scaleCluster.Size))
	// }

	// --- AddCluster from JSON File
	// jsonFile := "./cluster.json"
	// // Convert a JSON file to a Cluster struct
	// newCluster, err := client.ConvertJSONToCluster(jsonFile)
	// if err != nil {
	// 	fmt.Println("error:", err)
	// } else {
	// 	fmt.Println("Success")
	// }

	// fmt.Println("* New cluster name to create: " + *newCluster.Name)
	// createdCluster, err := client.AddCluster(newCluster)
	// if err != nil {
	// 	fmt.Println("Error from AddCluster:")
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("* Cluster sent to API: " + *createdCluster.Name)

	// --- Delete Cluster
	// clustername := string("romoss-testcp01-tenant04")
	// cluster, err := client.GetClusterByName(clustername)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// err = client.DeleteCluster(*cluster.UUID)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// --- create the Cluster struct, for later sending to the AddCluster function
	// // https://stackoverflow.com/questions/51916592/fill-a-struct-which-contains-slices
	// ccpsshuser := ccp.String("ccpadmin")
	// ccpsshkey := ccp.String("ecdsa-sha2-nistp521 AAAAE2VjZHNhLXNoYTItbmlzdHA1MjEAAAAIbmlzdHA1MjEAAACFBAFXQk0bZlFiFV6FD5DT0HdVJ2TsL9wlciD3UkcFI+/kpIj2AfOqqoQjt0BYZKzNJ6z4a25nkIueQJFog04S0/+PkQGX/Hc2DVccatAOWMRCedwukdgfoURLHyEdgl9EeCmiyqnUe6XVxiqcX9dkqXuI1KsP/oRir8ZAui3nXvdyUm8TGA== ccpadmin@galaxy.cisco.com")
	// ccptemplateimg := ccp.String("hx1-ccp-tenant-image-1.16.3-ubuntu18-6.1.1-pre")
	// kubernetesversion := ccp.String("1.16.3")
	// newCluster := &ccp.Cluster{
	// 	Name: ccp.String("romoss-testcp01-tenant03"),
	// 	Type: ccp.String("vsphere"),
	// 	// WorkerNodePool: newWorkers,
	// 	WorkerNodePool: &[]ccp.WorkerNodePool{
	// 		// first worker node pool
	// 		ccp.WorkerNodePool{
	// 			Name:              ccp.String("node-pool"), // default name
	// 			Size:              ccp.Int64(1),
	// 			VCPUs:             ccp.Int64(8),
	// 			Memory:            ccp.Int64(32768),
	// 			Template:          ccptemplateimg,
	// 			SSHUser:           ccpsshuser,
	// 			SSHKey:            ccpsshkey,
	// 			KubernetesVersion: kubernetesversion,
	// 		},
	// 	},
	// 	MasterNodePool: &ccp.MasterNodePool{
	// 		Name:              ccp.String("master-group"),
	// 		Size:              ccp.Int64(1),
	// 		VCPUs:             ccp.Int64(2),
	// 		Memory:            ccp.Int64(16384),
	// 		Template:          ccptemplateimg,
	// 		SSHUser:           ccpsshuser,
	// 		SSHKey:            ccpsshkey,
	// 		KubernetesVersion: kubernetesversion,
	// 	},
	// 	Infra: &ccp.Infra{
	// 		Datastore:  ccp.String("GFFA-HX1-CCPInstallTest01"),
	// 		Datacenter: ccp.String("GFFA-DC"),
	// 		Networks:   &[]string{"DV_VLAN1060"},
	// 		Cluster:    ccp.String("GFFA-HX1-Cluster"),
	// 	},
	// 	KubernetesVersion:  kubernetesversion,
	// 	InfraProviderUUID:  infraProvider.UUID,
	// 	SubnetUUID:         networkProviderSubnet.UUID,
	// 	LoadBalancerIPNum:  ccp.Int64(2),
	// 	IPAllocationMethod: ccp.String("ccpnet"),
	// 	AWSIamEnabled:      ccp.Bool(false),
	// 	NetworkPlugin: &ccp.NetworkPlugin{
	// 		Name: ccp.String("calico"),
	// 		Details: &ccp.NetworkPluginDetails{
	// 			PodCIDR: ccp.String("192.168.0.0/16"),
	// 		},
	// 	},
	// }
	// fmt.Println(newCluster)

	// // now create the cluster
	// // fmt.Println("* New cluster name to create: " + *newCluster.Name)
	// createdCluster, err := client.AddCluster(newCluster)
	// if err != nil {
	// 	fmt.Println("Error from AddCluster:")
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("* Cluster sent to API: " + *createdCluster.Name)

	// // ---- GetInfraProviderByName
	// infraProvider, err := client.GetInfraProviderByName("vsphere")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// // Print out the providerClientConfig details
	// fmt.Println("* Provider Config name: " + *infraProvider.Name + " hostname: " + *infraProvider.Address + " UUID: " + *infraProvider.UUID)

	cluster, err := client.GetClusterByName("romoss-testcp01-tenant02")
	if err != nil {
		fmt.Println(err)
		return
	}
	// // GetAddons uses UUID
	// addons, err := client.GetAddonsCatalogue(*cluster.UUID)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	//
	// fmt.Println(addons)

	// err = client.InstallAddonIstioOp(*cluster.UUID)
	// if err != nil {
	// 	ccp.Debug(1, "Error:")
	// 	fmt.Println(err)
	// 	return
	// }
	// time.Sleep(2 * time.Seconds)

	// err = client.InstallAddonIstioInstance(*cluster.UUID)
	// if err != nil {
	// 	ccp.Debug(1, "Error:")
	// 	fmt.Println(err)
	// 	return
	// }
	// time.Sleep(2 * time.Seconds)

	// err = client.InstallAddonIstio(*cluster.UUID)
	// if err != nil {
	// 	ccp.Debug(1, "Error:")
	// 	fmt.Println(err)
	// 	return
	// }
	// time.Sleep(2 * time.Seconds)

	// err = client.InstallAddonMonitoring(*cluster.UUID)
	// if err != nil {
	// 	ccp.Debug(1, "Error:")
	// 	fmt.Println(err)
	// 	return
	// }
	// time.Sleep(2 * time.Seconds)

	// err = client.InstallAddonLogging(*cluster.UUID)
	// if err != nil {
	// 	ccp.Debug(1, "Error:")
	// 	fmt.Println(err)
	// 	return
	// }
	// time.Sleep(2 * time.Second)

	// err = client.InstallAddonHarborOp(*cluster.UUID)
	// if err != nil {
	// 	ccp.Debug(1, "Error:")
	// 	fmt.Println(err)
	// 	return
	// }
	// time.Sleep(2 * time.Seconds)

	// err = client.InstallAddonHarborInstance(*cluster.UUID)
	// if err != nil {
	// 	ccp.Debug(1, "Error:")
	// 	fmt.Println(err)
	// 	return
	// }
	// time.Sleep(2 * time.Seconds)

	// err = client.InstallAddonHarbor(*cluster.UUID)
	// if err != nil {
	// 	ccp.Debug(1, "Error:")
	// 	fmt.Println(err)
	// 	return
	// }
	// time.Sleep(2 * time.Seconds)

	// ---- delete addons
	// err = client.DeleteAddOnLogging(*cluster.UUID)
	// if err != nil {
	// 	ccp.Debug(1, "Error:")
	// 	fmt.Println(err)
	// 	return
	// }

	err = client.DeleteAddOnMonitor(*cluster.UUID)
	if err != nil {
		ccp.Debug(1, "Error:")
		fmt.Println(err)
		return
	}

	// err = client.DeleteAddOnIstio(*cluster.UUID)
	// if err != nil {
	// 	ccp.Debug(1, "Error:")
	// 	fmt.Println(err)
	// 	return
	// }

	// err = client.DeleteAddOnHarbor(*cluster.UUID)
	// if err != nil {
	// 	ccp.Debug(1, "Error:")
	// 	fmt.Println(err)
	// 	return
	// }

	fmt.Printf("* Closed\n")

	// todo:
	// GetAddons from /catalogue/
	// InstallAddon<AddonName>
	// - Istio
	// - HX-CSI
	// - Logging
	// - Monitoring
	// - Dashboard
	// - Kubeflow
}

/* toDo
- Create JSON config
- Make connection to CCP CP via Proxy (optional)
- Set defaults: image, sshkey, sshuser, provider, network
- Log in to CCP using X-Auth-Token
- Create functions to:
-- Get kubernetes version for deployments
-- Fetch provider by name -> uuid
-- Fetch subnet by name -> uuid
-- Create Cluster (Calico, vSphere)
-- Scale Cluster (Worker nodes)
-- Delete Cluster

v2 todo
- Create functions to:
-- Install Add-Ons
--- Istio
--- Harbor
--- HX-CSI
--- Monitoring
--- Logging
*/
