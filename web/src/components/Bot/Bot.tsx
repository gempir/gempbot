import { useEffect, useState } from 'react';
import { useBotConfig } from '../../hooks/useBotConfig';
import { useSubscribtions } from '../../hooks/useSubscriptions';
import { Toggle } from './Toggle';

export function Bot() {
    const [subscribe, unsubscribe, subscriptionsStatus, loading] = useSubscribtions();

    const [predictionsAnnouncements, setPredictionAnnouncements] = useState(false);

    useEffect(() => {
        setPredictionAnnouncements(subscriptionsStatus.predictions);
    }, [subscriptionsStatus.predictions]);

    const handlePredictionAnnouncementChange = (value: boolean) => {
        setPredictionAnnouncements(value);
        if (value) {
            subscribe();
        } else {
            unsubscribe();
        }
    };

    const [botConfig, setBotConfig, loadingUserConfig] = useBotConfig();
    const handlePredictionCommandsChange = (value: boolean) => {
        if (botConfig) {
            setBotConfig({ ...botConfig, JoinBot: value });
        }
    };

    return <div className={"p-4"}>
        <div className={"bg-gray-800 rounded shadow relative p-4 " + (loading ? "animate-pulse pointer-events-none" : "")}>
            <div className="flex items-start justify-between">
                <div>
                    <h3 className="font-bold text-xl">Prediction Announcements</h3>
                    <div className="p-2 text-gray-200 mx-0 px-0">
                        Announces when predictions
                        <ul className="list-disc pl-6 mt-2">
                            <li>are made</li>
                            <li>locked</li>
                            <li>canceled</li>
                            <li>resolved</li>
                        </ul>
                        <img src={"/images/announcement.png"} className="mt-2" />
                    </div>
                </div>
                <Toggle checked={predictionsAnnouncements} onChange={handlePredictionAnnouncementChange} />
            </div>
        </div>
        <div className={"bg-gray-800 rounded shadow relative p-4 mt-4 " + (loadingUserConfig ? "animate-pulse pointer-events-none" : "")}>
            <div className="flex items-start justify-between">
                <div>
                    <h3 className="font-bold text-xl">Prediction Commands</h3>
                    <div className="p-2 text-gray-200 mx-0 px-0">
                        Commands to manage predictions.<br />
                        <strong>Default:</strong> <span className="font-mono">1m;yes;no</span>
                        <ul className="list-disc pl-6 font-mono mt-2">
                            <li>!prediction Will she win</li>
                            <li>!prediction Will she win;maybe</li>
                            <li>!prediction Who will win this game?;5m;NymN;Yabbe;Zoil;Pokimane;gempir;Forsen;xQc;pajlada</li>
                            <li className="mt-2">!prediction lock</li>
                            <li>!prediction cancel</li>
                            <li className="mt-2">!outcome 1</li>
                            <li>!outcome 2</li>
                            <li>!outcome 10</li>
                            <li>!outcome yabbe</li>
                        </ul>
                    </div>
                </div>
                <Toggle checked={!!botConfig?.JoinBot} onChange={handlePredictionCommandsChange} />
            </div>
        </div>
    </div >;
}


