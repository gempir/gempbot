export enum RewardTypes {
  Bttv = "bttv",
  SevenTv = "seventv",
}

export interface ChannelPointReward {
  OwnerTwitchID: string;
  ApproveOnly: boolean;
  Type: RewardTypes;
  Title: string;
  Cost: number;
  Prompt: string;
  BackgroundColor: string;
  IsMaxPerStreamEnabled: boolean;
  MaxPerStream: number;
  IsUserInputRequired: boolean;
  IsMaxPerUserPerStreamEnabled: boolean;
  MaxPerUserPerStream: number;
  IsGlobalCooldownEnabled: boolean;
  GlobalCooldownSeconds: number;
  ShouldRedemptionsSkipRequestQueue: boolean;
  Enabled: boolean;
  RewardID?: string;
  AdditionalOptionsParsed: BttvAdditionalOptions;
}

export interface RawBttvChannelPointReward extends ChannelPointReward {
  AdditionalOptions: string;
}

export interface BttvAdditionalOptions {
  Slots: number;
}
