import { useEffect, useRef, useState } from "react";
import ReactPlayer from 'react-player';
import { useWs, WsAction } from "../../hooks/useWs";
import { useStore } from "../../store";
import { PlayerState } from "../../types/Media";

export function MediaPlayer({ channel }: { channel: string }): JSX.Element {
    const player = useRef<ReactPlayer | null>(null);

    const tokenContent = useStore(state => state.scTokenContent);
    const isChannelOwner = useRef(tokenContent?.Login === channel || channel === "");

    useEffect(() => {
        isChannelOwner.current = tokenContent?.Login === channel || channel === "";
    }, [channel, tokenContent?.Login]);

    const [url, setUrl] = useState("https://www.youtube.com/watch?v=wzE2nsjsHhg");
    const [playing, setPlaying] = useState(false);

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
        if (data.action === WsAction.DEBUG) {
            console.log(data);
        }
    }

    useEffect(() => {
        sendJsonMessage({ action: WsAction.JOIN, channel: channel });
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

    return <div>
        <ReactPlayer controls={true} ref={player} volume={0} pip={true} url={url} playing={playing} onPause={handlePause} onSeek={handleSeek} onPlay={handlePlay} />
    </div>
}