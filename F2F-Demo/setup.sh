#!/bin/sh


show_yaml() {
    if command -v bat &> /dev/null
    then
        pe "bat ${1}"
    else
        pe "cat ${1}"
    fi
}

#shellcheck source=demo-magic.sh
source ./demo-magic.sh -n
clear 

p "Start kcp ..."
echo "using KUBECONFIG=$KUBECONFIG"
echo "Setting up the environment"
echo "Setting up services workspaces ..."
kubectl kcp workspace root

kubectl kcp workspace create redhat --enter --ignore-existing

kubectl kcp workspace create appstudio --enter --ignore-existing
kubectl apply -f redhat/appstudio/apiresourceschema-pipelines.appstudio.yaml
show_yaml "redhat/appstudio/appstudio-export.yaml"
kubectl apply -f redhat/appstudio/appstudio-export.yaml
kubectl kcp workspace ..

kubectl kcp workspace create hacbs --enter --ignore-existing
kubectl apply -f redhat/hacbs/apiresourceschema-pipelines.hacbs.yaml
show_yaml "redhat/hacbs/hacbs-export.yaml"
kubectl apply -f redhat/hacbs/hacbs-export.yaml

echo "Setting up Users/Consumers workspaces ..."

kubectl kcp workspace root
kubectl apply -f aspian/aspian-employees.yaml
kubectl kcp workspace create aspian --ignore-existing
kubectl kcp workspace aspian
kubectl apply -f aspian/entitlement.yaml
kubectl apply -f aspian/platform-team.yaml
kubectl kcp workspace create platform --ignore-existing
kubectl kcp workspace ..

kubectl kcp workspace create bspian --ignore-existing

echo ""
echo "Abigail gives access to team 'platform', part of aspian org in workspace called platform-ws for HACBS"
echo "Aspian is entitled to HACBS and Abigail give entitlement to the platform-ws workspace"
echo ""

show_yaml "./aspian/team/hacbs-binding.yaml"
pe kubectl kcp ws root:aspian:platform
pe kubectl apply -f ./aspian/team/hacbs-binding.yaml --token abigail
echo "Lets look at the bindings. We expect the hacbs binding to show up"
pe kubectl get apibindings


echo ""
echo "Abigail tries to access to team 'platform', part of aspian org in workspace called platform-ws for AppStudio"
echo "Aspian is not entitled to AppStudio and Abigail fails to give entitlement to the platform-ws workspace"
echo ""
show_yaml "./aspian/team/appstudio-binding.yaml"
pe kubectl kcp ws root:aspian:platform
pe kubectl apply -f ./aspian/team/appstudio-binding.yaml --token abigail
echo "Lets look at the bindings. We expect nothing to be bound."
pe kubectl get apibindings
echo ""
echo "Now we want to create a hacbs pipeline instance as we are bound to hacbs."
show_yaml "aspian/team/pipeline1.yaml"

pe kubectl apply -f ./aspian/team/pipeline1.yaml --token abigail
echo ""
echo "We have a hacbs basic subscription, allowing us one instance only. Let's create a second one."
show_yaml "aspian/team/pipeline2.yaml"
pe kubectl apply -f ./aspian/team/pipeline2.yaml --token abigail
echo "Oh dang, quota reached! But did it really not create the 2nd pipeline?"
pe kubectl get pipelines
echo "Yip, quota applied."

