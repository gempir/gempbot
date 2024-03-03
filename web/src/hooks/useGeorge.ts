import { Method } from '../service/doFetch';
import { useStore } from '../store';

type FetchStreamResponse = {
    data: string | null;
    error: Error | null;
    loading: boolean;
};

interface GeorgeRequest {
    channel: string;
    username: string;
    year: number;
    month: number;
    day: number;
    query: string;
    model: string;
}

export const useGeorge = (): (req: GeorgeRequest, onText: (text: string) => void) => void => {
    const apiBaseUrl = useStore(state => state.apiBaseUrl);
    const scToken = useStore(state => state.scToken);

    const request = async (req: GeorgeRequest, onText: (text: string) => void) => {
        const response = await fetch(apiBaseUrl + "/api/george", {
            method: Method.POST,
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${scToken}`
            },
            body: JSON.stringify(req)
        });

        const reader = response.body?.getReader();
        if (!reader) {
            return "";
        }

        let result = '';
        while (true) {
            const { done, value } = await reader.read();
            if (done) break;
            result += new TextDecoder().decode(value);

            onText(result);
        }
    };

    return request;
};
