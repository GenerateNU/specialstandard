import { useRouter } from 'next/navigation'
import { Avatar } from '@/components/ui/avatar'
import { cn, getAvatarVariant } from '@/lib/utils'

interface RecentStudentCardProps {
  id: string
  firstName: string
  lastName: string
  grade?: number | null
  className?: string
}

export default function RecentStudentCard({
  id,
  firstName,
  lastName,
  grade,
  className,
}: RecentStudentCardProps) {
  const router = useRouter()
  const fullName = `${firstName} ${lastName}`
  const avatarVariant = getAvatarVariant(id)

  return (
    <button
      onClick={() => router.push(`/student/${id}`)}
      className={cn(
        'w-full flex items-center gap-3 p-3 cursor-pointer rounded-lg hover:bg-pink-disabled transition-colors text-left',
        className,
      )}
    >
      {/* Avatar */}
      <Avatar
        name={fullName + id}
        variant={avatarVariant}
        className="w-13 h-13 shrink-0"
      />

      {/* Student Info */}
      <div className="flex-1 min-w-0">
        <p className="text-md font-bold truncate text-foreground">
          {firstName}
          {' '}
          {lastName}
        </p>
        {grade !== undefined && (
          <p className="text-sm text-muted-foreground">
            Grade
            {' '}
            {grade}
          </p>
        )}
      </div>
    </button>
  )
}
