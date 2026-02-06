# Azure Tagger API

**Production-Style Go REST API deployed to Azure Container Apps**

A cloud-native REST API built with Go that registers Azure resources and applies tags using the Azure SDK.

This project demonstrates real-world backend engineering practices including:

* Clean architecture
* Dependency injection
* Unit testing with table-driven tests
* Swagger/OpenAPI documentation
* Docker containerization
* Azure Container Registry integration
* Azure Container Apps deployment
* CI/CD readiness
* Cloud debugging & environment configuration

---

## Live Architecture

```
Client (REST / Swagger UI)
        ↓
Azure Container Apps (public endpoint)
        ↓
Docker Container (Go API)
        ↓
Azure SDK (DefaultAzureCredential)
        ↓
Azure Resource Manager (Tag Updates)
```

---

## Core Features

### REST API

* Create resource metadata entries
* List resources
* Get resource by ID
* Delete resource
* Apply tags directly to Azure resources

### Cloud Integration

* Azure SDK for Go
* DefaultAzureCredential authentication
* Service Principal configuration
* Real Azure Resource Manager tag updates

### Documentation

* Swagger UI
* OpenAPI spec generation via swaggo

### Testing

* Unit tests for handlers and store
* Table-driven tests
* Coverage reporting
* Mocked Azure interface for isolation

---

## Tech Stack

**Backend**

* Go
* Chi Router
* Swaggo (Swagger)
* Azure SDK for Go

**Cloud**

* Azure Container Registry
* Azure Container Apps
* Azure Log Analytics
* Azure RBAC

**DevOps**

* Docker
* GitHub Actions (CI-ready)
* Azure CLI automation

---

## Local Development

### Run locally

```bash
go run ./cmd/api
```

API available at:

