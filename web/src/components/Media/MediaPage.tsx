import { useWs, WsAction } from "../../hooks/useWs";
import YouTube, { YouTubeProps } from "react-youtube";

export function MediaPage(): JSX.Element {
    const { lastJsonMessage, sendJsonMessage } = useWs();

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

    return <div className="p-4 w-full max-h-screen flex gap-4">
        <div className="p-4 bg-gray-800 rounded shadow relative">
            {JSON.stringify(lastJsonMessage)}
            <YouTube className="my-4" videoId="TjAa0wOe5k4" opts={opts} onReady={onPlayerReady} onPlay={onPlay} onStateChange={onStateChange} />
        </div>
    </div>;
}