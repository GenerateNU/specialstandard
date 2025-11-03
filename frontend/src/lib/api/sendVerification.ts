import { createRouteHandlerClient } from "@supabase/auth-helpers-nextjs";
import { cookies } from "next/headers";
import type { NextRequest } from "next/server";
import { NextResponse } from "next/server";
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

// This block gets the users email from the cookies
export async function POST(_request: NextRequest) {
  try {
    const supabase = createRouteHandlerClient({ cookies });

    const {
      data: { user },
      error: userError,
    } = await supabase.auth.getUser();

    if (userError || !user) {
      return NextResponse.json(
        { success: false, error: "Not authenticated" },
        { status: 401 }
      );
    }

    const userEmail = user.email;

    // Generating random code for email verif
    const code = Math.floor(100000 + Math.random() * 900000).toString();

    // Place the code in database
    const expiresAt = new Date();
    expiresAt.setMinutes(expiresAt.getMinutes() + 10);

    const { error: dbError } = await supabase
      .from("verification_codes")
      .insert({
        user_id: user.id,
        code,
        expires_at: expiresAt.toISOString(),
      });

    if (dbError) {
      console.error("Database error:", dbError);
      return NextResponse.json(
        { success: false, error: "Failed to store verification code" },
        { status: 500 }
      );
    }

    const { data, error } = await resend.emails.send({
      from: "Kevin Matula <matulakevin91@gmail.com>",
      to: userEmail!,
      subject: "The Special Standard Verification Code",
      html: getEmailTemplate(code),
    });

    if (error) {
      console.error("Resend error:", error);
      return NextResponse.json(
        { success: false, error: "Failed to send email" },
        { status: 500 }
      );
    }

    return NextResponse.json({
      success: true,
      messageId: data?.id,
    });
  } catch (error) {
    console.error("Unexpected error:", error);
    return NextResponse.json(
      { success: false, error: "An unexpected error occurred" },
      { status: 500 }
    );
  }
}
