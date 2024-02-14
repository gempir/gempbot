'use client';

import { YDocProvider } from '@y-sweet/react';
import dynamic from 'next/dynamic';
import Head from 'next/head';
import { useParams } from 'next/navigation';
import { useOverlayByRoomId } from '../../hooks/useOverlays';

const Editor = dynamic(async () => (await import('./Editor')).CustomEditor, { ssr: false })

export function IframeOverlayPage() {
    const params = useParams<{ roomId: string }>();
    const [overlayAuth] = useOverlayByRoomId(params.roomId);

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

                    .tl-loading {
                        display: none !important;
                    }
                `}</style>
            </Head>
            {overlayAuth &&
                <YDocProvider clientToken={overlayAuth.auth}>
                    <Editor hideUi readonly />
                </YDocProvider>
            }
        </div>
    );
}