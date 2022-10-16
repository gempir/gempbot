import { useWs } from "../../hooks/useWs";
import React from "react";

export function MediaPage(): JSX.Element {
    const { lastJsonMessage } = useWs();

    return <div className="p-4 w-full max-h-screen flex gap-4">
        <div className="p-4 bg-gray-800 rounded shadow relative">
            {lastJsonMessage}
        </div>
    </div>;
}