* [http://localhost:8080/health](http://localhost:8080/health)
* [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
* [http://localhost:8080/v1/resources](http://localhost:8080/v1/resources)

---

### Environment Variables (.env)

```env
PORT=8080
AZURE_SUBSCRIPTION_ID=...
AZURE_TENANT_ID=...
AZURE_CLIENT_ID=...
AZURE_CLIENT_SECRET=...
```

Loaded locally via PowerShell script.

---

## Docker Build

```bash
docker build -t taggeracr35316.azurecr.io/azure-tagger-api:1.2 .
docker push taggeracr35316.azurecr.io/azure-tagger-api:1.2
```

---

## Azure Deployment (Container Apps)

```bash
az containerapp update \
  -n azure-tagger-api \
  -g rg-azure-tagger-dev \
  --image taggeracr35316.azurecr.io/azure-tagger-api:1.2
```

Public endpoint:

```
https://azure-tagger-api.<region>.azurecontainerapps.io
```

---

## End-to-End Cloud Test

1. Register an Azure resource ID
2. Call `/v1/resources/{id}/apply-tags`
3. Confirm tags via Azure CLI:

```bash
az resource show --ids <RESOURCE_ID> --query tags
```

Tags are successfully applied using the Azure SDK from within the container.

---

## Real-World Engineering Challenges Solved

### Routing Bug (Chi Router)

Accidentally used `router.*` instead of `r.*` inside `router.Route("/v1", ...)`.

Impact:

* Swagger showed `/v1` routes
* Runtime did not mount them correctly
* Caused 404 inconsistencies

Resolution:

* Correct route scoping
* Standardized API base path

---

### Azure Free Trial Limitations

ACR Tasks were blocked in Free Trial subscription.

Resolution:

* Switched to local Docker build
* Pushed directly to ACR
* Updated Container App manually

---

### In-Memory Store Reset

Container revision updates wiped memory store.

Learning:

* Stateless containers require persistent storage
* Next iteration: PostgreSQL or Azure Cosmos DB

---

### Environment Variable Issues

Container App lacked `AZURE_SUBSCRIPTION_ID`.

Resolution:

* Configured secrets via:

  ```
  az containerapp secret set
  az containerapp update --set-env-vars
  ```
* Verified via container logs

---

### PowerShell Docker Tag Parsing Issue

Variable expansion failed due to `:` parsing.

Resolution:
Used:

```powershell
"${ACR_SERVER}/${IMAGE_NAME}:${IMAGE_TAG}"
```

---

## Testing Strategy

```bash
go test ./...
```

Generate coverage:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

Tested:

* Store logic
* Handler validation
* Azure interface (mocked)

---

## What This Project Demonstrates

This is not a simple CRUD API.

It demonstrates:

* Cloud authentication flows
* Secure secret management
* Container-based deployment
* Azure RBAC configuration
* Production-style routing design
* Stateless service behavior
* CI/CD-ready structure
* Real-world debugging process

---

## Next Iterations

Planned improvements:

* Replace MemoryStore with PostgreSQL
* Add Managed Identity (remove client secret)
* Add structured logging (Zap)
* Add health probes
* Add request validation middleware
* Add integration tests
* Add staging environment
* Implement GitHub Actions full CI/CD

---

## Why This Project Matters

This project simulates how a real backend service would:

* Run inside containers
* Authenticate to cloud providers
* Handle configuration securely
* Be deployed through versioned images
* Debug real production issues
* Scale statelessly

It reflects a transition from junior-level CRUD coding to cloud-ready backend engineering.

---


## Extra problems

## README — Common Azure ACR + Go Docker Build Issues (What We Hit & How We Fixed It)

This project was containerized and pushed to **Azure Container Registry (ACR)**. Along the way we ran into a few common problems. This short report explains what went wrong, why, and how to avoid it next time.

---


## 2) Go version mismatch during Docker build

### Symptom

Build failed with:

* `go.mod requires go >= 1.25.6 (running go 1.22.12; GOTOOLCHAIN=local)`

### Cause

The Dockerfile used:

* `FROM golang:1.22-alpine`

…but `go.mod` requires a newer Go version. Since toolchain auto-download is disabled in containers (`GOTOOLCHAIN=local`), it won’t fetch a newer compiler automatically.

### Fix options

**Recommended:** update Dockerfile to match required Go version:

* `FROM golang:1.25-alpine AS build` (or whatever version `go.mod` needs)

**Alternative:** lower the Go version in `go.mod` **only if the code/dependencies truly support it**.

### Prevention checklist

* Keep Docker `golang:X` image version aligned with the `go` version in `go.mod`.

---

## 3) Push to ACR fails: `authentication required`

### Symptom

Docker push failed with:

* `error from registry: authentication required`

### Cause

Docker was not authenticated to the ACR registry, or the Azure user didn’t have permission to push.

### Fix

Login using Azure CLI:

```powershell
az login
az acr login --name taggeracr35316
docker push taggeracr35316.azurecr.io/azure-tagger-api:0
```

If login succeeds but push still fails, the account likely lacks the **AcrPush** role on the registry.

### Alternative approach (often easiest)

Let ACR build and push the image:

```powershell
az acr build -r taggeracr35316 -t azure-tagger-api:0 .
```

### Prevention checklist

* Run `az acr login --name <registryName>` before pushing.
* Make sure your account has **AcrPush** role assigned in the registry IAM.

---

## Quick “Known Good” Commands

### Build locally

```powershell
docker build -t taggeracr35316.azurecr.io/azure-tagger-api:0 .
```

### Login + push

```powershell
az login
az acr login --name taggeracr35316
docker push taggeracr35316.azurecr.io/azure-tagger-api:0
```

### Or build+push in ACR

```powershell
az acr build -r taggeracr35316 -t azure-tagger-api:0 .
```

---

## Summary

Most issues came from:

1. **Bad image tags** (missing image name / invalid tag text)
2. **Go compiler version mismatch** (Docker base image too old vs `go.mod`)
3. **ACR authentication/permissions** (not logged in or missing AcrPush role)



I keep these three in mind and you’ll avoid 90% of container-to-ACR pain ;D
Also, need to implement other tests..






## Author

Thiago Scheffer
Junior Backend Developer -> | Go | Cloud | REST APIs

---
