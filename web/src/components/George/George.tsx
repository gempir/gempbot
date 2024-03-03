import { useState } from "react";
import { useGeorge } from "../../hooks/useGeorge";

export function George() {
    const [resp, setResp] = useState<string>("");

    const doReq = useGeorge();

    const onSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        setResp(".");
        const load = setInterval(() => {
            setResp(prev => prev + ".");
        }, 1000);

        e.preventDefault();
        const form = e.currentTarget;
        const formData = new FormData(form);
        const req = {
            channel: formData.get("channel") as string,
            username: formData.get("username") as string,
            year: parseInt(formData.get("year") as string),
            month: parseInt(formData.get("month") as string),
            query: formData.get("query") as string
        };
        doReq(req, (text: string) => {
            clearInterval(load);
            setResp(text.trim());
        });
    }


    return <div className={"p-4 w-full"}>
        <div className={"bg-gray-800 rounded shadow relative p-4 w-full"}>
            <div className="flex items-start justify-between w-full">
                <form className="w-full" onSubmit={onSubmit}>
                    <div className="flex gap-2 w-full">
                        <div>
                            <div className="flex gap-2">
                                <input type="text" placeholder="Channel" name="channel" className="w-full bg-gray-700 p-2 rounded mt-2" />
                                <input type="text" placeholder="Username" name="username" className="w-full bg-gray-700 p-2 rounded mt-2" />
                            </div>

                            <div className="flex gap-2">
                                <input type="number" placeholder="Year" name="year" className="w-full bg-gray-700 p-2 rounded mt-2" />
                                <input type="number" placeholder="Month" name="month" className="w-full bg-gray-700 p-2 rounded mt-2" />
                            </div>
                        </div>
                        <input type="text" placeholder="Query" name="query" className="w-full bg-gray-700 p-2 rounded mt-2 resize-none" />
                        <input type="submit" value="Send" className="bg-blue-500 py-2 px-5 rounded mt-2" />
                    </div>
                </form>
            </div>
        </div>
        <div className={"bg-gray-800 rounded shadow relative p-4 w-full mt-2"}>
            <div className="flex items-start justify-between w-full">
                <textarea readOnly value={resp} placeholder="Response" name="response" className="w-full min-h-[700px] bg-gray-900 p-2 border-none select-none rounded focus:outline-none focus:ring-0 resize-none" />
            </div>
        </div>
    </div >;
}


