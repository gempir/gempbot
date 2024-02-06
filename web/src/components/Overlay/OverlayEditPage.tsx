const Editor = dynamic(async () => (await import('./Editor')).Editor, { ssr: false })
import dynamic from "next/dynamic";
import { useOverlay } from "../../hooks/useOverlays";
import { useParams } from "next/navigation";

export function OverlayEditPage() {
    const params = useParams<{ overlayId: string }>();
    const [overlay] = useOverlay(params.overlayId);

    console.log("Joining", overlay?.RoomID);

    return <div className="relative w-full h-[100vh]">
        {overlay?.RoomID && <Editor roomId={overlay.RoomID} />}
    </div>;
}