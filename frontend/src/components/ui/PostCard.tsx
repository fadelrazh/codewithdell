'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { Button } from './Button';
import { Card, CardContent, CardHeader } from './Card';
import { apiClient } from '@/lib/api';
import { Post, User } from '@/types/api';
import { AuthManager } from '@/lib/auth';

interface PostCardProps {
  post: Post;
  showAuthor?: boolean;
  showExcerpt?: boolean;
  showStats?: boolean;
  className?: string;
}

export const PostCard: React.FC<PostCardProps> = ({
  post,
  showAuthor = true,
  showExcerpt = true,
  showStats = true,
  className = "",
}) => {
  const [isLiked, setIsLiked] = useState(false);
  const [isBookmarked, setIsBookmarked] = useState(false);
  const [likeCount, setLikeCount] = useState(post.like_count);
  const [loading, setLoading] = useState(false);
  const user = AuthManager.getUser();

  useEffect(() => {
    if (user) {
      checkUserInteraction();
    }
  }, [user, post.id]);

  const checkUserInteraction = async () => {
    try {
      const response = await apiClient.checkUserInteraction(post.id);
      setIsLiked(response.is_liked);
      setIsBookmarked(response.is_bookmarked);
    } catch (error) {
      console.error('Failed to check user interaction:', error);
    }
  };

  const handleLike = async () => {
    if (!user) {
      // Redirect to login
      window.location.href = '/login';
      return;
    }

    try {
      setLoading(true);
      if (isLiked) {
        await apiClient.unlikePost(post.id);
        setLikeCount(prev => prev - 1);
        setIsLiked(false);
      } else {
        await apiClient.likePost(post.id);
        setLikeCount(prev => prev + 1);
        setIsLiked(true);
      }
    } catch (error) {
      console.error('Failed to like/unlike post:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleBookmark = async () => {
    if (!user) {
      // Redirect to login
      window.location.href = '/login';
      return;
    }

    try {
      setLoading(true);
      if (isBookmarked) {
        await apiClient.removeBookmark(post.id);
        setIsBookmarked(false);
      } else {
        await apiClient.bookmarkPost(post.id);
        setIsBookmarked(true);
      }
    } catch (error) {
      console.error('Failed to bookmark/unbookmark post:', error);
    } finally {
      setLoading(false);
    }
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    const now = new Date();
    const diffTime = Math.abs(now.getTime() - date.getTime());
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

    if (diffDays === 1) return 'Today';
    if (diffDays === 2) return 'Yesterday';
    if (diffDays < 7) return `${diffDays - 1} days ago`;
    if (diffDays < 30) return `${Math.floor(diffDays / 7)} weeks ago`;
    if (diffDays < 365) return `${Math.floor(diffDays / 30)} months ago`;
    return `${Math.floor(diffDays / 365)} years ago`;
  };

  const formatNumber = (num: number) => {
    if (num >= 1000000) return `${(num / 1000000).toFixed(1)}M`;
    if (num >= 1000) return `${(num / 1000).toFixed(1)}K`;
    return num.toString();
  };

  return (
    <Card className={`hover:shadow-lg transition-shadow duration-300 ${className}`}>
      {/* Featured Image */}
      {post.featured_image && (
        <div className="relative h-48 overflow-hidden rounded-t-lg">
          <Image
            src={post.featured_image}
            alt={post.title}
            fill
            className="object-cover"
          />
          <div className="absolute top-2 right-2">
            <span className={`px-2 py-1 text-xs rounded-full ${
              post.status === 'published' 
                ? 'bg-green-100 text-green-800' 
                : 'bg-yellow-100 text-yellow-800'
            }`}>
              {post.status}
            </span>
          </div>
        </div>
      )}

      <CardHeader className="pb-2">
        {/* Categories */}
        {post.categories && post.categories.length > 0 && (
          <div className="flex flex-wrap gap-1 mb-2">
            {post.categories.slice(0, 2).map((category) => (
              <Link
                key={category.id}
                href={`/blog/category/${category.slug}`}
                className="px-2 py-1 text-xs bg-blue-100 text-blue-800 rounded-full hover:bg-blue-200 transition-colors"
              >
                {category.name}
              </Link>
            ))}
            {post.categories.length > 2 && (
              <span className="px-2 py-1 text-xs text-gray-500">
                +{post.categories.length - 2} more
              </span>
            )}
          </div>
        )}

        {/* Title */}
        <Link href={`/blog/${post.slug}`}>
          <h3 className="text-xl font-semibold text-gray-900 hover:text-blue-600 transition-colors line-clamp-2">
            {post.title}
          </h3>
        </Link>

        {/* Excerpt */}
        {showExcerpt && post.excerpt && (
          <p className="text-gray-600 text-sm line-clamp-3 mt-2">
            {post.excerpt}
          </p>
        )}

        {/* Author */}
        {showAuthor && (
          <div className="flex items-center space-x-2 mt-3">
            <div className="w-8 h-8 rounded-full bg-gradient-to-r from-blue-500 to-purple-600 flex items-center justify-center text-white text-sm font-bold">
              {post.author.first_name.charAt(0)}
            </div>
            <div className="flex-1">
              <div className="text-sm font-medium text-gray-900">
                {post.author.first_name} {post.author.last_name}
              </div>
              <div className="text-xs text-gray-500">
                {formatDate(post.created_at)}
              </div>
            </div>
          </div>
        )}
      </CardHeader>

      <CardContent className="pt-0">
        {/* Tags */}
        {post.tags && post.tags.length > 0 && (
          <div className="flex flex-wrap gap-1 mb-4">
            {post.tags.slice(0, 3).map((tag) => (
              <Link
                key={tag.id}
                href={`/blog/tag/${tag.slug}`}
                className="px-2 py-1 text-xs bg-gray-100 text-gray-700 rounded-full hover:bg-gray-200 transition-colors"
              >
                #{tag.name}
              </Link>
            ))}
            {post.tags.length > 3 && (
              <span className="px-2 py-1 text-xs text-gray-500">
                +{post.tags.length - 3} more
              </span>
            )}
          </div>
        )}

        {/* Stats and Actions */}
        <div className="flex items-center justify-between">
          {/* Stats */}
          {showStats && (
            <div className="flex items-center space-x-4 text-sm text-gray-500">
              <div className="flex items-center space-x-1">
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>
                <span>{formatNumber(post.view_count)}</span>
              </div>
              <div className="flex items-center space-x-1">
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
                </svg>
                <span>{formatNumber(likeCount)}</span>
              </div>
              <div className="flex items-center space-x-1">
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                </svg>
                <span>{formatNumber(post.comment_count)}</span>
              </div>
            </div>
          )}

          {/* Action Buttons */}
          <div className="flex items-center space-x-2">
            <Button
              variant="ghost"
              size="sm"
              onClick={handleLike}
              disabled={loading}
              className={`flex items-center space-x-1 ${
                isLiked ? 'text-red-600 hover:text-red-700' : 'text-gray-500 hover:text-gray-700'
              }`}
            >
              <svg 
                className={`w-4 h-4 ${isLiked ? 'fill-current' : 'fill-none'}`} 
                stroke="currentColor" 
                viewBox="0 0 24 24"
              >
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
              </svg>
              <span className="text-xs">{formatNumber(likeCount)}</span>
            </Button>

            <Button
              variant="ghost"
              size="sm"
              onClick={handleBookmark}
              disabled={loading}
              className={`flex items-center space-x-1 ${
                isBookmarked ? 'text-blue-600 hover:text-blue-700' : 'text-gray-500 hover:text-gray-700'
              }`}
            >
              <svg 
                className={`w-4 h-4 ${isBookmarked ? 'fill-current' : 'fill-none'}`} 
                stroke="currentColor" 
                viewBox="0 0 24 24"
              >
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
              </svg>
            </Button>

            <Button
              variant="outline"
              size="sm"
              href={`/blog/${post.slug}`}
              className="text-xs"
            >
              Read More
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}; 