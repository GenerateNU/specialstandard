interface SchoolTagProps {
  schoolName: string
}
export default function SchoolTag({schoolName}: SchoolTagProps) {
  // Hash function to generate consistent color based on school name
  const getSchoolColor = (name: string) => {
    let hash = 0
    for (let i = 0; i < name.length; i++) {
      hash = name.charCodeAt(i) + ((hash << 5) - hash)
    }
    
    // Get index from 0-2 for three colors
    const colorIndex = Math.abs(hash) % 3
    
    const colors = [
      'bg-blue text-black',
      'bg-orange text-black',
      'bg-pink text-black',
    ]
    
    return colors[colorIndex]
  }

  const colorClass = getSchoolColor(schoolName)

  return (
    <span className={`p-2 rounded-full text-base font-medium whitespace-nowrap ${colorClass}`}>
      {schoolName}
    </span>
  )
}
