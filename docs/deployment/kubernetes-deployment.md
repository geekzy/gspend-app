# Kubernetes Deployment

## Overview

gSpend is Kubernetes-ready with manifests for production deployment.

---

## Prerequisites

- Kubernetes cluster (1.25+)
- kubectl configured
- Helm (optional)

---

## Manifest Structure

```
devops/kubernetes/
├── namespace.yaml
├── configmaps/
│   ├── frontend-config.yaml
│   ├── auth-service-config.yaml
│   └── financial-service-config.yaml
├── secrets/
│   └── app-secrets.yaml
├── deployments/
│   ├── frontend-deployment.yaml
│   ├── auth-service-deployment.yaml
│   ├── financial-service-deployment.yaml
│   ├── mongodb-deployment.yaml
│   └── redis-deployment.yaml
├── services/
│   ├── frontend-service.yaml
│   ├── auth-service.yaml
│   ├── financial-service.yaml
│   ├── mongodb-service.yaml
│   └── redis-service.yaml
├── ingress/
│   └── ingress.yaml
└── volumes/
    ├── mongodb-pvc.yaml
    └── redis-pvc.yaml
```

---

## Deployment Steps

### 1. Create Namespace
```bash
kubectl apply -f devops/kubernetes/namespace.yaml
```

### 2. Create Secrets
```bash
# Edit with your values first
kubectl apply -f devops/kubernetes/secrets/
```

### 3. Create ConfigMaps
```bash
kubectl apply -f devops/kubernetes/configmaps/
```

### 4. Create Persistent Volumes
```bash
kubectl apply -f devops/kubernetes/volumes/
```

### 5. Deploy Services
```bash
kubectl apply -f devops/kubernetes/deployments/
kubectl apply -f devops/kubernetes/services/
```

### 6. Configure Ingress
```bash
kubectl apply -f devops/kubernetes/ingress/
```

---

## Example Deployment Manifest

```yaml
# auth-service-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: auth-service
  namespace: gspend
spec:
  replicas: 2
  selector:
    matchLabels:
      app: auth-service
  template:
    metadata:
      labels:
        app: auth-service
    spec:
      containers:
      - name: auth-service
        image: gspend/auth-service:latest
        ports:
        - containerPort: 8081
        - containerPort: 9091
        env:
        - name: MONGODB_URI
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: mongodb-uri
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: jwt-secret
        resources:
          limits:
            cpu: "500m"
            memory: "256Mi"
          requests:
            cpu: "100m"
            memory: "128Mi"
        livenessProbe:
          httpGet:
            path: /api/v1/auth/health
            port: 8081
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /api/v1/auth/health
            port: 8081
          initialDelaySeconds: 5
          periodSeconds: 5
```

---

## Scaling

```bash
# Manual scaling
kubectl scale deployment auth-service --replicas=3 -n gspend

# Horizontal Pod Autoscaler
kubectl autoscale deployment auth-service \
  --cpu-percent=70 \
  --min=2 \
  --max=10 \
  -n gspend
```

---

## Monitoring

```bash
# View pods
kubectl get pods -n gspend

# View logs
kubectl logs -f deployment/auth-service -n gspend

# Describe deployment
kubectl describe deployment auth-service -n gspend
```

---

## Rollback

```bash
# View history
kubectl rollout history deployment/auth-service -n gspend

# Rollback to previous
kubectl rollout undo deployment/auth-service -n gspend
```
