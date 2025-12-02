// File: components/statistics/GamePerformanceCharts.tsx
import React from 'react'
import {
  Bar,
  CartesianGrid,
  Cell,
  ComposedChart,
  Legend,
  Line,
  Pie,
  PieChart,
  ReferenceLine,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import type {
  CategoryPerformance,
  GamePerformanceStats,
  PerformanceByType,
  PerformanceTrend,
} from '@/hooks/useStudentGameResults'

interface GameOverallStatsProps {
  stats: GamePerformanceStats
}

export function GameOverallStats({ stats }: GameOverallStatsProps) {
  const StatCard = ({
    label,
    value,
    unit = '',
  }: {
    label: string
    value: number | string
    unit?: string
  }) => (
    <div className="flex flex-col gap-2 p-4 rounded-2xl bg-card border-2 border-default hover:border-primary/50 transition-colors">
      <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">{label}</span>
      <div className="flex items-baseline gap-1">
        <span className="text-3xl font-bold text-primary">{value}</span>
        {unit && <span className="text-sm text-muted-foreground font-medium">{unit}</span>}
      </div>
    </div>
  )

  return (
    <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-3">
      <StatCard label="Total Games" value={stats.totalGames} />
      <StatCard
        label="Completed"
        value={stats.completedGames}
        unit={`/ ${stats.totalGames}`}
      />
      <StatCard label="Completion" value={stats.completionRate.toFixed(0)} unit="%" />
      <StatCard label="Avg Time" value={stats.avgTimeSeconds} unit="s" />
      <StatCard label="Total Errors" value={stats.totalErrors} />
      <StatCard label="Avg Errors" value={stats.avgErrors} />
    </div>
  )
}

interface PerformanceByTypeChartProps {
  data: PerformanceByType[]
  title: string
}

export function PerformanceByTypeChart({
  data,
  title,
}: PerformanceByTypeChartProps) {
  if (!data || data.length === 0) {
    return (
      <Card className="bg-card border-2 border-default rounded-4xl shadow-lg">
        <CardContent className="flex items-center justify-center h-64 text-muted-foreground">
          No data available
        </CardContent>
      </Card>
    )
  }

  const avgCompletionRate = data.reduce((sum, item) => sum + item.completionRate, 0) / data.length

  const CustomTooltip = ({ active, payload }: any) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload
      return (
        <div className="bg-card border-2 border-primary/20 rounded-xl p-4 shadow-xl">
          <p className="font-semibold text-primary mb-2">{data.label}</p>
          <div className="space-y-1 text-sm">
            <p className="flex justify-between gap-4">
              <span className="text-muted-foreground">Total Attempts:</span>
              <span className="font-semibold">{data.total}</span>
            </p>
            <p className="flex justify-between gap-4">
              <span className="text-muted-foreground">Completed:</span>
              <span className="font-semibold text-green-600">{data.completed}</span>
            </p>
            <p className="flex justify-between gap-4">
              <span className="text-muted-foreground">Completion Rate:</span>
              <span className="font-semibold text-blue-600">{data.completionRate.toFixed(1)}%</span>
            </p>
            <p className="flex justify-between gap-4">
              <span className="text-muted-foreground">Avg Time:</span>
              <span className="font-semibold">{data.avgTimeSeconds}s</span>
            </p>
            <p className="flex justify-between gap-4">
              <span className="text-muted-foreground">Avg Errors:</span>
              <span className="font-semibold text-red-600">{data.avgErrors}</span>
            </p>
          </div>
        </div>
      )
    }
    return null
  }

  return (
    <Card className="bg-card border-2 border-default rounded-4xl shadow-lg hover:shadow-xl transition-shadow">
      <CardHeader className="pb-4">
        <div className="flex items-center justify-between">
          <CardTitle className="text-xl font-bold text-primary">{title}</CardTitle>
          <div className="text-sm text-muted-foreground">
            Avg: <span className="font-semibold text-primary">{avgCompletionRate.toFixed(1)}%</span>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <ResponsiveContainer width="100%" height={350}>
          <ComposedChart data={data} margin={{ top: 10, right: 20, left: 0, bottom: 5 }}>
            <CartesianGrid strokeDasharray="3 3" stroke="hsl(var(--border))" opacity={0.3} />
            <XAxis 
              dataKey="label" 
              stroke="hsl(var(--muted-foreground))"
              tick={{ fill: 'hsl(var(--muted-foreground))', fontSize: 12 }}
              angle={-45}
              textAnchor="end"
              height={80}
            />
            <YAxis 
              yAxisId="left"
              stroke="hsl(var(--muted-foreground))"
              tick={{ fill: 'hsl(var(--muted-foreground))', fontSize: 12 }}
              label={{ value: 'Count', angle: -90, position: 'insideLeft', style: { textAnchor: 'middle' } }}
            />
            <YAxis 
              yAxisId="right"
              orientation="right"
              stroke="hsl(var(--muted-foreground))"
              tick={{ fill: 'hsl(var(--muted-foreground))', fontSize: 12 }}
              label={{ value: 'Rate (%)', angle: 90, position: 'insideRight', style: { textAnchor: 'middle' } }}
            />
            <Tooltip content={<CustomTooltip />} />
            <Legend 
              wrapperStyle={{ paddingTop: '20px' }}
              iconType="rect"
            />
            <Bar 
              yAxisId="left"
              dataKey="total" 
              fill="hsl(var(--muted))" 
              name="Total Attempts"
              radius={[8, 8, 0, 0]}
              opacity={0.6}
            />
            <Bar 
              yAxisId="left"
              dataKey="completed" 
              fill="hsl(142 76% 36%)" 
              name="Completed"
              radius={[8, 8, 0, 0]}
            />
            <Line 
              yAxisId="right"
              type="monotone" 
              dataKey="completionRate" 
              stroke="hsl(217 91% 60%)" 
              strokeWidth={3}
              name="Completion Rate (%)"
              dot={{ r: 5, fill: 'hsl(217 91% 60%)', strokeWidth: 2, stroke: 'white' }}
              activeDot={{ r: 7 }}
            />
            <ReferenceLine 
              yAxisId="right"
              y={avgCompletionRate} 
              stroke="hsl(var(--muted-foreground))" 
              strokeDasharray="5 5"
              label={{ value: 'Average', position: 'right', fill: 'hsl(var(--muted-foreground))' }}
            />
          </ComposedChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  )
}

