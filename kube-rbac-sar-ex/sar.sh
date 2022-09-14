kubectl create -f - -o yaml << EOF
apiVersion: authorization.k8s.io/v1
kind: SubjectAccessReview
spec:
  #user:system:serviceaccount:<namepsace of SA>:<SA Name>
  user: system:serviceaccount:default:kafka-client-1
  # This is needed for now as there is a hardcoded check in the authz module of KCP
  # This might not be needed in future if the bug is resolved
  groups: ['system:authenticated']
  extra:
    # Use the ws below for cluster-name
    'authentication.kubernetes.io/cluster-name': ['<your ws>']
  resourceAttributes:
    group: kafka.io
    resource: topics/test/abc
    verb: create
    # Target Namespace
    namespace: default
EOF

