import { cn } from '@/lib/utils';

interface RaIconProps {
  name: string;
  className?: string;
}

export function RaIcon({ name, className }: RaIconProps) {
  return <i className={cn('ra', `ra-${name}`, className)} aria-hidden="true" />;
}
