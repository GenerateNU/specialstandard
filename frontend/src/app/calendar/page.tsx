"use client";

import { useSessions } from "@/hooks/useSessions";
import moment from "moment";
import { useState } from "react";
import { Calendar, momentLocalizer, View } from "react-big-calendar";
import "react-big-calendar/lib/css/react-big-calendar.css";
import "./override-calendar.css";

const localizer = momentLocalizer(moment);

// Here, we are defining a calendar event type
interface CalendarEvent {
  id: string;
  title: string;
  start: Date;
  end: Date;
}

export default function MyCalendar() {
  const { sessions, isLoading, error } = useSessions();
  const [date, setDate] = useState(new Date());
  const [view, setView] = useState<View>("week");

  // Transform sessions into calendar events
  const events: CalendarEvent[] = sessions.map((session) => ({
    id: session.id,
    title: "session",
    start: new Date(session.start_datetime),
    end: new Date(session.end_datetime),
  }));

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div>Loading sessions...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div>Error loading sessions: {error}</div>
      </div>
    );
  }

  return (
    <div className="flex items-center justify-center h-screen">
      <Calendar
        localizer={localizer}
        events={events}
        startAccessor="start"
        endAccessor="end"
        style={{ height: "80vh", width: "90vw" }}
        date={date}
        view={view}
        onNavigate={setDate}
        onView={setView}
        views={["week", "day", "month"]}
      />
    </div>
  );
}
