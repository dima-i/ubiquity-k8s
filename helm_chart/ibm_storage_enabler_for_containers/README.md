# Ubiquity
* [(PRODUCTNAME)](https://<PRODUCTURL>) is ... brief sentence regarding product
* Add "-Beta" as suffix if beta version - beta versions are generally < 1.0.0
* Don't include versions of charts or products

## Introduction
IBM Storage Enabler for Containers allows IBM storage systems to be used as persistent volumes for stateful applications running in Kubernetes clusters.
Thus, the containers can be used with stateful microservices, such as database applications (MongoDB, PostgreSQL etc).
IBM Storage Enabler for Containers uses Kubernetes dynamic provisioning for creating and deleting volumes on IBM storage systems.
In addition, IBM Storage Enabler for Containers utilizes the full set of Kubernetes FlexVolume APIs for volume operations on a host.
The operations include initiation, attachment/detachment, mounting/unmounting etc.

## Chart Details
This chart includes:
* A Ubiquity server Deployment used as the server of Kubernetes Dynamic Provisioner and FlexVolume.
* A Ubiquity database Deployment used to store the persistent data of Ubiquity server.
* A Kubernetes Dynamic Provisioner Deployment for creation  storage volumes on-demand, using Kubernetes storage classes based on Spectrum Connect storage services.
* A Kubernetes FlexVolume DaemonSet for attaching and mounting storage volumes into a pod within a Kubernetes node.

## Prerequisites
Follow these steps to prepare your environment for installing the IBM Storage Enabler for Containers in the Kubernetes cluster that requires persistent volumes for stateful containers.
1. Contact your storage administrator and make sure IBM Storage Enabler for Containers interface has been added to active Spectrum Connect instance and at least one storage service has been delegated to it.
2. Verify that there is a proper communication link between Spectrum Connect and Kubernetes cluster.
3. Perform these steps for each worker node in Kubernetes cluster:
  * Install the following Linux packages to ensure Fibre Channel and iSCSI connectivity. Skip this step, if the packages are already installed.
    * RHEL: `sg3_utils` and `iscsi-initiator-utils` (if iSCSI connection is required).
    ```bash
    sudo yum -y install sg3_utils
    sudo yum -y install iscsi-initiator-utils
    ```
    * Ubuntu: `scsitools` and `open-iscsi` (if iSCSI connection is required).
    ```bash
    sudo apt-get install scsitools
    sudo apt-get install open-iscsi
    ```
 * Configure Linux multipath devices on the host. Create and set the relevant storage system parameters in the `/etc/multipath.conf file`. You can also use the default `multipath.conf` file located in the `/usr/share/doc/device-mapper-multipath-*` directory. 
 Verify that the `systemctl` status `multipathd` output indicates that the multipath status is active and error-free.
   * RHEL: 
    ```bash
    yum install device-mapper-multipath
    sudo modprobe dm-multipath
    systemctl start multipathd
    systemctl status multipathd
    multipath -ll
    ```
   * Ubuntu: 
    ```bash
    apt-get install multipath-tools
    sudo modprobe dm-multipath
    systemctl start multipathd
    systemctl status multipathd
    multipath -ll
    ```  
   Important: When configuring Linux multipath devices, verify that the `find_multipaths` parameter in the `multipath.conf` file is disabled.
   * RHEL: Remove the `find_multipaths yes` string from the `multipath.conf` file.
   * Ubuntu: Add the `find_multipaths no` string to the `multipath.conf` file, see below:
    ```bash
    defaults {
   find_multipaths no
   }
    ```
* Configure storage system connectivity.
  * Define the hostname of each Kubernetes node on the relevant storage systems with the valid WWPN or IQN of the node. The hostname on the storage system must be the same as the hostname defined in the Kubernetes cluster. Use the `$> kubectl get nodes` command to display hostname, as illustrated below. In this example, the `k8s-worker-node1` and the `k8s-worker-node2` hostnames must be defined on a storage system.
  
  Note: In most cases, the local hostname of the node is the same as the Kubernetes node hostname as displayed in the `kubectl get nodes` command      output. However, if the names are different, make sure to use the Kubernetes node name, as it appears in the command output.
```bash
root@k8s-user-v18-master:~# kubectl get nodes
NAME               STATUS   ROLES      AGE       VERSION
k8s-master         Ready     master    34d       v1.8.4
k8s-worker-node1   Ready     <none>    34d       v1.8.4
k8s-worker-node2   Ready     <none>    34d       v1.8.4
```
 * After the node hostnames are defined, log into Spectrum Connect UI and refresh the relevant storage systems in the Storage System pane.
 * For iSCSI, perform these three steps.
   * Make sure that the login used to log in to the iSCSI targets is permanent and remains available after a reboot of the worker node. To do this, verify that the `node.startup` in the `/etc/iscsi/iscsid.conf` file is set to `automatic`. If not, set it as required and then restart the `iscsid` service (`$> service iscsid restart`).
   * Discover and log into at least two iSCSI targets on the relevant storage systems.
   ```bash
    $> iscsiadm -m discoverydb -t st -p ${storage system iSCSI port IP}:3260
    --discover
    $> iscsiadm -m node  -p ${storage system iSCSI port IP/hostname} --login
    ```
    * Verify that the login was successful and display all targets that you logged in. The `portal` value must be the iSCSI target IP address.
    ```bash
    $> iscsiadm -m session --rescan
    Rescanning session [sid: 1, target: {storage system IQN},
    portal: {storage system iSCSI port IP},{port number}
    ```
## Resources Required
* Describes Minimum System Resources Required

## Installing the Chart

To install the chart with the release name `my-release`:

```bash
$ helm install --name my-release --namespace ubiquity stable/ibm_storage_enabler_for_containers
```

The command deploys <Chart name> on the Kubernetes cluster in the default configuration. The [configuration](#configuration) section lists the parameters that can be configured during installation.


> **Tip**: List all releases using `helm list`

### Verifying the Chart
You can check the status by running:
```bash
$ helm status my-release
```

If all statuses are free of errors, you can run sanity test by:
```bash
$ helm test my-release
```

### Uninstalling the Chart

To uninstall/delete the `my-release` release:

```bash
$ helm delete my-release --purge
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the <Ubiquity> chart and their default values.

| Parameter                  | Description                                     | Default                                                    |
| -----------------------    | ---------------------------------------------   | ---------------------------------------------------------- |
| `images.ubiquity`                                                                               | Image for Ubiquity server | `ibmcom/ibm-storage-enabler-for-containers:2.0.0` |
| `images.ubiquitydb`                                                                             | Image for Ubiquity database | `ibmcom/ibm-storage-enabler-for-containers-db:2.0.0` |
| `images.provisioner`                                                                            | Image for Kubernetes Dynamic Provisioner | `ibmcom/ibm-storage-dynamic-provisioner-for-kubernetes:2.0.0` |
| `images.flex`                                                                                   | Image for Kubernetes FlexVolume | `ibmcom/ibm-storage-flex-volume-for-kubernetes:2.0.0` |
| `spectrumConnect.connectionInfo.fqdn`                                                           | IP\FQDN of Spectrum Connect server. | ` ` |
| `spectrumConnect.connectionInfo.port`                                                           | Port of Spectrum Connect server. | ` ` |
| `spectrumConnect.connectionInfo.username`                                                       | Username defined for IBM Storage Enabler for Containers interface in Spectrum Connect. | ` ` |
| `spectrumConnect.connectionInfo.password`                                                       | Password defined for IBM Storage Enabler for Containers interface in Spectrum Connect. | ` ` |
| `spectrumConnect.connectionInfo.sslMode`                                                        | SSL verification mode. Allowed values: require (no validation is required) and verify-full (user-provided certificates) | `require` |
| `spectrumConnect.backendConfig.instanceName`                                                    | A prefix for any new volume created on the storage system | ` ` |
| `spectrumConnect.backendConfig.skipRescanIscsi`                                                 | Allowed values: true or false. Set to true if the nodes have FC connectivity | `false` |
| `spectrumConnect.backendConfig.DefaultStorageService`                                           | Default Spectrum Connect storage service to be used, if not specified by the storage class | ` ` |
| `spectrumConnect.backendConfig.newVolumeDefaults.fsType`                                        | File system type. Allowed values: ext4 or xfs | `ext4` |
| `spectrumConnect.backendConfig.newVolumeDefaults.size`                                          | The default volume size (in GB) if not specified by the user when creating a new volume | `1` |
| `spectrumConnect.backendConfig.dbPvConfig.ubiquityDbPvName`                                     | Ubiquity database PV name. For Spectrum Virtualize and Spectrum Accelerate, use default value "ibm-ubiquity-db". For DS8000 Family, use "ibmdb" instead and make sure UBIQUITY_INSTANCE_NAME_VALUE value length does not exceed 8 chars | `ibm-ubiquity-db`                                                        |
| `spectrumConnect.backendConfig.dbPvConfig.storageClassForDbPv.storageClassName`                 | Parameters to create the first Storage Class that also be used by ubiquity for ibm-ubiquity-db PVC | ` `|
| `spectrumConnect.backendConfig.dbPvConfig.storageClassForDbPv.params.spectrumConnectServiceName`| Storage Class profile parameter should point to the Spectrum Connect storage service name | ` ` |
| `spectrumConnect.backendConfig.dbPvConfig.storageClassForDbPv.params.fsType`                    | Storage Class file-system type, Allowed values: ext4 or xfs | `ext4` |
| `genericConfig.ubiquityIpAddress`                                                               | The IP address of the ubiquity service object | ` ` |
| `genericConfig.logging.logLevel`                                                                | Log level. Allowed values: debug, info, error | `info` |
| `genericConfig.logging.flexLogDir`                                                              | Flex log directory. If you change the default, then make the new path exist on all the nodes and update the Flex daemonset hostpath according | `/var/log` |
| `genericConfig.ubiquityDbCredentials.username`                                                  | Username for the deployment of ubiquity-db database. Note : Do not use the "postgres" username, because it already exists | ` ` |
| `genericConfig.ubiquityDbCredentials.password`                                                  | Password for the deployment of ubiquity-db database | ` ` |

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`.

Alternatively, a YAML file that specifies the values for the parameters can be provided while installing the chart.

## Storage
* Define how storage works with the workload
* Dynamic vs PV pre-created
* Considerations if using hostpath, local volume, empty dir
* Loss of data considerations
* Any special quality of service or security needs for storage

## Limitations
* Deployment limits - can you deploy more than once, can you deploy into different namespace
* List specific limitations such as platforms, security, replica's, scaling, upgrades etc.. - noteworthy limits identified
* List deployment limitations such as : restrictions on deploying more than once or into custom namespaces.
* Not intended to provide chart nuances, but more a state of what is supported and not - key items in simple bullet form.
* Does it support IBM Cloud Kubernetes Service in addition to IBM Cloud Private?

## Documentation
* Can have as many supporting links as necessary for this specific workload however don't overload the consumer with unnecessary information.
* Can be links to special procedures in the knowledge center.
