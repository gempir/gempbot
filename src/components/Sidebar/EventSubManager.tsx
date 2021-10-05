import { useSubscribtions } from "../../hooks/useSubscriptions";

export function EventSubManager() {
    const [subscribe, remove, loadingSubscribe, loadingRemove] = useSubscribtions();


    return <div className="flex mt-2 gap-2">
        <div onClick={remove} className={"py-4 w-full text-center bg-red-900 rounded opacity-10 hover:opacity-100 cursor-pointer" + (loadingRemove ? " animate-pulse pointer-events-none" : "")}>
            remove
        </div>
        <div onClick={subscribe} className={"py-4 w-full text-center bg-green-900 rounded opacity-10 hover:opacity-100 cursor-pointer" + (loadingSubscribe ? " animate-pulse pointer-events-none" : "")}>
            sub
        </div>
    </div>
}