interface CompletionRateChartProps {
  data: PerformanceByType[]
  title: string
}

export function CompletionRateChart({
  data,
  title,
}: CompletionRateChartProps) {
  if (!data || data.length === 0) {
    return (
      <Card className="bg-card border-2 border-default rounded-4xl shadow-lg">
        <CardContent className="flex items-center justify-center h-64 text-muted-foreground">
          No data available
        </CardContent>
      </Card>
    )
  }

  const pieData = data.map((item) => ({
    name: item.label,
    value: Math.round(item.completionRate),
    total: item.total,
    completed: item.completed,
    avgTime: item.avgTimeSeconds,
    avgErrors: item.avgErrors,
  }))

  const COLORS = [
    'hsl(142 76% 36%)',
    'hsl(217 91% 60%)',
    'hsl(330 81% 60%)',
    'hsl(25 95% 53%)',
    'hsl(262 83% 58%)',
  ]

  const CustomTooltip = ({ active, payload }: any) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload
      return (
        <div className="bg-card border-2 border-primary/20 rounded-xl p-4 shadow-xl">
          <p className="font-semibold text-primary mb-2">{data.name}</p>
          <div className="space-y-1 text-sm">
            <p className="flex justify-between gap-4">
              <span className="text-muted-foreground">Completion Rate:</span>
              <span className="font-semibold text-blue-600">{data.value}%</span>
            </p>
            <p className="flex justify-between gap-4">
              <span className="text-muted-foreground">Completed:</span>
              <span className="font-semibold text-green-600">{data.completed} / {data.total}</span>
            </p>
            <p className="flex justify-between gap-4">
              <span className="text-muted-foreground">Avg Time:</span>
              <span className="font-semibold">{data.avgTime}s</span>
            </p>
            <p className="flex justify-between gap-4">
              <span className="text-muted-foreground">Avg Errors:</span>
              <span className="font-semibold text-red-600">{data.avgErrors}</span>
            </p>
          </div>
        </div>
      )
    }
    return null
  }

  const renderLabel = ({ name, value, percent }: any) => {
    return `${name}\n${value}% (${(percent * 100).toFixed(1)}%)`
  }

  return (
    <Card className="bg-card border-2 border-default rounded-4xl shadow-lg hover:shadow-xl transition-shadow">
      <CardHeader className="pb-4">
        <CardTitle className="text-xl font-bold text-primary">{title}</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="flex items-center justify-center">
          <ResponsiveContainer width="100%" height={350}>
            <PieChart>
              <Pie
                data={pieData}
                cx="50%"
                cy="50%"
                labelLine={false}
                label={renderLabel}
                outerRadius={100}
                innerRadius={40}
                fill="hsl(217 91% 60%)"
                dataKey="value"
                paddingAngle={2}
              >
                {pieData.map((entry, index) => (
                  <Cell 
                    key={`cell-${index}`} 
                    fill={COLORS[index % COLORS.length]}
                    stroke="hsl(var(--card))"
                    strokeWidth={2}
                  />
                ))}
              </Pie>
              <Tooltip content={<CustomTooltip />} />
              <Legend 
                wrapperStyle={{ paddingTop: '20px' }}
                iconType="circle"
                formatter={(value, entry: any) => (
                  <span className="text-foreground" style={{ color: entry.color } as React.CSSProperties}>
                    {value}: {entry.payload.value}%
                  </span>
                )}
              />
            </PieChart>
          </ResponsiveContainer>
        </div>
      </CardContent>
    </Card>
  )
}

