import Link from 'next/link';
import { Button } from '@/components/ui/Button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/Card';
import { AuthManager } from '@/lib/auth';

export default function Home() {
  const user = AuthManager.getUser();

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      {/* Hero Section */}
      <section className="bg-white dark:bg-gray-800 py-20">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center">
            <div className="flex justify-center mb-8">
              <div 
                className="w-20 h-20 rounded-full flex items-center justify-center"
                style={{
                  background: `linear-gradient(135deg, var(--primary), var(--accent))`
                }}
              >
                <span className="text-white font-bold text-2xl">C</span>
              </div>
            </div>
            <h1 className="text-4xl font-bold text-gray-900 dark:text-white sm:text-5xl md:text-6xl">
              Welcome to{' '}
              <span 
                className="bg-clip-text text-transparent"
                style={{
                  background: `linear-gradient(to right, var(--primary), var(--accent))`,
                  WebkitBackgroundClip: 'text',
                  WebkitTextFillColor: 'transparent'
                }}
              >
                CodeWithDell
              </span>
            </h1>
            <p className="mt-6 text-xl text-gray-600 dark:text-gray-300 max-w-3xl mx-auto">
              🎥 YouTube Project Showcase - Explore my latest programming projects, tutorials, and 
              real-world applications. From web development to mobile apps, discover practical 
              code examples and learn through hands-on experience.
            </p>
            <div className="mt-10 flex flex-col sm:flex-row gap-4 justify-center">
              <Link href="/blog">
                <Button 
                  size="lg" 
                  className="w-full sm:w-auto"
                  style={{
                    background: `linear-gradient(to right, var(--primary), var(--accent))`
                  }}
                >
                  🚀 Explore Blog
                </Button>
              </Link>
              <a 
                href="https://youtube.com/@codewithdell" 
                target="_blank" 
                rel="noopener noreferrer"
                className="inline-flex items-center px-6 py-3 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
              >
                <svg className="w-5 h-5 mr-2 text-red-600" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M23.498 6.186a3.016 3.016 0 0 0-2.122-2.136C19.505 3.545 12 3.545 12 3.545s-7.505 0-9.377.505A3.017 3.017 0 0 0 .502 6.186C0 8.07 0 12 0 12s0 3.93.502 5.814a3.016 3.016 0 0 0 2.122 2.136c1.871.505 9.376.505 9.376.505s7.505 0 9.377-.505a3.015 3.015 0 0 0 2.122-2.136C24 15.93 24 12 24 12s0-3.93-.502-5.814zM9.545 15.568V8.432L15.818 12l-6.273 3.568z"/>
                </svg>
                Watch on YouTube
              </a>
            </div>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-20">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl font-bold text-gray-900 dark:text-white">
              What You'll Discover
            </h2>
            <p className="mt-4 text-lg text-gray-600 dark:text-gray-300">
              A curated collection of my best YouTube projects and tutorials
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <Card className="hover:shadow-lg transition-shadow duration-300">
              <CardHeader>
                <div 
                  className="w-12 h-12 rounded-lg flex items-center justify-center mb-4"
                  style={{
                    background: `linear-gradient(135deg, var(--primary), #5855eb)`
                  }}
                >
                  <svg className="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                  </svg>
                </div>
                <CardTitle>Web Development</CardTitle>
                <CardDescription>
                  Full-stack applications, REST APIs, modern frameworks like React, Next.js, and backend technologies
                </CardDescription>
              </CardHeader>
            </Card>

            <Card className="hover:shadow-lg transition-shadow duration-300">
              <CardHeader>
                <div 
                  className="w-12 h-12 rounded-lg flex items-center justify-center mb-4"
                  style={{
                    background: `linear-gradient(135deg, var(--accent), #0d9488)`
                  }}
                >
                  <svg className="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.746 0 3.332.477 4.5 1.253v13C19.832 18.477 18.246 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
                  </svg>
                </div>
                <CardTitle>Step-by-Step Tutorials</CardTitle>
                <CardDescription>
                  Detailed guides with source code, explanations, and best practices for developers
                </CardDescription>
              </CardHeader>
            </Card>

            <Card className="hover:shadow-lg transition-shadow duration-300">
              <CardHeader>
                <div 
                  className="w-12 h-12 rounded-lg flex items-center justify-center mb-4"
                  style={{
                    background: `linear-gradient(135deg, var(--accent-alt), #d97706)`
                  }}
                >
                  <svg className="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                  </svg>
                </div>
                <CardTitle>Real Projects</CardTitle>
                <CardDescription>
                  Production-ready applications with downloadable source code and live demos
                </CardDescription>
              </CardHeader>
            </Card>
          </div>
        </div>
      </section>

      {/* Stats Section */}
      <section className="bg-white dark:bg-gray-800 py-16">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-8 text-center">
            <div>
              <div 
                className="text-3xl font-bold"
                style={{ color: 'var(--primary)' }}
              >
                50+
              </div>
              <div className="text-gray-600 dark:text-gray-300">Projects</div>
            </div>
            <div>
              <div 
                className="text-3xl font-bold"
                style={{ color: 'var(--accent)' }}
              >
                100+
              </div>
              <div className="text-gray-600 dark:text-gray-300">Tutorials</div>
            </div>
            <div>
              <div 
                className="text-3xl font-bold"
                style={{ color: 'var(--accent-alt)' }}
              >
                10K+
              </div>
              <div className="text-gray-600 dark:text-gray-300">Views</div>
            </div>
            <div>
              <div 
                className="text-3xl font-bold"
                style={{ color: 'var(--primary)' }}
              >
                24/7
              </div>
              <div className="text-gray-600 dark:text-gray-300">Support</div>
            </div>
          </div>
        </div>
      </section>
    </div>
  );
}
