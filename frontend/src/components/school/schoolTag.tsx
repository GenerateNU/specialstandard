import { getSchoolColor } from '@/lib/utils'

interface SchoolTagProps {
  schoolName: string
}

export default function SchoolTag({schoolName}: SchoolTagProps) {
  const colorClass = getSchoolColor(schoolName)

  return (
    <span className={`p-2 rounded-full text-base font-medium whitespace-nowrap ${colorClass}`}>
      {schoolName}
    </span>
  )
}
