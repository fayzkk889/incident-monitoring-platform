# 🧪 Testing Guide - Step by Step

> **Complete guide to test your Incident Monitoring Platform**

This guide will walk you through testing every feature of the platform using PowerShell.

---

## 📋 Table of Contents

1. [Before You Start](#before-you-start)
2. [Quick Test (Automated)](#quick-test-automated)
3. [Manual Testing](#manual-testing)
4. [Complete Test Flow](#complete-test-flow)
5. [Testing Slack Notifications](#testing-slack-notifications)
6. [Troubleshooting](#troubleshooting)

---

## Before You Start

### ✅ Prerequisites Checklist

- [ ] Docker Desktop is running
- [ ] All services are started (`docker compose up`)
- [ ] You can see the dashboard at http://localhost:5173
- [ ] You're using PowerShell (Windows)

### 🚀 Start Services (If Not Running)

```powershell
# Navigate to project folder
cd Incident_Monitoring_Project

# Start all services
docker compose up --build
```

**Wait until you see:**
- ✅ `Go API listening on :8080`
- ✅ `Application startup complete` (Python ML)
- ✅ Frontend accessible

**This usually takes 30-60 seconds.**

---

## Quick Test (Automated)

### 🎯 Run the Automated Test Script

The easiest way to test everything:

```powershell
.\test-api.ps1
```

**What this does:**
1. ✅ Checks if services are healthy
2. ✅ Sends test logs
3. ✅ Triggers anomaly detection
4. ✅ Creates incidents
5. ✅ Gets AI summaries
6. ✅ Shows you results

**Expected output:**
```
========================================
  Testing Incident Monitoring Platform
========================================

--- Step 1: Health Checks ---
Testing: Go API Health Check
  ✓ Success

Testing: Python ML API Health Check
  ✓ Success

--- Step 2: Ingest Test Logs ---
Testing: Ingest Multiple Logs
  ✓ Success

... (more steps)

✓ Testing complete! Open http://localhost:5173 to view the dashboard
```

**⏱️ Takes about 1-2 minutes to complete**

---

## Manual Testing

Follow these steps to test each feature individually.

---

### Step 1: Health Checks ✅

**Purpose:** Verify all services are running correctly

#### Test Go API Health

```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/health" -Method GET
```

**Expected Response:**
```json
{
  "status": "ok",
  "checks": {
    "db": true
  }
}
```

**✅ Success if:** You see `"status": "ok"` and `"db": true`

#### Test Python ML API Health

```powershell
Invoke-RestMethod -Uri "http://localhost:8000/health" -Method GET
```

**Expected Response:**
```json
{
  "status": "ok"
}
```

**✅ Success if:** You see `"status": "ok"`

---

### Step 2: Send Logs 📥

**Purpose:** Test log ingestion

#### Option A: Send Single Log Entry

```powershell
$body = @{
    logs = @(
        @{
            service = "api-server"
            level = "error"
            message = "Database connection timeout"
            metadata = @{
                user_id = 12345
                endpoint = "/api/users"
            }
        }
    )
} | ConvertTo-Json -Depth 10

Invoke-RestMethod -Uri "http://localhost:8080/api/logs" -Method POST -Body $body -ContentType "application/json"
```

**Expected Response:**
```json
{
  "status": "accepted",
  "count": 1
}
```

**✅ Success if:** You see `"status": "accepted"` and `"count": 1`

#### Option B: Send Multiple Logs

```powershell
$body = @{
    logs = @(
        @{
            service = "api-server"
            level = "error"
            message = "Database connection timeout"
            metadata = @{ user_id = 12345 }
        },
        @{
            service = "auth-service"
            level = "warn"
            message = "Rate limit approaching"
            metadata = @{ ip = "192.168.1.1" }
        },
        @{
            service = "payment-service"
            level = "info"
            message = "Payment processed successfully"
            metadata = @{ transaction_id = "txn_123" }
        }
    )
} | ConvertTo-Json -Depth 10

Invoke-RestMethod -Uri "http://localhost:8080/api/logs" -Method POST -Body $body -ContentType "application/json"
```

**Expected Response:**
```json
{
  "status": "accepted",
  "count": 3
}
```

#### Option C: Use Test Script (Easiest)

```powershell
# Send 20 test logs
.\scripts\send_test_logs.ps1 -Count 20
```

**✅ Success if:** You see "Sent log X/20" messages

---

### Step 3: Detect Anomalies 🔍

**Purpose:** Trigger the system to find problems in logs

```powershell
Invoke-RestMethod -Uri "http://localhost:8000/detect_anomalies" -Method POST
```

**Expected Response (if anomalies found):**
```json
{
  "anomalies_detected": 2,
  "incidents_created": [1, 2],
  "anomalies": [
    {
      "type": "error_rate_spike",
      "severity": "high",
      "description": "Error rate spike detected..."
    }
  ]
}
```

**Expected Response (if no anomalies):**
```json
{
  "anomalies_detected": 0,
  "message": "No anomalies detected"
}
```

**✅ Success if:** 
- You get a response (even if no anomalies)
- If anomalies found, incidents are created
- Check Slack for notifications (if configured)

**💡 Tip:** Send some error logs first (Step 2) to ensure anomalies are detected

---

### Step 4: List Incidents 📋

**Purpose:** See all detected incidents

```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/incidents" -Method GET
```

**Pretty print the response:**
```powershell
$incidents = Invoke-RestMethod -Uri "http://localhost:8080/api/incidents" -Method GET
$incidents | ConvertTo-Json -Depth 5
```

**Expected Response:**
```json
[
  {
    "id": 1,
    "created_at": "2026-02-13T10:30:00Z",
    "status": "open",
    "severity": "high",
    "description": "Error rate spike detected...",
    "summary": null,
    "root_cause": null,
    "resolved_at": null
  }
]
```

**✅ Success if:** You see a list of incidents (may be empty if none created yet)

**View specific fields:**
```powershell
$incidents = Invoke-RestMethod -Uri "http://localhost:8080/api/incidents" -Method GET
$incidents | ForEach-Object { 
    Write-Host "Incident #$($_.id): $($_.description) - Severity: $($_.severity)" 
}
```

---

### Step 5: Get AI Summary 🤖

**Purpose:** Get AI-powered explanation of an incident

**⚠️ Note:** This requires OpenAI API key and may take 10-30 seconds

#### Get Summary for Specific Incident

```powershell
# Replace 1 with your incident ID
Invoke-RestMethod -Uri "http://localhost:8080/api/summary/1" -Method GET
```

**Expected Response:**
```json
{
  "id": 1,
  "created_at": "2026-02-13T10:30:00Z",
  "status": "open",
  "severity": "high",
  "description": "Error rate spike detected...",
  "summary": "The system experienced a significant increase in error rates...",
  "root_cause": "The root cause appears to be database connectivity issues...",
  "resolved_at": null
}
```

**✅ Success if:** You see `summary` and `root_cause` fields populated

#### Automatically Get First Incident Summary

```powershell
$incidents = Invoke-RestMethod -Uri "http://localhost:8080/api/incidents" -Method GET
if ($incidents.Count -gt 0) {
    $firstId = $incidents[0].id
    Write-Host "Getting AI summary for incident #$firstId..." -ForegroundColor Yellow
    Write-Host "This may take 10-30 seconds..." -ForegroundColor Yellow
    Invoke-RestMethod -Uri "http://localhost:8080/api/summary/$firstId" -Method GET | ConvertTo-Json -Depth 5
} else {
    Write-Host "No incidents found. Create some incidents first!" -ForegroundColor Red
}
```

**✅ Success if:** AI provides summary and root cause analysis

---

## Complete Test Flow

Run this complete sequence to test everything:

```powershell
Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "  Complete Test Flow" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

# Step 1: Health Checks
Write-Host "Step 1: Checking service health..." -ForegroundColor Yellow
try {
    $goHealth = Invoke-RestMethod -Uri "http://localhost:8080/api/health"
    $mlHealth = Invoke-RestMethod -Uri "http://localhost:8000/health"
    Write-Host "✓ All services healthy" -ForegroundColor Green
} catch {
    Write-Host "✗ Services not ready. Start them first!" -ForegroundColor Red
    exit
}

# Step 2: Send Test Logs
Write-Host "`nStep 2: Sending test logs..." -ForegroundColor Yellow
.\scripts\send_test_logs.ps1 -Count 25
Write-Host "✓ Test logs sent" -ForegroundColor Green

# Step 3: Detect Anomalies
Write-Host "`nStep 3: Detecting anomalies..." -ForegroundColor Yellow
$anomalyResult = Invoke-RestMethod -Uri "http://localhost:8000/detect_anomalies" -Method POST
Write-Host "Anomalies detected: $($anomalyResult.anomalies_detected)" -ForegroundColor Cyan
Write-Host "Incidents created: $($anomalyResult.incidents_created.Count)" -ForegroundColor Cyan
Write-Host "✓ Anomaly detection complete" -ForegroundColor Green

# Step 4: List Incidents
Write-Host "`nStep 4: Listing incidents..." -ForegroundColor Yellow
$incidents = Invoke-RestMethod -Uri "http://localhost:8080/api/incidents" -Method GET
Write-Host "Total incidents: $($incidents.Count)" -ForegroundColor Cyan
if ($incidents.Count -gt 0) {
    $incidents | ForEach-Object {
        Write-Host "  - Incident #$($_.id): $($_.description) (Severity: $($_.severity))" -ForegroundColor Gray
    }
    Write-Host "✓ Incidents listed" -ForegroundColor Green
} else {
    Write-Host "⚠ No incidents found" -ForegroundColor Yellow
}

# Step 5: Get AI Summary (if incidents exist)
if ($incidents.Count -gt 0) {
    Write-Host "`nStep 5: Getting AI summary..." -ForegroundColor Yellow
    Write-Host "This may take 10-30 seconds..." -ForegroundColor Yellow
    $firstId = $incidents[0].id
    try {
        $summary = Invoke-RestMethod -Uri "http://localhost:8080/api/summary/$firstId" -Method GET
        Write-Host "`nAI Summary:" -ForegroundColor Cyan
        Write-Host "  Summary: $($summary.summary)" -ForegroundColor Gray
        Write-Host "  Root Cause: $($summary.root_cause)" -ForegroundColor Gray
        Write-Host "✓ AI analysis complete" -ForegroundColor Green
    } catch {
        Write-Host "⚠ AI analysis failed: $($_.Exception.Message)" -ForegroundColor Yellow
        Write-Host "  Make sure OPENAI_API_KEY is set in .env" -ForegroundColor Yellow
    }
} else {
    Write-Host "`nStep 5: Skipped (no incidents to analyze)" -ForegroundColor Yellow
}

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "  Test Complete!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "`nNext Steps:" -ForegroundColor Yellow
Write-Host "  1. Open dashboard: http://localhost:5173" -ForegroundColor White
Write-Host "  2. View incidents in the UI" -ForegroundColor White
Write-Host "  3. Check Slack for notifications (if configured)" -ForegroundColor White
Write-Host "`n"
```

**⏱️ Total time: 2-3 minutes**

---

## Testing Slack Notifications

### Prerequisites

- ✅ `ALERT_WEBHOOK_URL` is set in `.env` file
- ✅ Slack webhook URL is valid

### Test Steps

1. **Send logs with errors:**
   ```powershell
   # Send multiple error logs
   for ($i = 1; $i -le 10; $i++) {
       $body = @{
           logs = @(
               @{
                   service = "api-server"
                   level = "error"
                   message = "Test error #$i"
                   metadata = @{ test = $true }
               }
           )
       } | ConvertTo-Json -Depth 10
       
       Invoke-RestMethod -Uri "http://localhost:8080/api/logs" -Method POST -Body $body -ContentType "application/json" | Out-Null
       Start-Sleep -Milliseconds 200
   }
   ```

2. **Trigger anomaly detection:**
   ```powershell
   Invoke-RestMethod -Uri "http://localhost:8000/detect_anomalies" -Method POST
   ```

3. **Check Slack:**
   - ✅ You should receive a notification in your Slack channel
   - ✅ Message includes incident ID, severity, and description
   - ✅ Includes links to dashboard and AI analysis

**✅ Success if:** You see a Slack notification with incident details

---

## Testing Error Scenarios

### Test 1: Invalid Log Payload

```powershell
# Try to send empty logs array
$invalidBody = @{ logs = @() } | ConvertTo-Json
try {
    Invoke-RestMethod -Uri "http://localhost:8080/api/logs" -Method POST -Body $invalidBody -ContentType "application/json"
    Write-Host "✗ Should have failed!" -ForegroundColor Red
} catch {
    Write-Host "✓ Correctly rejected invalid payload" -ForegroundColor Green
    Write-Host "  Error: $($_.Exception.Message)" -ForegroundColor Gray
}
```

**Expected:** Error response with `"error": "no logs provided"`

### Test 2: Invalid Incident ID

```powershell
# Try to get summary for non-existent incident
try {
    Invoke-RestMethod -Uri "http://localhost:8080/api/summary/99999" -Method GET
    Write-Host "✗ Should have failed!" -ForegroundColor Red
} catch {
    Write-Host "✓ Correctly handled invalid incident ID" -ForegroundColor Green
    Write-Host "  Error: $($_.Exception.Message)" -ForegroundColor Gray
}
```

**Expected:** Error response with `"error": "incident not found"`

---

## Viewing Responses

### Pretty Print JSON

```powershell
$response = Invoke-RestMethod -Uri "http://localhost:8080/api/incidents"
$response | ConvertTo-Json -Depth 5
```

### Save Response to File

```powershell
$response = Invoke-RestMethod -Uri "http://localhost:8080/api/incidents"
$response | ConvertTo-Json -Depth 5 | Out-File -FilePath "incidents.json"
```

### View Specific Fields

```powershell
$incidents = Invoke-RestMethod -Uri "http://localhost:8080/api/incidents"
$incidents | ForEach-Object { 
    Write-Host "Incident #$($_.id)" -ForegroundColor Cyan
    Write-Host "  Description: $($_.description)" -ForegroundColor Gray
    Write-Host "  Severity: $($_.severity)" -ForegroundColor Gray
    Write-Host "  Status: $($_.status)" -ForegroundColor Gray
    Write-Host ""
}
```

---

## Troubleshooting

### ❌ Problem: Service Not Responding

**Symptoms:**
- Connection refused errors
- Timeout errors

**Solutions:**
```powershell
# Check if services are running
docker compose ps

# Check logs for errors
docker compose logs go-api
docker compose logs python-ml

# Restart services
docker compose restart
```

### ❌ Problem: Port Already in Use

**Symptoms:**
- `bind: address already in use`
- Can't start services

**Solutions:**
```powershell
# Find what's using the port
netstat -ano | findstr :8080

# Stop the conflicting service
# Or change ports in docker-compose.yml
```

### ❌ Problem: Database Connection Failed

**Symptoms:**
- `failed to connect to database`
- Health check shows `"db": false`

**Solutions:**
```powershell
# Check if PostgreSQL is running
docker compose ps postgres

# Check PostgreSQL logs
docker compose logs postgres

# Restart PostgreSQL
docker compose restart postgres
```

### ❌ Problem: AI Analysis Not Working

**Symptoms:**
- Summary returns empty or error
- Takes too long

**Solutions:**
```powershell
# Check if OpenAI key is set
# Look in .env file for OPENAI_API_KEY

# Check Python ML logs
docker compose logs python-ml

# Verify OpenAI key is valid
# Make sure you have credits in OpenAI account
```

### ❌ Problem: Slack Notifications Not Working

**Symptoms:**
- No notifications in Slack
- Incidents created but no alerts

**Solutions:**
```powershell
# Check if webhook URL is set in .env
# Verify ALERT_WEBHOOK_URL is correct

# Check Python ML logs for errors
docker compose logs python-ml | Select-String -Pattern "Slack"

# Test webhook URL manually
$webhookUrl = "YOUR_WEBHOOK_URL"
$testPayload = @{ text = "Test message" } | ConvertTo-Json
Invoke-RestMethod -Uri $webhookUrl -Method POST -Body $testPayload -ContentType "application/json"
```

### ❌ Problem: No Anomalies Detected

**Symptoms:**
- `"anomalies_detected": 0` even with error logs

**Solutions:**
```powershell
# Make sure you have enough logs (at least 10)
# Send more error logs
.\scripts\send_test_logs.ps1 -Count 30

# Wait a bit, then try again
Start-Sleep -Seconds 5
Invoke-RestMethod -Uri "http://localhost:8000/detect_anomalies" -Method POST
```

---

## Frontend Testing

### Access the Dashboard

1. Open browser: **http://localhost:5173**
2. You should see:
   - System health status
   - List of incidents
   - Service statistics

### Test Dashboard Features

1. **View Incidents:**
   - ✅ See all incidents listed
   - ✅ See severity badges (high/medium/low)
   - ✅ See status (open/resolved)

2. **Generate AI Analysis:**
   - ✅ Click "Generate AI Analysis" button
   - ✅ Wait 10-30 seconds
   - ✅ See summary and root cause appear

3. **Monitor Health:**
   - ✅ See system health status
   - ✅ See number of open incidents
   - ✅ See high severity count

---

## Performance Testing

### Send Many Logs Quickly

```powershell
Write-Host "Sending 100 logs..." -ForegroundColor Yellow
$startTime = Get-Date

for ($i = 1; $i -le 100; $i++) {
    $body = @{
        logs = @(
            @{
                service = "service-$($i % 5)"
                level = @("info", "warn", "error")[$i % 3]
                message = "Test log #$i"
                metadata = @{ iteration = $i }
            }
        )
    } | ConvertTo-Json -Depth 10
    
    try {
        Invoke-RestMethod -Uri "http://localhost:8080/api/logs" -Method POST -Body $body -ContentType "application/json" | Out-Null
        if ($i % 10 -eq 0) { 
            Write-Host "  Sent $i logs..." -ForegroundColor Gray 
        }
    } catch {
        Write-Host "  Failed at log $i" -ForegroundColor Red
    }
}

$endTime = Get-Date
$duration = ($endTime - $startTime).TotalSeconds
Write-Host "✓ Sent 100 logs in $duration seconds" -ForegroundColor Green
```

**Expected:** Should complete in under 10 seconds

---

## Quick Reference

### Most Common Commands

```powershell
# Health check
Invoke-RestMethod -Uri "http://localhost:8080/api/health"

# Send test logs
.\scripts\send_test_logs.ps1 -Count 20

# Detect anomalies
Invoke-RestMethod -Uri "http://localhost:8000/detect_anomalies" -Method POST

# List incidents
Invoke-RestMethod -Uri "http://localhost:8080/api/incidents"

# Get AI summary
Invoke-RestMethod -Uri "http://localhost:8080/api/summary/1"

# Run full test
.\test-api.ps1
```

---

## Need More Help?

- 📖 Check [README.md](./README.md) for general information
- 🐛 Check service logs: `docker compose logs [service-name]`
- 🔍 View running services: `docker compose ps`
- 💬 Check if Slack webhook is working by testing it manually

---

**Happy Testing! 🎉**
