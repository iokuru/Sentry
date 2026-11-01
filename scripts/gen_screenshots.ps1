Add-Type -AssemblyName System.Drawing

$outDir = "testdata/screenshots"
if (!(Test-Path $outDir)) {
    New-Item -ItemType Directory -Path $outDir -Force | Out-Null
}

$specs = @(
    @{
        id = "screen-01-delivery-sms-light"
        type = "sms"
        theme = "light"
        sender = "USPS Tracking"
        time = "10:24 AM"
        text = "USPS: Your package is on hold due to an unpaid redelivery fee of $1.99. Pay here: https://bit.ly/usps-fee-update"
        blur = $false
        empty = $false
    },
    @{
        id = "screen-02-bank-sms-dark"
        type = "sms"
        theme = "dark"
        sender = "Chase Alert"
        time = "02:15 PM"
        text = "Chase Alert: Your debit card has been suspended due to suspicious activity. Verify identity immediately at https://chase-security-verify.com"
        blur = $false
        empty = $false
    },
    @{
        id = "screen-03-family-whatsapp-light"
        type = "whatsapp"
        theme = "light"
        sender = "+1 (555) 019-2831"
        time = "04:45 PM"
        text = "Hi Mom, I dropped my phone in water and this is my new number. I have an urgent bill due today, could you transfer $650 for me?"
        blur = $false
        empty = $false
    },
    @{
        id = "screen-04-upi-whatsapp-dark"
        type = "whatsapp"
        theme = "dark"
        sender = "Marketplace Buyer"
        time = "11:10 AM"
        text = "Scan this QR code on PhonePe or GooglePay and enter your secret UPI PIN to receive the Rs 15000 payment in your bank account."
        blur = $false
        empty = $false
    },
    @{
        id = "screen-05-geek-squad-email-light"
        type = "email"
        theme = "light"
        sender = "Geek Squad Billing <billing@service-support.net>"
        time = "09:00 AM"
        text = "Geek Squad Total Support auto renewed for $399.99 today. To dispute or cancel this charge, call our billing helpdesk at 1-800-555-0144."
        blur = $false
        empty = $false
    },
    @{
        id = "screen-06-norton-email-dark"
        type = "email"
        theme = "dark"
        sender = "Norton LifeLock <notice@security-renewals.com>"
        time = "01:30 PM"
        text = "Norton Antivirus 3-year subscription renewed for $449.00. Debited from your registered account. Call 1-888-555-0183 to cancel."
        blur = $false
        empty = $false
    },
    @{
        id = "screen-07-otp-sms-light"
        type = "sms"
        theme = "light"
        sender = "Security Desk"
        time = "03:12 PM"
        text = "A wire transfer cancellation was requested. Please share the 6-digit verification code you just received to reverse the charge."
        blur = $false
        empty = $false
    },
    @{
        id = "screen-08-prize-sms-dark"
        type = "sms"
        theme = "dark"
        sender = "Mega Lottery"
        time = "06:00 PM"
        text = "CONGRATULATIONS! You won $50,000 in the National Cash Draw! Pay $100 processing fee via Apple gift card to claim funds today."
        blur = $false
        empty = $false
    },
    @{
        id = "screen-09-techsupport-popup-light"
        type = "popup"
        theme = "light"
        sender = "Microsoft Windows Alert"
        time = "Alert"
        text = "VIRUS DETECTED: Spyware trojan found on your system. Do not restart computer. Call Windows Support immediately at 1-800-555-0199."
        blur = $false
        empty = $false
    },
    @{
        id = "screen-10-defender-popup-dark"
        type = "popup"
        theme = "dark"
        sender = "Windows Defender Security"
        time = "Security Warning"
        text = "CRITICAL BREACH: Unauthorized network access detected. Remote session required to clean system. Call 1-800-555-0122 now."
        blur = $false
        empty = $false
    },
    @{
        id = "screen-11-electricity-sms-light"
        type = "sms"
        theme = "light"
        sender = "Power Dept"
        time = "08:30 AM"
        text = "Electricity Disconnection Notice: Service scheduled for shutoff at 6 PM today due to unpaid bill. Call 1-800-555-0188 immediately."
        blur = $false
        empty = $false
    },
    @{
        id = "screen-12-digital-arrest-whatsapp-light"
        type = "whatsapp"
        theme = "light"
        sender = "Cyber Crime Branch"
        time = "01:15 PM"
        text = "Customs narcotics wing Delhi: Parcel seized under your Aadhaar. Digital arrest warrant active. Connect on Skype video call immediately."
        blur = $false
        empty = $false
    },
    @{
        id = "screen-13-crypto-sms-dark"
        type = "sms"
        theme = "dark"
        sender = "VIP Trading"
        time = "12:00 PM"
        text = "Join our automated Binance trading group. Guaranteed 15% daily return on Bitcoin. Deposit 500 USDT to start earning immediately."
        blur = $false
        empty = $false
    },
    @{
        id = "screen-14-romance-chat-light"
        type = "sms"
        theme = "light"
        sender = "Dr. James"
        time = "10:45 PM"
        text = "My love, my medical trunk is held at customs immigration. Please send $1,200 via wire transfer so I can get my flight to see you."
        blur = $false
        empty = $false
    },
    @{
        id = "screen-15-job-task-sms-light"
        type = "sms"
        theme = "light"
        sender = "Global Media"
        time = "11:20 AM"
        text = "Work from home part-time: Earn $300 to $500 daily evaluating products online. Pay $150 refundable training deposit to start tasks."
        blur = $false
        empty = $false
    },
    @{
        id = "screen-16-kyc-pan-sms-light"
        type = "sms"
        theme = "light"
        sender = "SBI Alert"
        time = "05:10 PM"
        text = "Dear Customer, your SBI account KYC expired today. Netbanking will be blocked. Update PAN now by downloading: bit.ly/sbi-pan-app"
        blur = $false
        empty = $false
    },
    @{
        id = "screen-17-refund-overpay-sms-dark"
        type = "sms"
        theme = "dark"
        sender = "Accounts Billing"
        time = "03:40 PM"
        text = "Overpayment notice: We accidentally credited $2,400 instead of $240 to your account. Wire back the excess $2,160 immediately."
        blur = $false
        empty = $false
    },
    @{
        id = "screen-18-craigslist-chat-light"
        type = "sms"
        theme = "light"
        sender = "Craigslist Buyer"
        time = "02:05 PM"
        text = "I mailed a cashier cheque for $1,200 for your dining table ($400). Please wire the $800 difference to my movers before tomorrow."
        blur = $false
        empty = $false
    },
    @{
        id = "screen-19-delivery-blurry-light"
        type = "sms"
        theme = "light"
        sender = "FedEx Express"
        time = "01:22 PM"
        text = "FedEx: Delivery attempted. $2.50 redelivery fee required. Confirm address and payment at https://bit.ly/fdx-parcel-hold"
        blur = $true
        empty = $false
    },
    @{
        id = "screen-20-bank-blurry-dark"
        type = "sms"
        theme = "dark"
        sender = "Wells Fargo"
        time = "04:18 PM"
        text = "Wells Fargo: Unusual debit card activity detected. Your card is temporarily locked. Verify now at bit.ly/wf-unlock-card"
        blur = $true
        empty = $false
    },
    @{
        id = "screen-21-family-blurry-light"
        type = "whatsapp"
        theme = "light"
        sender = "Daughter"
        time = "09:15 AM"
        text = "Hey Dad, my phone broke so texting from this friend number. I need $400 for emergency car repair please wire to my account."
        blur = $true
        empty = $false
    },
    @{
        id = "screen-22-amazon-blurry-light"
        type = "email"
        theme = "light"
        sender = "Amazon Support"
        time = "10:00 AM"
        text = "Amazon Prime Membership renewed for $499.00 today. If this was not authorized by you, call customer desk at 1-800-555-0166."
        blur = $true
        empty = $false
    },
    @{
        id = "screen-23-safe-chase-2fa-light"
        type = "sms"
        theme = "light"
        sender = "Chase Shortcode 24273"
        time = "12:30 PM"
        text = "Your Chase verification code is 492011. Do not share this code with anyone. Chase will never call or text to ask for this code."
        blur = $false
        empty = $false
    },
    @{
        id = "screen-24-safe-google-2fa-dark"
        type = "sms"
        theme = "dark"
        sender = "Google 22000"
        time = "07:45 PM"
        text = "G-821940 is your Google verification code. Do not share your code with anyone."
        blur = $false
        empty = $false
    },
    @{
        id = "screen-25-safe-usps-delivered-light"
        type = "sms"
        theme = "light"
        sender = "USPS Tracking"
        time = "02:15 PM"
        text = "USPS: Package 9400111899562544 was delivered at 2:15 PM on your front porch. Thank you for shipping with USPS."
        blur = $false
        empty = $false
    },
    @{
        id = "screen-26-safe-cvs-rx-light"
        type = "sms"
        theme = "light"
        sender = "CVS Pharmacy"
        time = "01:00 PM"
        text = "CVS: Your prescription for Lisinopril is ready for pickup at 145 Main Street. Pharmacy closes today at 7:00 PM."
        blur = $false
        empty = $false
    },
    @{
        id = "screen-27-safe-doctor-appt-light"
        type = "sms"
        theme = "light"
        sender = "Dr. Miller Dental"
        time = "03:30 PM"
        text = "Dr. Miller Office: Reminder for your dental cleaning tomorrow at 10:30 AM. Reply C to confirm or call 555-0122 to reschedule."
        blur = $false
        empty = $false
    },
    @{
        id = "screen-28-safe-family-dinner-dark"
        type = "sms"
        theme = "dark"
        sender = "Grandson Alex"
        time = "05:15 PM"
        text = "Hey Grandma, hope you had a good week! Are we still having family dinner this Sunday at 1:00 PM at your place?"
        blur = $false
        empty = $false
    },
    @{
        id = "screen-29-damage-control-card-light"
        type = "sms"
        theme = "light"
        sender = "Security Alert"
        time = "11:50 AM"
        text = "Netflix Billing: Payment failed. Update credit card at http://bit.ly/netflix-billing-fix. Entered card details on site."
        blur = $false
        empty = $false
    },
    @{
        id = "screen-30-damage-control-wire-dark"
        type = "sms"
        theme = "dark"
        sender = "Federal Tax Division"
        time = "04:30 PM"
        text = "IRS Settlement: Outstanding tax warrant of $3,500. Wired settlement from credit union an hour ago."
        blur = $false
        empty = $false
    },
    @{
        id = "screen-31-no-text-coffee-photo"
        type = "photo"
        theme = "light"
        sender = ""
        time = ""
        text = ""
        blur = $false
        empty = $true
    },
    @{
        id = "screen-32-no-text-landscape-photo"
        type = "photo"
        theme = "dark"
        sender = ""
        time = ""
        text = ""
        blur = $false
        empty = $true
    }
)

