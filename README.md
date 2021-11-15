# SMI Compliant Controller for Consul Service Mesh

## Cert manager

```shell
kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.1.0/cert-manager.crds.yaml
helm repo add jetstack https://charts.jetstack.io
helm install cert-manager jetstack/cert-manager --namespace cert-manager --create-namespace
```

## Secrets 

```shell
kubectl get secret consul-server-cert -o json \
 | jq 'del(.metadata["namespace","creationTimestamp","resourceVersion","selfLink","uid"])' \
 | kubectl apply -n smi -f -
```

```shell
kubectl get secret consul-controller-acl-token -o json \
 | jq 'del(.metadata["namespace","creationTimestamp","resourceVersion","selfLink","uid"])' \
 | kubectl apply -n smi -f -
```

## Service

If you are deploying the SMI controller to a different namespace from the main consul instance, you will need to create a service
as the CA certificate is setup with the name of the service.

```
echo "
kind: Service
apiVersion: v1
metadata:
  name: consul-server
  namespace: smi
spec:
  type: ExternalName
  externalName: consul-server.default.svc.cluster.local
" | kubectl apply -f -
```

## Helm Install

```shell
helm repo add smi-controller https://servicemeshinterface.github.io/smi-controller-sdk/
helm install smi-controller ./ --namespace smi --create-namespace --values ./consul-values.yaml
```

```yaml
controller:
  enabled: true
  env:
    - name: CONSUL_CAPATH
      value: /tmp/consul/ca/tls.crt
    - name: CONSUL_HTTP_TOKEN
      valueFrom: 
        secretKeyRef:
          name: consul-controller-acl-token
          key: token
    - name: HOST_IP
      valueFrom:
        fieldRef:
          fieldPath: status.hostIP
    - name: CONSUL_HTTP_ADDR
      value: https://$(HOST_IP):8501

  image:
    repository: "nicholasjackson/smi-controller-consul"
    tag: "dev"
    pullPolicy: IfNotPresent

  volumeMounts:
    - name: consul-ca
      mountPath: /tmp/consul/ca
  
  volumes:
    - name: consul-ca
      secret:
        secretName: consul-server-cert

webhook:  
  enabled: true
  service: "smi-webhook"
  port: 443
  additionalDNSNames:
    - localhost
```