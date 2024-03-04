import { useRef } from 'react';
import { Method } from '../service/doFetch';
import { useStore } from '../store';
import { on } from 'events';

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

type RequestFunc = (req: GeorgeRequest, onText: (text: string) => void, onQuery: (text: string) => void) => void;

export const useGeorge = (): [RequestFunc, AbortController] => {
    const apiBaseUrl = useStore(state => state.apiBaseUrl);
    const scToken = useStore(state => state.scToken);

    const controller = useRef<AbortController>(new AbortController());

    const request = async (req: GeorgeRequest, onText: (text: string) => void, onQuery: (text: string) => void) => {
        const response = await fetch(apiBaseUrl + "/api/george", {
            method: Method.POST,
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${scToken}`
            },
            body: JSON.stringify(req),
            signal: controller.current.signal
        });

        const reader = response.body?.getReader();
        if (!reader) {
            return "";
        }

        let result = '';
        let queryDone = false;

        while (true) {
            const { done, value } = await reader.read();
            if (done) {
                onText("@DONE");
                break;
            }
            const textValue = new TextDecoder().decode(value);
            result += textValue;
            
            if (!queryDone && textValue.includes("==QUERYDONE==")) {
                onQuery(result);
                result = "";
                queryDone = true;
                continue;
            }

            if (queryDone) {
                onText(result);
            }
        }

        // something failed
        if (!queryDone) {
            onText(result);
        }
    };

    return [request, controller.current];
};
