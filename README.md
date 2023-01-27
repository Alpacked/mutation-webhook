# Mutation admission webhook

## What is mutation webhook
A mutation webhook is a Kubernetes admission controller that allows you to modify resources before they are persisted in the API server. Mutation webhooks are used to enforce specific policies or to add additional functionality to the Kubernetes API. They are called by the API server during the admission process, and they can modify the request object before it is persisted.

For example, you can use a mutation webhook to automatically add a specific label or annotation to all pods created in a namespace, to automatically add a specific security context to all containers, or to validate that all resources have specific labels or annotations before they are persisted in the API server.

Mutation webhooks are configured as CustomResourceDefinitions(CRDs) in Kubernetes, and they must be registered with the API server before they can be used. Once registered, the API server will call the webhook for each request that matches the specified rules in the webhook configuration.

## Admission Logic
A set of mutations are implemented in this extensible repository. Those happen on the fly when a pod is deployed and no further resources are tracked and updated (ie. no controller logic).

### Mutating Webhooks
#### Implemented
- [inject env](pkg/mutation/inject_env.go): inject environment variables into the pod such as `KUBE: true`
- [inject securityContext](pkg/mutation/security_context.go): inject PodSecurityContext and SecurityContext for Containers and InitContainers. Webhook checks if existing SecurityContext isn't empty on all levels in spec and add only if it's empty.

#### How to add a new pod mutation
To add a new pod mutation, create a file `pkg/mutation/MUTATION_NAME.go`, then create a new struct implementing the `mutation.podMutator` interface and write your logic for mutation.

## Installation
This project includes helm chart for deploying mutation webhook.

### How to build and push image
```bash
docker build . -t [REPOSITORY]/[IMAGE_NAME]
docker push [REPOSITORY]/[IMAGE_NAME]
```

### How to deploy
```bash
    helm upgrade --install [RELEASE_NAME] ./webhook-chart
```

### Requirements
- cert-manager (1.11 and later)

When you configure a mutation webhook in Kubernetes, the API server needs to communicate with the webhook over a secure connection. To ensure that the connection is secure, you need to provide the API server with a certificate and key that it can use to authenticate the webhook.

The certificate and key are used to establish a secure connection (TLS) between the API server and the webhook. This ensures that the request and response between the API server and the webhook are encrypted and cannot be intercepted or tampered with by a third party.

Additionally, the certificate contains the webhook's public key, which the API server can use to verify that the webhook is who it claims to be. This is important, because it ensures that the webhook is actually the one that you expect it to be, and not a rogue service that is attempting to impersonate the webhook.

It's important to note that the certificate must be signed by a trusted certificate authority (CA) or self-signed. When the certificate is self-signed, you must configure the API server to trust the certificate by providing the root CA of the certificate.

We decided to manage those certificates with cert-manager and it's required for correct work of webhook.

## Testing
Unit tests can be run with the following command:
```bash
go test ./pkg/mutation/
ok      github.com/alpacked/mutation-webhook/pkg/mutation       0.354s

go test ./pkg/admission
ok      github.com/alpacked/mutation-webhook/pkg/admission      0.487s
```

## Examples of work
Let's assume that our webhook is deployed into cluster.
So to show how webhook works let's create simple pod

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
  namespace: apps
spec:
  containers:
    - args:
        - sleep
        - "3600"
      image: busybox
      name: test-container
  restartPolicy: Always
```

And we have next live manifest of running pod

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
  namespace: apps
spec:
  containers:
    - name: test-container
      image: busybox
      args:
        - sleep
        - '3600'
      env:
        - name: KUBE
          value: 'true'
      imagePullPolicy: Always
      securityContext:
        capabilities:
          drop:
            - ALL
        runAsNonRoot: true
        readOnlyRootFilesystem: true
        allowPrivilegeEscalation: false
  restartPolicy: Always
  securityContext:
    runAsUser: 1000
    runAsGroup: 1000
    runAsNonRoot: true
...
```

So we see that our webhook added SecurityContext on Pod level and also set it for Container. Also it injected env var `KUBE` with value `true`.

Logs from webhook:
```log
time="2023-01-27T14:53:27Z" level=debug msg="received mutation request" uri="/mutate-pods?timeout=5s"
time="2023-01-27T14:53:27Z" level=debug msg="pod env injected {KUBE true nil}" mutation=inject_env pod_name=test-pod
time="2023-01-27T14:53:27Z" level=info msg="setting pod security context" inj_sec_context="&PodSecurityContext{SELinuxOptions:nil,RunAsUser:*1000,RunAsNonRoot:*true,SupplementalGroups:[],FSGroup:nil,RunAsGroup:*1000,Sysctls:[]Sysctl{},WindowsOptions:nil,FSGroupChangePolicy:nil,SeccompProfile:nil,}" mutation=inj_sec_context pod_name=test-pod
time="2023-01-27T14:53:27Z" level=info msg="setting security context for container: test-container" inj_sec_context="&SecurityContext{Capabilities:&Capabilities{Add:[],Drop:[ALL],},Privileged:nil,SELinuxOptions:nil,RunAsUser:nil,RunAsNonRoot:*true,ReadOnlyRootFilesystem:*true,AllowPrivilegeEscalation:*false,RunAsGroup:nil,ProcMount:nil,WindowsOptions:nil,SeccompProfile:nil,}" mutation=inj_sec_context pod_name=test-pod
time="2023-01-27T14:53:27Z" level=debug msg="sending response" uri="/mutate-pods?timeout=5s"
time="2023-01-27T14:53:27Z" level=debug msg="{\"kind\":\"AdmissionReview\",\"apiVersion\":\"admission.k8s.io/v1\",\"response\":{\"uid\":\"6c3dc0ed-174f-460b-bfcb-6b8998f7b38a\",\"allowed\":true,\"patch\":\"W3sib3AiOiJhZGQiLCJwYXRoIjoiL3NwZWMvY29udGFpbmVycy8wL2VudiIsInZhbHVlIjpbeyJuYW1lIjoiS1VCRSIsInZhbHVlIjoidHJ1ZSJ9XX0seyJvcCI6ImFkZCIsInBhdGgiOiIvc3BlYy9jb250YWluZXJzLzAvc2VjdXJpdHlDb250ZXh0IiwidmFsdWUiOnsiYWxsb3dQcml2aWxlZ2VFc2NhbGF0aW9uIjpmYWxzZSwiY2FwYWJpbGl0aWVzIjp7ImRyb3AiOlsiQUxMIl19LCJyZWFkT25seVJvb3RGaWxlc3lzdGVtIjp0cnVlLCJydW5Bc05vblJvb3QiOnRydWV9fSx7Im9wIjoiYWRkIiwicGF0aCI6Ii9zcGVjL3NlY3VyaXR5Q29udGV4dC9ydW5Bc0dyb3VwIiwidmFsdWUiOjEwMDB9LHsib3AiOiJhZGQiLCJwYXRoIjoiL3NwZWMvc2VjdXJpdHlDb250ZXh0L3J1bkFzTm9uUm9vdCIsInZhbHVlIjp0cnVlfSx7Im9wIjoiYWRkIiwicGF0aCI6Ii9zcGVjL3NlY3VyaXR5Q29udGV4dC9ydW5Bc1VzZXIiLCJ2YWx1ZSI6MTAwMH1d\",\"patchType\":\"JSONPatch\"}}" uri="/mutate-pods?timeout=5s"
```
