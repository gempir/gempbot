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
        if (!confirm("This will re-subscribe all eventsub subscriptions to your channel.\nOnly do this if you know what you are doing")) {
            return
        }

        executeSubscriptions(Method.PUT);
    }
    const remove = () => {
        if (!confirm("This will remove all eventsub subscriptions to your channel.\nOnly do this if you know what you are doing")) {
            return
        }

        executeSubscriptions(Method.DELETE);
    }

    return [subscribe, remove];
}