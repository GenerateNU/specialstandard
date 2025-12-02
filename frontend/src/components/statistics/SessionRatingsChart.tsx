import { CartesianGrid, Line, LineChart, XAxis, YAxis } from "recharts";
import { ChartContainer, ChartTooltip } from "@/components/ui/chart";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { Input } from "@/components/ui/input";
import React, { useEffect, useState } from "react";

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
  // Calculate initial date range - ensure yyyy-mm-dd format
  const getDateRange = () => {
    const dates = chartData
      .filter(d => categories.some(cat => d[cat.key] !== null))
      .map(d => {
        const dateStr = d.session as string;
        // Extract just the date part if it's a full ISO datetime
        return dateStr.includes('T') ? dateStr.split('T')[0] : dateStr.slice(0, 10);
      })
      .sort();
    return dates.length > 0 ? { min: dates[0], max: dates[dates.length - 1] } : { min: '', max: '' };
  };

  // State to track which categories are visible
  const [visibleCategories, setVisibleCategories] = useState<Record<string, boolean>>(
    categories.reduce((acc, cat) => ({ ...acc, [cat.key]: true }), {})
  );

  // State for date range filtering - initialize with calculated range
  const [startDate, setStartDate] = useState<string>('');
  const [endDate, setEndDate] = useState<string>('');

  // Update date range when chartData changes
  useEffect(() => {
    const { min, max } = getDateRange();
    if (min) {
      setStartDate(min);
    }
    // Add a day to max date to compensate for timezone offset
    if (max) {
      const maxDate = new Date(max);
      maxDate.setDate(maxDate.getDate() + 1);
      setEndDate(maxDate.toISOString().split('T')[0]);
    }
  }, [chartData, categories]);

  // Toggle category visibility
  const toggleCategory = (key: string) => {
    setVisibleCategories(prev => ({ ...prev, [key]: !prev[key] }));
  };

  // Build chart config from categories
  const chartConfig = categories.reduce((config, category) => {
    config[category.key] = {
      label: category.label,
      color: category.color,
    };
    return config;
  }, {} as Record<string, { label: string; color: string }>);

  // Filter out sessions where all specified categories are null, apply date range filter, then sort chronologically
  const sortedData = [...chartData]
    .filter(d => categories.some(cat => d[cat.key] !== null))
    .filter(d => {
      const sessionDate = d.session as string;
      if (startDate && sessionDate < startDate) return false;
      if (endDate && sessionDate > endDate) return false;
      return true;
    })
    .sort((a, b) => a.session.localeCompare(b.session));

  // Filter categories to only show visible ones
  const activeCategories = categories.filter(cat => visibleCategories[cat.key]);

  return (
  <Card className="bg-white border-0">
    <CardHeader className="space-y-4 pb-4">
      <div className="flex flex-row items-center justify-between">
        <CardTitle>{title}</CardTitle>
        <div className="flex gap-4">
        {categories.length > 1 &&
        categories.map(category => (
          <div key={category.key} className="flex items-center gap-2">
            <Checkbox
              id={`${category.key}-checkbox`}
              checked={visibleCategories[category.key]}
              onCheckedChange={() => toggleCategory(category.key)}
              style={visibleCategories[category.key] ? { 
                backgroundColor: category.color, 
                borderColor: category.color,
                cursor: 'pointer'
              } : {cursor: 'pointer'}}
            />
            <label
              htmlFor={`${category.key}-checkbox`}
              className="text-sm font-medium cursor-pointer select-none"
              style={{ color: category.color }}
            >
              {category.label}
            </label>
          </div>
        ))}
        </div>
      </div>
      <div className="flex items-center gap-3">
        <span className="text-sm font-medium text-muted-foreground">Date Range:</span>
        <Input
          type="date"
          value={startDate}
          onChange={(e) => setStartDate(e.target.value)}
          className="w-40"
        />
        <span className="text-sm text-muted-foreground">to</span>
        <Input
          type="date"
          value={endDate}
          onChange={(e) => setEndDate(e.target.value)}
          className="w-40"
        />
      </div>
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
                content={<CustomTooltip categories={activeCategories} />}
              />
              {activeCategories.map(category => (
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