interface PerformanceTrendChartProps {
  data: PerformanceTrend[]
}

export function PerformanceTrendChart({ data }: PerformanceTrendChartProps) {
  if (!data || data.length === 0) {
    return (
      <Card className="bg-card border-2 border-default rounded-4xl shadow-lg">
        <CardContent className="flex items-center justify-center h-64 text-muted-foreground">
          No data available
        </CardContent>
      </Card>
    )
  }

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr)
    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
  }

  const chartData = data.map((item) => ({
    ...item,
    formattedDate: formatDate(item.date),
    total: item.completed + item.failed,
    successRate: item.completed + item.failed > 0 
      ? (item.completed / (item.completed + item.failed)) * 100 
      : 0,
  }))

  const avgSuccessRate = chartData.reduce((sum, item) => sum + item.successRate, 0) / chartData.length

  const CustomTooltip = ({ active, payload }: any) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload
      return (
        <div className="bg-card border-2 border-primary/20 rounded-xl p-4 shadow-xl">
          <p className="font-semibold text-primary mb-2">{data.formattedDate}</p>
          <div className="space-y-1 text-sm">
            <p className="flex justify-between gap-4">
              <span className="text-muted-foreground">Total Sessions:</span>
              <span className="font-semibold">{data.total}</span>
            </p>
            <p className="flex justify-between gap-4">
              <span className="text-muted-foreground">Completed:</span>
              <span className="font-semibold text-green-600">{data.completed}</span>
            </p>
            <p className="flex justify-between gap-4">
              <span className="text-muted-foreground">Failed:</span>
              <span className="font-semibold text-red-600">{data.failed}</span>
            </p>
            <p className="flex justify-between gap-4">
              <span className="text-muted-foreground">Success Rate:</span>
              <span className="font-semibold text-blue-600">{data.successRate.toFixed(1)}%</span>
            </p>
            <p className="flex justify-between gap-4">
              <span className="text-muted-foreground">Avg Time:</span>
              <span className="font-semibold">{data.avgTimeSeconds}s</span>
            </p>
            <p className="flex justify-between gap-4">
              <span className="text-muted-foreground">Total Errors:</span>
              <span className="font-semibold text-orange-600">{data.totalErrors}</span>
            </p>
          </div>
        </div>
      )
    }
    return null
  }

  return (
    <Card className="bg-card border-2 border-default rounded-4xl shadow-lg hover:shadow-xl transition-shadow">
      <CardHeader className="pb-4">
        <div className="flex items-center justify-between">
          <CardTitle className="text-xl font-bold text-primary">Performance Trend Over Time</CardTitle>
          <div className="text-sm text-muted-foreground">
            Avg Success: <span className="font-semibold text-primary">{avgSuccessRate.toFixed(1)}%</span>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <ResponsiveContainer width="100%" height={400}>
          <ComposedChart data={chartData} margin={{ top: 10, right: 20, left: 0, bottom: 5 }}>
            <CartesianGrid strokeDasharray="3 3" stroke="hsl(var(--border))" opacity={0.3} />
            <XAxis 
              dataKey="formattedDate" 
              stroke="hsl(var(--muted-foreground))"
              tick={{ fill: 'hsl(var(--muted-foreground))', fontSize: 11 }}
              angle={-45}
              textAnchor="end"
              height={80}
            />
            <YAxis 
              yAxisId="left"
              stroke="hsl(var(--muted-foreground))"
              tick={{ fill: 'hsl(var(--muted-foreground))', fontSize: 12 }}
              label={{ value: 'Count', angle: -90, position: 'insideLeft', style: { textAnchor: 'middle' } }}
            />
            <YAxis 
              yAxisId="right"
              orientation="right"
              stroke="hsl(var(--muted-foreground))"
              tick={{ fill: 'hsl(var(--muted-foreground))', fontSize: 12 }}
              label={{ value: 'Rate (%)', angle: 90, position: 'insideRight', style: { textAnchor: 'middle' } }}
            />
            <Tooltip content={<CustomTooltip />} />
            <Legend 
              wrapperStyle={{ paddingTop: '20px' }}
              iconType="line"
            />
            <Bar 
              yAxisId="left"
              dataKey="total" 
              fill="hsl(var(--muted))" 
              name="Total Sessions"
              radius={[4, 4, 0, 0]}
              opacity={0.3}
            />
            <Bar 
              yAxisId="left"
              dataKey="completed" 
              fill="hsl(142 76% 36%)" 
              name="Completed"
              radius={[4, 4, 0, 0]}
            />
            <Bar 
              yAxisId="left"
              dataKey="failed" 
              fill="hsl(0 84% 60%)" 
              name="Failed"
              radius={[4, 4, 0, 0]}
            />
            <Line 
              yAxisId="right"
              type="monotone" 
              dataKey="successRate" 
              stroke="hsl(217 91% 60%)" 
              strokeWidth={3}
              name="Success Rate (%)"
              dot={{ r: 5, fill: 'hsl(217 91% 60%)', strokeWidth: 2, stroke: 'white' }}
              activeDot={{ r: 7 }}
            />
            <ReferenceLine 
              yAxisId="right"
              y={avgSuccessRate} 
              stroke="hsl(var(--muted-foreground))" 
              strokeDasharray="5 5"
              label={{ value: 'Average', position: 'right', fill: 'hsl(var(--muted-foreground))' }}
            />
          </ComposedChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  )
}

