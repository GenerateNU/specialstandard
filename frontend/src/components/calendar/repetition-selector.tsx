// components/repetition-selector.tsx

import React from "react";
import { ChevronDown, ChevronUp, Repeat } from "lucide-react";
import { Switch } from "@/components/ui/switch";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
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
  { label: "Sun", value: 0 },
  { label: "Mon", value: 1 },
  { label: "Tue", value: 2 },
  { label: "Wed", value: 3 },
  { label: "Thu", value: 4 },
  { label: "Fri", value: 5 },
  { label: "Sat", value: 6 },
];

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
    if (new Date(end) <= new Date(sessionDate)) return;

    const config: RepetitionConfig = {
      recur_start: new Date(`${sessionDate  }T${  sessionTime}`).toISOString(),
      recur_end: new Date(`${end  }T${  sessionTime}`).toISOString(),
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
    <div className="border rounded-lg overflow-hidden bg-white">
      {/* Header */}
      <button
        type="button"
        onClick={() => setIsExpanded(!isExpanded)}
        className="w-full flex items-center justify-between p-3 bg-slate-50 hover:bg-slate-100 transition-colors"
      >
        <label className="flex items-center gap-3 cursor-pointer flex-1">
          <Repeat className="w-4 h-4" />
          <div className="text-left">
            <p className="font-semibold text-sm">Recurring Session</p>
            {isRecurring && selectedDays.length > 0 && (
              <p className="text-xs text-gray-600">{getRecurrenceSummary()}</p>
            )}
          </div>
          <Switch
            checked={isRecurring}
            onCheckedChange={handleToggleRecurring}
            onClick={(e) => e.stopPropagation()}
            className={cn(
              isRecurring ? "bg-blue-500" : "bg-gray-300",
              "transition-colors"
            )}
          />
        </label>
        {isExpanded ? (
          <ChevronUp className="w-4 h-4 text-gray-500 ml-2" />
        ) : (
          <ChevronDown className="w-4 h-4 text-gray-500 ml-2" />
        )}
      </button>

      {/* Expanded Content */}
      {isExpanded && isRecurring && (
        <div className="p-4 space-y-4 border-t bg-white">
          {/* Every N Weeks */}
          <div className="space-y-2">
            <Label className="text-sm font-medium text-gray-700">
              Frequency
            </Label>
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
                className="w-16 px-2 py-1 border border-gray-300 rounded text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              <span className="text-sm">week{everyNWeeks !== 1 ? "s" : ""}</span>
            </div>
          </div>

          {/* Days of Week Selection */}
          <div className="space-y-2">
            <Label className="text-sm font-medium text-gray-700">
              Days of Week
            </Label>
            <div className="grid grid-cols-7 gap-2">
              {DAYS_OF_WEEK.map((day) => (
                <button
                  key={day.value}
                  type="button"
                  onClick={() => handleDayToggle(day.value)}
                  className={`py-2 px-1 text-xs font-semibold rounded transition-all border ${
                    selectedDays.includes(day.value)
                      ? "bg-blue-500 text-white border-blue-600 shadow-sm"
                      : "bg-white text-gray-700 border-gray-200 hover:border-gray-300"
                  }`}
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
          <div className="space-y-2">
            <Label htmlFor="endDate" className="text-sm font-medium text-gray-700">
              Recurrence End Date
            </Label>
            <Input
              id="endDate"
              type="date"
              value={endDate}
              onChange={(e) => setEndDate(e.target.value)}
              min={sessionDate}
              className="w-full px-3 py-2 border border-gray-300 rounded text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            {new Date(endDate) <= new Date(sessionDate) && (
              <p className="text-xs text-red-500">
                End date must be after the session date
              </p>
            )}
          </div>

          {/* Summary Card */}
          {selectedDays.length > 0 &&
            new Date(endDate) > new Date(sessionDate) && (
              <div className="p-3 bg-blue-50 border border-blue-200 rounded space-y-2">
                <p className="text-sm font-semibold text-blue-900">Summary:</p>
                <ul className="text-xs text-blue-800 space-y-1">
                  <li>
                    • First session: {new Date(sessionDate).toLocaleDateString()}
                  </li>
                  <li>
                    • Repeats every {everyNWeeks} week
                    {everyNWeeks !== 1 ? "s" : ""} on{" "}
                    {selectedDays
                      .map((d) => DAYS_OF_WEEK.find((day) => day.value === d)?.label)
                      .join(", ")}
                  </li>
                  <li>
                    • Until: {new Date(endDate).toLocaleDateString()}
                  </li>
                </ul>
              </div>
            )}
        </div>
      )}
    </div>
  );
}