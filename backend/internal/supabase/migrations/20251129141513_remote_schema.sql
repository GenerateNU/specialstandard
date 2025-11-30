alter table "public"."game_content" drop column "applicable_game_types";

alter table "public"."game_content" drop column "exercise_type";

alter table "public"."resource" drop column "date";

alter table "public"."resource" add column "week" integer not null default 1;

alter table "public"."session" alter column "session_parent_id" drop not null;

drop type "public"."exercise_type";

drop type "public"."game_type";

alter table "public"."resource" add constraint "resource_week_check" CHECK (((week >= 1) AND (week <= 4))) not valid;

alter table "public"."resource" validate constraint "resource_week_check";


