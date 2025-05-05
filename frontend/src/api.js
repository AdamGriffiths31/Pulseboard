import axios from "axios";

const API_BASE_URL = "http://localhost:8080"; // Backend URL

// Fetch latency metrics with startDate and endDate
export const fetchLatencyMetrics = async (startDate, endDate) => {
  try {
    console.log("Requesting latency metrics with params:", { startDate, endDate });
    const response = await axios.get(`${API_BASE_URL}/getlatency`, {
      params: {
        startDate: startDate,
        endDate: endDate,    
      },
    });
    return response.data; // Return the fetched data
  } catch (error) {
    console.error("Error fetching latency metrics:", error);
    return [];
  }
};

// Fetch status code distribution metrics (no date range needed for this example)
export const fetchStatusCodeDistributionMetrics = async (startDate, endDate) => {
  try {
    console.log("Requesting status code distribution metrics");
    const response = await axios.get(`${API_BASE_URL}/statuscodedistribution`, {
      params: {
        startDate: startDate,
        endDate: endDate,    
      },
    });
    return response.data; // Return the fetched data
  } catch (error) {
    console.error("Error fetching status code distribution:", error);
    return {};
  }
};
