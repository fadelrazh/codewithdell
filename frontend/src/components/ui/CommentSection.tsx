'use client';

import { useState, useEffect } from 'react';
import { Button } from './Button';
import { Card, CardContent, CardHeader } from './Card';
import { apiClient } from '@/lib/api';
import { Comment, CreateCommentRequest } from '@/types/api';
import { AuthManager } from '@/lib/auth';

interface CommentSectionProps {
  postId?: number;
  projectId?: number;
}

interface CommentItemProps {
  comment: Comment;
  onReply: (commentId: number) => void;
  onEdit: (comment: Comment) => void;
  onDelete: (commentId: number) => void;
  replyToId?: number;
  setReplyToId: (id: number | null) => void;
}

const CommentItem: React.FC<CommentItemProps> = ({
  comment,
  onReply,
  onEdit,
  onDelete,
  replyToId,
  setReplyToId,
}) => {
  const user = AuthManager.getUser();
  const isOwner = user?.id === comment.user_id;
  const isAdmin = user?.role === 'admin';

  return (
    <Card className="mb-4">
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <div className="w-8 h-8 rounded-full bg-gradient-to-r from-blue-500 to-purple-600 flex items-center justify-center text-white text-sm font-bold">
              {comment.user.first_name.charAt(0)}
            </div>
            <div>
              <div className="font-semibold text-sm">
                {comment.user.first_name} {comment.user.last_name}
              </div>
              <div className="text-xs text-gray-500">
                {new Date(comment.created_at).toLocaleDateString()}
              </div>
            </div>
          </div>
          <div className="flex items-center space-x-2">
            {comment.status === 'pending' && (
              <span className="px-2 py-1 text-xs bg-yellow-100 text-yellow-800 rounded-full">
                Pending
              </span>
            )}
            {(isOwner || isAdmin) && (
              <div className="flex space-x-1">
                {isOwner && (
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => onEdit(comment)}
                    className="text-xs"
                  >
                    Edit
                  </Button>
                )}
                {(isOwner || isAdmin) && (
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => onDelete(comment.id)}
                    className="text-xs text-red-600 hover:text-red-700"
                  >
                    Delete
                  </Button>
                )}
              </div>
            )}
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <p className="text-sm text-gray-700 mb-3">{comment.content}</p>
        <Button
          variant="ghost"
          size="sm"
          onClick={() => onReply(comment.id)}
          className="text-xs text-blue-600 hover:text-blue-700"
        >
          Reply
        </Button>
        
        {/* Nested comments */}
        {comment.children && comment.children.length > 0 && (
          <div className="mt-4 ml-6 space-y-3">
            {comment.children.map((child) => (
              <CommentItem
                key={child.id}
                comment={child}
                onReply={onReply}
                onEdit={onEdit}
                onDelete={onDelete}
                replyToId={replyToId}
                setReplyToId={setReplyToId}
              />
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
};

export const CommentSection: React.FC<CommentSectionProps> = ({
  postId,
  projectId,
}) => {
  const [comments, setComments] = useState<Comment[]>([]);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [newComment, setNewComment] = useState('');
  const [replyToId, setReplyToId] = useState<number | null>(null);
  const [editingComment, setEditingComment] = useState<Comment | null>(null);
  const [totalComments, setTotalComments] = useState(0);

  const user = AuthManager.getUser();

  useEffect(() => {
    loadComments();
  }, [postId, projectId]);

  const loadComments = async () => {
    try {
      setLoading(true);
      const response = await apiClient.getComments(postId, projectId);
      setComments(response.comments);
      setTotalComments(response.total);
    } catch (error) {
      console.error('Failed to load comments:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newComment.trim() || !user) return;

    try {
      setSubmitting(true);
      const commentData: CreateCommentRequest = {
        content: newComment.trim(),
        post_id: postId,
        project_id: projectId,
        parent_id: replyToId || undefined,
      };

      const response = await apiClient.createComment(commentData);
      
      // Add new comment to the list
      if (replyToId) {
        // Add as child comment
        setComments(prev => prev.map(comment => {
          if (comment.id === replyToId) {
            return {
              ...comment,
              children: [...(comment.children || []), response.comment],
            };
          }
          return comment;
        }));
      } else {
        // Add as top-level comment
        setComments(prev => [response.comment, ...prev]);
      }
      
      setTotalComments(prev => prev + 1);
      setNewComment('');
      setReplyToId(null);
    } catch (error) {
      console.error('Failed to submit comment:', error);
    } finally {
      setSubmitting(false);
    }
  };

  const handleReply = (commentId: number) => {
    setReplyToId(commentId);
    setEditingComment(null);
  };

  const handleEdit = (comment: Comment) => {
    setEditingComment(comment);
    setNewComment(comment.content);
    setReplyToId(null);
  };

  const handleDelete = async (commentId: number) => {
    if (!confirm('Are you sure you want to delete this comment?')) return;

    try {
      await apiClient.deleteComment(commentId);
      
      // Remove comment from the list
      setComments(prev => {
        const removeComment = (comments: Comment[]): Comment[] => {
          return comments.filter(comment => {
            if (comment.id === commentId) {
              return false;
            }
            if (comment.children) {
              comment.children = removeComment(comment.children);
            }
            return true;
          });
        };
        return removeComment(prev);
      });
      
      setTotalComments(prev => prev - 1);
    } catch (error) {
      console.error('Failed to delete comment:', error);
    }
  };

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingComment || !newComment.trim()) return;

    try {
      setSubmitting(true);
      const response = await apiClient.updateComment(editingComment.id, {
        content: newComment.trim(),
      });
      
      // Update comment in the list
      const updateComment = (comments: Comment[]): Comment[] => {
        return comments.map(comment => {
          if (comment.id === editingComment.id) {
            return response.comment;
          }
          if (comment.children) {
            comment.children = updateComment(comment.children);
          }
          return comment;
        });
      };
      
      setComments(updateComment(comments));
      setNewComment('');
      setEditingComment(null);
    } catch (error) {
      console.error('Failed to update comment:', error);
    } finally {
      setSubmitting(false);
    }
  };

  const cancelEdit = () => {
    setEditingComment(null);
    setNewComment('');
    setReplyToId(null);
  };

  if (loading) {
    return (
      <div className="space-y-4">
        <h3 className="text-lg font-semibold">Comments ({totalComments})</h3>
        <div className="animate-pulse space-y-4">
          {[1, 2, 3].map((i) => (
            <div key={i} className="bg-gray-200 h-20 rounded"></div>
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <h3 className="text-lg font-semibold">Comments ({totalComments})</h3>
      
      {/* Comment form */}
      {user ? (
        <Card>
          <CardContent className="pt-6">
            <form onSubmit={editingComment ? handleUpdate : handleSubmit}>
              <div className="space-y-4">
                <div className="flex items-start space-x-3">
                  <div className="w-8 h-8 rounded-full bg-gradient-to-r from-blue-500 to-purple-600 flex items-center justify-center text-white text-sm font-bold">
                    {user.first_name.charAt(0)}
                  </div>
                  <div className="flex-1">
                    <textarea
                      value={newComment}
                      onChange={(e) => setNewComment(e.target.value)}
                      placeholder={
                        replyToId
                          ? "Write a reply..."
                          : editingComment
                          ? "Edit your comment..."
                          : "Write a comment..."
                      }
                      className="w-full p-3 border border-gray-300 rounded-lg resize-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      rows={3}
                      required
                    />
                    {(replyToId || editingComment) && (
                      <div className="mt-2 text-sm text-gray-500">
                        {replyToId && "Replying to a comment"}
                        {editingComment && "Editing your comment"}
                      </div>
                    )}
                  </div>
                </div>
                
                <div className="flex justify-between items-center">
                  <div className="flex space-x-2">
                    <Button
                      type="submit"
                      disabled={submitting || !newComment.trim()}
                      className="bg-blue-600 hover:bg-blue-700"
                    >
                      {submitting ? 'Submitting...' : editingComment ? 'Update' : 'Submit'}
                    </Button>
                    {(replyToId || editingComment) && (
                      <Button
                        type="button"
                        variant="outline"
                        onClick={cancelEdit}
                        disabled={submitting}
                      >
                        Cancel
                      </Button>
                    )}
                  </div>
                  {editingComment && (
                    <div className="text-sm text-gray-500">
                      Editing comment
                    </div>
                  )}
                </div>
              </div>
            </form>
          </CardContent>
        </Card>
      ) : (
        <Card>
          <CardContent className="pt-6">
            <div className="text-center py-8">
              <p className="text-gray-600 mb-4">
                Please log in to leave a comment
              </p>
              <Button href="/login" className="bg-blue-600 hover:bg-blue-700">
                Log In
              </Button>
            </div>
          </CardContent>
        </Card>
      )}
      
      {/* Comments list */}
      <div className="space-y-4">
        {comments.length === 0 ? (
          <div className="text-center py-8 text-gray-500">
            No comments yet. Be the first to comment!
          </div>
        ) : (
          comments.map((comment) => (
            <CommentItem
              key={comment.id}
              comment={comment}
              onReply={handleReply}
              onEdit={handleEdit}
              onDelete={handleDelete}
              replyToId={replyToId}
              setReplyToId={setReplyToId}
            />
          ))
        )}
      </div>
    </div>
  );
}; 