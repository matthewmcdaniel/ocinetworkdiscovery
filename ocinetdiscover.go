package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/common/auth"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/oracle/oci-go-sdk/v65/example/helpers"
)

type Vnic struct {
	MacAddr         string `json:"macAddr"`
	PrivateIP       string `json:"privateIp"`
	SubnetCidrBlock string `json:"subnetCidrBlock"`
	VirtualRouterIP string `json:"virtualRouterIp"`
	VlanTag         int    `json:"vlanTag"`
	VnicId          string `json:"vnicId"`
}

const (
	instance_metadata_url = "http://169.254.169.254/opc/v1/vnics/"
)

var provider common.ConfigurationProvider
var err error

func main() {

	for _, val := range os.Args {
		if strings.EqualFold(val, "--help") || strings.EqualFold(val, "-h") {
			fmt.Println("Lists networking information of node.")
			fmt.Println("Usage: ./ocinetdiscovery <instance_principal | workload_identity> <args>")
			fmt.Println("Arguments:\n--publicipv4 - Lists public IPv4 address of primary VNIC.")
			fmt.Println("--privateipv4 - Lists private IPv4 address of primary VNIC.")
			return
		}
	}

	if strings.EqualFold(os.Args[1], "instance_principal") { // Dynamic group required
		provider, err = auth.InstancePrincipalConfigurationProvider()
		_ = provider // Avoid compilation error
		helpers.FatalIfError(err)
	} else if strings.EqualFold(os.Args[1], "workload_identity") { // Workload identity required
		provider, err = auth.OkeWorkloadIdentityConfigurationProvider()
		_ = provider // Avoid compilation error
		helpers.FatalIfError(err)
	} else {
		fmt.Println("No valid authentication method, please include authentication method i.e instance_principal or workload_identity")
		os.Exit(1)
	}

	// Request for instance metadata
	instance_metadata_response, err := http.Get(instance_metadata_url)
	helpers.FatalIfError(err)

	defer instance_metadata_response.Body.Close()

	// Deserialize response JSON
	var vn []Vnic
	err = json.NewDecoder(instance_metadata_response.Body).Decode(&vn)
	helpers.FatalIfError(err)

	// Initialize network client
	virtual_network_client, err := core.NewVirtualNetworkClientWithConfigurationProvider(provider)
	helpers.FatalIfError(err)

	// Send request for VNIC using VnicId acquired from instance metadata
	vnic_request := core.GetVnicRequest{
		VnicId: common.String(vn[0].VnicId),
	}

	// Send request for VNIC using VnicId acquired from instance metadata
	vnic_response, err := virtual_network_client.GetVnic(context.Background(), vnic_request)
	helpers.FatalIfError(err)

	for _, val := range os.Args {
		if strings.EqualFold(val, "--privateIPv4") {
			_, err := fmt.Println(*vnic_response.PrivateIp)
			helpers.FatalIfError(err)
		}
		if strings.EqualFold(val, "--publicIPv4") {
			_, err := fmt.Println(*vnic_response.PublicIp)
			helpers.FatalIfError(err)
		}
	}
	return
}
