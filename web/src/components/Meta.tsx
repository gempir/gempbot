

export function Meta({ activeChannels, joinedChannels }: { [key: string]: number }) {
    return <div className="flex flex-row">
        <div className="rounded shadow bg-gray-700 p-5 w-64 m-5">
            Joined Channels: <strong>{joinedChannels}</strong>
        </div>
        <div className="rounded shadow bg-gray-700 p-5 w-64 m-5 -ml-0">
            Active Channels: <strong>{activeChannels}</strong>
        </div>
    </div>
}