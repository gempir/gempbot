

export function Meta({ activeChannels, joinedChannels }: { [key: string]: number }) {
    return <div className="flex flex-row">
        <div className="rounded shadow bg-gray-800 p-4 w-64 m-4">
            Joined Channels: <strong>{joinedChannels}</strong>
        </div>
        <div className="rounded shadow bg-gray-800 p-4 w-64 m-4 -ml-0">
            Active Channels: <strong>{activeChannels}</strong>
        </div>
    </div>
}