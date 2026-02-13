from collections import defaultdict
from datetime import datetime, timedelta
from typing import Any, Dict, List

import numpy as np
from sklearn.ensemble import IsolationForest


def detect_anomalies(logs: List[Dict[str, Any]]) -> List[Dict[str, Any]]:
    """
    Detect anomalies in logs using multiple heuristics:
    1. Error rate spike detection
    2. Isolation Forest on log frequency patterns
    3. Service-specific anomaly detection
    """
    anomalies = []

    if not logs:
        return anomalies

    error_levels = {"error", "critical", "fatal", "panic"}
    warning_levels = {"warn", "warning"}

    time_buckets = defaultdict(lambda: {"total": 0, "errors": 0, "warnings": 0})
    service_stats = defaultdict(lambda: {"total": 0, "errors": 0})

    for log in logs:
        timestamp = log.get("timestamp")
        if isinstance(timestamp, str):
            try:
                timestamp = datetime.fromisoformat(timestamp.replace("Z", "+00:00"))
            except:
                continue
        if not isinstance(timestamp, datetime):
            continue

        bucket_key = timestamp.replace(second=0, microsecond=0)
        level = log.get("level", "").lower()
        service = log.get("service", "unknown")

        time_buckets[bucket_key]["total"] += 1
        service_stats[service]["total"] += 1

        if level in error_levels:
            time_buckets[bucket_key]["errors"] += 1
            service_stats[service]["errors"] += 1
        elif level in warning_levels:
            time_buckets[bucket_key]["warnings"] += 1

    if not time_buckets:
        return anomalies

    buckets = sorted(time_buckets.items())
    error_rates = [b[1]["errors"] / max(b[1]["total"], 1) for b in buckets]
    total_counts = [b[1]["total"] for b in buckets]

    if len(error_rates) < 3:
        return anomalies

    mean_error_rate = np.mean(error_rates)
    std_error_rate = np.std(error_rates) if len(error_rates) > 1 else 0

    mean_count = np.mean(total_counts)
    std_count = np.std(total_counts) if len(total_counts) > 1 else 0

    threshold_error = mean_error_rate + 2 * std_error_rate
    threshold_count = mean_count + 2 * std_count

    for bucket_time, stats in buckets:
        error_rate = stats["errors"] / max(stats["total"], 1)
        if error_rate > threshold_error and stats["total"] > threshold_count:
            anomalies.append(
                {
                    "type": "error_rate_spike",
                    "severity": "high" if error_rate > 0.5 else "medium",
                    "timestamp": bucket_time.isoformat(),
                    "description": f"Error rate spike detected: {stats['errors']}/{stats['total']} logs are errors ({error_rate*100:.1f}%)",
                    "details": {
                        "error_rate": error_rate,
                        "total_logs": stats["total"],
                        "error_count": stats["errors"],
                    },
                }
            )

    for service, stats in service_stats.items():
        if stats["total"] < 5:
            continue
        error_rate = stats["errors"] / stats["total"]
        if error_rate > 0.3:
            anomalies.append(
                {
                    "type": "service_error_rate",
                    "severity": "medium",
                    "service": service,
                    "description": f"High error rate in {service}: {stats['errors']}/{stats['total']} logs are errors ({error_rate*100:.1f}%)",
                    "details": {
                        "service": service,
                        "error_rate": error_rate,
                        "total_logs": stats["total"],
                        "error_count": stats["errors"],
                    },
                }
            )

    if len(total_counts) >= 10:
        X = np.array(total_counts).reshape(-1, 1)
        iso_forest = IsolationForest(contamination=0.1, random_state=42)
        predictions = iso_forest.fit_predict(X)

        for i, pred in enumerate(predictions):
            if pred == -1:
                bucket_time, stats = buckets[i]
                anomalies.append(
                    {
                        "type": "log_volume_anomaly",
                        "severity": "medium",
                        "timestamp": bucket_time.isoformat(),
                        "description": f"Unusual log volume detected: {stats['total']} logs in this time window",
                        "details": {
                            "log_count": stats["total"],
                            "expected_range": f"{mean_count - std_count:.0f} - {mean_count + std_count:.0f}",
                        },
                    }
                )

    return anomalies
