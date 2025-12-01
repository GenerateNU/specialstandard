import { CSS } from '@dnd-kit/utilities'
import { useSortable } from '@dnd-kit/sortable'


interface DraggableImageProps {
  id: string
  filename: string
  url: string
}


// Draggable image component
export default function DraggableImage({ id, filename, url }: DraggableImageProps) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id })

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  }

  return (
    <div
      ref={setNodeRef}
      style={style}
      {...attributes}
      {...listeners}
      className="bg-card rounded-lg shadow-md p-3 cursor-grab active:cursor-grabbing"
    >
      <div className="aspect-square rounded-lg overflow-hidden">
        <img 
          src={url} 
          alt={filename}
          className="w-full h-full object-contain"
          onError={(e) => {
            console.error('Failed to load image:', url)
            e.currentTarget.src = 'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" width="100" height="100"%3E%3Crect fill="%23ddd" width="100" height="100"/%3E%3Ctext x="50%25" y="50%25" dominant-baseline="middle" text-anchor="middle"%3EError%3C/text%3E%3C/svg%3E'
          }}
        />
      </div>
    </div>
  )
}

