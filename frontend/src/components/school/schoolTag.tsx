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
  
    
    const colors = [
      'bg-blue text-black',
      'bg-orange text-black',
      'bg-pink text-black',
      'bg-blue-light text-black',
      'bg-orange-light text-black',
      'bg-pink-light text-black',


    ]

    const colorIndex = Math.abs(hash) % colors.length
    
    return colors[colorIndex]
  }

  const colorClass = getSchoolColor(schoolName)

  return (
    <span className={`p-2 rounded-full text-base font-medium whitespace-nowrap ${colorClass}`}>
      {schoolName}
    </span>
  )
}
