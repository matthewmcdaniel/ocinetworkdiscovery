# ocinetdiscover

This CLI tool can be used to introspect an OCI cloud VMs private and public IPv4 address. This tool relies on [instance principals](https://docs.oracle.com/en-us/iaas/Content/Identity/Tasks/callingservicesfrominstances.htm) for Compute-based workloads, and [workload identity](https://docs.oracle.com/en-us/iaas/Content/ContEng/Tasks/contenggrantingworkloadaccesstoresources.htm) for OKE-based workloads (Oracle Kubernetes Engine). Using manually configured credentials is not supported (or recommended).

Supported services:
* OCI Compute 
* OKE w/ managed nodes
* OKE w/ self-managed nodes

Not supported/Not tested:
* OCI Container Instances
* Functions
* OKE w/ virtual nodes

## Prerequisites

If you are using this on OCI Compute, you will need a dynamic group with the following rule

```
instance.compartment.id='compartment_ocid'
```

Assign a policy that will allow the dynamic group to access the instance's VNIC.

```
Allow dynamic-group <dyn-group> to read vnics in compartment <compartment-name>
```

See [Dynamic Groups](https://docs.oracle.com/en-us/iaas/Content/Identity/Tasks/managingdynamicgroups.htm) for more information.

If you are using this on OKE, I recommend you use [Workload Identities](https://docs.oracle.com/en-us/iaas/Content/ContEng/Tasks/contenggrantingworkloadaccesstoresources.htm). The method described above will also work with OKE, but will not grant you pod-level isolation. This is where workload identity can help.

## Usage

```
go build
```

```
./ocinetdiscover <instance_principal | workload_identity> <arg>
```

If you are using instance principals, use `instance_principal`.

If you are using workload identity, use `workload_identity`.
## Supported Arguments

```
--privateipv4 
```
Lists primary private IPv4 address of primary VNIC.
```
--publicipv4
```
Lists primary public IPv4 address of primary VNIC.  

