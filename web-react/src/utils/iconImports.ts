/**
 * Optimized icon imports from lucide-react
 * Individual imports enable tree-shaking and reduce bundle size
 * 
 * Usage:
 * import { AlertTriangle, Shield } from '@/utils/iconImports';
 */

// Export commonly used icons individually for tree-shaking
export {
  AlertTriangle,
  Shield,
  TrendingUp,
  TrendingDown,
  Activity,
  Server,
  Zap,
  Eye,
  RefreshCw,
  Download,
  Settings,
  BarChart3,
  Clock,
  Monitor,
  Wifi,
  MapPin,
  Cpu,
  Tag,
  ArrowRight,
  Check,
  X,
  ChevronDown,
  ChevronUp,
  ChevronLeft,
  ChevronRight,
  Search,
  Filter,
  MoreVertical,
  Plus,
  Minus,
  Edit,
  Trash2,
  Copy,
  ExternalLink,
  Info,
  HelpCircle,
  CheckCircle,
  XCircle,
  AlertCircle,
  Loader2,
  Menu,
  X as Close,
} from 'lucide-react';

// Re-export types if needed
export type { LucideIcon, LucideProps } from 'lucide-react';

/**
 * Icon component wrapper for consistent sizing and styling
 */
import React from 'react';
import { LucideIcon } from 'lucide-react';
import { cn } from '@/lib/utils';

interface IconProps {
  icon: LucideIcon;
  size?: number | string;
  className?: string;
  strokeWidth?: number;
}

export const Icon: React.FC<IconProps> = ({ 
  icon: IconComponent, 
  size = 20, 
  className,
  strokeWidth = 2 
}) => {
  return (
    <IconComponent
      size={size}
      className={cn(className)}
      strokeWidth={strokeWidth}
    />
  );
};

