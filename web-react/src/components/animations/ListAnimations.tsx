/**
 * List item animations
 */
import React from 'react';
import { motion } from 'framer-motion';

interface ListItemAnimationProps {
  children: React.ReactNode;
  index: number;
  layout?: boolean;
}

export const AnimatedListItem: React.FC<ListItemAnimationProps> = ({
  children,
  index,
  layout = true,
}) => {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, x: -20 }}
      transition={{
        duration: 0.3,
        delay: index * 0.05,
      }}
      layout={layout}
      whileHover={{ scale: 1.02 }}
      whileTap={{ scale: 0.98 }}
    >
      {children}
    </motion.div>
  );
};

/**
 * Card animation
 */
export const AnimatedCard: React.FC<{
  children: React.ReactNode;
  index?: number;
}> = ({ children, index = 0 }) => {
  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.9 }}
      animate={{ opacity: 1, scale: 1 }}
      transition={{
        duration: 0.3,
        delay: index * 0.1,
      }}
      whileHover={{ y: -4, boxShadow: '0 10px 20px rgba(0,0,0,0.1)' }}
    >
      {children}
    </motion.div>
  );
};

