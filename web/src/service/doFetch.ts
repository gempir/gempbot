import { store } from "../store";

export enum Method {
    GET = "GET",
    POST = "POST",
    DELETE = "POST"
}

export async function doFetch(method: Method, path: string, body: any = undefined) {
    const {apiBaseUrl, scToken} = store.getRawState();

    const headers: Record<string, string> = { 'content-type': 'application/json' }
    if (scToken) {
        headers.Authorization = `Bearer ${scToken}`
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

    return window.fetch(`${apiBaseUrl}${path}`, config)
        .then(async response => {
            if (response.status === 401) {
                window.location.assign("/")
                return
            }
            if (response.ok) {
                return await response.json()
            } else {
                const errorMessage = await response.text()
                return Promise.reject(new Error(errorMessage))
            }
        })
}