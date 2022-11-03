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

kubectl kcp workspace create redhat --ignore-existing
kubectl kcp workspace redhat

kubectl kcp workspace create appstudio --ignore-existing
kubectl kcp workspace appstudio
kubectl apply -f redhat/appstudio/apiresourceschema-pipelines.appstudio.yaml
show_yaml "redhat/appstudio/appstudio-export.yaml"
kubectl apply -f redhat/appstudio/appstudio-export.yaml
kubectl kcp workspace ..

kubectl kcp workspace create hacbs --ignore-existing
kubectl kcp workspace hacbs
kubectl apply -f redhat/hacbs/apiresourceschema-pipelines.hacbs.yaml
show_yaml "redhat/hacbs/hacbs-export.yaml"
kubectl apply -f redhat/hacbs/hacbs-export.yaml

echo "Setting up Users/Consumers workspaces ..."

kubectl kcp workspace root
kubectl apply -f aspian/aspian-employees.yaml
kubectl kcp workspace create aspian --ignore-existing
kubectl kcp workspace aspian
kubectl apply -f aspian/hacbs-entitlement.yaml
kubectl apply -f aspian/platform-team.yaml
kubectl kcp workspace create platform --ignore-existing
kubectl kcp workspace ..

kubectl kcp workspace create bspian --ignore-existing

echo ""
echo "Abigail gives access to team 'platform', part of aspian org in workspace called platform for HACBS"
echo "Aspian is entitled to HACBS and Abigail binds the API to the platform workspace"
echo ""

show_yaml "./aspian/team/hacbs-binding.yaml"
pe "kubectl kcp ws root:aspian:platform"
pe "kubectl apply -f ./aspian/team/hacbs-binding.yaml --token abigail"
echo "Lets look at the bindings. We expect the hacbs binding to show up"
pe "kubectl get apibindings"
echo ""
echo "Now we want to create a hacbs pipeline instance to use with hacbs."
show_yaml "aspian/team/pipeline1.yaml"

pe "kubectl apply -f ./aspian/team/pipeline1.yaml --token abigail"
echo ""
echo "We have a hacbs basic subscription, allowing us one instance only. Let's try to create a second one, which should fail."
show_yaml "aspian/team/pipeline2.yaml"
pe "kubectl apply -f ./aspian/team/pipeline2.yaml --token abigail"
echo "Let's look at the existing pipelines. There should not be any new ones."
pe "kubectl get pipelines"

echo ""
echo "Abigail tries to bind AppStudio for the platform team in the aspian org"
echo "Aspian is not entitled to AppStudio and Abigail fails to bind this API to the platform workspace"
echo ""
show_yaml "./aspian/team/appstudio-binding.yaml"
pe "kubectl kcp ws root:aspian:platform"
pe "kubectl apply -f ./aspian/team/appstudio-binding.yaml --token abigail"
echo "Lets look at the bindings. We expect nothing new to be bound."
pe "kubectl get apibindings"
echo ""
echo "Entitlements are stored at the organization level, appear to be plain old kube objects and can be viewed as normal:"
echo ""
pe "kubectl ws root:aspian"
pe "kubectl get entitlement hacbs -oyaml"
echo "What if she were to make some superficial changes to this data and create an appstudio entitlement like the following?"
echo ""
show_yaml aspian/appstudio-entitlement.yaml
echo "Let's give it a try."
pe "kubectl create -f aspian/appstudio-entitlement.yaml --token abigail"
echo "We should have an error here. Entitlement objects are protected from modification by normal users"