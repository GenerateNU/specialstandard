import { createPagesServerClient } from "@supabase/auth-helpers-nextjs";
import { createClient } from "@supabase/supabase-js";
import type { NextApiRequest, NextApiResponse } from "next";
import process from "node:process";
import { Resend } from "resend";

const resend = new Resend(process.env.RESEND_API_KEY!);

const getEmailTemplate = (code: string) => `
  <div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
    <h2>Your verification code is:</h2>
    <div style="background: #f4f4f4; padding: 20px; text-align: center; border-radius: 8px;">
      <h1 style="font-size: 36px; letter-spacing: 8px; margin: 0;">${code}</h1>
    </div>
    <p style="margin-top: 20px;">This code will expire in 10 minutes.</p>
    <p style="color: #666;">If you didn't request this code, please ignore this email.</p>
  </div>
`;

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse
) {
  if (req.method !== "POST") {
    return res.status(405).json({ error: "Method not allowed" });
  }

  try {
    const authHeader = req.headers.authorization;
    const token = authHeader?.replace("Bearer ", "");

    if (!token) {
      return res.status(401).json({
        success: false,
        error: "No authentication token provided",
      });
    }

    // Client for auth verification (respects RLS)
    const supabaseAuth = createPagesServerClient({ req, res });

    // Verify the user
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

    const userEmail = user.email;

    // Create service role client for database insert (bypasses RLS)
    const supabaseAdmin = createClient(
      process.env.NEXT_PUBLIC_SUPABASE_URL!,
      process.env.SUPABASE_SERVICE_ROLE_KEY! // This bypasses RLS
    );

    // Generating random code for email verif
    const code = Math.floor(100000 + Math.random() * 900000).toString();

    const expiresAt = new Date();
    expiresAt.setMinutes(expiresAt.getMinutes() + 10);

    // Use admin client to insert
    const { error: dbError } = await supabaseAdmin
      .from("verification_codes")
      .insert({
        user_id: user.id,
        code,
        expires_at: expiresAt.toISOString(),
      });

    if (dbError) {
      console.error("Database error:", dbError);
      return res.status(500).json({
        success: false,
        error: "Failed to store verification code",
      });
    }
    const { data, error } = await resend.emails.send({
      from: "Kevin Matula <matulakevin91@gmail.com>",
      to: userEmail!,
      subject: "The Special Standard Verification Code",
      html: getEmailTemplate(code),
    });

    if (error) {
      console.error("Resend error:", error);
      return res.status(500).json({
        success: false,
        error: "Failed to send email",
      });
    }

    return res.status(200).json({
      success: true,
      messageId: data?.id,
    });
  } catch (error) {
    console.error("Unexpected error:", error);
    return res.status(500).json({
      success: false,
      error: "An unexpected error occurred",
    });
  }
}
