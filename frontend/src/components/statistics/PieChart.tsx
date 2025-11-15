import { Cell, Pie, PieChart, ResponsiveContainer } from 'recharts'

interface CustomPieChartProps {
  percentage: number
  title?: string
  color?: string
  className?: string
  size?: number
}

export default function CustomPieChart({
  percentage,
  title,
  color = 'var(--color-orange)',
  className = '',
  size = 120,
}: CustomPieChartProps) {
    const data = [
        { name: 'Completed', value: percentage },
        { name: 'Remaining', value: 100 - percentage },
    ]
    return (
        <div className={`flex items-center gap-6 ${className}`}>
        {/* Donut Chart */}
        <div style={{ width: size, height: size }}>
            <ResponsiveContainer width="100%" height="100%">
            <PieChart>
                <Pie
                data={data}
                cx="50%"
                cy="50%"
                innerRadius={size * 0.35}
                outerRadius={size * 0.45}
                dataKey="value"
                startAngle={90}
                endAngle={450}
                strokeWidth={0}
                >
                <Cell fill={color} />
                <Cell fill={'var(--color-white-hover)'} />
                </Pie>
            </PieChart>
            </ResponsiveContainer>
        </div>

        {/* Text Content */}
        <div className="flex flex-col justify-center">
            {title && (
            <div className="text-xl font-normal text-primary leading-tight mb-1">
                {title}
            </div>
            )}
            <div className="text-4xl font-bold text-primary">
            {percentage}%
            </div>
        </div>
        </div>
    )
}
