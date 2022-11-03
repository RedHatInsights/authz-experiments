echo "Start kcp ..."  
read kubeconfig

#export KUBECONFIG={$kubeconfig}

kubectl kcp workspace root

kubectl kcp workspace create redhat --enter --ignore-existing
kubectl kcp workspace create appstudio --ignore-existing

kubectl kcp workspace create hacbs --enter --ignore-existing
kubectl apply -f redhat/hacbs/hacbs-export.yaml

kubectl kcp workspace root
kubectl apply -f aspian/aspian-employees.yaml
kubectl kcp workspace create aspian --enter --ignore-existing
kubectl apply -f aspian/entitlement.yaml
kubectl apply -f aspian/team.yaml
kubectl kcp workspace create team --ignore-existing
kubectl kcp workspace ..

kubectl kcp workspace create bspian --ignore-existing

