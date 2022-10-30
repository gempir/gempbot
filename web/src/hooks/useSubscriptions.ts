import { useEffect, useState } from "react";
import { doFetch, Method } from "../service/doFetch";
import { useStore } from "../store";

interface SubscriptionStatus {
    predictions: boolean;
}

export function useSubscribtions(): [() => void, () => void, SubscriptionStatus, boolean] {
    const managing = useStore(state => state.managing);
    const [loading, setLoading] = useState(true);
    const [subscriptionStatus, setSubscriptionStatus] = useState<SubscriptionStatus>({ predictions: false });
    const apiBaseUrl = useStore(state => state.apiBaseUrl);
    const scToken = useStore(state => state.scToken);

    const executeSubscriptions = (method: Method) => {
        const endPoint = "/api/subscriptions";
        const searchParams = new URLSearchParams();
        if (managing) {
            searchParams.append("managing", managing);
        }

        return doFetch({apiBaseUrl, managing, scToken }, method, endPoint, searchParams);
    };

    const subscribe = () => {
        setLoading(true);

        executeSubscriptions(Method.PUT).then(() => setSubscriptionStatus({predictions: true})).then(() => setLoading(false)).catch(err => {
            console.error(err);
            setLoading(false);
        });
    }
    const remove = () => {
        setLoading(true);

        executeSubscriptions(Method.DELETE).then(() => setSubscriptionStatus({predictions: false})).then(() => setLoading(false)).catch(err => {
            console.error(err);
            setLoading(false);
        });
    }
    const status = () => {
        setLoading(true);

        executeSubscriptions(Method.GET).then(setSubscriptionStatus).then(() => setLoading(false)).catch(err => {
            console.error(err);
            setLoading(false);
        });
    }

    useEffect(() => {
        status();
    }, [managing]);

    return [subscribe, remove, subscriptionStatus, loading];
}