import { useEffect, useState } from "react";
import { useStore } from "../store";
import { doFetch, Method } from "../service/doFetch";

export interface BotConfig {
  predictionAnnouncements: boolean;
}

export function useBotConfig() {
  const [config, setConfig] = useState<BotConfig | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const managing = useStore((state) => state.managing);
  const apiBaseUrl = useStore((state) => state.apiBaseUrl);
  const scToken = useStore((state) => state.scToken);

  const fetchConfig = async () => {
    try {
      setLoading(true);
      const searchParams = new URLSearchParams();
      if (managing) {
        searchParams.append("channelId", managing);
      }
      const data = await doFetch(
        { apiBaseUrl, managing, scToken },
        Method.GET,
        "/api/bot/config",
        searchParams,
      );
      setConfig(data);
      setError(null);
    } catch (err) {
      setError("Error loading bot configuration");
    } finally {
      setLoading(false);
    }
  };

  const updateConfig = async (updates: Partial<BotConfig>) => {
    try {
      const body: any = { ...updates };
      if (managing) {
        body.channelId = managing;
      }

      const data = await doFetch(
        { apiBaseUrl, managing, scToken },
        Method.POST,
        "/api/bot/config",
        undefined,
        body,
      );
      setConfig(data);
      setError(null);
    } catch (err) {
      setError("Error updating bot configuration");
      throw err;
    }
  };

  useEffect(() => {
    fetchConfig();
  }, [managing]);

  return {
    config,
    loading,
    error,
    updateConfig,
    refetch: fetchConfig,
  };
}
