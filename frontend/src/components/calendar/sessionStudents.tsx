import { Avatar } from '@/components/ui/avatar'
import { useSessionStudentsForSession } from '@/hooks/useSessionStudents'
import { getAvatarName, getAvatarVariant } from '@/lib/avatarUtils'

export default function SessionStudents({ sessionId }: { sessionId: string }) {
  const { students, isLoading } = useSessionStudentsForSession(sessionId)

  if (isLoading || students.length === 0) {
    return null
  }

  return (
    <div className="grid grid-cols-2 gap-1 px-1">
      {students.map((student, index) => {
        const avatarVariant = getAvatarVariant(student.id)
        const isLastAndOdd = index === students.length - 1 && students.length % 2 !== 0
        return (
          <div
            key={student.id}
            className={`flex flex-row justify-center items-center gap-1 ${isLastAndOdd ? 'col-span-2' : ''}`}
          >
            <Avatar
              name={getAvatarName(student.first_name, student.last_name, student.id)}
              variant={avatarVariant}
              className="w-8 h-8"
              title={`${student.first_name} ${student.last_name}`}
            />
            <div>
              {student.first_name.charAt(0).toUpperCase()}
              {student.last_name.charAt(0).toUpperCase()}
            </div>
          </div>
        )
      })}
    </div>
  )
}
