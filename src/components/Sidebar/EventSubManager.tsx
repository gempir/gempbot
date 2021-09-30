import { useSubscribtions } from "../../hooks/useSubscriptions";

export function EventSubManager() {
    const [subscribe, remove] = useSubscribtions();


    return <div className="flex mt-2 gap-2">
        <div onClick={remove} className="py-4 w-full text-center bg-red-900 rounded opacity-10 hover:opacity-100 cursor-pointer">
            remove
        </div>
        <div onClick={subscribe} className="py-4 w-full text-center bg-green-900 rounded opacity-10 hover:opacity-100 cursor-pointer">
            sub
        </div>
    </div>
}

