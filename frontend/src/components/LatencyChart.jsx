import React from "react";
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from "recharts";

function LatencyChart({ data }) {
  return (
    <ResponsiveContainer width="100%" height={300}>
      <LineChart data={data}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="timestamp" />
        <YAxis />
        <Tooltip
          contentStyle={{
            backgroundColor: "#2d3748", // Dark gray background
            borderColor: "#4a5568", // Slightly lighter border
            color: "#f7fafc", // Light text color
          }}
          itemStyle={{
            color: "#f7fafc", // Light text color for items
          }}
          formatter={(value, name, props) => {
            if (name === "timestamp") {
              return [new Date(value).toLocaleString(), name];
            }
            return [value, name];
          }}
        />
        <Legend />
        <Line type="monotone" dataKey="latency" stroke="#8884d8" />
      </LineChart>
    </ResponsiveContainer>
  );
}

export default LatencyChart;
