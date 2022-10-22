import { useStore } from "../store";

export enum Method {
    GET = "GET",
    POST = "POST",
    DELETE = "DELETE",
    PUT = "PUT",
    PATCH = "PATCH"
}

export enum RejectReason {
    NotFound = '{"message":"Not Found"}'
}

export async function doFetch(method: Method, path: string, params: URLSearchParams = new URLSearchParams(), body: any = undefined) {
    const { apiBaseUrl, managing, scToken } = useStore.getState();

    const headers: Record<string, string> = { 'content-type': 'application/json' }
    if (scToken) {
        headers['Authorization'] = `Bearer ${scToken}`
    }

    const config: RequestInit = {
        method: method,
        headers: {
            ...headers,
        },
    }

    if (body) {
        config.body = JSON.stringify(body)
    }

    const url = new URL(path, apiBaseUrl);
    if (managing) {
        url.searchParams.append('managing', managing);
    }

    params.forEach((value, key) => url.searchParams.append(key, value));

    return window.fetch(url.toString(), config)
        .then(async response => {
            if (response.status === 401) {
                useStore.setState({ scToken: undefined });
            }
            if (response.ok) {
                const text = await response.text();
                return text !== "" ? JSON.parse(text) : null
            } else {
                if (response.status === 404) {
                    return Promise.reject(RejectReason.NotFound);
                }

                const errorMessage = await response.text()
                return Promise.reject(errorMessage)
            }
        })
}