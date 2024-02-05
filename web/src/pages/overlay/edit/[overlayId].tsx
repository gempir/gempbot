'use client';

import dynamic from "next/dynamic";
import { useParams } from "next/navigation";
const Editor = dynamic(async () => (await import('../../../components/Overlay/Editor')).Editor, { ssr: false })

export default function OverlaysEditPage() {
    const params = useParams<{ overlayId: string }>();

    return <div className="relative w-full h-[100vh]">
        <Editor overlayId={params.overlayId} />
    </div>;
}