import { useEffect, useRef, useState } from "react";
import ReactPlayer from 'react-player';
import { useWs, WsAction } from "../../hooks/useWs";
import { useStore } from "../../store";
import { PlayerState } from "../../types/Media";

export function MediaPlayer({ channel }: { channel: string }): JSX.Element {
    const [url, setUrl] = useState("https://www.youtube.com/watch?v=TjAa0wOe5k4");
    const [playing, setPlaying] = useState(false);
    const player = useRef<ReactPlayer | null>(null);

    const tokenContent = useStore(state => state.scTokenContent);
    const isChannelOwner = useRef(tokenContent?.Login === channel || channel === "");

    useEffect(() => {
        isChannelOwner.current = tokenContent?.Login === channel || channel === "";
    }, [channel, tokenContent?.Login]);

    const handleWsMessage = (event: MessageEvent<any>): void => {
        const data = JSON.parse(event.data);

        if (data.action === WsAction.PLAYER_STATE) {
            if (player.current) {
                player.current.seekTo(data.currentTime, "seconds");

                if (data.state === PlayerState.PLAYING) {
                    setPlaying(true);
                }
                if (data.state === PlayerState.PAUSED) {
                    setPlaying(false);
                }
            }
        }
        if (data.action === WsAction.DEBUG) {
            console.log(data);
        }
    };

    const { sendJsonMessage } = useWs(handleWsMessage);

    useEffect(() => {
        sendJsonMessage({ action: WsAction.JOIN, channel: channel });
    }, []);


    const handlePause = () => {
        if (!isChannelOwner.current) {
            return;
        }

        const time = player.current?.getCurrentTime()

        sendJsonMessage({ action: WsAction.PLAYER_STATE, seconds: time, videoId: "", state: PlayerState.PAUSED });
    }

    const handleSeek = (seconds: number) => {
        if (!isChannelOwner.current) {
            return;
        }

        sendJsonMessage({ action: WsAction.PLAYER_STATE, seconds: seconds, videoId: "", state: PlayerState.PLAYING });
    }

    return <div>
        <ReactPlayer ref={player} url={url} playing={playing} onPause={handlePause} onSeek={handleSeek} controls={isChannelOwner.current} />
    </div>
}