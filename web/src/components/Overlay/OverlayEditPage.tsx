const Editor = dynamic(async () => (await import('./Editor')).CustomEditor, { ssr: false })
import dynamic from "next/dynamic";
import { useParams } from "next/navigation";
import { useOverlay } from "../../hooks/useOverlays";
import { useStore } from "../../store";


export function OverlayEditPage() {
    const params = useParams<{ overlayId: string }>();
    const [overlayResponse] = useOverlay(params.overlayId);
    const scTokenContent = useStore(state => state.scTokenContent);

    return <div className="relative w-full h-[100vh]">
        {overlayResponse &&
            <Editor roomID={overlayResponse.overlay.RoomID} userID={scTokenContent?.UserID} username={scTokenContent?.Login} />
        }
    </div>;
}