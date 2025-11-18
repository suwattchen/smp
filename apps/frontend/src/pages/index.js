// File: apps/frontend/src/pages/index.js

import { useState, useEffect } from 'react';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api";

export default function Home() {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    async function fetchData() {
      try {
        // *** สมมติว่ามี Endpoint ทดสอบสุขภาพระบบ (Health Check)
        const response = await fetch(`${API_BASE_URL}/health`); 
        
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const result = await response.json();
        setData(result);
      } catch (e) {
        setError(e.message);
      } finally {
        setLoading(false);
      }
    }
    fetchData();
  }, []);

  if (loading) return <div style={{ padding: 20 }}>Loading Backend status...</div>;
  if (error) return (
    <div style={{ padding: 20, color: 'red', border: '1px solid red' }}>
      <h1>Connection Failed! ❌</h1>
      <p>Could not connect to Backend at: {API_BASE_URL}/health</p>
      <p>Error: {error}</p>
    </div>
  );

  return (
    <div style={{ padding: 20, maxWidth: 600, margin: '0 auto' }}>
      <h1>Project Sunmart Frontend</h1>
      <p>Welcome to your Next.js application.</p>
      <hr />
      <h2>Backend Status (Core-Go) ✅</h2>
      <p>Status received from **{API_BASE_URL}/health**:</p>
      <pre style={{ backgroundColor: '#f0f0f0', padding: 10 }}>
        {JSON.stringify(data, null, 2)}
      </pre>
      <p>Now connected and ready to develop!</p>
    </div>
  );
}
