import { NextRequest, NextResponse } from 'next/server';

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
