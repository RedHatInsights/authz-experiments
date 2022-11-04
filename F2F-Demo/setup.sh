#!/bin/sh

RED="\e[31m"
GREEN="\e[32m"
YELLOW="\e[91m" # 33 is regular yellow
ENDCOLOR="\e[0m"

echo_yellow() {
    echo -e "${YELLOW}${1}${ENDCOLOR}"
}

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
echo_yellow "using KUBECONFIG=$KUBECONFIG"
echo_yellow "Ready to go!"
echo_yellow "First, we set up a fictional redhat organization using the KCP kubectl plugin:"

wait 

kubectl kcp ws root

kubectl ws create redhat --ignore-existing
kubectl ws redhat

echo_yellow "Then we set up the service provider workspaces for AppStudio and HACBS inside RedHat:"

kubectl ws create appstudio --ignore-existing
kubectl ws appstudio
kubectl apply -f redhat/appstudio/apiresourceschema-pipelines.appstudio.yaml
echo_yellow "Workspace AppStudio is set up. Now let us look into the APIExport:"

show_yaml "redhat/appstudio/appstudio-export.yaml"
echo_yellow "... and apply it for others to bind:"
kubectl apply -f redhat/appstudio/appstudio-export.yaml

echo_yellow "Same goes for the HACBS service provider:"
kubectl ws ..

kubectl ws create hacbs --ignore-existing
kubectl ws hacbs
kubectl apply -f redhat/hacbs/apiresourceschema-pipelines.hacbs.yaml
show_yaml "redhat/hacbs/hacbs-export.yaml"
kubectl apply -f redhat/hacbs/hacbs-export.yaml

echo_yellow "Now that the Service providers are online, we set up two consumer organizations: Aspian and Bspian."
echo_yellow "Note that these are created using an admin-user:"
echo_yellow ""
kubectl ws root
echo_yellow ""
echo_yellow "Apply the ClusterRole 'aspian-employee' and its binding 'aspian-employees'"
pe "kubectl apply -f aspian/aspian-employees.yaml"
echo_yellow ""
echo_yellow "...and create the organizational workspace aspian under root."
pe "kubectl ws create aspian --ignore-existing"
echo_yellow ""
echo_yellow "... then join the workspace and apply the entitlement to bind HACBS to the aspian org."
kubectl kcp ws aspian
pe "kubectl apply -f aspian/hacbs-entitlement.yaml"
echo_yellow ""
echo_yellow "... create a 'platform' team role/binding for Aspian:"
pe "kubectl apply -f aspian/platform-team.yaml"

echo_yellow "... and last but not least create the platform team workspace, where HACBS should get bound and used."
pe "kubectl kcp workspace create platform --ignore-existing"
kubectl kcp workspace ..

echo_yellow "Same goes for Bspian and AppStudio"
kubectl apply -f bspian/bspian-employees.yaml
kubectl ws create bspian --ignore-existing
kubectl ws bspian
kubectl apply -f bspian/appstudio-entitlement.yaml
kubectl apply -f bspian/platform-team.yaml
kubectl ws create platform --ignore-existing
kubectl ws ..
echo_yellow "All Service providers and orgs are set up! Ready to go: "
### Aspian part ####
echo_yellow ""
echo_yellow "Scenario 1: Abigail gives access to the team 'platform', part of aspian org in a workspace called platform for HACBS"
echo_yellow "We just entitled Aspian to bind HACBS, as they bought a subscription. Abigail binds the API to the platform workspace"
echo_yellow ""

show_yaml "aspian/team/hacbs-binding.yaml"
pe "kubectl apply -f ./aspian/team/hacbs-binding.yaml --token abigail"
echo_yellow "Lets look at the bindings for this workspace. We expect the hacbs binding to show up:"
pe "kubectl get apibindings"
echo_yellow ""
echo_yellow "Now let's create a HACBS pipeline instance to use HACBS in the platform workspace:"
show_yaml "aspian/team/pipeline1.yaml"

