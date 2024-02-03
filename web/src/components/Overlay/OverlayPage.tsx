'use client';

import dynamic from 'next/dynamic';

const Tldraw = dynamic(async () => (await import('@tldraw/tldraw')).Tldraw, { ssr: false })
import '@tldraw/tldraw/tldraw.css'

export function OverlayPage() {
    return (
        <div className="relative w-full">
            <Tldraw inferDarkMode />
        </div>
    );
}