import { PieChart, Pie, Cell, Tooltip, Legend, ResponsiveContainer } from 'recharts';

const StatusCodePieChart = ({ data }) => {
  const COLORS = data.map(item => (item.status_code >= 200 && item.status_code < 300 ? '#22c55e' : '#ef4444'));

  return (
    <ResponsiveContainer width="100%" height={250}>
      <PieChart>
        <Pie
          data={data}
          dataKey="count"
          nameKey="status_code"
          cx="50%"
          cy="50%"
          outerRadius={80}
          label={({ status_code }) => status_code}
        >
          {data.map((entry, index) => (
            <Cell key={`cell-${index}`} fill={COLORS[index]} />
          ))}
        </Pie>
        <Tooltip />
        <Legend />
      </PieChart>
    </ResponsiveContainer>
  );
};

export default StatusCodePieChart;