pe "kubectl apply -f ./aspian/team/pipeline1.yaml --token abigail"
echo_yellow ""
echo_yellow "Aspian bought a HACBS basic subscription, allowing them to use one pipeline instance only. Let's try to create a second one, which should fail:"
show_yaml "aspian/team/pipeline2.yaml"
pe "kubectl apply -f ./aspian/team/pipeline2.yaml --token abigail"
echo_yellow "Let's look at the existing pipelines. There should not be any new ones:"
pe "kubectl get pipelines"

wait

echo_yellow ""
echo_yellow "Scenario 2: Abigail heard of AppStudio and wants to use it. She tries to bind AppStudio for the platform workspace, not knowing if her org purchased a subscription:"
echo_yellow ""
pe "kubectl kcp ws root:aspian:platform"
show_yaml "./aspian/team/appstudio-binding.yaml"
pe "kubectl apply -f ./aspian/team/appstudio-binding.yaml --token abigail"
echo_yellow "Aspian is not entitled to bind AppStudio. Abigail fails to bind this API to the platform workspace."
echo_yellow "Lets look at the bindings. We expect nothing new to be bound:"
pe "kubectl get apibindings"
echo_yellow ""
echo_yellow "Entitlements are stored at the organization level, they are be plain old kube objects and can be viewed as normal:"
echo_yellow ""
pe "kubectl ws root:aspian"
pe "kubectl get entitlement hacbs -oyaml"
echo_yellow ""
echo_yellow "What if she tries to make superficial changes to this data and create an appstudio entitlement like the following?"
echo_yellow ""
show_yaml "aspian/appstudio-entitlement.yaml"
echo_yellow "Let's give it a try:"
pe "kubectl create -f aspian/appstudio-entitlement.yaml --token abigail"
echo_yellow "As we see, an error shows up. That's because Entitlements are protected from modification by normal users."

wait
### Bspian part ###
echo_yellow ""
echo_yellow "Scenario 3: Bspian bought an AppStudio subscription."
echo_yellow "Ben gives access to team 'platform' to use AppStudio. The team is part of the Bspian organisation, in a workspace called platform."
echo_yellow "Bspian is entitled to bind AppStudio, so Ben binds the API to the platform workspace:"
echo_yellow ""

pe "kubectl kcp ws root:bspian:platform"
show_yaml "bspian/team/appstudio-binding.yaml"
### TODO: should we have another user who cannot apply these bindings cause their token is missing the respective group? ###
pe "kubectl apply -f ./bspian/team/appstudio-binding.yaml --token ben"
echo_yellow "Lets look at the bindings. We expect the app-studio binding to show up:"
pe "kubectl get apibindings"
echo_yellow ""
echo_yellow "Now Ben wants to create an AppStudio App instance to be used:"
show_yaml "bspian/team/appstudio1.yaml"

pe "kubectl apply -f ./bspian/team/appstudio1.yaml --token ben"
echo_yellow ""
echo_yellow "Bspian has an AppStudio basic subscription, allowing them to use one App only. Let's try to create a second one, which should fail:"
show_yaml "bspian/team/appstudio2.yaml"
pe "kubectl apply -f ./bspian/team/appstudio2.yaml --token ben"
echo_yellow "Let us look at the existing apps. There should not be any new ones:"
pe "kubectl get apps"

echo_yellow ""
echo_yellow "Ben also heard of HACBS, not knowing if Bspian bought a subscription. He tries to bind HACBS for Bspians platform team: "
echo_yellow ""
pe "kubectl ws root:bspian:platform"
show_yaml "bspian/team/hacbs-binding.yaml"
pe "kubectl apply -f ./bspian/team/hacbs-binding.yaml --token ben"
echo_yellow "Bspian is not entitled to Hacbs and ben fails to bind this API to the platform workspace"
echo_yellow "Lets look at the bindings. We expect nothing new to be bound."
pe "kubectl get apibindings"
echo_yellow ""
echo_yellow "That was it for now. Thanks for your time! Questions?"
