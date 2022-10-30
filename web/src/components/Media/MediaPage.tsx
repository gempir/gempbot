import { useEffect, useRef } from "react";
import YouTube, { YouTubeProps } from "react-youtube";
import { useWs, WsAction } from "../../hooks/useWs";
import { useStore } from "../../store";
import { PlayerState } from "../../types/Media";
import { MediaQueue } from "./MediaQueue";

export function MediaPage({ channel = "" }: { channel?: string }): JSX.Element {
    const tokenContent = useStore(state => state.scTokenContent);
    const isChannelOwner = useRef(tokenContent?.Login === channel || channel === "");

    useEffect(() => {
        isChannelOwner.current = tokenContent?.Login === channel || channel === "";
    }, [channel, tokenContent?.Login]);

    const player = useRef<YouTube | null>(null);

    const onPlayerReady: YouTubeProps['onReady'] = (event) => {
        if (!isChannelOwner.current || [PlayerState.BUFFERING, PlayerState.CUED, PlayerState.ENDED, PlayerState.UNSTARTED].includes(event.target.getPlayerState())) {
            return;
        }

        sendJsonMessage({ action: WsAction.PLAYER_STATE, seconds: event.target.getCurrentTime(), videoId: event.target.getVideoData()['video_id'], state: event.target.getPlayerState() });
    }

    const onPlay: YouTubeProps['onPlay'] = (event) => {
        if (!isChannelOwner.current || [PlayerState.BUFFERING, PlayerState.CUED, PlayerState.ENDED, PlayerState.UNSTARTED].includes(event.target.getPlayerState())) {
            return;
        }

        sendJsonMessage({ action: WsAction.PLAYER_STATE, seconds: event.target.getCurrentTime(), videoId: event.target.getVideoData()['video_id'], state: event.target.getPlayerState() });
    }

    const onStateChange: YouTubeProps['onStateChange'] = (event) => {
        if (!isChannelOwner.current || [PlayerState.BUFFERING, PlayerState.CUED, PlayerState.ENDED, PlayerState.UNSTARTED].includes(event.target.getPlayerState())) {
            return;
        }

        sendJsonMessage({ action: WsAction.PLAYER_STATE, seconds: event.target.getCurrentTime(), videoId: event.target.getVideoData()['video_id'], state: event.target.getPlayerState() });
    }

    const opts: YouTubeProps['opts'] = {
        height: '450',
        width: '800',
        playerVars: {
            // https://developers.google.com/youtube/player_parameters
            modestbranding: 1,
            autoplay: 1,
            mute: 1, // mute for now because this allows autoplay
        },
    };

    const handleWsMessage = (event: MessageEvent<any>): void => {
        const data = JSON.parse(event.data);

        if (data.action === WsAction.PLAYER_STATE) {
            if (player.current) {
                player.current.getInternalPlayer().seekTo(data.currentTime);

                if (data.state === PlayerState.PLAYING) {
                    player.current.getInternalPlayer().playVideo();
                }
                if (data.state === PlayerState.PAUSED) {
                    player.current.getInternalPlayer().pauseVideo();
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

    return <div className="p-4 w-full flex gap-4">
        <YouTube ref={player} className="" videoId="TjAa0wOe5k4" opts={opts} onReady={onPlayerReady} onPlay={onPlay} onStateChange={onStateChange} />
        <MediaQueue />
    </div>;
}