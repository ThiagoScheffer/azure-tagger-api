


$SUB="<e9582386-74c7-4f59-89cb-6f09aaa7db41>"

az ad sp create-for-rbac `
  --name "azure-tagger-api-sp" `
  --role "Contributor" `
  --scopes "/subscriptions/$SUB"

  $RG="rg-azure-tagger-dev"
$LOC="eastus"
$APP="azure-tagger-api"

az containerapp up `
  --name $APP `
  --resource-group $RG `
  --location $LOC `
  --source . `
  --ingress external `
  --target-port 8080



  az containerapp secret set `
  --name $APP `
  --resource-group $RG `
  --secrets `
    azure-client-secret="<CLIENT_SECRET>"



    $create = @{
  name    = "test-vm"
  azureId = "/subscriptions/<SUB>/resourceGroups/<RG>/providers/Microsoft.Compute/virtualMachines/<VM>"
  tags    = @{ env="dev" }
} | ConvertTo-Json -Depth 5

$r = Invoke-RestMethod "$BASE/v1/resources" -Method Post -ContentType "application/json" -Body $create
$r