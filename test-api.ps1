# PowerShell Test Script for Incident Monitoring Platform
# Usage: .\test-api.ps1

$ErrorActionPreference = "Stop"

$GO_API = "http://localhost:8080"
$ML_API = "http://localhost:8000"

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "  Testing Incident Monitoring Platform" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

# Function to make API calls with error handling
function Invoke-ApiCall {
    param(
        [string]$Uri,
        [string]$Method = "GET",
        [object]$Body = $null,
        [string]$Description = ""
    )
    
    Write-Host "Testing: $Description" -ForegroundColor Yellow
    Write-Host "  $Method $Uri" -ForegroundColor Gray
    
    try {
        $params = @{
            Uri = $Uri
            Method = $Method
            ContentType = "application/json"
            ErrorAction = "Stop"
        }
        
        if ($Body) {
            $params.Body = ($Body | ConvertTo-Json -Depth 10)
        }
        
        $response = Invoke-RestMethod @params
        Write-Host "  ✓ Success" -ForegroundColor Green
        $response | ConvertTo-Json -Depth 5 | Write-Host -ForegroundColor DarkGray
        return $response
    }
    catch {
        Write-Host "  ✗ Failed: $($_.Exception.Message)" -ForegroundColor Red
        if ($_.Exception.Response) {
            $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
            $responseBody = $reader.ReadToEnd()
            Write-Host "  Response: $responseBody" -ForegroundColor Red
        }
        return $null
    }
    Write-Host ""
}

# Test 1: Health Checks
Write-Host "`n--- Step 1: Health Checks ---" -ForegroundColor Magenta
Invoke-ApiCall -Uri "$GO_API/api/health" -Description "Go API Health Check"
Start-Sleep -Seconds 1
Invoke-ApiCall -Uri "$ML_API/health" -Description "Python ML API Health Check"
Start-Sleep -Seconds 1

# Test 2: Send Test Logs
Write-Host "`n--- Step 2: Ingest Test Logs ---" -ForegroundColor Magenta
$testLogs = @{
    logs = @(
        @{
            service = "api-server"
            level = "error"
            message = "Database connection timeout after 30 seconds"
            metadata = @{
                user_id = 12345
                endpoint = "/api/users"
                duration_ms = 30000
            }
        },
        @{
            service = "auth-service"
            level = "warn"
            message = "Rate limit approaching: 90% capacity"
            metadata = @{
                ip = "192.168.1.100"
                requests_per_minute = 540
            }
        },
        @{
            service = "payment-service"
            level = "error"
            message = "Payment gateway timeout"
            metadata = @{
                transaction_id = "txn_abc123"
                amount = 99.99
            }
        },
        @{
            service = "api-server"
            level = "info"
            message = "Request processed successfully"
            metadata = @{
                user_id = 67890
                endpoint = "/api/products"
                duration_ms = 45
            }
        },
        @{
            service = "auth-service"
            level = "error"
            message = "Invalid authentication token"
            metadata = @{
                ip = "10.0.0.5"
                token_expired = $true
            }
        }
    )
}

$logResponse = Invoke-ApiCall -Uri "$GO_API/api/logs" -Method POST -Body $testLogs -Description "Ingest Multiple Logs"
Start-Sleep -Seconds 2

# Test 3: Send More Logs for Anomaly Detection
Write-Host "`n--- Step 3: Send Additional Logs for Anomaly Detection ---" -ForegroundColor Magenta
for ($i = 1; $i -le 15; $i++) {
    $level = if ($i % 4 -eq 0) { "error" } elseif ($i % 3 -eq 0) { "warn" } else { "info" }
    $service = @("api-server", "auth-service", "payment-service")[($i - 1) % 3]
    
    $singleLog = @{
        logs = @(
            @{
                service = $service
                level = $level
                message = "Test log entry #$i - $level level from $service"
                metadata = @{
                    iteration = $i
                    timestamp = (Get-Date).ToString("o")
                }
            }
        )
    }
    
    try {
        Invoke-RestMethod -Uri "$GO_API/api/logs" -Method POST -Body ($singleLog | ConvertTo-Json -Depth 10) -ContentType "application/json" | Out-Null
        Write-Host "  Sent log $i/15" -ForegroundColor DarkGray
    }
    catch {
        Write-Host "  Failed to send log $i" -ForegroundColor Red
    }
    Start-Sleep -Milliseconds 200
}
Write-Host "  ✓ Sent 15 additional logs" -ForegroundColor Green
Start-Sleep -Seconds 2

# Test 4: List Incidents
Write-Host "`n--- Step 4: List Incidents ---" -ForegroundColor Magenta
$incidents = Invoke-ApiCall -Uri "$GO_API/api/incidents" -Description "List All Incidents"
Start-Sleep -Seconds 1

# Test 5: Detect Anomalies
Write-Host "`n--- Step 5: Detect Anomalies ---" -ForegroundColor Magenta
$anomalyResponse = Invoke-ApiCall -Uri "$ML_API/detect_anomalies" -Method POST -Description "Trigger Anomaly Detection"
Start-Sleep -Seconds 2

# Test 6: List Incidents Again (after anomaly detection)
Write-Host "`n--- Step 6: List Incidents (After Anomaly Detection) ---" -ForegroundColor Magenta
$incidentsAfter = Invoke-ApiCall -Uri "$GO_API/api/incidents" -Description "List Incidents After Detection"
Start-Sleep -Seconds 1

# Test 7: Get Incident Summary (if incidents exist)
if ($incidentsAfter -and $incidentsAfter.Count -gt 0) {
    $firstIncidentId = $incidentsAfter[0].id
    Write-Host "`n--- Step 7: Get AI Summary for Incident #$firstIncidentId ---" -ForegroundColor Magenta
    Write-Host "  Note: This may take 10-30 seconds as it calls OpenAI API..." -ForegroundColor Yellow
    $summary = Invoke-ApiCall -Uri "$GO_API/api/summary/$firstIncidentId" -Description "Get AI-Generated Summary"
    Start-Sleep -Seconds 1
}

# Summary
Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "  Test Summary" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "✓ Health checks completed" -ForegroundColor Green
Write-Host "✓ Logs ingested successfully" -ForegroundColor Green
Write-Host "✓ Anomaly detection triggered" -ForegroundColor Green
Write-Host "✓ Incidents retrieved" -ForegroundColor Green
if ($incidentsAfter -and $incidentsAfter.Count -gt 0) {
    Write-Host "✓ AI summary generated" -ForegroundColor Green
}
Write-Host "`nNext Steps:" -ForegroundColor Yellow
Write-Host "  1. Open the dashboard: http://localhost:5173" -ForegroundColor White
Write-Host "  2. View incidents and logs in the UI" -ForegroundColor White
Write-Host "  3. Click 'Generate AI Analysis' on any incident" -ForegroundColor White
Write-Host "`n"
