import { createLoginUrl } from "../factory/createLoginUrl";
import { store } from "../store";


export function Navbar() {
    const scToken = store.useState(s => s.scToken);

    let buttons = <Login className="p-4 rounded shadow bg-purple-800 hover:bg-purple-600" />

    if (scToken) {
        buttons = <LoggedIn />;
    }

    return <div className="m-4 flex flex-row justify-end gap-4">
        {buttons}
    </div>;
}

function LoggedIn() {
    return <>
        <Login className={`p-4 rounded shadow opacity-25 bg-purple-800 hover:bg-purple-600`} />
    </>;
}

function Login({ className }: { className?: string }) {
    const { apiBaseUrl, twitchClientId } = store.getRawState();
    const url = createLoginUrl(apiBaseUrl, twitchClientId);

    return <a className={className} href={url.toString()}>Login</a>
}