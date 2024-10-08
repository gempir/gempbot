'use client';

import dynamic from 'next/dynamic';
import Head from 'next/head';
import { useParams } from 'next/navigation';

const Editor = dynamic(async () => (await import('./Editor')).CustomEditor, { ssr: false })

export function IframeOverlayPage() {
    const params = useParams<{ roomId: string }>();

    return (
        <div className="relative w-full h-[100vh]">
            <Head>
                <style>{`
                    body {
                        background-color: transparent !important;
                    }

                    .tl-background__wrapper, .tl-background, .tl-canvas {
                        background-color: transparent !important;
                    }

                    .tl-loading {
                        display: none !important;
                    }

                    .tl-cursor, .tl-nametag, .tl-nametag-title, .tl-nametag-chat {
                        display: none !important;
                    }

                    #cursor_hint, .tl-collaborator__cursor-hint {
                        display: none !important;
                        opacity: 0 !important;
                    }

                    .tl-overlays__item {
                        opacity: 0 !important;
                    }

                    /** Please don't hate me tldraw, I can't show this in the overlay, that would suck for the stream. But it's still visible for the editors. **/
                    
                    .tl-watermark_SEE-LICENSE { 
                        display: none !important;
                    }
                `}</style>
            </Head>
            <Editor hideUi readonly roomID={params.roomId} components={{Cursor: null}} />
        </div>
    );
}