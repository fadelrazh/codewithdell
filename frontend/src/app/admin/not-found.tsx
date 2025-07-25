import Link from 'next/link';
import { Button } from '@/components/ui/Button';

export default function AdminNotFound() {
  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
      <div className="sm:mx-auto sm:w-full sm:max-w-md">
        <div className="flex justify-center">
          <div className="w-12 h-12 bg-red-600 rounded-lg flex items-center justify-center">
            <span className="text-white font-bold text-lg">!</span>
          </div>
        </div>
        <h2 className="mt-6 text-center text-3xl font-bold text-gray-900 dark:text-white">
          Admin Page Not Found
        </h2>
        <p className="mt-2 text-center text-sm text-gray-600 dark:text-gray-300">
          The admin page you're looking for doesn't exist
        </p>
      </div>

      <div className="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
        <div className="bg-white dark:bg-gray-800 py-8 px-4 shadow sm:rounded-lg sm:px-10 text-center">
          <div className="space-y-4">
            <Link href="/admin">
              <Button className="w-full">
                Go to Dashboard
              </Button>
            </Link>
            <Link href="/">
              <Button variant="outline" className="w-full">
                Back to Site
              </Button>
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
} 