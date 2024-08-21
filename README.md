# Create Cluster
```
k3d cluster create demo \
--api-port 6550 \
-p "80:80@loadbalancer" \
-p "443:443@loadbalancer" \
-p "8080:8080@loadbalancer" \
--k3s-arg "--disable=traefik@server:*"
```
#  Install Gateway CRD

```
kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.1.0/standard-install.yaml
```

# Install Gateway Class
```
echo "
---
apiVersion: gateway.networking.k8s.io/v1
kind: GatewayClass
metadata:
 name: kong
 annotations:
   konghq.com/gatewayclass-unmanaged: 'true'

spec:
 controllerName: konghq.com/kic-gateway-controller
---
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
 name: kong
spec:
 gatewayClassName: kong
 listeners:
 - name: proxy
   hostname: www.example.com
   port: 80
   protocol: HTTP
   allowedRoutes:
      namespaces:
         from: All
 - name: proxy-ssl
   hostname: www.example.com
   port: 443
   protocol: HTTPS
   tls:
      mode: Terminate
      certificateRefs:
         - kind: Secret
           name: www-secret
   allowedRoutes:
      namespaces:
         from: All
" | kubectl apply -f -

```

# Install Kong
```
helm install kong kong/ingress -n kong --create-namespace 
```

## Check by
```
curl -s -I -HHost:www.example.com http://localhost

---- Edit the hostname to resolver and 
open browser: https://www.example.com
```

# Install Cert-Manager
```
helm install \
  cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version v1.15.3 \
  --set crds.enabled=true
```

# Create SelfSigned
```
echo "
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: selfsigned-issuer
spec:
  selfSigned: {}
---
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: my-selfsigned-ca
spec:
  isCA: true
  commonName: my-selfsigned-ca
  secretName: root-secret
  privateKey:
    algorithm: RSA
    encoding: PKCS1
    size: 4096
  issuerRef:
    name: selfsigned-issuer
    kind: ClusterIssuer
    group: cert-manager.io
---
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: my-ca-issuer
spec:
  ca:
    secretName: root-secret
" | kubectl apply -f -
```

# Create www -selfSigned
```
echo "
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: www-selfsigned
spec:
  isCA: false
  secretName: www-secret
  privateKey:
    algorithm: RSA
    encoding: PKCS1
    size: 4096
  dnsNames:
    - 'example.com'
    - '*.example.com'
    - 'www.example.com'
  usages:
    - server auth
    - client auth
  issuerRef:
    name: my-ca-issuer
    kind: Issuer
    group: cert-manager.io
" | kubectl apply -f -
```

# Create Cert file from Root CA
```
kubectl get secret root-secret -o json -o=jsonpath="{.data.tls\.crt}" | base64 -d > selfsigned-root-secret.cer
```

# Create Deployment and Service
```
kubectl create deployment hello-node --image=nginx
kubectl expose deployment hello-node --type=ClusterIP --port=80 --target-port=80
```

# Create HTTPRoute
```
echo "
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
 name: nginx
 annotations:
   konghq.com/strip-path: 'true'
spec:
 parentRefs:
 - name: kong
   sectionName: proxy
 hostnames:
 - www.example.com
 rules:
 - matches:
   - path:
       type: PathPrefix
       value: /
   backendRefs:
   - name: hello-node
     kind: Service
     port: 80

---
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
 name: nginx-ssl
 annotations:
   konghq.com/strip-path: 'true'
spec:
 parentRefs:
 - name: kong
   sectionName: proxy-ssl
 hostnames:
 - www.example.com
 rules:
 - matches:
   - path:
       type: PathPrefix
       value: /
   backendRefs:
   - name: hello-node
     kind: Service
     port: 80
" | kubectl apply -f -

```
