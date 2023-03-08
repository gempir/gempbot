export enum PlayerState {
    PLAYING = "PLAYING",
    PAUSED = "PAUSED",
}

export type Queue = QueueItem[];

export interface QueueItem {
    ID:              string;
    ChannelTwitchId: string;
    Url:             string;
    Approved:        boolean;
    Author:          string;
    Approver:        string;
    CreatedAt:       Date;
    UpdatedAt:       Date;
}