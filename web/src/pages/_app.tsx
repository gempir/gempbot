import Head from "next/head";
import 'tailwindcss/tailwind.css';
import { Sidebar } from "../components/Sidebar/Sidebar";
import { StoreProvider, useCreateStore } from "../store";

export default function App({ Component, pageProps }: { Component: any; pageProps: any }) {
    const createStore = useCreateStore(pageProps.store);

    return <StoreProvider createStore={createStore}>
        <Head>
            <title>gempbot</title>
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
    </StoreProvider>
}