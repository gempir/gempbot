import { SpeakerWaveIcon, SpeakerXMarkIcon } from '@heroicons/react/24/solid';
import { useEffect, useRef, useState } from "react";
import ReactPlayer from 'react-player';
import { useWs, WsAction } from "../../hooks/useWs";
import { useStore } from "../../store";
import { PlayerState, Queue } from "../../types/Media";
import { MediaQueue } from './MediaQueue';

export function MediaPlayer({ channel }: { channel: string }): JSX.Element {
    const player = useRef<ReactPlayer | null>(null);

    const tokenContent = useStore(state => state.scTokenContent);
    const isChannelOwner = useRef(tokenContent?.Login === channel || channel === "");
    const [queue, setQueue] = useState<Queue>([]);

    useEffect(() => {
        isChannelOwner.current = tokenContent?.Login === channel || channel === "";
    }, [channel, tokenContent?.Login]);

    const [url, setUrl] = useState();
    const [playing, setPlaying] = useState(false);
    const [volume, setVolume] = useState(0);

    const { sendJsonMessage } = useWs(handleWsMessage);

    function handleWsMessage(event: MessageEvent<any>): void {
        const data = JSON.parse(event.data);

        if (data.action === WsAction.PLAYER_STATE) {
            if (player.current) {
                player.current.seekTo(data.time, "seconds");

                if (data.state === PlayerState.PLAYING) {
                    setPlaying(true);
                }
                if (data.state === PlayerState.PAUSED) {
                    setPlaying(false);
                }
                setUrl(data.url);
            }
        }
        if (data.action === WsAction.QUEUE_STATE) {
            const queue: Queue = data.queue;

            setQueue(queue);
        }
        if (data.action === WsAction.DEBUG) {
            console.log(data);
        }
    }

    useEffect(() => {
        sendJsonMessage({ action: WsAction.JOIN, channel: channel });
        sendJsonMessage({ action: WsAction.GET_QUEUE, channel: channel });
    }, []);

    const handlePause = () => {
        if (!isChannelOwner.current) {
            return;
        }

        const time = player.current?.getCurrentTime()
        sendJsonMessage({ action: WsAction.PLAYER_STATE, time: time, url: url, state: PlayerState.PAUSED });
    }

    const handlePlay = () => {
        if (!isChannelOwner.current) {
            return;
        }

        const time = player.current?.getCurrentTime()
        sendJsonMessage({ action: WsAction.PLAYER_STATE, time: time, url: url, state: PlayerState.PLAYING });
    }

    const handleSeek = (seconds: number) => {
        if (!isChannelOwner.current) {
            return;
        }

        sendJsonMessage({ action: WsAction.PLAYER_STATE, time: seconds, url: url, state: PlayerState.PLAYING });
    }

    return <div className="flex gap-4 w-full h-full">
        <div>
            <div className="mb-4">
                <div className="p-4 bg-gray-800 rounded shadow flex items-center">
                    <div className="flex gap-2 items-center">
                        <div onClick={() => setVolume(vol => vol > 0 ? 0 : 1)} className="cursor-pointer hover:text-blue-500">
                            {volume > 0 ? <SpeakerWaveIcon className="h-6 w-6" /> : <SpeakerXMarkIcon className="h-6 w-6" />}
                        </div>
                        <input
                            type="range"
                            className="cursor-pointer"
                            min={0}
                            max={1}
                            step={0.01}
                            value={volume}
                            onChange={event => {
                                setVolume(event.target.valueAsNumber)
                            }}
                        />
                        {Math.round(volume * 100)}%
                    </div>
                </div>
            </div>
            <ReactPlayer controls={true} ref={player} muted={volume === 0} volume={volume} pip={true} url={url} playing={playing} onPause={handlePause} onSeek={handleSeek} onPlay={handlePlay} />
        </div>
        <MediaQueue queue={queue} />
    </div>
}