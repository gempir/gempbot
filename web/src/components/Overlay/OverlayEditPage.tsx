const Editor = dynamic(async () => (await import('./Editor')).CustomEditor, { ssr: false })
import dynamic from "next/dynamic";
import { useParams } from "next/navigation";
import { useOverlay } from "../../hooks/useOverlays";


export function OverlayEditPage() {
    const params = useParams<{ overlayId: string }>();
    const [overlayResponse] = useOverlay(params.overlayId);

    return <div className="relative w-full h-[100vh]">
        {overlayResponse &&
            <Editor roomID={overlayResponse.overlay.RoomID} />
        }
    </div>;
}