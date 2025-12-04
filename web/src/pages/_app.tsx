import { AppShell, ColorSchemeScript, MantineProvider } from "@mantine/core";
import { Notifications } from "@mantine/notifications";
import dayjs from "dayjs";
import * as localizedFormat from "dayjs/plugin/localizedFormat";
import Head from "next/head";
import { useEffect } from "react";
import { Sidebar } from "../components/Sidebar/Sidebar";
import { initializeStore } from "../store";
import { theme } from "../theme";

import "@mantine/core/styles.css";
import "@mantine/notifications/styles.css";

// @ts-ignore
dayjs.extend(localizedFormat);

export default function App({
  Component,
  pageProps,
}: {
  Component: any;
  pageProps: any;
}) {
  const renderFullLayout = pageProps.renderFullLayout ?? true;

  // Initialize store with server state on mount
  useEffect(() => {
    if (pageProps.store) {
      initializeStore(pageProps.store);
    }
  }, [pageProps.store]);

  return (
    <>
      <Head>
        <title>gempbot - Twitch Bot & Overlay Manager</title>
        <link rel="icon" href="/favicon.ico" />
        <meta
          name="viewport"
          content="minimum-scale=1, initial-scale=1, width=device-width, user-scalable=no"
        />
        <ColorSchemeScript />
      </Head>

      <MantineProvider theme={theme} forceColorScheme="dark">
        <Notifications position="top-right" />

        {renderFullLayout ? (
          <AppShell navbar={{ width: 280, breakpoint: "sm" }} padding="md">
            <AppShell.Navbar p="md">
              <Sidebar />
            </AppShell.Navbar>

            <AppShell.Main>
              <Component {...pageProps} />
            </AppShell.Main>
          </AppShell>
        ) : (
          <Component {...pageProps} />
        )}
      </MantineProvider>
    </>
  );
}
