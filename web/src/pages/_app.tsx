import dayjs from 'dayjs';
import * as localizedFormat from 'dayjs/plugin/localizedFormat';
import Head from "next/head";
import Link from "next/link";
import 'tailwindcss/tailwind.css';
import { Sidebar } from "../components/Sidebar/Sidebar";
import { StoreProvider, useCreateStore } from "../store";

// @ts-ignore
dayjs.extend(localizedFormat);

export default function App({ Component, pageProps }: { Component: any; pageProps: any }) {
    const createStore = useCreateStore(pageProps.store);

    const renderFullLayout = Component.name !== "Overlay";

    return (
        <StoreProvider createStore={createStore}>
            <Head>
                <title>gempbot</title>
                <link rel="icon" href="/favicon.ico" />
            </Head>
            <style jsx global>{`
                body {
                    --tw-bg-opacity: 1;
                    min-height: 100vh;
                    background-color: ${renderFullLayout ? "rgba(17, 24, 39, var(--tw-bg-opacity))" : "transparent"};
                    line-height: 1.25;
                    --tw-text-opacity: 1;
                    color: rgba(209, 213, 219, var(--tw-text-opacity));
                }
            `}</style>
            {renderFullLayout && <>
                <main>
                    <div className="flex" style={{ scrollbarGutter: "stable" }}>
                        <Sidebar />
                        <Component {...pageProps} />
                    </div>
                    <div className="absolute bottom-3 text-center right-3 mx-auto hover:text-blue-500">
                        <Link href="/privacy">
                            Privacy
                        </Link>
                    </div>
                </main>
            </>}
            {!renderFullLayout && <>
                <Component {...pageProps} />
            </>}
        </StoreProvider>
    );
}