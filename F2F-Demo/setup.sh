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
echo "Ready to go!"
echo "First, we set up a fictional redhat organization using the KCP kubectl plugin:"
kubectl kcp ws root

kubectl ws create redhat --ignore-existing
kubectl ws redhat

echo "Then we set up the service provider workspaces for AppStudio and HACBS inside RedHat:"

kubectl ws create appstudio --ignore-existing
kubectl ws appstudio
kubectl apply -f redhat/appstudio/apiresourceschema-pipelines.appstudio.yaml
echo "Workspace AppStudio is set up. Now let us look into the APIExport:"

show_yaml "redhat/appstudio/appstudio-export.yaml"
echo "... and apply it for others to bind:"
kubectl apply -f redhat/appstudio/appstudio-export.yaml

echo "Same goes for the HACBS service provider:"
kubectl ws ..

kubectl ws create hacbs --ignore-existing
kubectl ws hacbs
kubectl apply -f redhat/hacbs/apiresourceschema-pipelines.hacbs.yaml
show_yaml "redhat/hacbs/hacbs-export.yaml"
kubectl apply -f redhat/hacbs/hacbs-export.yaml

echo "Now that the Service providers are online, we set up two consumer organizations: Aspian and Bspian."
echo "Note that these are created using an admin-user:"
echo ""
kubectl ws root
echo ""
echo "Apply the ClusterRole 'aspian-employee' and its binding 'aspian-employees'"
pe "kubectl apply -f aspian/aspian-employees.yaml"
echo ""
echo "...and create the organizational workspace aspian under root."
pe "kubectl ws create aspian --ignore-existing"
echo ""
echo "... then join the workspace and apply the entitlement to bind HACBS to the aspian org."
kubectl kcp ws aspian
pe "kubectl apply -f aspian/hacbs-entitlement.yaml"
echo ""
echo "... create a 'platform' team role/binding for Aspian:"
pe "kubectl apply -f aspian/platform-team.yaml"

echo "... and last but not least create the platform team workspace, where HACBS should get bound and used."
pe "kubectl kcp workspace create platform --ignore-existing"
kubectl kcp workspace ..

echo "Same goes for Bspian and AppStudio"
kubectl apply -f bspian/bspian-employees.yaml
kubectl ws create bspian --ignore-existing
kubectl ws bspian
kubectl apply -f bspian/appstudio-entitlement.yaml
kubectl apply -f bspian/platform-team.yaml
kubectl ws create platform --ignore-existing
kubectl ws ..
echo "All Service providers and orgs are set up! Ready to go: "
### Aspian part ####
echo ""
echo "Scenario 1: Abigail gives access to the team 'platform', part of aspian org in a workspace called platform for HACBS"
echo "We just entitled Aspian to bind HACBS, as they bought a subscription. Abigail binds the API to the platform workspace"
echo ""

show_yaml "aspian/team/hacbs-binding.yaml"
pe "kubectl ws root:aspian:platform"
pe "kubectl apply -f ./aspian/team/hacbs-binding.yaml --token abigail"
echo "Lets look at the bindings for this workspace. We expect the hacbs binding to show up:"
pe "kubectl get apibindings"
echo ""
echo "Now let's create a HACBS pipeline instance to use HACBS in the platform workspace:"
show_yaml "aspian/team/pipeline1.yaml"

pe "kubectl apply -f ./aspian/team/pipeline1.yaml --token abigail"
echo ""
echo "Aspian bought a HACBS basic subscription, allowing them to use one pipeline instance only. Let's try to create a second one, which should fail:"
show_yaml "aspian/team/pipeline2.yaml"
pe "kubectl apply -f ./aspian/team/pipeline2.yaml --token abigail"
echo "Let's look at the existing pipelines. There should not be any new ones:"
pe "kubectl get pipelines"

echo ""
echo "Scenario 2: Abigail heard of AppStudio and wants to use it. She tries to bind AppStudio for the platform workspace, not knowing if her org purchased a subscription:"
echo ""
show_yaml "./aspian/team/appstudio-binding.yaml"
pe "kubectl kcp ws root:aspian:platform"
pe "kubectl apply -f ./aspian/team/appstudio-binding.yaml --token abigail"
echo "Aspian is not entitled to bind AppStudio. Abigail fails to bind this API to the platform workspace."
echo "Lets look at the bindings. We expect nothing new to be bound:"
pe "kubectl get apibindings"
echo ""
echo "Entitlements are stored at the organization level, appear to be plain old kube objects and can be viewed as normal:"
echo ""
pe "kubectl ws root:aspian"
pe "kubectl get entitlement hacbs -oyaml"
echo ""
echo "What if she tries to make superficial changes to this data and create an appstudio entitlement like the following?"
echo ""
show_yaml "aspian/appstudio-entitlement.yaml"
echo "Let's give it a try:"
pe "kubectl create -f aspian/appstudio-entitlement.yaml --token abigail"
echo "As we see, an error shows up. That's because Entitlements are protected from modification by normal users."


### Bspian part ###
echo ""
echo "Scenario 3: Bspian bought an AppStudio subscription."
echo "Ben gives access to team 'platform' to use AppStudio. The team is part of the Bspian organisation, in a workspace called platform."
echo "Bspian is entitled to bind AppStudio, so Ben binds the API to the platform workspace:"
echo ""

show_yaml "bspian/team/appstudio-binding.yaml"
pe "kubectl kcp ws root:bspian:platform"
### TODO: should we have another user who cannot apply these bindings cause their token is missing the respective group? ###
pe "kubectl apply -f ./bspian/team/appstudio-binding.yaml --token ben"
echo "Lets look at the bindings. We expect the app-studio binding to show up:"
pe "kubectl get apibindings"
echo ""
echo "Now Ben wants to create an AppStudio App instance to be used:"
show_yaml "bspian/team/appstudio1.yaml"

pe "kubectl apply -f ./bspian/team/appstudio1.yaml --token ben"
echo ""
echo "Bspian has an AppStudio basic subscription, allowing them to use one App only. Let's try to create a second one, which should fail:"
show_yaml "bspian/team/appstudio2.yaml"
pe "kubectl apply -f ./bspian/team/appstudio2.yaml --token ben"
echo "Let us look at the existing apps. There should not be any new ones:"
pe "kubectl get apps"

echo ""
echo "Ben also heard of HACBS, not knowing if Bspian bought a subscription. He tries to bind HACBS for Bspians platform team: "
echo ""
show_yaml "bspian/team/hacbs-binding.yaml"
pe "kubectl ws root:bspian:platform"
pe "kubectl apply -f ./bspian/team/hacbs-binding.yaml --token ben"
echo "Bspian is not entitled to Hacbs and ben fails to bind this API to the platform workspace"
echo "Lets look at the bindings. We expect nothing new to be bound."
pe "kubectl get apibindings"
echo ""
echo "That was it for now. Thanks for your time! Questions?"