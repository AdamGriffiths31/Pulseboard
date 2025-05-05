import { useEffect, useState } from 'react';
import { fetchLatencyMetrics, fetchStatusCodeDistributionMetrics } from './api';
import useWebSocketMetrics from './hooks/useWebSocketMetrics';
import LatencyChart from './components/LatencyChart';
import StatusCodePieChart from './components/StatusCodePieChart';
import './tailwind.css';

function App() {
  const [metrics, setMetrics] = useState([]);
  const [statusCodes, setStatusCodes] = useState({});
  const [loading, setLoading] = useState(true);
  const liveData = useWebSocketMetrics();

  // Set default start and end times for the last 24 hours
  const [startDate, setStartDate] = useState(new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString().slice(0, 19));
  const [endDate, setEndDate] = useState(new Date().toISOString().slice(0, 19));

  // Fetch historical data
  useEffect(() => {
    setLoading(true);

    const formattedStartDate = new Date(startDate).toISOString();
    const formattedEndDate = new Date(endDate).toISOString();

    Promise.all([
      fetchLatencyMetrics(formattedStartDate, formattedEndDate),
      fetchStatusCodeDistributionMetrics(formattedStartDate, formattedEndDate),
    ])
      .then(([latencyData, statusCodeData]) => {
        const groupedLatency = latencyData.reduce((acc, metric) => {
          if (!acc[metric.endpoint_id]) {
            acc[metric.endpoint_id] = {
              url: metric.url,
              data: [],
            };
          }
          acc[metric.endpoint_id].data.push({
            timestamp: new Date(metric.timestamp).toLocaleString(),
            latency: metric.latency_ms,
          });
          return acc;
        }, {});

        setMetrics(groupedLatency);
        setStatusCodes(statusCodeData);
        setLoading(false);
      })
      .catch(err => {
        console.error("Failed to fetch metrics:", err);
        setLoading(false);
      });
  }, [startDate, endDate]);

  const allUrls = [
    ...new Set([
      ...Object.values(metrics).map(m => m.url),
      ...Object.keys(statusCodes),
      ...Object.keys(liveData),
    ])
  ].sort();

  return (
    <div className="p-6 bg-gray-900 text-gray-100 min-h-screen">
      <h1 className="text-4xl font-extrabold mb-6 text-center text-gray-100 bg-gradient-to-r from-gray-800 via-gray-700 to-gray-800 p-4 rounded-lg shadow-lg border border-gray-600">
        Pulseboard
      </h1>

      {loading ? (
        <p className="text-center text-lg">Loading metrics...</p>
      ) : (
        <div>
          {/* Live Data Section */}
          <h2 className="text-2xl font-bold mb-6 text-left">Live Data</h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {Object.keys(liveData).length ? (
              Object.entries(liveData).map(([url, data]) => (
                <div key={url} className="bg-gray-800 p-6 rounded-lg shadow-lg">
                  <h4 className="text-lg font-semibold mb-4 text-center">{url}</h4>
                  <LatencyChart data={data} />
                </div>
              ))
            ) : (
              <p className="text-center text-lg col-span-full">Waiting for live data...</p>
            )}
          </div>

          {/* Historical Data Section */}
          <h2 className="text-2xl font-bold my-6 text-left">Historical Data</h2>

          <div className="mb-6">
            <label htmlFor="startDate" className="mr-2">Start Date:</label>
            <input
              type="datetime-local"
              id="startDate"
              value={startDate}
              onChange={(e) => setStartDate(e.target.value)}
              className="bg-gray-700 text-white p-2 rounded"
            />
            <label htmlFor="endDate" className="ml-4 mr-2">End Date:</label>
            <input
              type="datetime-local"
              id="endDate"
              value={endDate}
              onChange={(e) => setEndDate(e.target.value)}
              className="bg-gray-700 text-white p-2 rounded"
            />
          </div>

          {/* Latency Section */}
          <h3 className="text-xl font-bold mb-6 text-left">Latency</h3>
          <hr className="border-gray-700 mb-6" />
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {allUrls.map((url) => {
              const endpointId = Object.keys(metrics).find(id => metrics[id].url === url);
              return (
                <div key={url} className="bg-gray-800 p-6 rounded-lg shadow-lg">
                  <h4 className="text-lg font-semibold mb-4 text-center">{url}</h4>
                  <LatencyChart data={endpointId ? metrics[endpointId].data : []} />
                </div>
              );
            })}
          </div>

          {/* Status Code Section */}
          <h3 className="text-xl font-bold my-6 text-left">Status Code Distribution</h3>
          <hr className="border-gray-700 mb-6" />
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {allUrls.map((url) => (
              <div key={url} className="bg-gray-800 p-6 rounded-lg shadow-lg">
                <h4 className="text-lg font-semibold mb-4 text-center">{url}</h4>
                <StatusCodePieChart data={statusCodes[url] || []} />
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}

export default App;
