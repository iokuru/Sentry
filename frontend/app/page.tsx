"use client";

import Link from "next/link";
import { Shell } from "@/components/shell";
import { ArrowRight, Check, Clock3, Image, MessageSquareText, UploadCloud, UsersRound } from "lucide-react";

const steps = [
  { number: "01", title: "Bring the message", text: "Use a screenshot, paste the text, or try a sample message." },
  { number: "02", title: "Get a clear read", text: "See what deserves a pause, explained in plain language." },
  { number: "03", title: "Choose a safe next step", text: "Use a practical action and the right official contact for your region." },
];

export default function HomePage() {
  return (
    <Shell activeSection="landing">
      <div className="landing-page">
        <section className="landing-hero" aria-labelledby="hero-title">
          <div className="hero-copy">
            <p className="hero-kicker"><span /> A calmer way to handle unexpected messages</p>
            <h1 id="hero-title">Pause before<br /><em>you reply.</em></h1>
            <p className="hero-lede">A message can feel urgent. You don’t have to answer it that way. Get a clear second look before you click, pay, or share.</p>
            <div className="hero-actions">
              <Link href="/check" className="primary-action">Check a message <ArrowRight size={17} /></Link>
              <Link href="/caregiver" className="secondary-action"><UsersRound size={16} /> For family</Link>
            </div>
            <div className="hero-footnote"><span className="footnote-check"><Check size={12} /></span> No account needed <span className="footnote-divider">/</span> Your screenshot is deleted after its check</div>
          </div>

          <div className="hero-art" aria-label="Example text message history">
            <div className="art-topline"><span>INCOMING MESSAGE</span><span className="art-time">11:42</span></div>
            <div className="message-window">
              <div className="window-header"><span className="window-avatar">M</span><div><strong>Delivery update</strong><small>Text message · today</small></div><span className="window-dots">•••</span></div>
              <div className="conversation" aria-label="Conversation history about a suspicious parcel fee">
                <div className="conversation-turn incoming"><span className="turn-label">DELIVERY UPDATE · 11:42</span><p className="chat-bubble incoming-bubble">Your parcel is on hold. Pay a small redelivery fee to receive it today.<br /><span>Track delivery: parcel-update.example</span></p></div>
                <div className="conversation-turn outgoing"><span className="turn-label">YOU · 11:43</span><p className="chat-bubble outgoing-bubble">Why is a redelivery fee required? I already paid for delivery.</p></div>
                <div className="conversation-turn incoming"><span className="turn-label">DELIVERY UPDATE · 11:44</span><p className="chat-bubble incoming-bubble">Pay within 30 minutes or the parcel will be returned.</p></div>
              </div>
            </div>
            <div className="art-index">A closer look before the next tap <span>01 / 03</span></div>
          </div>
        </section>

        <section className="landing-gallery" aria-labelledby="start-title">
          <div className="gallery-heading"><div><p className="section-kicker">START WITH WHAT YOU HAVE</p><h2 id="start-title">Bring the message in.</h2><p>Choose a format to open the right part of the checker.</p></div><span className="gallery-count">01 <i /> 03</span></div>
          <div className="entry-gallery">
            <article className="entry-card">
              <div className="entry-head"><span className="entry-icon"><Image size={17} /></span><span>01</span></div>
              <div className="entry-preview upload-preview"><UploadCloud size={24} /><strong>Drop an image here</strong><small>PNG · JPG · HEIC</small></div>
              <div className="entry-copy"><div><h3>Screenshot or photo</h3><p>Upload a message or use a screenshot from your camera roll.</p></div><Link href="/check?mode=image" aria-label="Check a screenshot" className="entry-link"><ArrowRight size={17} /></Link></div>
            </article>
            <article className="entry-card">
              <div className="entry-head"><span className="entry-icon"><MessageSquareText size={17} /></span><span>02</span></div>
              <div className="entry-preview text-entry-preview"><span className="text-lines"><i /><i /><i /></span><span className="text-caret" /></div>
              <div className="entry-copy"><div><h3>Paste the message</h3><p>Include the wording, sender, or link that made you pause.</p></div><Link href="/check?mode=text" aria-label="Paste message text" className="entry-link"><ArrowRight size={17} /></Link></div>
            </article>
            <article className="entry-card">
              <div className="entry-head"><span className="entry-icon"><Clock3 size={17} /></span><span>03</span></div>
              <div className="entry-preview practice-entry-preview"><div className="practice-bubble">Your account will close today. Verify your details to keep access.</div><small>Unknown sender · 9:18 AM</small></div>
              <div className="entry-copy"><div><h3>Try a sample</h3><p>See how a message review works before using your own.</p></div><Link href="/check?mode=practice" aria-label="Try a sample message" className="entry-link"><ArrowRight size={17} /></Link></div>
            </article>
          </div>
        </section>

        <section className="how-panel" aria-labelledby="how-title">
          <div className="how-copy"><p className="section-kicker">A SMALL PAUSE GOES A LONG WAY</p><h2 id="how-title">A clear read.<br />A useful next step.</h2><p>Review the details at your own pace. You stay in control of what happens next.</p></div>
          <ol className="steps-list">{steps.map((step) => <li key={step.number}><span>{step.number}</span><div><h3>{step.title}</h3><p>{step.text}</p></div><ArrowRight size={16} /></li>)}</ol>
        </section>

        <section className="family-band">
          <div className="family-symbol"><UsersRound size={19} /></div>
          <div className="family-copy"><p className="section-kicker">FOR FAMILY</p><h2>Someone you trust can stay in the loop.</h2><p>Family Circle gives a caregiver a place to review shared alerts and follow up.</p></div>
          <Link href="/caregiver" className="family-cta">Open family space <ArrowRight size={15} /></Link>
        </section>
        <footer className="landing-footer"><span>Take a moment before you reply.</span><Link href="/check">Message check <ArrowRight size={13} /></Link></footer>
      </div>
      <style jsx>{`
        .landing-page { width:100%; max-width:1220px; margin:0 auto; color:var(--text-main); }
        .landing-hero { display:grid; grid-template-columns:minmax(0,1.04fr) minmax(380px,.96fr); align-items:center; gap:clamp(34px,6vw,90px); padding:35px 0 55px; border-bottom:1px solid var(--border-light); }
        .hero-copy { max-width:570px; }
        .hero-kicker { display:flex; align-items:center; gap:9px; color:var(--text-muted); font-size:13px; }
        .hero-kicker span { width:8px; height:8px; border-radius:50%; background:#8fba38; }
        h1 { margin:24px 0 18px; font-size:clamp(48px,6.2vw,80px); line-height:.99; letter-spacing:-.065em; font-weight:650; }
        h1 em { color:var(--color-primary); font-style:normal; }
        .hero-lede { max-width:480px; color:var(--text-muted); font-size:18px; line-height:1.58; }
        .hero-actions { display:flex; align-items:center; gap:18px; margin-top:30px; }
        :global(.primary-action) { display:inline-flex; align-items:center; gap:20px; min-height:50px; padding:0 18px; border-radius:var(--radius-sm); background:var(--color-primary); color:var(--button-on-primary); font-size:14px; font-weight:650; text-decoration:none; transition:background .15s,transform .15s; }
        :global(.primary-action:hover) { color:var(--button-on-primary); background:var(--color-primary-hover); transform:translateY(-1px); }
        :global(.secondary-action) { display:inline-flex; align-items:center; gap:8px; color:var(--text-main); font-size:14px; font-weight:560; text-decoration:none; }
        :global(.secondary-action:hover) { color:var(--color-primary); }
        .hero-footnote { display:flex; align-items:center; gap:8px; margin-top:24px; color:var(--text-subtle); font-size:12px; }
        .footnote-check { display:grid; place-items:center; width:17px; height:17px; border-radius:50%; color:#fff; background:#6c9f33; }.footnote-divider { color:var(--border-medium); }
        .hero-art { position:relative; min-height:470px; padding:18px 0 18px 30px; border-left:1px solid var(--border-light); }
        .art-topline { display:flex; justify-content:space-between; color:var(--text-subtle); font:11px var(--font-mono); letter-spacing:.12em; }
        .art-time { letter-spacing:0; }
        .message-window { width:min(100%,430px); margin:25px 0 0 auto; padding:16px; border:1px solid var(--border-light); border-radius:var(--radius-md); background:var(--gallery-card); box-shadow:0 14px 36px rgba(0,0,0,.08); }
        .window-header { display:flex; align-items:center; gap:10px; padding:2px 0 13px; border-bottom:1px solid var(--border-light); }.window-header > div { display:grid; gap:3px; }.window-header strong { font-size:13px; font-weight:620; }.window-header small { color:var(--text-subtle); font-size:11px; }.window-avatar { display:grid; place-items:center; width:30px; height:30px; border-radius:50%; background:var(--bg-muted); color:var(--text-main); font-size:12px; font-weight:650; }.window-dots { margin-left:auto; color:var(--text-subtle); letter-spacing:2px; }
        .conversation { display:grid; gap:13px; padding:18px 0 4px; }.conversation-turn { display:grid; gap:5px; width:90%; }.conversation-turn.incoming { justify-self:start; }.conversation-turn.outgoing { justify-self:end; justify-items:end; }.turn-label { color:var(--text-subtle); font:9px var(--font-mono); letter-spacing:.07em; }.chat-bubble { width:fit-content; max-width:100%; margin:0; padding:11px 13px; border-radius:var(--radius-md); color:var(--sample-ink); font-size:12px; line-height:1.5; overflow-wrap:anywhere; }.incoming-bubble { border-bottom-left-radius:3px; background:var(--sample-bubble); }.outgoing-bubble { border-bottom-right-radius:3px; background:var(--outgoing-bubble); color:var(--outgoing-ink); }.incoming-bubble span { display:inline-block; margin-top:7px; color:#567629; font-size:11px; text-decoration:underline; text-underline-offset:2px; }
        .art-index { display:flex; justify-content:space-between; margin-top:16px; color:var(--text-subtle); font:10px var(--font-mono); }.art-index span { color:var(--text-muted); }
        .landing-gallery { padding:44px 0 51px; border-bottom:1px solid var(--border-light); }.gallery-heading { display:flex; justify-content:space-between; align-items:end; gap:20px; margin-bottom:22px; }.section-kicker { color:var(--text-subtle); font:10px var(--font-mono); letter-spacing:.12em; }.gallery-heading h2 { margin-top:8px; font-size:27px; font-weight:630; letter-spacing:-.035em; }.gallery-heading > div > p:last-child { margin-top:6px; color:var(--text-muted); font-size:14px; }.gallery-count { display:flex; align-items:center; gap:9px; color:var(--text-subtle); font:11px var(--font-mono); }.gallery-count i { width:30px; height:1px; background:var(--border-medium); }
        .entry-gallery { display:grid; grid-template-columns:repeat(3,minmax(0,1fr)); gap:16px; }.entry-card { display:grid; grid-template-rows:auto 140px 1fr; gap:13px; min-width:0; padding:15px; border:1px solid var(--border-light); border-radius:var(--radius-md); background:var(--gallery-card); transition:transform .16s,border-color .16s; }.entry-card:hover { transform:translateY(-2px); border-color:var(--border-medium); }.entry-head { display:flex; align-items:center; justify-content:space-between; color:var(--text-subtle); font:11px var(--font-mono); }.entry-icon { display:grid; place-items:center; width:31px; height:31px; border:1px solid var(--border-light); border-radius:var(--radius-sm); color:var(--text-muted); }.entry-preview { display:flex; align-items:center; justify-content:center; overflow:hidden; min-width:0; border-radius:var(--radius-sm); background:var(--preview-bg); color:var(--preview-ink); }.upload-preview { flex-direction:column; gap:9px; border:1px dashed var(--preview-border); color:var(--text-subtle); }.upload-preview strong { color:var(--preview-ink); font-size:13px; font-weight:560; }.upload-preview small { font-size:10px; color:var(--text-subtle); }.text-entry-preview { position:relative; align-items:flex-start; justify-content:flex-start; padding:25px; border:1px solid var(--preview-border); }.text-lines { display:grid; gap:10px; width:88%; }.text-lines i { height:6px; border-radius:4px; background:var(--preview-line); }.text-lines i:nth-child(2) { width:81%; }.text-lines i:nth-child(3) { width:58%; }.text-caret { position:absolute; top:67px; left:calc(25px + 59%); width:1px; height:15px; background:var(--color-primary); }.practice-entry-preview { flex-direction:column; align-items:flex-end; gap:9px; padding:18px; }.practice-bubble { max-width:94%; padding:11px 13px; border-radius:var(--radius-md) var(--radius-md) 3px var(--radius-md); background:var(--sample-bubble); color:var(--sample-ink); font-size:12px; line-height:1.45; }.practice-entry-preview small { color:var(--text-subtle); font-size:10px; }.entry-copy { display:flex; align-items:center; justify-content:space-between; gap:12px; }.entry-copy h3 { font-size:16px; font-weight:620; letter-spacing:-.015em; }.entry-copy p { margin-top:5px; color:var(--text-muted); font-size:12px; line-height:1.5; }:global(.entry-link) { display:grid; place-items:center; width:36px; height:36px; flex:none; border:1px solid var(--border-light); border-radius:var(--radius-sm); color:var(--text-main); background:var(--work-card); text-decoration:none; }:global(.entry-link:hover) { color:var(--color-primary); border-color:var(--color-primary); }
        .how-panel { display:grid; grid-template-columns:.85fr 1.15fr; gap:clamp(28px,5vw,68px); align-items:center; margin:45px 0; padding:28px clamp(20px,3vw,38px); border:1px solid var(--border-light); border-radius:var(--radius-lg); background:var(--bg-muted); }.how-copy h2 { margin-top:12px; font-size:clamp(26px,3vw,37px); line-height:1.08; letter-spacing:-.04em; font-weight:630; }.how-copy > p:last-child { max-width:340px; margin-top:13px; color:var(--text-muted); font-size:14px; line-height:1.55; }.steps-list { list-style:none; display:grid; }.steps-list li { display:grid; grid-template-columns:36px 1fr 18px; align-items:center; gap:12px; padding:15px 0; border-bottom:1px solid var(--border-light); }.steps-list li:last-child { border-bottom:0; }.steps-list li > span { color:var(--color-primary); font:11px var(--font-mono); }.steps-list h3 { font-size:14px; font-weight:620; }.steps-list p { margin-top:4px; color:var(--text-muted); font-size:12px; line-height:1.45; }.steps-list li > svg { color:var(--text-subtle); }
        .family-band { display:flex; align-items:center; gap:16px; padding:21px 23px; border:1px solid var(--border-light); border-radius:var(--radius-md); background:var(--gallery-card); }.family-symbol { display:grid; place-items:center; width:43px; height:43px; flex:none; border:1px solid var(--border-light); border-radius:var(--radius-sm); color:var(--color-primary); background:var(--bg-muted); }.family-copy .section-kicker { margin-bottom:6px; }.family-band h2 { font-size:16px; font-weight:620; }.family-band p:not(.section-kicker) { margin-top:4px; color:var(--text-muted); font-size:13px; }:global(.family-cta) { display:inline-flex; align-items:center; gap:12px; margin-left:auto; padding:10px 13px; border:1px solid var(--border-light); border-radius:var(--radius-sm); color:var(--text-main); background:var(--work-card); text-decoration:none; font-size:12px; font-weight:560; white-space:nowrap; }:global(.family-cta:hover) { color:var(--color-primary); border-color:var(--color-primary); }
        .landing-footer { display:flex; justify-content:space-between; gap:18px; padding:23px 0 4px; color:var(--text-subtle); font-size:11px; }:global(.landing-footer a) { display:flex; align-items:center; gap:6px; color:var(--text-muted); text-decoration:none; }:global(.landing-footer a:hover) { color:var(--color-primary); }
        @media(max-width:1050px) { .landing-hero { grid-template-columns:minmax(0,1fr); gap:30px; }.hero-art { width:100%; max-width:640px; margin:0 auto; } }
        @media(max-width:800px) { .entry-gallery { grid-template-columns:1fr; }.entry-card { grid-template-columns:36px minmax(100px,.7fr) 1.4fr; grid-template-rows:auto auto; align-items:center; gap:5px 12px; }.entry-head { grid-row:1 / 3; align-self:start; }.entry-preview { grid-row:1 / 3; grid-column:2; height:95px; }.entry-copy { grid-column:3; }.how-panel { grid-template-columns:1fr; gap:18px; }.how-copy h2 br { display:none; } }
        @media(max-width:680px) { .landing-hero { grid-template-columns:1fr; padding:20px 0 34px; gap:27px; }.hero-art { min-height:0; padding:18px 0 8px 17px; }.message-window { margin-top:19px; }.hero-footnote { align-items:flex-start; flex-wrap:wrap; }.landing-gallery { padding:31px 0 38px; }.gallery-heading h2 { font-size:23px; }.entry-card { grid-template-columns:31px minmax(92px,.65fr) 1.35fr; gap:5px 9px; padding:12px; }.entry-copy h3 { font-size:14px; }.entry-copy p { font-size:11px; }:global(.entry-link) { width:31px; height:31px; }.how-panel { margin:30px 0; padding:20px 17px; border-radius:var(--radius-md); }.steps-list li { grid-template-columns:29px 1fr 16px; gap:9px; }.family-band { align-items:flex-start; padding:15px; }.family-copy { min-width:0; }.family-band h2 { font-size:14px; }.family-band p:not(.section-kicker) { line-height:1.45; font-size:11px; }:global(.family-cta) { align-self:center; padding:9px; font-size:0; gap:0; }.family-cta svg { width:16px; height:16px; }.landing-footer { font-size:10px; } }
      `}</style>
    </Shell>
  );
}
