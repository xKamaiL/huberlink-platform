'use client'
import { Icons } from '@/components/icons'
import { cn } from '@/lib/utils'
import Link from 'next/link'
import { usePathname } from 'next/navigation'

const menuClass = (active: boolean) =>
  cn(
    `w-full focus:text-primary hover:primary justify-center inline-block text-center pt-2 pb-1 text-base`,
    active ? 'text-indigo-500' : 'text-slate-900'
  )

const BottomNavigation = () => {
  const path = usePathname()
  const currentHome = '121'
  return (
    <div
      className={cn`block fixed inset-x-0 bottom-0 z-10 bg-white shadow border-t mb-0 p-safe`}
    >
      <div className="flex justify-between  container mx-auto">
        <Link href={`/h/${currentHome}`} className={menuClass(false)}>
          <>
            <Icons.home className="w-5 h-4 mx-auto" />
            <span className="block text-xs">Home</span>
          </>
        </Link>
        <Link
          href={`/h/${currentHome}/automation`}
          className={menuClass(false)}
        >
          <>
            <Icons.bot className="w-5 h-5 mx-auto" />
            <span className="block text-xs">Automation</span>
          </>
        </Link>
        <Link href={`/h/${currentHome}/devices`} className={menuClass(false)}>
          <>
            <Icons.network className="w-5 h-4 mx-auto" />
            <span className="block text-xs">Devices</span>
          </>
        </Link>
        <Link href="/account" className={menuClass(path === '/account')}>
          <>
            <Icons.settings className="w-5 h-4 mx-auto" />
            <span className="block text-xs">Account</span>
          </>
        </Link>
      </div>
    </div>
  )
}

export default BottomNavigation

BottomNavigation.displayName = 'BottomNavigation'