import json
from typing import Any, Dict, List

from openai import OpenAI


async def summarize_incident(
    incident_description: str, recent_logs: List[Dict[str, Any]], api_key: str
) -> tuple[str, str]:
    """
    Use OpenAI API to generate a summary and root cause analysis for an incident.
    Returns (summary, root_cause) tuple.
    """
    if not api_key:
        return (
            "LLM analysis unavailable: API key not configured",
            "Please configure OPENAI_API_KEY environment variable",
        )

    client = OpenAI(api_key=api_key)

    error_logs = [log for log in recent_logs if log.get("level", "").lower() in {"error", "critical", "fatal", "panic"}]
    relevant_logs = error_logs[:20] if error_logs else recent_logs[:20]

    log_context = "\n".join(
        [
            f"[{log.get('timestamp', 'N/A')}] {log.get('service', 'unknown')} [{log.get('level', 'info')}]: {log.get('message', '')}"
            for log in relevant_logs
        ]
    )

    prompt = f"""You are a DevOps engineer analyzing an incident. Based on the incident description and recent logs, provide:

1. A concise summary of what happened (2-3 sentences)
2. The most likely root cause (1-2 sentences)

Incident Description:
{incident_description}

Recent Logs:
{log_context}

Respond in JSON format with "summary" and "root_cause" fields.
"""

    try:
        response = client.chat.completions.create(
            model="gpt-4o-mini",
            messages=[
                {
                    "role": "system",
                    "content": "You are a helpful DevOps engineer analyzing system incidents. Always respond with valid JSON.",
                },
                {"role": "user", "content": prompt},
            ],
            temperature=0.3,
            max_tokens=500,
        )

        content = response.choices[0].message.content.strip()
        if content.startswith("```json"):
            content = content[7:]
        if content.startswith("```"):
            content = content[3:]
        if content.endswith("```"):
            content = content[:-3]
        content = content.strip()

        result = json.loads(content)
        summary = result.get("summary", "Unable to generate summary")
        root_cause = result.get("root_cause", "Unable to determine root cause")

        return summary, root_cause

    except Exception as e:
        return (
            f"Error during LLM analysis: {str(e)}",
            "Please check API key and network connectivity",
        )
