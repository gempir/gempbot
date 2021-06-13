export enum RewardTypes {
    Bttv = "bttv",
    Timeout = "timeout"
}

export interface ChannelPointReward {
    OwnerTwitchID: string
    Type: RewardTypes
    Title: string
    Cost: number
    Prompt: string
    BackgroundColor: string
    IsMaxPerStreamEnabled: boolean
    MaxPerStream: number
    IsUserInputRequired: boolean
    IsMaxPerUserPerStreamEnabled: boolean
    MaxPerUserPerStream: number
    IsGlobalCooldownEnabled: boolean
    GlobalCooldownSeconds: number
    ShouldRedemptionsSkipRequestQueue: boolean
    Enabled: boolean
    RewardID?: string
}