import React from "react";
import { createLoginUrl } from "../factory/createLoginUrl";
import { store } from "../store";

export function Login({ className, children }: { className?: string, children?: React.ReactNode }) {
    const { apiBaseUrl, twitchClientId } = store.getRawState();
    const url = createLoginUrl(apiBaseUrl, twitchClientId);

    return <a className={"p-4 rounded shadow opacity-25 bg-purple-800 hover:bg-purple-600 text-center"} href={url.toString()}>Login Again</a>
}