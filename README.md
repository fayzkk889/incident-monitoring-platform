# 🔥 Incident Monitoring & Prediction Platform

> **A smart monitoring system that watches your applications, detects problems automatically, and uses AI to explain what went wrong.**

---

## 📖 What Is This?

Think of this platform as a **smart security guard** for your applications:

- 👁️ **Watches** all your application logs
- 🚨 **Detects** problems automatically when something goes wrong
- 🤖 **Explains** what happened using AI
- 📊 **Shows** everything in an easy-to-use dashboard
- 💬 **Alerts** you on Slack when incidents occur

Perfect for **DevOps teams**, **startups**, and **SaaS companies** who want intelligent log analysis without the complexity.

---

## ✨ Key Features

| Feature | What It Does |
|---------|--------------|
| **📥 Log Ingestion** | Collects logs from all your services in one place |
| **🔍 Anomaly Detection** | Automatically finds unusual patterns and problems |
| **🤖 AI Analysis** | Uses OpenAI to explain incidents and find root causes |
| **📊 Real-time Dashboard** | Beautiful web interface to view everything |
| **💬 Slack Alerts** | Get notified instantly when problems are detected |
| **⚡ High Performance** | Built with Go for fast log processing |

---

## 🏗️ How It Works

```
Your Apps → Send Logs → Go API → Database
                              ↓
                    Python ML Service → Detects Problems
                              ↓
                    Creates Incidents → Sends Slack Alert
                              ↓
                    Dashboard Shows Everything
```

### The 4 Main Components

1. **Go API** (Port 8080) - Receives and stores logs
2. **Python ML Service** (Port 8000) - Detects problems and uses AI
3. **PostgreSQL Database** (Port 5432) - Stores everything
4. **React Dashboard** (Port 5173) - Shows you what's happening

---

## 🚀 Quick Start Guide

### Step 1: Prerequisites

Make sure you have:
- ✅ **Docker Desktop** installed and running
- ✅ **OpenAI API Key** (optional, but recommended for AI features)
  - Get one at: https://platform.openai.com/api-keys
- ✅ **Slack Webhook URL** (optional, for notifications)
  - Create one at: https://api.slack.com/messaging/webhooks

### Step 2: Setup

1. **Navigate to the project folder:**
   ```powershell
   cd Incident_Monitoring_Project
   ```

2. **Configure environment variables:**
   
   Edit the `.env` file and add your keys:
   ```env
   # Add your OpenAI API key (for AI features)
   OPENAI_API_KEY=sk-your-key-here
   
   # Add your Slack webhook URL (for notifications)
   ALERT_WEBHOOK_URL=https://hooks.slack.com/services/YOUR/WEBHOOK/URL
   ```

3. **Start everything:**
   ```powershell
    
   ```
   
   ⏳ **Wait 30-60 seconds** for all services to start. You'll see:
   - ✅ `Go API listening on :8080`
   - ✅ `Application startup complete` (Python ML)
   - ✅ Frontend ready on port 5173

4. **Open the dashboard:**
   - Open your browser: **http://localhost:5173**

🎉 **That's it!** Your platform is now running.

---

## 📝 How to Use

### 1️⃣ Send Logs to the System

Your applications can send logs like this:

**Using PowerShell:**
```powershell
$body = @{
    logs = @(
        @{
            service = "api-server"
            level = "error"
            message = "Database connection timeout"
            metadata = @{ user_id = 12345 }
        }
    )
} | ConvertTo-Json -Depth 10

Invoke-RestMethod -Uri "http://localhost:8080/api/logs" -Method POST -Body $body -ContentType "application/json"
```

**Or use the test script:**
```powershell
.\scripts\send_test_logs.ps1 -Count 20
```

### 2️⃣ Detect Problems

The system automatically detects problems, or you can trigger it manually:

```powershell
Invoke-RestMethod -Uri "http://localhost:8000/detect_anomalies" -Method POST
```

**What happens:**
- ✅ System analyzes recent logs
- ✅ Creates incidents if problems found
- ✅ Sends Slack notification (if configured)
- ✅ Shows incidents in dashboard

### 3️⃣ Get AI Analysis

When an incident is created, get an AI-powered explanation:

```powershell
# Get summary for incident #1
Invoke-RestMethod -Uri "http://localhost:8080/api/summary/1"
```

**The AI will tell you:**
- 📝 What happened (summary)
- 🔍 Why it happened (root cause)
- 💡 What to do about it (insights)

### 4️⃣ View Everything in Dashboard

- Open **http://localhost:5173**
- See all incidents
- Click "Generate AI Analysis" on any incident
- Monitor system health

---

## 🧪 Testing

**Quick test everything:**
```powershell
.\test-api.ps1
```

