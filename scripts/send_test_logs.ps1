# PowerShell script to send test logs to the Go API
# Usage: .\send_test_logs.ps1 [count]

param(
    [int]$Count = 10
)

$API_URL = "http://localhost:8080/api/logs"

Write-Host "Sending $Count test log entries..." -ForegroundColor Cyan

for ($i = 1; $i -le $Count; $i++) {
    $level = "info"
    if ($i % 5 -eq 0) {
        $level = "error"
    } elseif ($i % 3 -eq 0) {
        $level = "warn"
    }

    $service = "service-$((Get-Random) % 3 + 1)"
    
    $body = @{
        logs = @(
            @{
                service = $service
                level = $level
                message = "Test log message #$i - $level level"
                metadata = @{
                    test = $true
                    iteration = $i
                }
            }
        )
    } | ConvertTo-Json -Depth 10

    try {
        Invoke-RestMethod -Uri $API_URL -Method Post -Body $body -ContentType "application/json" | Out-Null
        Write-Host "Sent log $i/$Count" -ForegroundColor Green
    } catch {
        Write-Host "Failed to send log $i : $_" -ForegroundColor Red
    }

    Start-Sleep -Milliseconds 100
}

Write-Host "`nDone! Check the dashboard at http://localhost:5173" -ForegroundColor Cyan
