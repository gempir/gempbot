import { Dayjs } from "dayjs"

export interface Election {
    ChannelTwitchID: string
    Hours: number
    NominationCost: number
    EmoteAmount: number
    MaxNominationPerUser: number
    CreatedAt: Dayjs
    UpdatedAt: Dayjs
    StartedRunAt?: Dayjs
    SpecificTime?: Dayjs
}