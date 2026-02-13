import { useState, useEffect } from 'react'

const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:8080'

function App() {
  const [health, setHealth] = useState(null)
  const [incidents, setIncidents] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  const fetchHealth = async () => {
    try {
      const res = await fetch(`${API_BASE}/api/health`)
      const data = await res.json()
      setHealth(data)
    } catch (err) {
      console.error('Health check failed:', err)
    }
  }

  const fetchIncidents = async () => {
    try {
      setLoading(true)
      setError(null)
      const res = await fetch(`${API_BASE}/api/incidents`)
      if (!res.ok) throw new Error('Failed to fetch incidents')
      const data = await res.json()
      setIncidents(data || [])
    } catch (err) {
      setError(err.message)
      console.error('Failed to fetch incidents:', err)
    } finally {
      setLoading(false)
    }
  }

  const fetchIncidentSummary = async (incidentId) => {
    try {
      const res = await fetch(`${API_BASE}/api/summary/${incidentId}`)
      if (!res.ok) throw new Error('Failed to fetch summary')
      const data = await res.json()
      setIncidents(prev => prev.map(inc => 
        inc.id === incidentId ? { ...inc, summary: data.summary, root_cause: data.root_cause } : inc
      ))
    } catch (err) {
      console.error('Failed to fetch summary:', err)
    }
  }

  useEffect(() => {
    fetchHealth()
    fetchIncidents()
    const interval = setInterval(() => {
      fetchHealth()
      fetchIncidents()
    }, 30000)
    return () => clearInterval(interval)
  }, [])

  const openIncidents = incidents.filter(i => i.status === 'open').length
  const highSeverity = incidents.filter(i => i.severity === 'high').length

  return (
    <div className="container">
      <div className="header">
        <h1>🔥 Incident Monitoring & Prediction Platform</h1>
        <p>Real-time system health monitoring and AI-powered incident analysis</p>
      </div>

      <div className="stats-grid">
        <div className="stat-card">
          <h3>System Health</h3>
          <div className="value">
            {health?.checks?.db ? '✓' : '✗'}
          </div>
          <span className={`status ${health?.checks?.db ? 'ok' : 'error'}`}>
            {health?.checks?.db ? 'Operational' : 'Degraded'}
          </span>
        </div>

        <div className="stat-card">
          <h3>Open Incidents</h3>
          <div className="value">{openIncidents}</div>
          <span className={`status ${openIncidents === 0 ? 'ok' : openIncidents > 5 ? 'error' : 'warning'}`}>
            {openIncidents === 0 ? 'All Clear' : openIncidents > 5 ? 'Critical' : 'Active'}
          </span>
        </div>

        <div className="stat-card">
          <h3>High Severity</h3>
          <div className="value">{highSeverity}</div>
          <span className={`status ${highSeverity === 0 ? 'ok' : 'error'}`}>
            {highSeverity === 0 ? 'None' : 'Alert'}
          </span>
        </div>

        <div className="stat-card">
          <h3>Total Incidents</h3>
          <div className="value">{incidents.length}</div>
          <span className="status ok">Tracked</span>
        </div>
      </div>

      <div className="incidents-section">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
          <h2>Recent Incidents</h2>
          <button className="refresh-btn" onClick={fetchIncidents} disabled={loading}>
            {loading ? 'Refreshing...' : 'Refresh'}
          </button>
        </div>

        {error && <div className="error">Error: {error}</div>}

        {loading && incidents.length === 0 ? (
          <div className="loading">Loading incidents...</div>
        ) : incidents.length === 0 ? (
          <div className="loading">No incidents found. System is healthy! 🎉</div>
        ) : (
          <div className="incident-list">
            {incidents.map(incident => (
              <div key={incident.id} className="incident-card">
                <div className="incident-header">
                  <div>
                    <span className="incident-id">Incident #{incident.id}</span>
                    <span className={`incident-severity ${incident.severity}`}>
                      {incident.severity.toUpperCase()}
                    </span>
                  </div>
                  <div className="incident-time">
                    {new Date(incident.created_at).toLocaleString()}
                  </div>
                </div>
                <div className="incident-description">{incident.description}</div>
                {incident.status === 'open' && (
                  <div style={{ marginTop: '8px' }}>
                    <span style={{ 
                      padding: '4px 8px', 
                      background: '#fee2e2', 
                      color: '#991b1b', 
                      borderRadius: '4px', 
                      fontSize: '12px',
                      fontWeight: '600'
                    }}>
                      OPEN
                    </span>
                  </div>
                )}
                {incident.summary && (
                  <div className="incident-summary">
                    <h4>AI Summary:</h4>
                    <p>{incident.summary}</p>
                    {incident.root_cause && (
                      <>
                        <h4 style={{ marginTop: '12px' }}>Root Cause:</h4>
                        <p>{incident.root_cause}</p>
                      </>
                    )}
                  </div>
                )}
                {!incident.summary && incident.status === 'open' && (
                  <button 
                    onClick={() => fetchIncidentSummary(incident.id)}
                    style={{ marginTop: '12px', fontSize: '12px', padding: '6px 12px' }}
                  >
                    Generate AI Analysis
                  </button>
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}

export default App
