import { useEffect, useState } from 'react';

export default function useWebSocketMetrics() {
  const [liveData, setLiveData] = useState({});

  useEffect(() => {
    const socket = new WebSocket('ws://localhost:8080/ws');

    socket.onopen = () => console.log('WebSocket connected');

    socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data); // should be an array of metrics
        const grouped = {};

        for (const metric of data) {
          const { url, timestamp, latency_ms } = metric;

          if (!grouped[url]) {
            grouped[url] = [];
          }

          grouped[url].push({
            timestamp: new Date(timestamp).toLocaleString(),
            latency: latency_ms,
          });
        }

        // Merge new data into state per URL
        setLiveData(prev => {
          const updated = { ...prev };
          for (const url in grouped) {
            if (!updated[url]) {
              updated[url] = [];
            }
            updated[url] = [...updated[url], ...grouped[url]].slice(-50); // keep last 50 points
          }
          return updated;
        });

      } catch (err) {
        console.error('Failed to parse message:', err);
      }
    };

    socket.onerror = (err) => console.error('WebSocket error:', err);
    socket.onclose = () => console.log('WebSocket closed');

    return () => socket.close();
  }, []);

  return liveData;
}
