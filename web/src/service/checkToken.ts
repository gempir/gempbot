export function checkToken(setScToken: (scToken: string | null) => void, response: Response) {
    if (response.status === 403) {
        setScToken(null);
    }
}