$fontFamily = [System.Drawing.FontFamily]::GenericSansSerif
$fontBody = [System.Drawing.Font]::new($fontFamily, 13, [System.Drawing.FontStyle]::Regular)
$fontHeader = [System.Drawing.Font]::new($fontFamily, 12, [System.Drawing.FontStyle]::Bold)
$fontMeta = [System.Drawing.Font]::new($fontFamily, 10, [System.Drawing.FontStyle]::Regular)

$manifest = @()

foreach ($s in $specs) {
    $width = 720
    $height = 400
    $bmp = [System.Drawing.Bitmap]::new($width, $height)
    $g = [System.Drawing.Graphics]::FromImage($bmp)
    $g.SmoothingMode = [System.Drawing.Drawing2D.SmoothingMode]::AntiAlias
    $g.TextRenderingHint = [System.Drawing.Text.TextRenderingHint]::ClearTypeGridFit

    if ($s.empty) {
        # Render a photographic gradient / pattern with zero text
        $bgBrush = [System.Drawing.Drawing2D.LinearGradientBrush]::new(
            [System.Drawing.Point]::new(0, 0),
            [System.Drawing.Point]::new($width, $height),
            [System.Drawing.Color]::FromArgb(40, 80, 120),
            [System.Drawing.Color]::FromArgb(180, 140, 90)
        )
        $g.FillRectangle($bgBrush, 0, 0, $width, $height)
        $bgBrush.Dispose()

        # Draw decorative circles/shapes like a photo
        $pen = [System.Drawing.Pen]::new([System.Drawing.Color]::FromArgb(80, 255, 255, 255), 4)
        $g.DrawEllipse($pen, 120, 80, 240, 240)
        $g.DrawEllipse($pen, 360, 140, 180, 180)
        $pen.Dispose()
    } else {
        $isDark = ($s.theme -eq "dark")
        $bgColor = if ($isDark) { [System.Drawing.Color]::FromArgb(24, 26, 31) } else { [System.Drawing.Color]::FromArgb(242, 244, 247) }
        $bubbleBg = if ($isDark) { [System.Drawing.Color]::FromArgb(44, 48, 58) } else { [System.Drawing.Color]::FromArgb(255, 255, 255) }
        $textBrush = if ($isDark) { [System.Drawing.Brushes]::White } else { [System.Drawing.Brushes]::Black }
        $metaBrush = if ($isDark) { [System.Drawing.Brushes]::LightGray } else { [System.Drawing.Brushes]::DarkGray }
        $headerBg = if ($isDark) { [System.Drawing.Color]::FromArgb(32, 35, 42) } else { [System.Drawing.Color]::FromArgb(230, 233, 238) }

        # Background
        $g.Clear($bgColor)

        # App bar / header
        $headerBrush = [System.Drawing.SolidBrush]::new($headerBg)
        $g.FillRectangle($headerBrush, 0, 0, $width, 60)
        $headerBrush.Dispose()

        # Sender title in header
        $g.DrawString($s.sender, $fontHeader, $textBrush, 28, 20)
        $g.DrawString($s.time, $fontMeta, $metaBrush, ($width - 120), 22)

        # Bubble
        $bubbleBrush = [System.Drawing.SolidBrush]::new($bubbleBg)
        $bubbleRect = [System.Drawing.RectangleF]::new(28, 90, 664, 260)
        $g.FillRectangle($bubbleBrush, $bubbleRect)
        $bubbleBrush.Dispose()

        # Message text inside bubble
        $textRect = [System.Drawing.RectangleF]::new(48, 110, 624, 220)
        $format = [System.Drawing.StringFormat]::new()
        $format.LineAlignment = [System.Drawing.StringAlignment]::Near
        $format.Alignment = [System.Drawing.StringAlignment]::Near
        $g.DrawString($s.text, $fontBody, $textBrush, $textRect, $format)
        $format.Dispose()
    }

    $filePath = Join-Path $outDir "$($s.id).png"
    $bmp.Save($filePath, [System.Drawing.Imaging.ImageFormat]::Png)
    $g.Dispose()
    $bmp.Dispose()

    $manifest += [PSCustomObject]@{
        id = $s.id
        file = "$($s.id).png"
        type = $s.type
        theme = $s.theme
        expected_text = $s.text
        is_blurry = $s.blur
        is_blank = $s.empty
    }
}

$fontBody.Dispose()
$fontHeader.Dispose()
$fontMeta.Dispose()

$manifestJson = $manifest | ConvertTo-Json -Depth 4
$manifestPath = Join-Path $outDir "manifest.json"
Set-Content -Path $manifestPath -Value $manifestJson -Encoding utf8

Write-Output "Successfully generated $($specs.Count) screenshots in $outDir"
