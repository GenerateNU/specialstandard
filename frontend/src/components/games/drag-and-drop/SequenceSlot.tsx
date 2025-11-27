
import { useDroppable } from "@dnd-kit/core"

export interface SequenceSlotProps {
  position: number
  filename: string | null
  url?: string
}

// Droppable sequence slot
export default function SequenceSlot({ position, filename, url, slotId }: SequenceSlotProps & { slotId: string }) {
  const { setNodeRef, isOver } = useDroppable({
    id: slotId,
  })

  return (
    <div 
      ref={setNodeRef}
      className={`relative bg-card-hover border-2 border-dashed rounded-lg p-4 min-h-[150px] flex items-center justify-center transition-colors ${
        isOver ? 'border-blue bg-blue/10' : 'border-border'
      }`}
    >
      {filename && url ? (
        <div className="w-full h-full">
          <div className="aspect-square rounded-lg overflow-hidden">
            <img 
              src={url} 
              alt={`Position ${position}`}
              className="w-full h-full object-contain"
            />
          </div>
        </div>
      ) : (
        <span className="text-secondary text-sm">Drop here ({position})</span>
      )}
    </div>
  )
}