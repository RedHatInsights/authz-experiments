# Checking subject access review in KCP

This experiment is try a SAR for a service account using KCP

 # Create Role
`kubectl create -f kafka_clusteradmin_role.yaml`

# Create Service account
`kubectl create -f sa-kafka-client-1.yaml`

# Try subject access review

Try the SAR using the sar.sh 

It should return
``` allowed: false``` - Since there is no rolebinding at this stage

# Create Role Binding 
`kubectl create -f kafka_clusteradmin_rolebinding.yaml`

# Try subject access review again
Execute the sar using sar.sh

It should return
``` allowed: true``` - Since rolebinding is present now
