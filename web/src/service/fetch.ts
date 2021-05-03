

export function handleResponse(response: Response) {
    if (!response.ok) { throw response }
    return response.json()
}