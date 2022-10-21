import Head from "next/head";
import React, { useEffect } from "react";
import { Sidebar } from "../components/Sidebar/Sidebar";
import 'tailwindcss/tailwind.css';
import { useStore } from "../store";

export default function App({ Component, pageProps }: { Component: any; pageProps: any }) {
    useEffect(() => {
        useStore.setState(pageProps.store);
    }, []);

    return <>
        <Head>
            <title>gempbot - fly.io</title>
            <link rel="icon" href="/favicon.ico" />
        </Head>
        <style jsx global>{`
            body {
                --tw-bg-opacity: 1;
                background-color: rgba(17, 24, 39, var(--tw-bg-opacity));
                line-height: 1.25;
                --tw-text-opacity: 1;
                color: rgba(209, 213, 219, var(--tw-text-opacity));
            }
        `}</style>
        <main>
            <div className="flex" style={{scrollbarGutter: "stable"}}>
                <Sidebar />
                <Component {...pageProps} />
            </div>
        </main>
    </>
}