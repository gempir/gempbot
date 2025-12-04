import { useEffect, useState } from 'react';
import { useStore } from '../store';
import { doFetch } from '../service/doFetch';

export interface BotConfig {
    predictionAnnouncements: boolean;
}

export function useBotConfig() {
    const [config, setConfig] = useState<BotConfig | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const managing = useStore((state) => state.managing);
    const apiBaseUrl = useStore((state) => state.apiBaseUrl);

    const fetchConfig = async () => {
        if (!managing) {
            setLoading(false);
            return;
        }

        try {
            setLoading(true);
            const response = await doFetch(`${apiBaseUrl}/api/bot/config?channelId=${managing}`, {
                method: 'GET',
            });

            if (response.ok) {
                const data = await response.json();
                setConfig(data);
            } else {
                setError('Failed to load bot configuration');
            }
        } catch (err) {
            setError('Error loading bot configuration');
        } finally {
            setLoading(false);
        }
    };

    const updateConfig = async (updates: Partial<BotConfig>) => {
        if (!managing) return;

        try {
            const response = await doFetch(`${apiBaseUrl}/api/bot/config`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    channelId: managing,
                    ...updates,
                }),
            });

            if (response.ok) {
                const data = await response.json();
                setConfig(data);
            } else {
                throw new Error('Failed to update configuration');
            }
        } catch (err) {
            setError('Error updating bot configuration');
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
