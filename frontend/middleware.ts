import { NextRequest, NextResponse } from 'next/server';

/**
 * Next.js middleware for route protection.
 * Redirects unauthenticated users to login page when accessing protected routes.
 * Protected routes include /forms, /dashboard, and /admin paths.
 *
 * @param req - The incoming request
 * @returns Redirect response to login page or allows request to continue
 */
export function middleware(req: NextRequest) {
    const token = req.cookies.get('auth_token')?.value;
    const pathname = req.nextUrl.pathname;

    const isProtected =
        pathname.startsWith('/forms') ||
        pathname.startsWith('/dashboard') ||
        pathname.startsWith('/admin');

    if (isProtected && !token) {
        const loginUrl = new URL('/login', req.url);
        loginUrl.searchParams.set('next', pathname);
        return NextResponse.redirect(loginUrl);
    }

    return NextResponse.next();
}

export const config = {
    matcher: ['/forms/:path*', '/dashboard/:path*', '/admin/:path*'],
};
