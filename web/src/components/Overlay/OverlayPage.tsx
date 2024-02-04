'use client';

import dynamic from 'next/dynamic';

const Editor = dynamic(async () => (await import('./Editor')).Editor, { ssr: false })

export function OverlayPage() {
    return (
        <div className="relative w-full">
            <Editor />
        </div>
    );
}