
export interface EventMessage {
    records: Array<Record>
    joinedChannels: number
    activeChannels: number
}

export interface Record {
    title: string
    scores: Array<Score>
}

export interface Score {
    id: string
    score: number
}