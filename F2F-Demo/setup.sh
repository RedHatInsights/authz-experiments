echo "Start kcp ..."  
read kubeconfig

#export KUBECONFIG={$kubeconfig}

kubectl kcp workspace root

kubectl kcp workspace create redhat --enter --ignore-existing
kubectl kcp workspace create hacbs --ignore-existing
kubectl kcp workspace create appstudio --ignore-existing

kubectl kcp workspace root
kubectl kcp workspace create aspian --ignore-existing
kubectl kcp workspace create bspian --ignore-existing

