import type { Metadata } from "next";
import "../styles/globals.css";

export const metadata: Metadata = {
  title: "Take a moment — Check before you reply",
  description: "A calm, private second look at unexpected messages before you click, pay, or share.",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" data-theme="dark">
      <body>{children}</body>
    </html>
  );
}
