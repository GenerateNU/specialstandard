import { Button } from "../ui/button"

// Reusable ResourceButton component
export function ResourceButton({ resource, icon: Icon }: { resource: any, icon: any }) {
  const buttonContent = (
    <>
      <Icon className='!w-6 !h-6 flex-shrink-0' strokeWidth={2}/>
      <div className='flex flex-col items-start'>
        <p>{resource.title || 'Untitled'}</p>
      </div>
    </>
  )

  const buttonClassName = 'flex flex-row items-center justify-baseline gap-4 px-5 py-8  rounded-2xl shadow-sm border border-border'

  return (
    <Button
      variant='secondary'
      className={buttonClassName}
      asChild={!!resource.presigned_url}
    >
      {resource.presigned_url ? (
        <a href={resource.presigned_url} target='_blank' rel='noopener noreferrer'>
          {buttonContent}
        </a>
      ) : (
        buttonContent
      )}
    </Button>
  )
}