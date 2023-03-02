// middleware.ts
import type { NextRequest } from 'next/server'
import { NextResponse } from 'next/server'
import { ResponseCode } from './lib/types'
import { doRefreshToken } from './services/rawFetch'

// This function can be marked `async` if using `await` inside
export async function middleware(req: NextRequest) {
  const baseURL = req.url
  const accessToken = req.cookies.get('accessToken')
  if (!accessToken) {
    return NextResponse.redirect(
      new URL('/auth/sign-in?redirect=' + req.nextUrl.pathname, baseURL)
    )
  }

  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/me`, {
    headers: {
      Authorization: `Bearer ${accessToken.value}`,
    },
  }).then((r) => r.json()) // <--- This is the problem
  if (!res.success) {
    if (res.code === ResponseCode.TokenExpired) {
      console.log('[INFO] Refreshing token')
      const refreshResult = await doRefreshToken(req)
      if (refreshResult) {
        console.log('[INFO] Refresh token success')
        req.cookies.set('accessToken', refreshResult.token)
        req.cookies.set('refreshToken', refreshResult.refreshToken)
        return
      }
      console.log('[INFO] Refresh token failed', {
        refreshResult,
      })
    }

    return NextResponse.redirect(
      new URL('/auth/sign-in?redirect=' + req.nextUrl.pathname, baseURL)
    )
  }
}

// See "Matching Paths" below to learn more
export const config = {
  matcher: ['/h/:path*', '/account/:path*'],
}
