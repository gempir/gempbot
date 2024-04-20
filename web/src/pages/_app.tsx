import { ColorSchemeScript, MantineProvider, createTheme } from '@mantine/core';
import dayjs from 'dayjs';
import * as localizedFormat from 'dayjs/plugin/localizedFormat';
import Head from "next/head";
import Link from "next/link";
import 'tailwindcss/tailwind.css';
import { Sidebar } from "../components/Sidebar/Sidebar";
import { StoreProvider, useCreateStore } from "../store";

import '@mantine/core/styles.css';
import '@mantine/dates/styles.css';
import '@mantine/dropzone/styles.css';

// @ts-ignore
dayjs.extend(localizedFormat);


const theme = createTheme({
});

export default function App({ Component, pageProps }: { Component: any; pageProps: any }) {
    const createStore = useCreateStore(pageProps.store);

    const renderFullLayout = pageProps.renderFullLayout ?? true;

    return (
        <StoreProvider createStore={createStore}>
            <Head>
                <title>gempbot</title>
                <link rel="icon" href="/favicon.ico" />
                <ColorSchemeScript />
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
            <MantineProvider theme={theme} forceColorScheme="dark">
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
            </MantineProvider>
        </StoreProvider>
    );
}