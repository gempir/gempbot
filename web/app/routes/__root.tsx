import { AppShell, ColorSchemeScript, MantineProvider } from "@mantine/core";
import { Notifications } from "@mantine/notifications";
import { NavigationProgress, nprogress } from "@mantine/nprogress";
import { createRootRoute, Outlet } from "@tanstack/react-router";
import dayjs from "dayjs";
import localizedFormat from "dayjs/plugin/localizedFormat";
import { jwtDecode } from "jwt-decode";
import { useEffect } from "react";
import { Sidebar } from "../components/Sidebar/Sidebar";
import { parseCookie } from "../service/cookie";
import type { ScTokenContent } from "../store";
import { initializeStore } from "../store";
import { theme } from "../theme";

import "@mantine/core/styles.css";
import "@mantine/notifications/styles.css";
import "@mantine/nprogress/styles.css";
import "../styles/globals.css";

dayjs.extend(localizedFormat);

export const Route = createRootRoute({
  component: RootComponent,
});

function RootComponent() {
  // Initialize store with client-side cookie data on mount
  useEffect(() => {
    nprogress.start();

    const cookies = parseCookie(document.cookie);

    let scTokenContent: ScTokenContent | null = null;
    if (cookies.scToken) {
      try {
        scTokenContent = jwtDecode<ScTokenContent>(cookies.scToken);
      } catch (e) {
        console.error("Failed to decode scToken:", e);
      }
    }

    initializeStore({
      scTokenContent,
      scToken: cookies.scToken || null,
      managing: cookies.managing || null,
      twitchClientId: import.meta.env.VITE_TWITCH_CLIENT_ID || "",
      apiBaseUrl: import.meta.env.VITE_API_BASE_URL || "",
      yjsWsUrl: import.meta.env.VITE_YJS_WS_URL || "",
      baseUrl: import.meta.env.VITE_BASE_URL || "",
    });

    nprogress.complete();
  }, []);

  return (
    <>
      <ColorSchemeScript forceColorScheme="dark" />
      <MantineProvider theme={theme} forceColorScheme="dark">
        <NavigationProgress />
        <Notifications position="top-right" />

        <AppShell navbar={{ width: 280, breakpoint: "sm" }} padding="md">
          <AppShell.Navbar p="md">
            <Sidebar />
          </AppShell.Navbar>

          <AppShell.Main>
            <Outlet />
          </AppShell.Main>
        </AppShell>
      </MantineProvider>
    </>
  );
}