This automated script will:
1. ✅ Check if services are running
2. ✅ Send test logs
3. ✅ Detect anomalies
4. ✅ Create incidents
5. ✅ Get AI summaries

**For detailed testing instructions, see [TESTING-GUIDE.md](./TESTING-GUIDE.md)**

---

## 🔌 API Endpoints

### Go API (http://localhost:8080)

| Endpoint | Method | What It Does |
|----------|--------|--------------|
| `/api/logs` | POST | Send logs to the system |
| `/api/health` | GET | Check if API is working |
| `/api/incidents` | GET | Get list of all incidents |
| `/api/summary/:id` | GET | Get AI analysis for an incident |

### Python ML API (http://localhost:8000)

| Endpoint | Method | What It Does |
|----------|--------|--------------|
| `/health` | GET | Check if ML service is working |
| `/detect_anomalies` | POST | Analyze logs and find problems |
| `/analyze_incident` | POST | Get AI analysis (internal use) |

---

## ⚙️ Configuration

### Environment Variables

Create a `.env` file in the project root:

```env
# Go API Configuration
DATABASE_URL=postgres://incident:incidentpassword@postgres:5432/incidentdb?sslmode=disable
ML_SERVICE_URL=http://python-ml:8000
ALERT_WEBHOOK_URL=https://hooks.slack.com/services/YOUR/WEBHOOK/URL

# Python ML Service Configuration
OPENAI_API_KEY=sk-your-openai-key-here

# Frontend Configuration
VITE_API_BASE=http://localhost:8080
```

### What Each Variable Does

- **`OPENAI_API_KEY`** - Required for AI features (get from OpenAI)
- **`ALERT_WEBHOOK_URL`** - Optional, for Slack notifications
- **`DATABASE_URL`** - Usually don't need to change this
- **`ML_SERVICE_URL`** - Usually don't need to change this

---

## 🎯 Common Use Cases

### Use Case 1: Monitor Your Application
- Send logs from your app to the Go API
- System automatically detects problems
- Get Slack alerts when issues occur
- View everything in the dashboard

### Use Case 2: Debug Issues
- When something breaks, check the dashboard
- Click "Generate AI Analysis" on the incident
- AI explains what happened and why
- Fix the issue faster

### Use Case 3: Track System Health
- Monitor error rates across services
- See which services have the most problems
- Get alerts before users notice issues

---

## 🛠️ Troubleshooting

### Problem: Services won't start

**Solution:**
```powershell
# Stop everything
docker compose down

# Remove old containers
docker compose down -v

# Start fresh
docker compose up --build
```

### Problem: Port already in use

**Solution:**
- Check if ports 8080, 8000, 5173, or 5432 are already used
- Stop the service using those ports
- Or change ports in `docker-compose.yml`

### Problem: Can't connect to database

**Solution:**
```powershell
# Check if PostgreSQL is running
docker compose ps postgres

# Check logs
docker compose logs postgres
```

### Problem: AI analysis not working

**Solution:**
- Make sure `OPENAI_API_KEY` is set in `.env`
- Check if you have credits in your OpenAI account
- Check logs: `docker compose logs python-ml`

### Problem: Slack notifications not working

**Solution:**
- Verify `ALERT_WEBHOOK_URL` is correct in `.env`
- Test your webhook URL manually
- Check logs: `docker compose logs python-ml`

---

## 📊 How It Works (Technical)

1. **Log Ingestion**
   - Your apps send logs → Go API receives them → Stores in PostgreSQL

2. **Anomaly Detection**
   - Python ML service analyzes logs using:
     - Error rate spike detection
     - Service-specific patterns
     - Volume anomaly detection

3. **Incident Creation**
   - When problems found → Creates incident record
   - Sends Slack notification (if configured)

4. **AI Analysis**
   - On-demand → Uses OpenAI to:
     - Summarize what happened
     - Identify root causes
     - Provide actionable insights

---

## 🛠️ Development Mode

### Running Without Docker

**Go API:**
```powershell
cd go-api
go mod download
go run cmd/server/main.go
```

**Python ML:**
```powershell
cd python-ml
pip install -r requirements.txt
uvicorn main:app --reload
```

**Frontend:**
```powershell
cd frontend
npm install
npm run dev
```

---

## 📦 Tech Stack

- **Backend**: Go (Echo framework) + PostgreSQL
- **ML Service**: Python (FastAPI) + scikit-learn + OpenAI
- **Frontend**: React + Vite
- **Infrastructure**: Docker + Docker Compose

---

## 🤝 Need Help?

- 📖 Check [TESTING-GUIDE.md](./TESTING-GUIDE.md) for detailed testing instructions
- 🐛 Check service logs: `docker compose logs [service-name]`
- 🔍 View all running services: `docker compose ps`

---

## 📄 License

MIT License - feel free to use this project for your own needs!

---

**Made with ❤️ for developers who want smarter monitoring**