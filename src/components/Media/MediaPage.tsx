import { useWs } from "../../hooks/useWs";
import YouTube, { YouTubeProps } from "react-youtube";

export function MediaPage(): JSX.Element {
    const {lastJsonMessage} = useWs();

    const onPlayerReady: YouTubeProps['onReady'] = (event) => {
        // access to player in all event handlers via event.target
        event.target.pauseVideo();
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
            <YouTube className="my-4" videoId="TjAa0wOe5k4" opts={opts} onReady={onPlayerReady} />
        </div>
    </div>;
}