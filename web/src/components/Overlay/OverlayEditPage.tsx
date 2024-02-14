const Editor = dynamic(async () => (await import('./Editor')).CustomEditor, { ssr: false })
import { YDocProvider } from '@y-sweet/react';
import dynamic from "next/dynamic";
import { useParams } from "next/navigation";
import { useOverlay } from "../../hooks/useOverlays";


export function OverlayEditPage() {
    const params = useParams<{ overlayId: string }>();
    const [overlayAuth] = useOverlay(params.overlayId);


    return <div className="relative w-full h-[100vh]">
        {overlayAuth &&
            <YDocProvider clientToken={overlayAuth.auth}>
                <Editor />
            </YDocProvider>
        }
    </div>;
}