import { Dayjs } from "dayjs"

export interface Election {
    ChannelTwitchID?: string
    Hours: number
    NominationCost: number
    CreatedAt: Dayjs
    UpdatedAt: Dayjs
    StartedRunAt: Dayjs
}