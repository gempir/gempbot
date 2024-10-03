import { useStore } from "../store";

export type Asset = {
    id: string;
    isAnimated: boolean;
    isVideo: boolean;
    mimeType: string;
    url: string;
}


export function useAssetUploader(): (file: File) => Promise<Asset> {
    const apiBaseUrl = useStore(state => state.apiBaseUrl);
    const scToken = useStore(state => state.scToken);

    const upload = (file: File): Promise<Asset> => {
        const endPoint = "/api/asset";
        
        const formData = new FormData();
        formData.append("file", file);

        return fetch(apiBaseUrl + endPoint, {
            method: "POST",
            headers: {
                Authorization: `Bearer ${scToken}`,
            },
            body: formData,
        }).then(response => response.json() as unknown as Asset);
    }

    return upload;
}