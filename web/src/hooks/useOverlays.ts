import { useEffect, useState } from "react";
import { Method, doFetch } from "../service/doFetch";
import { useStore } from "../store";

type Overlay = {
    ID: string;
    RoomID: string;
}

export function useOverlays(): [Overlay[], () => void, (id: string) => void, string | null, boolean] {
    const [overlays, setOverlays] = useState<Overlay[]>([]);
    const [errorMessage, setErrorMessage] = useState<string | null>(null);
    const [loading, setLoading] = useState(false);
    const managing = useStore(state => state.managing);
    const apiBaseUrl = useStore(state => state.apiBaseUrl);
    const scToken = useStore(state => state.scToken);

    const fetchOverlays = () => {
        setLoading(true);
        const endPoint = "/api/overlay";

        doFetch({ apiBaseUrl, managing, scToken }, Method.GET, endPoint).then(setOverlays).catch(setErrorMessage).finally(() => setLoading(false));
    }

    useEffect(fetchOverlays, []);

    const addOverlay = () => {
        setLoading(true);
        const endPoint = "/api/overlay";

        doFetch({ apiBaseUrl, managing, scToken }, Method.POST, endPoint).then(fetchOverlays).catch(setErrorMessage).finally(() => setLoading(false));
    }

    const deleteOverlay = (id: string) => {
        setLoading(true);
        const endPoint = "/api/overlay";

        doFetch({ apiBaseUrl, managing, scToken }, Method.DELETE, endPoint, new URLSearchParams({id})).then(() => setErrorMessage(null)).then(fetchOverlays).catch(setErrorMessage).finally(() => setLoading(false));
    }

    return [overlays, addOverlay, deleteOverlay, errorMessage, loading];
}


export function useOverlay(id: string): [Overlay|null, boolean] {
    const [overlay, setOverlay] = useState<Overlay | null>(null);
    const [loading, setLoading] = useState(false);
    const managing = useStore(state => state.managing);
    const apiBaseUrl = useStore(state => state.apiBaseUrl);
    const scToken = useStore(state => state.scToken);

    const fetchOverlay = () => {
        setLoading(true);
        const endPoint = "/api/overlay";

        doFetch({ apiBaseUrl, managing, scToken }, Method.GET, endPoint, new URLSearchParams({id})).then(setOverlay).finally(() => setLoading(false));
    }

    useEffect(fetchOverlay, [id]);

    return [overlay, loading];
}

export function useOverlayByRoomId(roomId: string): [Overlay|null, boolean] {
    const [overlay, setOverlay] = useState<Overlay | null>(null);
    const [loading, setLoading] = useState(false);
    const managing = useStore(state => state.managing);
    const apiBaseUrl = useStore(state => state.apiBaseUrl);
    const scToken = useStore(state => state.scToken);

    const fetchOverlay = () => {
        setLoading(true);
        const endPoint = "/api/overlay";

        doFetch({ apiBaseUrl, managing, scToken }, Method.GET, endPoint, new URLSearchParams({roomId})).then(setOverlay).finally(() => setLoading(false));
    }

    useEffect(fetchOverlay, [roomId]);

    return [overlay, loading];
}