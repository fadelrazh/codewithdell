import React from 'react';

interface ButtonProps {
  children: React.ReactNode;
  variant?: 'primary' | 'secondary' | 'outline' | 'ghost';
  size?: 'sm' | 'md' | 'lg';
  className?: string;
  onClick?: () => void;
  disabled?: boolean;
  type?: 'button' | 'submit' | 'reset';
  style?: React.CSSProperties;
}

export const Button: React.FC<ButtonProps> = ({
  children,
  variant = 'primary',
  size = 'md',
  className = '',
  onClick,
  disabled = false,
  type = 'button',
  style,
  ...props
}) => {
  const baseClasses = 'inline-flex items-center justify-center rounded-lg font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed';
  
  const sizeClasses = {
    sm: 'px-3 py-1.5 text-sm',
    md: 'px-4 py-2 text-base',
    lg: 'px-6 py-3 text-lg'
  };

  const variantClasses = {
    primary: 'text-white border-0',
    secondary: 'text-gray-900 border-0',
    outline: 'border-2 bg-transparent',
    ghost: 'bg-transparent border-0'
  };

  const getVariantStyles = () => {
    switch (variant) {
      case 'primary':
        return {
          backgroundColor: 'var(--primary)',
          color: 'white'
        };
      case 'secondary':
        return {
          backgroundColor: 'var(--secondary)',
          color: '#111827'
        };
      case 'outline':
        return {
          backgroundColor: 'transparent',
          borderColor: 'var(--primary)',
          color: 'var(--primary)'
        };
      case 'ghost':
        return {
          backgroundColor: 'transparent',
          color: 'var(--primary)'
        };
      default:
        return {};
    }
  };

  const classes = `${baseClasses} ${sizeClasses[size]} ${variantClasses[variant]} ${className}`;
  const variantStyles = getVariantStyles();
  const combinedStyles = { ...variantStyles, ...style };

  return (
    <button
      type={type}
      className={classes}
      style={combinedStyles}
      onClick={onClick}
      disabled={disabled}
      {...props}
    >
      {children}
    </button>
  );
}; 