interface CategoryPerformanceChartProps {
  data: CategoryPerformance[]
}

export function CategoryPerformanceChart({
  data,
}: CategoryPerformanceChartProps) {
  if (!data || data.length === 0) {
    return (
      <Card className="bg-card border-2 border-default rounded-4xl shadow-lg">
        <CardContent className="flex items-center justify-center h-64 text-muted-foreground">
          No data available
        </CardContent>
      </Card>
    )
  }

  const sortedData = [...data].sort((a, b) => b.total - a.total)
  const avgCompletionRate = data.reduce((sum, item) => sum + item.completionRate, 0) / data.length

  const CustomTooltip = ({ active, payload }: any) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload
      return (
        <div className="bg-card border-2 border-primary/20 rounded-xl p-4 shadow-xl">
          <p className="font-semibold text-primary mb-2">{data.label}</p>
          <div className="space-y-1 text-sm">
            <p className="flex justify-between gap-4">
              <span className="text-muted-foreground">Total Attempts:</span>
              <span className="font-semibold">{data.total}</span>
            </p>
            <p className="flex justify-between gap-4">
              <span className="text-muted-foreground">Completed:</span>
              <span className="font-semibold text-green-600">{data.completed}</span>
            </p>
            <p className="flex justify-between gap-4">
              <span className="text-muted-foreground">Completion Rate:</span>
              <span className="font-semibold text-blue-600">{data.completionRate.toFixed(1)}%</span>
            </p>
            <p className="flex justify-between gap-4">
              <span className="text-muted-foreground">Avg Time:</span>
              <span className="font-semibold">{data.avgTimeSeconds}s</span>
            </p>
            <p className="flex justify-between gap-4">
              <span className="text-muted-foreground">Avg Errors:</span>
              <span className="font-semibold text-red-600">{data.avgErrors}</span>
            </p>
          </div>
        </div>
      )
    }
    return null
  }

  return (
    <Card className="bg-card border-2 border-default rounded-4xl shadow-lg hover:shadow-xl transition-shadow">
      <CardHeader className="pb-4">
        <div className="flex items-center justify-between">
          <CardTitle className="text-xl font-bold text-primary">Performance by Category</CardTitle>
          <div className="text-sm text-muted-foreground">
            Avg: <span className="font-semibold text-primary">{avgCompletionRate.toFixed(1)}%</span>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <ResponsiveContainer width="100%" height={400}>
          <ComposedChart data={sortedData} margin={{ top: 10, right: 20, left: 0, bottom: 5 }}>
            <CartesianGrid strokeDasharray="3 3" stroke="hsl(var(--border))" opacity={0.3} />
            <XAxis 
              dataKey="label" 
              stroke="hsl(var(--muted-foreground))"
              tick={{ fill: 'hsl(var(--muted-foreground))', fontSize: 11 }}
              angle={-45}
              textAnchor="end"
              height={100}
            />
            <YAxis 
              yAxisId="left"
              stroke="hsl(var(--muted-foreground))"
              tick={{ fill: 'hsl(var(--muted-foreground))', fontSize: 12 }}
              label={{ value: 'Count', angle: -90, position: 'insideLeft', style: { textAnchor: 'middle' } }}
            />
            <YAxis 
              yAxisId="right"
              orientation="right"
              stroke="hsl(var(--muted-foreground))"
              tick={{ fill: 'hsl(var(--muted-foreground))', fontSize: 12 }}
              label={{ value: 'Rate (%)', angle: 90, position: 'insideRight', style: { textAnchor: 'middle' } }}
            />
            <Tooltip content={<CustomTooltip />} />
            <Legend 
              wrapperStyle={{ paddingTop: '20px' }}
              iconType="rect"
            />
            <Bar 
              yAxisId="left"
              dataKey="total" 
              fill="hsl(var(--muted))" 
              name="Total Attempts"
              radius={[8, 8, 0, 0]}
              opacity={0.6}
            />
            <Bar 
              yAxisId="left"
              dataKey="completed" 
              fill="hsl(142 76% 36%)" 
              name="Completed"
              radius={[8, 8, 0, 0]}
            />
            <Line 
              yAxisId="right"
              type="monotone" 
              dataKey="completionRate" 
              stroke="hsl(217 91% 60%)" 
              strokeWidth={3}
              name="Completion Rate (%)"
              dot={{ r: 5, fill: 'hsl(217 91% 60%)', strokeWidth: 2, stroke: 'white' }}
              activeDot={{ r: 7 }}
            />
            <ReferenceLine 
              yAxisId="right"
              y={avgCompletionRate} 
              stroke="hsl(var(--muted-foreground))" 
              strokeDasharray="5 5"
              label={{ value: 'Average', position: 'right', fill: 'hsl(var(--muted-foreground))' }}
            />
          </ComposedChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  )
}

