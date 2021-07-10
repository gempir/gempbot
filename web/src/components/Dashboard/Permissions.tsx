import { UserConfig } from "../../hooks/useUserConfig";

export function Permissions({ userConfig }: { userConfig: UserConfig }) {
    return <div className="m-4 p-4 bg-gray-800 rounded shadow w-96 relative">
        <h2 className="mb-4 text-xl">Permissions</h2>
        <table className="w-full">
            <thead>
                <tr className="border-b-8 border-transparent">
                    <th>User</th>
                    <th>Prediction</th>
                </tr>
            </thead>
            <tbody>
                {userConfig.Permissions.map((perm, index) => <tr className={index % 2 ? "bg-gray-900" : ""}>
                    <th>{perm.User}</th>
                    <th><input type="checkbox" disabled defaultChecked={perm.Prediction} /></th>
                </tr>)}
            </tbody>
        </table>
    </div>
}