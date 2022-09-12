kubectl create -f - -o yaml << EOF
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
