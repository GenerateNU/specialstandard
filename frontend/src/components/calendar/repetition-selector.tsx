import React from "react";
import { ChevronDown, ChevronUp, Repeat } from "lucide-react";
import { Switch } from "@/components/ui/switch";
import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";

export interface RepetitionConfig {
  recur_start: string; // ISO datetime string
  recur_end: string; // ISO datetime string
  every_n_weeks: number;
  days: number[]; // 0-6 representing Sun-Sat
}

interface RepetitionSelectorProps {
  value: RepetitionConfig | undefined;
  onChange: (config: RepetitionConfig | undefined) => void;
  sessionDate: string;
  sessionTime: string;
}

const DAYS_OF_WEEK = [
  { label: "Mon", value: 1 },
  { label: "Tue", value: 2 },
  { label: "Wed", value: 3 },
  { label: "Thu", value: 4 },
  { label: "Fri", value: 5 },
];

// Helper to create a date from YYYY-MM-DD string and HH:MM time string
// This avoids timezone parsing issues with ISO string format
function createLocalDate(dateStr: string, timeStr: string): Date {
  const [year, month, day] = dateStr.split('-').map(Number);
  const [hour, minute] = timeStr.split(':').map(Number);
  return new Date(year, month - 1, day, hour, minute, 0, 0);
}

// Find the first occurrence of any of the selected days on or after the given date
function getFirstOccurrence(dateStr: string, timeStr: string, selectedDays: number[]): Date {
  if (selectedDays.length === 0) {
    return createLocalDate(dateStr, timeStr);
  }

  const baseDate = createLocalDate(dateStr, timeStr);
  const baseWeekday = baseDate.getDay(); // 0 = Sunday, 1 = Monday, etc.
  
  // Sort days to find the earliest one
  const sortedDays = [...selectedDays].sort((a, b) => a - b);
  
  // Find the smallest positive offset to reach one of the selected days
  let minOffset = 7; // Start with a week (max offset)
  
  for (const targetDay of sortedDays) {
    let offset = targetDay - baseWeekday;
    if (offset < 0) {
      offset += 7; // Wrap to next week if the day has passed
    }
    if (offset < minOffset) {
      minOffset = offset;
    }
  }
  
  // If one of the selected days is today (offset = 0), use today
  // Otherwise, move to the nearest selected day
  const result = new Date(baseDate);
  result.setDate(result.getDate() + minOffset);
  return result;
}

