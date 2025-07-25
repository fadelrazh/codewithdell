import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { Header } from "@/components/layout/Header";
import { Footer } from "@/components/layout/Footer";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "CodeWithDell - YouTube Project Showcase",
  description: "🎥 YouTube Project Showcase - Explore my latest programming projects, tutorials, and real-world applications. From web development to mobile apps, discover practical code examples and learn through hands-on experience.",
  keywords: ["youtube", "programming", "tutorials", "web development", "projects", "code", "react", "nextjs", "go"],
  authors: [{ name: "CodeWithDell" }],
  creator: "CodeWithDell",
  openGraph: {
    type: "website",
    locale: "en_US",
    url: "https://codewithdell.com",
    title: "CodeWithDell - YouTube Project Showcase",
    description: "🎥 YouTube Project Showcase - Explore my latest programming projects, tutorials, and real-world applications.",
    siteName: "CodeWithDell",
  },
  twitter: {
    card: "summary_large_image",
    title: "CodeWithDell - YouTube Project Showcase",
    description: "🎥 YouTube Project Showcase - Explore my latest programming projects, tutorials, and real-world applications.",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <head>
        <script
          dangerouslySetInnerHTML={{
            __html: `
              (function() {
                try {
                  // Default to system theme if no theme saved
                  var theme = localStorage.getItem('theme') || 'system';
                  var isDark = theme === 'dark' || (theme === 'system' && window.matchMedia('(prefers-color-scheme: dark)').matches);
                  var root = document.documentElement;
                  
                  if (isDark) {
                    root.classList.add('dark');
                  } else {
                    root.classList.remove('dark');
                  }
                  
                  // Listen for system theme changes when using system theme
                  if (theme === 'system') {
                    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', function(e) {
                      if (localStorage.getItem('theme') === 'system') {
                        if (e.matches) {
                          root.classList.add('dark');
                        } else {
                          root.classList.remove('dark');
                        }
                      }
                    });
                  }
                } catch (e) {
                  console.error('Theme initialization error:', e);
                }
              })();
            `,
          }}
        />
      </head>
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased min-h-screen bg-gray-50 dark:bg-gray-900`}
      >
        <div className="flex flex-col min-h-screen">
          <Header />
          <main className="flex-1">
            {children}
          </main>
          <Footer />
        </div>
      </body>
    </html>
  );
}
