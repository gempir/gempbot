import { doFetch, Method } from "../service/doFetch";
import { useStore } from "../store";

export function useSubscribtions(): [() => void, () => void] {
    const managing = useStore(state => state.managing);

    const executeSubscriptions = (method: Method) => {
        const endPoint = "/api/subscriptions";
        const searchParams = new URLSearchParams();
        if (managing) {
            searchParams.append("managing", managing);
        }

        doFetch(method, endPoint, searchParams);
    };

    const subscribe = () => {
        if (!confirm("Subscribe to all prediction events, to log them here and write to chat")) {
            return
        }

        executeSubscriptions(Method.PUT);
    }
    const remove = () => {
        if (!confirm("Unsubscribe prediction events")) {
            return
        }

        executeSubscriptions(Method.DELETE);
    }

    return [subscribe, remove];
}