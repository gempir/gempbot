import { useState } from "react";
import { doFetch, Method } from "../service/doFetch";
import { useStore } from "../store";

export function useSubscribtions(): [() => void, () => void, boolean, boolean] {
    const managing = useStore(state => state.managing);
    const [loadingSubscribe, setLoadingSubscribe] = useState(false);
    const [loadingUnsubscribe, setLoadingUnsubscribe] = useState(false);

    const executeSubscriptions = (method: Method) => {
        const endPoint = "/api/subscriptions";
        const searchParams = new URLSearchParams();
        if (managing) {
            searchParams.append("managing", managing);
        }

        return doFetch(method, endPoint, searchParams);
    };

    const subscribe = () => {
        if (!confirm("Subscribe to all prediction events, to log them here and write to chat")) {
            return
        }
        setLoadingSubscribe(true);

        executeSubscriptions(Method.PUT).then(() => setLoadingSubscribe(false)).catch(() => setLoadingSubscribe(false));
    }
    const remove = () => {
        if (!confirm("Unsubscribe prediction events")) {
            return
        }
        setLoadingUnsubscribe(true);

        executeSubscriptions(Method.DELETE).then(() => setLoadingUnsubscribe(false)).catch(() => setLoadingSubscribe(false));
    }

    return [subscribe, remove, loadingSubscribe, loadingUnsubscribe];
}