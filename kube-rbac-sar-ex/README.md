# Checking subject access review


# Create Role
`kubectl create -f kakfa_clusteradmin_role.yaml`

# Create Service account
`kubectl create -f sa-kafka-client-1.yaml`


# Try subject access review

```kubectl create -f - -o yaml << EOF
apiVersion: authorization.k8s.io/v1
kind: SubjectAccessReview
spec:
  resourceAttributes:
    group: kafka.io
    resource: topics/test/abc
    verb: create
    namespace: default
  user: "system:serviceaccount:default:kafka-client-1"
EOF
```

It should return
``` allowed: false```

# Create Role Binding 
`kubectl create -f kakfa_clusteradmin_rolebinding.yaml`

# Try subject access review again
```kubectl create -f - -o yaml << EOF
apiVersion: authorization.k8s.io/v1
kind: SubjectAccessReview
spec:
  resourceAttributes:
    group: kafka.io
    resource: topics/test/abc
    verb: create
    namespace: default
  user: "system:serviceaccount:default:kafka-client-1"
EOF
```

It should return true (using on Kubernetes cluster)