#!/bin/bash -ex

# Export FILESYSTEM to gpfs filesystem, by default it will use gpfs_device

set -e
scripts=$(dirname $0)
. $scripts/acceptance_utils.sh
shopt -s expand_aliases
export PATH=$PATH:/usr/lpp/mmfs/bin/
PODName="write-pod-test"

if [[ -z "${FILESYSTEM}" ]]; then
        FileSystem="gpfs_device"
else
        FileSystem=${FILESYSTEM}
fi


function changefilesystem() {
	filename=$1
	from=$2
	to=$3

	sed -i "s/$from/$to/" $scripts/../deploy/$1
}

function cleanup() {
	kubectl delete storageclass --all
	kubectl delete pvc --all
	kubectl delete pods --all
	maxretry=20
	counter=0
	while [ $counter -lt $maxretry ]; do
		storageclass=`kubectl get storageclass --show-all`
		pvc=`kubectl get pvc --show-all`
		pods=`kubectl get pods --show-all`
		if [ "$storageclass" == "" ] &&  [ "$pvc" == "" ] && [ "$pods" == "" ]; then
			return
		fi 
		sleep 5
		counter=$(($counter+1))		
	done
}

function verify_delete()
{
	local lightweight=$1	
	local fileSystem=$2
	local fileSet=$3
	if [ $lightweight -eq 1 ]; then
		ls /gpfs/$fileSystem/$fileSet && rc=$? || rc=$?	
	else
		mmlsfileset $fileSystem $fileSet && rc=$? || rc=$?
	fi
	if [ $rc -ne 2 ]; then
		if [ $lightweight -eq 1 ]; then
			echo "/gpfs/$fileSystem/$fileSet exists. Exitting"
			exit -1
		else
			echo "$fileSystem fileset $fileSet exists. Exitiing"
			exit -1
		fi		
	fi	
	return
}

function verify_create()
{
	local lightweight=$1
	l_pvcName=$2
	local __pvname=$3
	local __fileset=$4
	local __filesystem=$5
	pvname=`kubectl get pvc $l_pvcName --no-headers -o custom-columns=name:spec.volumeName`
	fileset=`kubectl get pv -o json $pvname | grep -A15 flexVolume |grep fileset |awk '{print $2}' |cut -d, -f1`
	filesystemname=`kubectl get pv -o json $pvname | grep -A15 flexVolume |grep filesystem |awk '{print $2}' |cut -d, -f1`
	filesystemname=`sed -e 's/^"//' -e 's/"$//' <<<"$filesystemname"`
	fileset=`sed -e 's/^"//' -e 's/"$//' <<<"$fileset"`

	if [ $lightweight -eq 1 ];then
		ls /gpfs/$filesystemname/$fileset
	else
		mmlsfileset $filesystemname $fileset
	fi
	eval $__pvname=$pvname	
	eval $__filesystem=$filesystemname
	eval $__fileset=$fileset
}

function execute_scale_test()
{

	backend=$1
	pvName=''
	fileSet=''
	fileSystemName=''
	isLightWeight=0

	if [ "$backend" == "spectrum-scale" ]; then
		storageClassYml="storage_class_fileset.yml"
		pvcYml="pvc_fileset.yml"
		pvcName="ubiquity-claim-fileset"
		podYml="pod.yml"
	elif [ "$backend" == "spectrum-scale-nfs" ]; then
		storageClassYml="storage_class_fileset_nfs.yml"
		pvcYml="pvc_fileset_nfs.yml"
		pvcName="ubiquity-claim-fileset-nfs"
		podYml="pod-nfs.yml"
	else	
		storageClassYml="storage_class_lightweight.yml"
		pvcYml="pvc_lightweight.yml"
		pvcName="ubiquity-claim-lightweight"
		podYml="pod-lightweight.yml"
		isLightWeight=1
	fi
	

	echo "Creating Storage class...."
	kubectl create -f $scripts/../deploy/$storageClassYml

	echo "Listing Storage classes"
	kubectl get storageclass


	echo "Creating Persistent Volume Claim..."
	kubectl create -f $scripts/../deploy/$pvcYml
	sleep 2
	wait_for_item pvc $pvcName  ${PVC_GOOD_STATUS} 20 5

	echo "Listing Persistent Volume Claim..."
	kubectl get pvc

	echo "Listing Persistent Volume..."
	kubectl get pv

	echo "Verify GPFS Fileset.."
	verify_create $isLightWeight $pvcName pvName fileSet fileSystemName


	echo "Creating Test Pod"
	kubectl create -f $scripts/../deploy/$podYml
	wait_for_item pod $PODName Running 30 5


	echo "Listing pods"
	kubectl get pods

	echo "Writing success.txt to mounted volume"
	kubectl exec write-pod-test -c write-pod touch /mnt/success.txt

	echo "Reading from mounted volume"
	kubectl exec write-pod-test -c write-pod ls /mnt


	echo "Cleaning test environment"

	echo "Deleting Pod"
	kubectl delete -f $scripts/../deploy/$podYml --grace-period=0 --force
	wait_for_item_to_delete pod $PODName 30 5

	echo "Deleting Persistent Volume Claim"
	kubectl delete -f $scripts/../deploy/$pvcYml
	wait_for_item_to_delete pv $pvName 30 5

	echo "Listing PVC"
	kubectl get pvc

	echo "Listing PV"
	kubectl get pv

	verify_delete $isLightWeight $fileSystemName $fileSet

	echo "Deleting Storage Class"
	kubectl delete -f $scripts/../deploy/$storageClassYml

	echo "Listing Storage Classes"
	kubectl get storageclass
}
cleanup
changefilesystem "storage_class_fileset.yml" "gold" $FileSystem 
changefilesystem "storage_class_fileset_nfs.yml" "gold" $FileSystem 

echo "Execute spectrum-scale Backend"
execute_scale_test "spectrum-scale"

echo "Execute spectrum-scale-nfs Backend"
execute_scale_test "spectrum-scale-nfs"

#echo "Execute spectrum-scale Backend with type lightweight"
#execute_scale_test "lightweight"

echo "Successfully Executed Spectrum Scale Acceptance Tests"
