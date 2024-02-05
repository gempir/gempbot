'use client';

import dynamic from 'next/dynamic';
import Head from 'next/head';

const Editor = dynamic(async () => (await import('./Editor')).Editor, { ssr: false })

export function OverlayPage() {
    return (
        <div className="relative w-full h-[100vh]">
            <Head>
                <style>{`
                    body {
                        background-color: transparent !important;
                    }

                    .tl-background {
                        background-color: transparent !important;
                    }
                `}</style>
            </Head>
            <Editor hideUi />
        </div>
    );
}