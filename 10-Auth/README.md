# Create local registry
```
k3d registry create registry.localhost --port 4000
```

# Add host to file
```
edit /etc/hosts
---
add record with

127.0.0.1 k3d-registry.localhost
```

# Create k3d cluster with local registry
```
k3d cluster create demo \
--api-port 6550 \
-p "80:80@loadbalancer" \
-p "443:443@loadbalancer" \
-p "8080:8080@loadbalancer" \
--k3s-arg "--disable=metrics-server,traefik@server:*" \
--registry-use k3d-registry.localhost:4000
```

# Install Kong follow sheet 09-Auth

# Publish your computer to public with port 8080
```
docker run --net=host \
-it -e NGROK_AUTHTOKEN=<your_api_token> \ 
ngrok/ngrok:latest http 8080
```

# Install keycloak
```
kubectl apply -f keycloak.yaml
```

# Configure keycloak with 'kong' client and realms named 'test_app'

# Build test application
```
cd test_app
---
docker build -t k3d-registry.localhost:4000/krismt/test_app:latest .
```

# Push image to local registry
```
docker push  k3d-registry.localhost:4000/krismt/test_app:latest
```

# Deploy test application
```
kubectl apply -f test_app.yaml
```

# Deploy HTTPRoute 
```
kubectl apply -f httproute-test_app.yaml
```

# Create Kong Plugin
```
kubectl create configmap keycloak-introspection --from-file=kong-custom-plugin -n kong
```

# Update Kong with Plugin
```
helm upgrade kong kong/ingress -n kong --values kong-value.yaml
```

# Edit file 'kong-plugins.yaml' on URL Endpoint and Client Secret
```
Don't forgot to edit file and change Introspection Endpoint and Secret
---
kubectl apply -f kong-plugins.yaml
```

# Deploy HTTPRoute for secure path
```
kubectl apply -f httproute-test_app-secure.yaml
```

# Enjoy to test follow the sheet