interface DifficultyPerformanceChartProps {
  data: Array<{
    level: number
    total: number
    completed: number
    completionRate: number
    avgTimeSeconds: number
    avgErrors: number
  }>
}

export function DifficultyPerformanceChart({
  data,
}: DifficultyPerformanceChartProps) {
  if (!data || data.length === 0) {
    return (
      <Card className="bg-card border-2 border-default rounded-4xl shadow-lg">
        <CardContent className="flex items-center justify-center h-64 text-muted-foreground">
          No data available
        </CardContent>
      </Card>
    )
  }

  const chartData = data.map((item) => ({
    ...item,
    name: `Level ${item.level}`,
  }))

  const avgCompletionRate = data.reduce((sum, item) => sum + item.completionRate, 0) / data.length
  const avgErrors = data.reduce((sum, item) => sum + item.avgErrors, 0) / data.length

  const CustomTooltip = ({ active, payload }: any) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload
      return (
        <div className="bg-card border-2 border-primary/20 rounded-xl p-4 shadow-xl">
          <p className="font-semibold text-primary mb-2">{data.name}</p>
          <div className="space-y-1 text-sm">
            <p className="flex justify-between gap-4">
              <span className="text-muted-foreground">Total Attempts:</span>
              <span className="font-semibold">{data.total}</span>
            </p>
            <p className="flex justify-between gap-4">
              <span className="text-muted-foreground">Completed:</span>
              <span className="font-semibold text-green-600">{data.completed}</span>
            </p>
            <p className="flex justify-between gap-4">
              <span className="text-muted-foreground">Completion Rate:</span>
              <span className="font-semibold text-blue-600">{data.completionRate.toFixed(1)}%</span>
            </p>
            <p className="flex justify-between gap-4">
              <span className="text-muted-foreground">Avg Time:</span>
              <span className="font-semibold">{data.avgTimeSeconds}s</span>
            </p>
            <p className="flex justify-between gap-4">
              <span className="text-muted-foreground">Avg Errors:</span>
              <span className="font-semibold text-red-600">{data.avgErrors}</span>
            </p>
          </div>
        </div>
      )
    }
    return null
  }

  return (
    <Card className="bg-card border-2 border-default rounded-4xl shadow-lg hover:shadow-xl transition-shadow">
      <CardHeader className="pb-4">
        <div className="flex items-center justify-between">
          <CardTitle className="text-xl font-bold text-primary">Performance by Difficulty Level</CardTitle>
          <div className="text-sm text-muted-foreground space-x-4">
            <span>
              Avg Rate: <span className="font-semibold text-primary">{avgCompletionRate.toFixed(1)}%</span>
            </span>
            <span>
              Avg Errors: <span className="font-semibold text-primary">{avgErrors.toFixed(1)}</span>
            </span>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <ResponsiveContainer width="100%" height={400}>
          <ComposedChart data={chartData} margin={{ top: 10, right: 20, left: 0, bottom: 5 }}>
            <CartesianGrid strokeDasharray="3 3" stroke="hsl(var(--border))" opacity={0.3} />
            <XAxis 
              dataKey="name" 
              stroke="hsl(var(--muted-foreground))"
              tick={{ fill: 'hsl(var(--muted-foreground))', fontSize: 12 }}
            />
            <YAxis 
              yAxisId="left"
              stroke="hsl(var(--muted-foreground))"
              tick={{ fill: 'hsl(var(--muted-foreground))', fontSize: 12 }}
              label={{ value: 'Errors', angle: -90, position: 'insideLeft', style: { textAnchor: 'middle' } }}
            />
            <YAxis 
              yAxisId="right"
              orientation="right"
              stroke="hsl(var(--muted-foreground))"
              tick={{ fill: 'hsl(var(--muted-foreground))', fontSize: 12 }}
              label={{ value: 'Rate (%)', angle: 90, position: 'insideRight', style: { textAnchor: 'middle' } }}
            />
            <Tooltip content={<CustomTooltip />} />
            <Legend 
              wrapperStyle={{ paddingTop: '20px' }}
              iconType="line"
            />
            <Bar 
              yAxisId="left"
              dataKey="avgErrors" 
              fill="hsl(0 84% 60%)" 
              name="Avg Errors"
              radius={[8, 8, 0, 0]}
              opacity={0.7}
            />
            <Line 
              yAxisId="right"
              type="monotone" 
              dataKey="completionRate" 
              stroke="hsl(217 91% 60%)" 
              strokeWidth={3}
              name="Completion Rate (%)"
              dot={{ r: 6, fill: 'hsl(217 91% 60%)', strokeWidth: 2, stroke: 'white' }}
              activeDot={{ r: 8 }}
            />
            <Line 
              yAxisId="left"
              type="monotone" 
              dataKey="avgTimeSeconds" 
              stroke="hsl(25 95% 53%)" 
              strokeWidth={2}
              name="Avg Time (s)"
              dot={{ r: 4, fill: 'hsl(25 95% 53%)', strokeWidth: 2, stroke: 'white' }}
              strokeDasharray="5 5"
              opacity={0.7}
            />
            <ReferenceLine 
              yAxisId="right"
              y={avgCompletionRate} 
              stroke="hsl(var(--muted-foreground))" 
              strokeDasharray="5 5"
              label={{ value: 'Avg Rate', position: 'right', fill: 'hsl(var(--muted-foreground))' }}
            />
            <ReferenceLine 
              yAxisId="left"
              y={avgErrors} 
              stroke="hsl(var(--muted-foreground))" 
              strokeDasharray="5 5"
              label={{ value: 'Avg Errors', position: 'left', fill: 'hsl(var(--muted-foreground))' }}
            />
          </ComposedChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  )
}