
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
    user: User
    score: number
}

export interface User {
    id: string
    displayName: string
    profilePicture: string
}