import { CartesianGrid, Line, LineChart, XAxis, YAxis } from "recharts";
import { ChartContainer, ChartTooltip } from "@/components/ui/chart";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import React from "react";

const ratingLevelLabels = {
  1: "minimal",
  2: "low",
  3: "moderate",
  4: "high",
  5: "maximal",
};

type RatingCategory = "visual_cue" | "verbal_cue" | "gestural_cue" | "engagement";

interface CategoryConfig {
  key: RatingCategory;
  label: string;
  color: string;
}

interface SessionRatingsChartProps {
  title: string;
  chartData: Array<{ 
    session: string;
    [key: string]: number | null | string;
  }>;
  categories: CategoryConfig[];
}

const CustomTooltip = ({ active, payload, categories }: any) => {
  if (active && payload && payload.length) {
    const data = payload[0].payload;
    
    return (
      <div className="bg-white border border-border rounded-lg shadow-lg p-3">
        <p className="text-sm font-medium mb-2">Date: {data.session}</p>
        {categories.map((category: CategoryConfig, index: number) => (
          <div key={category.key} className={index < categories.length - 1 ? "mb-1" : ""}>
            <span className="text-sm font-medium" style={{ color: category.color }}>
              {category.label}:{" "}
            </span>
            <span className="text-sm">
              {data[category.key] !== null 
                ? `${data[category.key]} (${ratingLevelLabels[data[category.key] as keyof typeof ratingLevelLabels]})` 
                : "No data"}
            </span>
          </div>
        ))}
      </div>
    );
  }
  return null;
};

const SessionRatingsChart: React.FC<SessionRatingsChartProps> = ({ title, chartData, categories }) => {
  // Build chart config from categories
  const chartConfig = categories.reduce((config, category) => {
    config[category.key] = {
      label: category.label,
      color: category.color,
    };
    return config;
  }, {} as Record<string, { label: string; color: string }>);

  // Filter out sessions where all specified categories are null, then sort chronologically
  const sortedData = [...chartData]
    .filter(d => categories.some(cat => d[cat.key] !== null))
    .sort((a, b) => a.session.localeCompare(b.session));

  return (
  <Card className="bg-white border-0">
    <CardHeader>
      <CardTitle>{title}</CardTitle>
    </CardHeader>
    <CardContent>
      <ChartContainer config={chartConfig} className="h-[200px]">
        <LineChart
              accessibilityLayer
              data={sortedData}
              margin={{ left: 12, right: 12, top: 12 }}
            >
              <CartesianGrid vertical={false} />
              <XAxis
                dataKey="session"
                tickLine={false}
                axisLine={false}
                tickMargin={8}
                tickFormatter={(value) =>
                  typeof value === "string" ? value.slice(0, 10) : value
                }
              />
              <YAxis
                domain={[1, 5]}
                ticks={[1, 2, 3, 4, 5]}
                tickLine={false}
                axisLine={false}
                tickMargin={8}
                tickFormatter={(value: number) => (ratingLevelLabels[value as keyof typeof ratingLevelLabels] ?? String(value))}
              />
              <ChartTooltip
                cursor={false}
                content={<CustomTooltip categories={categories} />}
              />
              {categories.map(category => (
                <Line
                  key={category.key}
                  dataKey={category.key}
                  type="linear"
                  stroke={category.color}
                  strokeWidth={2}
                  dot={{ r: 3, fill: category.color }}
                  connectNulls
                />
              ))}
          </LineChart>
        </ChartContainer>
    </CardContent>
  </Card>
  );
};

export default SessionRatingsChart;
