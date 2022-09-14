# How the KCP API Export & Binding works

 export KUBECONFIG=.kcp/admin.kubeconfig

# Check the api resources
``` kubectl api-resources ```

# List workspaces 
``` kubectl get workspaces ```
# Create testwk1 workspace
``` kubectl kcp ws create testwk1 ```
# Create testwk2 workspace
``` kubectl kcp ws create testwk2 ```
# switch to workspace: testwk1 
```kubectl kcp ws testwk1```
# Check the apiexport
```kubectl get apiexports```
# Expose the kind:Widget (data.domain) schema
```kubectl apply -f dd_apischema.yaml```
# Create a API export for kind:Widget
```kubectl apply -f dd_apiexport.yaml```

# Switch to root workspace
```kubectl kcp ws root```
# Switch to workspace testwk2
```kubectl kcp ws create testwk2 ```
# Check the apiexports (data)
```kubectl get apiexports```

# Create a api binding
```kubectl apply -f dd_apibinding.yaml```
# kubectl get api
``` kubectl get apiresourceschemas.apis.kcp.dev```