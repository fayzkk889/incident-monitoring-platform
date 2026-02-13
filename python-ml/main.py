import os
from contextlib import asynccontextmanager
from datetime import datetime, timedelta
from typing import Optional

import httpx
import numpy as np
from dotenv import load_dotenv
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from sqlalchemy import create_engine, text
from sqlalchemy.engine import Engine

from anomaly import detect_anomalies
from llm_service import summarize_incident

load_dotenv()

DATABASE_URL = os.getenv(
    "DATABASE_URL", "postgresql+psycopg2://incident:incidentpassword@localhost:5432/incidentdb"
)
OPENAI_API_KEY = os.getenv("OPENAI_API_KEY", "")
ALERT_WEBHOOK_URL = os.getenv("ALERT_WEBHOOK_URL", "")

engine: Optional[Engine] = None


@asynccontextmanager
async def lifespan(app: FastAPI):
    global engine
    engine = create_engine(DATABASE_URL, pool_pre_ping=True)
    yield
    if engine:
        engine.dispose()


app = FastAPI(lifespan=lifespan)


async def send_slack_alert(incident_id: int, severity: str, description: str):
    """Send alert to Slack webhook when an incident is created."""
    if not ALERT_WEBHOOK_URL:
        return
    
    # Determine color based on severity
    color_map = {
        "high": "#FF0000",  # Red
        "critical": "#8B0000",  # Dark red
        "medium": "#FFA500",  # Orange
        "low": "#FFFF00",  # Yellow
    }
    color = color_map.get(severity.lower(), "#808080")  # Gray default
    
    # Create Slack message payload
    payload = {
        "text": f"🚨 New Incident Detected: #{incident_id}",
        "blocks": [
            {
                "type": "header",
                "text": {
                    "type": "plain_text",
                    "text": f"🚨 Incident #{incident_id} Detected",
                    "emoji": True
                }
            },
            {
                "type": "section",
                "fields": [
                    {
                        "type": "mrkdwn",
                        "text": f"*Severity:*\n{severity.upper()}"
                    },
                    {
                        "type": "mrkdwn",
                        "text": f"*Status:*\nOpen"
                    }
                ]
            },
            {
                "type": "section",
                "text": {
                    "type": "mrkdwn",
                    "text": f"*Description:*\n{description}"
                }
            },
            {
                "type": "context",
                "elements": [
                    {
                        "type": "mrkdwn",
                        "text": f"<http://localhost:5173|View in Dashboard> | <http://localhost:8080/api/summary/{incident_id}|Get AI Analysis>"
                    }
                ]
            }
        ],
        "attachments": [
            {
                "color": color,
                "footer": "Incident Monitoring Platform",
                "ts": int(datetime.now().timestamp())
            }
        ]
    }
    
    try:
        async with httpx.AsyncClient(timeout=5.0) as client:
            response = await client.post(ALERT_WEBHOOK_URL, json=payload)
            response.raise_for_status()
    except Exception as e:
        # Log error but don't fail the request
        print(f"Failed to send Slack alert: {e}")


class AnalyzeIncidentRequest(BaseModel):
    incident_id: int
    description: str


class AnalyzeIncidentResponse(BaseModel):
    summary: str
    root_cause: str


@app.get("/health")
async def health():
    return {"status": "ok"}


@app.post("/analyze_incident", response_model=AnalyzeIncidentResponse)
async def analyze_incident(req: AnalyzeIncidentRequest):
    """Analyze an incident using LLM to generate summary and root cause."""
    if not OPENAI_API_KEY:
        return AnalyzeIncidentResponse(
            summary="LLM analysis unavailable: OPENAI_API_KEY not configured",
            root_cause="Please configure OPENAI_API_KEY environment variable",
        )

    with engine.connect() as conn:
        result = conn.execute(
            text("""
                SELECT id, timestamp, service, level, message, metadata
                FROM logs
                WHERE timestamp >= NOW() - INTERVAL '1 hour'
                ORDER BY timestamp DESC
                LIMIT 100
            """)
        )
        logs = [
            {
                "id": row[0],
                "timestamp": row[1].isoformat() if row[1] else None,
                "service": row[2],
                "level": row[3],
                "message": row[4],
                "metadata": row[5] if row[5] else {},
            }
            for row in result
        ]

    summary, root_cause = await summarize_incident(
        incident_description=req.description,
        recent_logs=logs,
        api_key=OPENAI_API_KEY,
    )

    return AnalyzeIncidentResponse(summary=summary, root_cause=root_cause)


@app.post("/detect_anomalies")
async def detect_anomalies_endpoint():
    """Detect anomalies in recent logs and create incidents if found."""
    with engine.connect() as conn:
        result = conn.execute(
            text("""
                SELECT id, timestamp, service, level, message, metadata
                FROM logs
                WHERE timestamp >= NOW() - INTERVAL '1 hour'
                ORDER BY timestamp DESC
                LIMIT 1000
            """)
        )
        logs = [
            {
                "id": row[0],
                "timestamp": row[1],
                "service": row[2],
                "level": row[3],
                "message": row[4],
                "metadata": row[5] if row[5] else {},
            }
            for row in result
        ]

    if len(logs) < 10:
        return {"anomalies_detected": 0, "message": "Not enough logs for anomaly detection"}

    anomalies = detect_anomalies(logs)

    if not anomalies:
        return {"anomalies_detected": 0, "message": "No anomalies detected"}

    created_incidents = []
    with engine.begin() as conn:
        for anomaly in anomalies:
            severity = anomaly.get("severity", "medium")
            description = anomaly.get("description", "Anomaly detected in logs")
            
            result = conn.execute(
                text("""
                    INSERT INTO incidents (status, severity, description)
                    VALUES (:status, :severity, :description)
                    RETURNING id
                """),
                {
                    "status": "open",
                    "severity": severity,
                    "description": description,
                },
            )
            incident_id = result.scalar()
            created_incidents.append(incident_id)
            
            # Send Slack alert for each incident
            await send_slack_alert(incident_id, severity, description)

    return {
        "anomalies_detected": len(anomalies),
        "incidents_created": created_incidents,
        "anomalies": anomalies,
    }


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host="0.0.0.0", port=8000)
