import fs from "fs";
import path from "path";

const API_BASE = "https://ni3bvlk951.execute-api.ap-south-1.amazonaws.com";

const textCases = [
  {
    id: "live-delivery-fee",
    category: "delivery",
    expectedLevel: "high",
    country: "US",
    text: "USPS: Your package is on hold due to $2.35 unpaid postage fee. Click http://usps-redelivery-fee.xyz to pay.",
  },
  {
    id: "live-bank-suspended",
    category: "bank",
    expectedLevel: "high",
    country: "US",
    text: "CHASE ALERT: Your debit card has been suspended due to unusual activity. Verify immediately at https://bit.ly/chase-unblock-now",
  },
  {
    id: "live-otp-demand",
    category: "bank",
    expectedLevel: "high",
    country: "US",
    text: "This is bank security agent Davis. Tell me the 6-digit code sent to your phone immediately to stop the fraudulent charge.",
  },
  {
    id: "live-family-claim",
    category: "relative in need",
    expectedLevel: "high",
    country: "US",
    text: "Hi Mom, I dropped my phone in the sink. Need $450 urgently for groceries from this new number.",
  },
  {
    id: "live-prize-lottery",
    category: "prize",
    expectedLevel: "high",
    country: "US",
    text: "Congratulations! You won $50,000 in the Mega Lottery. Pay $150 processing fee via Target gift card to claim.",
  },
  {
    id: "live-tech-support",
    category: "tech support",
    expectedLevel: "high",
    country: "US",
    text: "Windows Defender: Trojan spyware detected on your computer! Call Microsoft support immediately at 1-800-555-0199 to clean device.",
  },
  {
    id: "live-digital-arrest",
    category: "tax",
    expectedLevel: "high",
    country: "IN",
    text: "This is CBI narcotics wing. A courier containing contraband in your name was seized. Digital arrest warrant issued, remain on video call.",
  },
  {
    id: "live-crypto-scheme",
    category: "other",
    expectedLevel: "high",
    country: "US",
    text: "Guaranteed 20% daily returns with automated cryptocurrency arbitrage algorithm. Join VIP telegram group to multiply savings.",
  },
  {
    id: "live-romance-transfer",
    category: "romance",
    expectedLevel: "medium",
    country: "US",
    text: "My love, my medical trunk is held up at airport customs, please wire $800 to release it so I can fly to meet you.",
  },
  {
    id: "live-fake-job",
    category: "other",
    expectedLevel: "high",
    country: "US",
    text: "Work from home rating hotels for $300 daily. Pay refundable equipment registration fee of $50 to start today.",
  },
  {
    id: "live-utility-shutoff",
    category: "utility",
    expectedLevel: "high",
    country: "US",
    text: "Urgent: Electricity shutoff scheduled today at 5pm due to unpaid balance. Call billing department immediately at 800-555-0144.",
  },
  {
    id: "live-upi-qr-cashback",
    category: "bank",
    expectedLevel: "high",
    country: "IN",
    text: "Scan this QR code and enter secret UPI PIN to receive 5000 rupees cashback instantly into your bank account.",
  },
  {
    id: "live-damage-control",
    category: "bank",
    expectedLevel: "high",
    country: "US",
    text: "I panicked and already typed my debit card numbers and pin into the link they sent.",
  },
  {
    id: "live-safe-2fa-sms",
    category: "other",
    expectedLevel: "low",
    country: "US",
    text: "Your Google verification code is 842190. Do not share this code with anyone.",
  },
  {
    id: "live-safe-appointment",
    category: "other",
    expectedLevel: "low",
    country: "US",
    text: "Reminder: Your dental cleaning is scheduled for Tuesday at 2:00 PM with Dr. Smith. Reply C to confirm.",
  },
];

const screenshotCases = [
  {
    id: "live-shot-bank-chase",
    file: "testdata/screenshots/screen-02-bank-sms-dark.png",
    category: "bank",
    expectedLevel: "high",
    country: "US",
    ocrText: "Chase Alert: Your debit card has been suspended due to suspicious activity. Verify identity immediately at https://chase-security-verify.com",
  },
  {
    id: "live-shot-delivery-usps",
    file: "testdata/screenshots/screen-01-delivery-sms-light.png",
    category: "delivery",
    expectedLevel: "high",
    country: "US",
    ocrText: "USPS: Your package is on hold due to an unpaid redelivery fee of $1.99. Pay here: https://bit.ly/usps-fee-update",
  },
  {
    id: "live-shot-digital-arrest",
    file: "testdata/screenshots/screen-12-digital-arrest-whatsapp-light.png",
    category: "tax",
    expectedLevel: "high",
    country: "IN",
    ocrText: "Customs narcotics wing Delhi: Parcel seized under your Aadhaar. Digital arrest warrant active. Connect on Skype video call immediately.",
  },
];

