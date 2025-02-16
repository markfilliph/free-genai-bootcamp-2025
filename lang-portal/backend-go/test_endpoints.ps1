# Test script for language portal API endpoints
$baseUrl = "http://localhost:8080/api"

Write-Host "`nStarting API endpoint tests...`n" -ForegroundColor Cyan

# Test 1: Get all groups (no pagination)
Write-Host "Test 1: Get all groups (no pagination)" -ForegroundColor Green
$groups = Invoke-RestMethod -Method GET -Uri "$baseUrl/groups"
$groups | ConvertTo-Json

# Test 2: Get groups with pagination
Write-Host "`nTest 2: Get groups with pagination (offset=0, limit=2)" -ForegroundColor Green
$pagedGroups = Invoke-RestMethod -Method GET -Uri "$baseUrl/groups?offset=0&limit=2"
$pagedGroups | ConvertTo-Json

# Test 3: Create a new group
Write-Host "`nTest 3: Create a new group" -ForegroundColor Green
$body = @{
    name = "Test Group $(Get-Date -Format 'yyyyMMddHHmmss')"
} | ConvertTo-Json
$headers = @{
    "Content-Type" = "application/json"
}
try {
    $newGroup = Invoke-RestMethod -Method POST -Uri "$baseUrl/groups" -Headers $headers -Body $body
    Write-Host "Created group with ID: $($newGroup.id)" -ForegroundColor Green
    $newGroup | ConvertTo-Json

    # Test 4: Get the created group
    Write-Host "`nTest 4: Get group details" -ForegroundColor Green
    $group = Invoke-RestMethod -Method GET -Uri "$baseUrl/groups/$($newGroup.id)"
    $group | ConvertTo-Json

    # Test 5: Get words in the created group
    Write-Host "`nTest 5: Get words in group" -ForegroundColor Green
    $words = Invoke-RestMethod -Method GET -Uri "$baseUrl/groups/$($newGroup.id)/words"
    $words | ConvertTo-Json

    # Test 6: Get study sessions for the created group
    Write-Host "`nTest 6: Get study sessions for group" -ForegroundColor Green
    $sessions = Invoke-RestMethod -Method GET -Uri "$baseUrl/groups/$($newGroup.id)/study-sessions"
    $sessions | ConvertTo-Json
}
catch {
    Write-Host "Error: $_" -ForegroundColor Red
    Write-Host $_.ErrorDetails.Message -ForegroundColor Red
}

Write-Host "`nAPI endpoint tests completed.`n" -ForegroundColor Cyan
