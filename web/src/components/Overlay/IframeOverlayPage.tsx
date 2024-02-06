'use client';

import dynamic from 'next/dynamic';
import Head from 'next/head';
import { useParams } from 'next/navigation';

const Editor = dynamic(async () => (await import('./Editor')).Editor, { ssr: false })

export function IframeOverlayPage() {
    const params = useParams<{ roomId: string }>();

    console.log("Joining", params.roomId);

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
            <Editor hideUi roomId={params.roomId} readonly />
        </div>
    );
}