export function RepetitionSelector({
  value,
  onChange,
  sessionDate,
  sessionTime,
}: RepetitionSelectorProps) {
  const [isExpanded, setIsExpanded] = React.useState(!!value);
  const [isRecurring, setIsRecurring] = React.useState(!!value);
  const [endDate, setEndDate] = React.useState(
    value?.recur_end.split("T")[0] || sessionDate
  );
  const [everyNWeeks, setEveryNWeeks] = React.useState(
    value?.every_n_weeks || 1
  );
  const [selectedDays, setSelectedDays] = React.useState<number[]>(
    value?.days || [new Date(sessionDate).getDay()]
  );

  const handleDayToggle = (day: number) => {
    const newDays = selectedDays.includes(day)
      ? selectedDays.filter((d) => d !== day)
      : [...selectedDays, day];
    setSelectedDays(newDays);
  };

  const updateRepetition = (
    days: number[],
    weeks: number,
    end: string
  ) => {
    if (days.length === 0) return;
    
    const endDateTime = createLocalDate(end, sessionTime);
    const sessionDateTime = createLocalDate(sessionDate, sessionTime);
    
    if (endDateTime <= sessionDateTime) return;

    // Calculate the first occurrence date based on selected days
    // This ensures the first session is on the selected weekday, not the original date
    const firstOccurrence = getFirstOccurrence(sessionDate, sessionTime, days);

    const config: RepetitionConfig = {
      recur_start: firstOccurrence.toISOString(),
      recur_end: endDateTime.toISOString(),
      every_n_weeks: weeks,
      days: days.sort((a, b) => a - b),
    };

    onChange(config);
  };

  const handleToggleRecurring = (checked: boolean) => {
    setIsRecurring(checked);
    setIsExpanded(checked);

    if (!checked) {
      onChange(undefined);
      return;
    }

    // Auto-apply when toggling on
    updateRepetition(selectedDays, everyNWeeks, endDate);
  };

  React.useEffect(() => {
    if (isRecurring) {
      updateRepetition(selectedDays, everyNWeeks, endDate);
    }
  }, [selectedDays, everyNWeeks, endDate, isRecurring]);

  const getRecurrenceSummary = () => {
    if (!isRecurring || selectedDays.length === 0) return "";

    const dayNames = selectedDays
      .map((d) => DAYS_OF_WEEK.find((day) => day.value === d)?.label)
      .join(", ");

    const endDateFormatted = new Date(endDate).toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric",
    });

    return `Every ${everyNWeeks} week${everyNWeeks !== 1 ? "s" : ""} on ${dayNames} until ${endDateFormatted}`;
  };

  return (
    <div className="border border-gray-300 rounded-lg overflow-hidden">
      {/* Header */}
      <button
        type="button"
        onClick={() => setIsExpanded(!isExpanded)}
        className="w-full flex items-center justify-between p-3 hover:bg-gray-50 border border-gray-300 transition-colors"
      >
        <label className="flex items-center gap-3 cursor-pointer flex-1">
          <Repeat className="w-4 h-4 text-secondary" />
          <div className="text-left">
            <p className="text-sm font-medium text-primary">Recurring Session</p>
            {isRecurring && selectedDays.length > 0 && (
              <p className="text-xs text-secondary mt-1">{getRecurrenceSummary()}</p>
            )}
          </div>
          <Switch
            checked={isRecurring}
            onCheckedChange={handleToggleRecurring}
            onClick={(e) => e.stopPropagation()}
            className={cn(
              "data-[state=checked]:bg-primary",
              "data-[state=unchecked]:bg-gray-300",
              "data-[state=unchecked]:border-gray-400"
            )}          />
        </label>
        {isExpanded ? (
          <ChevronUp className="w-4 h-4 text-secondary ml-2" />
        ) : (
          <ChevronDown className="w-4 h-4 text-secondary ml-2" />
        )}
      </button>

      {/* Expanded Content */}
      {isExpanded && isRecurring && (
        <div className="p-4 space-y-4 border-t">
          {/* Every N Weeks */}
          <div>
            <p className="text-xs text-secondary mb-1">Frequency</p>
            <div className="flex items-center gap-2">
              <span className="text-sm">Every</span>
              <Input
                type="number"
                min="1"
                max="12"
                value={everyNWeeks}
                onChange={(e) => {
                  const val = Math.max(1, Number.parseInt(e.target.value) || 1);
                  setEveryNWeeks(val);
                }}
                className="w-16 border border-gray-300"
              />
              <span className="text-sm">week{everyNWeeks !== 1 ? "s" : ""}</span>
            </div>
          </div>

          {/* Days of Week Selection */}
          <div>
            <p className="text-xs text-secondary mb-2">Days of Week</p>
            <div className="grid grid-cols-5 gap-1">
              {DAYS_OF_WEEK.map((day) => (
                <button
                  key={day.value}
                  type="button"
                  onClick={() => handleDayToggle(day.value)}
                  className={cn(
                    "py-2 px-1 text-xs font-medium rounded transition-all",
                    selectedDays.includes(day.value)
                      ? "bg-primary text-white"
                      : "bg-gray-50 text-secondary hover:bg-gray-100"
                  )}
                >
                  {day.label}
                </button>
              ))}
            </div>
            {selectedDays.length === 0 && (
              <p className="text-xs text-red-500 mt-1">
                Please select at least one day
              </p>
            )}
          </div>

          {/* End Date */}
          <div>
            <p className="text-xs text-secondary mb-1">Recurrence End Date</p>
            <Input
              type="date"
              value={endDate}
              onChange={(e) => setEndDate(e.target.value)}
              min={sessionDate}
              className="border border-gray-300"
            />
            {createLocalDate(endDate, sessionTime) <= createLocalDate(sessionDate, sessionTime) && (
              <p className="text-xs text-red-500 mt-1">
                End date must be after the session date
              </p>
            )}
          </div>

          {/* Summary Card */}
          {selectedDays.length > 0 &&
            createLocalDate(endDate, sessionTime) > createLocalDate(sessionDate, sessionTime) && (
              <div className="p-3 bg-gray-50 border border-gray-200 rounded">
                <p className="text-sm font-medium text-primary mb-2">Summary:</p>
                <ul className="text-xs text-secondary space-y-1">
                  <li>
                    • First session: {getFirstOccurrence(sessionDate, sessionTime, selectedDays).toLocaleDateString()}
                  </li>
                  <li>
                    • Repeats every {everyNWeeks} week
                    {everyNWeeks !== 1 ? "s" : ""} on{" "}
                    {selectedDays
                      .map((d) => DAYS_OF_WEEK.find((day) => day.value === d)?.label)
                      .join(", ")}
                  </li>
                  <li>
                    • Until: {createLocalDate(endDate, sessionTime).toLocaleDateString()}
                  </li>
                </ul>
              </div>
            )}
        </div>
      )}
    </div>
  );
}