async function presignUpload(filename) {
  const res = await fetch(`${API_BASE}/api/presign`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ content_type: "image/png", filename }),
  });
  if (!res.ok) throw new Error(`presign failed: ${res.status}`);
  return await res.json();
}

async function uploadFile(uploadUrl, filePath) {
  const fileBytes = fs.readFileSync(filePath);
  const res = await fetch(uploadUrl, {
    method: "PUT",
    headers: { "Content-Type": "image/png" },
    body: fileBytes,
  });
  if (!res.ok) throw new Error(`upload failed: ${res.status}`);
}

async function submitIntake(payload) {
  const res = await fetch(`${API_BASE}/api/intake`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });
  if (!res.ok) throw new Error(`intake failed: ${res.status}`);
  return await res.json();
}

async function pollIncident(incidentId, maxAttempts = 15) {
  for (let i = 0; i < maxAttempts; i++) {
    await new Promise((r) => setTimeout(r, 1200));
    const res = await fetch(`${API_BASE}/api/incidents?id=${incidentId}`);
    if (res.ok) {
      const data = await res.json();
      if (data.status === "analyzed" || data.status === "ANALYZED") {
        return data;
      }
    }
  }
  throw new Error(`timeout polling incident ${incidentId}`);
}

async function run() {
  console.log("=== SENTRY LIVE VERIFICATION ===");
  console.log(`Backend API: ${API_BASE}\n`);

  const results = [];

  // Run 15 text cases
  console.log("--- 15 Live Text Cases ---");
  for (const tc of textCases) {
    const t0 = Date.now();
    try {
      const intake = await submitIntake({
        senior_id: "live-eval-agent",
        user_context: tc.text,
        country: tc.country,
      });
      const incident = await pollIncident(intake.incident_id);
      const duration = Date.now() - t0;
      const level = (incident.risk_level || incident.severity || "").toLowerCase();
      const cat = (incident.category || "").toLowerCase();
      const output = incident.output || {};

      const levelMatch = level === tc.expectedLevel;
      results.push({
        id: tc.id,
        type: "text",
        expectedLevel: tc.expectedLevel,
        actualLevel: level,
        expectedCat: tc.category,
        actualCat: cat,
        pass: levelMatch,
        duration,
        output,
      });

      console.log(
        `[${levelMatch ? "PASS" : "FAIL"}] ${tc.id.padEnd(25)} | level: ${level.padEnd(6)} (want ${tc.expectedLevel.padEnd(6)}) | cat: ${cat.padEnd(16)} | ${duration}ms`
      );
    } catch (err) {
      console.log(`[ERR ] ${tc.id.padEnd(25)} | error: ${err.message}`);
      results.push({ id: tc.id, type: "text", pass: false, error: err.message });
    }
  }

  // Run 3 screenshot cases
  console.log("\n--- 3 Live Screenshot Cases ---");
  for (const sc of screenshotCases) {
    const t0 = Date.now();
    try {
      const presign = await presignUpload(path.basename(sc.file));
      await uploadFile(presign.upload_url, sc.file);
      const intake = await submitIntake({
        senior_id: "live-eval-agent",
        image_key: presign.image_key,
        user_context: sc.ocrText,
        country: sc.country,
      });
      const incident = await pollIncident(intake.incident_id);
      const duration = Date.now() - t0;
      const level = (incident.risk_level || incident.severity || "").toLowerCase();
      const cat = (incident.category || "").toLowerCase();
      const output = incident.output || {};

      const levelMatch = level === sc.expectedLevel;
      results.push({
        id: sc.id,
        type: "screenshot",
        expectedLevel: sc.expectedLevel,
        actualLevel: level,
        expectedCat: sc.category,
        actualCat: cat,
        pass: levelMatch,
        duration,
        output,
      });

      console.log(
        `[${levelMatch ? "PASS" : "FAIL"}] ${sc.id.padEnd(25)} | level: ${level.padEnd(6)} (want ${sc.expectedLevel.padEnd(6)}) | cat: ${cat.padEnd(16)} | ${duration}ms`
      );
    } catch (err) {
      console.log(`[ERR ] ${sc.id.padEnd(25)} | error: ${err.message}`);
      results.push({ id: sc.id, type: "screenshot", pass: false, error: err.message });
    }
  }

  const passed = results.filter((r) => r.pass).length;
  console.log(`\nLive Verification Summary: ${passed}/${results.length} passed (${((passed / results.length) * 100).toFixed(1)}%)`);

  fs.writeFileSync("testdata/live_verification_results.json", JSON.stringify(results, null, 2));
}

run().catch((e) => {
  console.error(e);
  process.exit(1);
});
