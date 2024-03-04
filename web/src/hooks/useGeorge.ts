import { Method } from '../service/doFetch';
import { useStore } from '../store';

interface GeorgeRequest {
    channel: string;
    username: string;
    year: number;
    month: number;
    day: number;
    query: string;
    model: string;
}

type RequestFunc = (req: GeorgeRequest, controller: AbortController, onText: (text: string) => void, onQuery: (text: string) => void) => void;

export const useGeorge = (): [RequestFunc] => {
    const apiBaseUrl = useStore(state => state.apiBaseUrl);
    const scToken = useStore(state => state.scToken);

    const request = async (req: GeorgeRequest, controller: AbortController, onText: (text: string) => void, onQuery: (text: string) => void) => {
        try {
            const response = await fetch(apiBaseUrl + "/api/george", {
                method: Method.POST,
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${scToken}`
                },
                body: JSON.stringify(req),
                signal: controller.signal
            });

            const reader = response.body?.getReader();
            if (!reader) {
                onText("Failed to read response body");
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

                if (!queryDone && textValue.includes("====QUERYDONE====")) {
                    onQuery(result.replace("====QUERYDONE====", ""));
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
        } catch (error) {
            const err = error as Error;
            onText(err?.message);
            return;
        }
    };

    return [request];
};
