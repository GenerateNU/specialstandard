import { createPagesServerClient } from "@supabase/auth-helpers-nextjs";
import { createClient } from "@supabase/supabase-js";
import type { NextApiRequest, NextApiResponse } from "next";
import process from "node:process";

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse
) {
  if (req.method !== "POST") {
    return res.status(405).json({ error: "Method not allowed" });
  }

  try {
    const { code } = req.body;

    if (!code || code.length !== 6) {
      return res.status(400).json({
        success: false,
        error: "Invalid code format",
      });
    }

    // Extract JWT from Authorization header
    const authHeader = req.headers.authorization;
    const token = authHeader?.replace("Bearer ", "");

    if (!token) {
      return res.status(401).json({
        success: false,
        error: "No authentication token provided",
      });
    }

    // Verify the user
    const supabaseAuth = createPagesServerClient({ req, res });
    const {
      data: { user },
      error: userError,
    } = await supabaseAuth.auth.getUser(token);

    if (userError || !user) {
      console.error("Auth error:", userError);
      return res.status(401).json({
        success: false,
        error: "Not authenticated",
      });
    }

    const supabaseAdmin = createClient(
      process.env.NEXT_PUBLIC_SUPABASE_URL!,
      process.env.SUPABASE_SERVICE_ROLE_KEY!
    );

    // Find the most recent valid code for this user
    const { data: verificationRecord, error: fetchError } = await supabaseAdmin
      .from("verification_codes")
      .select("*")
      .eq("user_id", user.id)
      .eq("code", code)
      .gt("expires_at", new Date().toISOString())
      .order("created_at", { ascending: false })
      .limit(1)
      .single();

    if (fetchError || !verificationRecord) {
      return res.status(400).json({
        success: false,
        verified: false,
        error: "Invalid or expired code",
      });
    }

    // Code is valid! Delete it so it can't be reused
    await supabaseAdmin
      .from("verification_codes")
      .delete()
      .eq("id", verificationRecord.id);

    return res.status(200).json({
      success: true,
      verified: true,
    });
  } catch (error) {
    console.error("Unexpected error:", error);
    return res.status(500).json({
      success: false,
      error: "An unexpected error occurred",
    });
  }
}
