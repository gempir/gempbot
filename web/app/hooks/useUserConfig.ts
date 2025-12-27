import { useCallback, useEffect, useState } from "react";
import { useDebounce } from "react-use";
import { doFetch, Method } from "../service/doFetch";
import { useStore } from "../store";

export interface UserConfig {
  Permissions: Record<string, Permission>;
  Protected: Protected;
}

export interface Permission {
  Editor: boolean;
  Prediction: boolean;
}

export interface Rewards {
  SevenTv: null | SevenTvReward;
}

export interface SevenTvReward {
  title: string;
  prompt?: string;
  cost: number;
  backgroundColor?: string;
  isMaxPerStreamEnabled?: boolean;
  maxPerStream?: number;
  isUserInputRequired?: boolean;
  isMaxPerUserPerStreamEnabled?: boolean;
  maxPerUserPerStream?: number;
  isGlobalCooldownEnabled?: boolean;
  globalCooldownSeconds?: number;
  shouldRedemptionsSkipRequestQueue?: boolean;
  enabled?: boolean;
  isDefault: boolean;
  ID?: string;
}

export interface Protected {
  EditorFor: Array<string>;
  CurrentUserID: string;
}

export type SetUserConfig = (userConfig: UserConfig | null) => void;

export function useUserConfig(): [
  UserConfig | null | undefined,
  SetUserConfig,
  () => void,
  boolean,
  string | undefined,
] {
  const [userConfig, setUserConfig] = useState<UserConfig | null | undefined>(
    undefined,
  );
  const [errorMessage, setError] = useState<string | undefined>();
  const [loading, setLoading] = useState(true);
  const [changeCounter, setChangeCounter] = useState(0);
  const managing = useStore((state) => state.managing);
  const scToken = useStore((state) => state.scToken);
  const apiBaseUrl = useStore((state) => state.apiBaseUrl);

  const fetchConfig = useCallback(() => {
    if (!scToken) {
      setLoading(false);
      return;
    }

    setLoading(true);
    const endPoint = "/api/userconfig";
    doFetch({ apiBaseUrl, managing, scToken }, Method.GET, endPoint)
      .then((userConfig) => setUserConfig(userConfig))
      .then(() => setLoading(false))
      .catch(setError);
  }, [scToken, managing, apiBaseUrl]);

  useEffect(() => {
    fetchConfig();
  }, [fetchConfig]);

  useDebounce(
    () => {
      if (changeCounter && userConfig) {
        const endPoint = "/api/userconfig";
        setLoading(true);
        doFetch(
          { apiBaseUrl, managing, scToken },
          Method.POST,
          endPoint,
          undefined,
          userConfig,
        )
          .then(fetchConfig)
          .then(() => setError(undefined))
          .catch((error) => setError(JSON.parse(error).message))
          .finally(() => setLoading(false));
      } else if (changeCounter && userConfig === null) {
        setLoading(true);
        doFetch(
          { apiBaseUrl, managing, scToken },
          Method.DELETE,
          "/api/userconfig",
        ).finally(() => setLoading(false));
      }
    },
    100,
    [changeCounter],
  );

  const setCfg = (config: UserConfig | null) => {
    setUserConfig(config);
    setChangeCounter(changeCounter + 1);
  };

  return [userConfig, setCfg, fetchConfig, loading, errorMessage];
}
