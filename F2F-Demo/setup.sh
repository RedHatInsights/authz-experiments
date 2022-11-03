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
source ./demo-magic.sh
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
kubectl kcp workspace create aspian --enter --ignore-existing
kubectl apply -f aspian/entitlement.yaml
kubectl apply -f aspian/team.yaml
kubectl kcp workspace create team --ignore-existing
kubectl kcp workspace ..

kubectl kcp workspace create bspian --ignore-existing

