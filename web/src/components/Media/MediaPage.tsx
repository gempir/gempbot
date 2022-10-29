import { useEffect, useRef } from "react";
import YouTube, { YouTubeProps } from "react-youtube";
import { useWs, WsAction } from "../../hooks/useWs";

export function MediaPage({ channel }: { channel: string }): JSX.Element {
    const player = useRef<YouTube | null>(null);

    const onPlayerReady: YouTubeProps['onReady'] = (event) => {
        sendJsonMessage({ action: WsAction.TIME_CHANGED, seconds: event.target.getCurrentTime(), videoId: event.target.getVideoData()['video_id'] });
    }

    const onPlay: YouTubeProps['onPlay'] = (event) => {
        sendJsonMessage({ action: WsAction.TIME_CHANGED, seconds: event.target.getCurrentTime(), videoId: event.target.getVideoData()['video_id'] });
    }

    const onStateChange: YouTubeProps['onStateChange'] = (event) => {
        sendJsonMessage({ action: WsAction.TIME_CHANGED, seconds: event.target.getCurrentTime(), videoId: event.target.getVideoData()['video_id'] });
    }

    const opts: YouTubeProps['opts'] = {
        height: '720',
        width: '1280',
        playerVars: {
            // https://developers.google.com/youtube/player_parameters
            autoplay: 0,
        },
    };

    const handleWsMessage = (event: MessageEvent<any>): void => {
        console.log("ev", event);
        const data = JSON.parse(event.data);

        if (data.action === WsAction.TIME_CHANGED) {
            if (player.current) {
                console.log(player.current.getInternalPlayer());
                player.current.getInternalPlayer().seekTo(data.currentTime);
            }
        }
    };

    const { lastJsonMessage, sendJsonMessage, getWebSocket } = useWs(handleWsMessage);

    useEffect(() => {
        sendJsonMessage({ action: WsAction.JOIN, channel: channel });
    }, []);

    return <div className="p-4 w-full max-h-screen flex gap-4">
        <div className="p-4 bg-gray-800 rounded shadow relative">
            {JSON.stringify(lastJsonMessage)}
            <YouTube ref={player} className="my-4" videoId="TjAa0wOe5k4" opts={opts} onReady={onPlayerReady} onPlay={onPlay} onStateChange={onStateChange} />
        </div>
    </div>;
}