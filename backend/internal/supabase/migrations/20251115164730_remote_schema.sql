drop extension if exists "pg_net";

alter table "public"."resource" drop constraint "resource_week_check";


  create table "public"."verification_codes" (
    "id" uuid not null default gen_random_uuid(),
    "user_id" uuid not null,
    "code" text not null,
    "created_at" timestamp with time zone default now(),
    "expires_at" timestamp with time zone not null,
    "used" boolean default false,
    "attempts" integer default 0
      );


alter table "public"."verification_codes" enable row level security;

alter table "public"."resource" drop column "week";

alter table "public"."resource" add column "date" date;

CREATE INDEX idx_verification_codes_code ON public.verification_codes USING btree (code);

CREATE INDEX idx_verification_codes_user_id ON public.verification_codes USING btree (user_id);

CREATE UNIQUE INDEX verification_codes_pkey ON public.verification_codes USING btree (id);

alter table "public"."verification_codes" add constraint "verification_codes_pkey" PRIMARY KEY using index "verification_codes_pkey";

alter table "public"."verification_codes" add constraint "verification_codes_user_id_fkey" FOREIGN KEY (user_id) REFERENCES auth.users(id) ON DELETE CASCADE not valid;

alter table "public"."verification_codes" validate constraint "verification_codes_user_id_fkey";

set check_function_bodies = off;

CREATE OR REPLACE FUNCTION public.log_table_access()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        INSERT INTO Audit_Log(access_type, previous_data, new_data, time_of_operation, person, table_accessed)
        VALUES ('CREATE', NULL, row_to_json(NEW), now(), NULL, TG_TABLE_NAME);
        RETURN NEW;
    ELSIF (TG_OP = 'UPDATE') THEN
        INSERT INTO Audit_Log(access_type, previous_data, new_data, time_of_operation, person, table_accessed)
        VALUES ('UPDATE', row_to_json(OLD), row_to_json(NEW), now(), NULL, TG_TABLE_NAME);
        RETURN NEW;
    ELSIF (TG_OP = 'DELETE') THEN
        INSERT INTO Audit_Log(access_type, previous_data, new_data, time_of_operation, person, table_accessed)
        VALUES ('DELETE', row_to_json(OLD), NULL, now(), NULL, TG_TABLE_NAME);
        RETURN OLD;
    END IF;
    RETURN NULL;
    END;
    $function$
;

grant delete on table "public"."verification_codes" to "anon";

grant insert on table "public"."verification_codes" to "anon";

grant references on table "public"."verification_codes" to "anon";

grant select on table "public"."verification_codes" to "anon";

grant trigger on table "public"."verification_codes" to "anon";

grant truncate on table "public"."verification_codes" to "anon";

grant update on table "public"."verification_codes" to "anon";

grant delete on table "public"."verification_codes" to "authenticated";

grant insert on table "public"."verification_codes" to "authenticated";

grant references on table "public"."verification_codes" to "authenticated";

grant select on table "public"."verification_codes" to "authenticated";

grant trigger on table "public"."verification_codes" to "authenticated";

grant truncate on table "public"."verification_codes" to "authenticated";

grant update on table "public"."verification_codes" to "authenticated";

grant delete on table "public"."verification_codes" to "service_role";

grant insert on table "public"."verification_codes" to "service_role";

grant references on table "public"."verification_codes" to "service_role";

grant select on table "public"."verification_codes" to "service_role";

grant trigger on table "public"."verification_codes" to "service_role";

grant truncate on table "public"."verification_codes" to "service_role";

grant update on table "public"."verification_codes" to "service_role";


  create policy "Users can delete their own verification codes"
  on "public"."verification_codes"
  as permissive
  for delete
  to authenticated
using ((auth.uid() = user_id));



  create policy "Users can insert their own verification codes"
  on "public"."verification_codes"
  as permissive
  for insert
  to authenticated
with check ((auth.uid() = user_id));



  create policy "Users can read their own verification codes"
  on "public"."verification_codes"
  as permissive
  for select
  to authenticated
using ((auth.uid() = user_id));



  create policy "Users can view own verification codes"
  on "public"."verification_codes"
  as permissive
  for select
  to public
using ((auth.uid() = user_id));



