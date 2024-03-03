import { useRef, useState } from "react";
import { useGeorge } from "../../hooks/useGeorge";

export function George() {
    const [resp, setResp] = useState<string>("");
    const loadRef = useRef<NodeJS.Timeout>();

    const doReq = useGeorge();

    const onSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        clearTimeout(loadRef.current);

        e.preventDefault();
        const form = e.currentTarget;
        const formData = new FormData(form);

        // when username is defined username is required
        if (!formData.get("username") && !formData.get("day")) {
            setResp("channel log requires day");
            return;
        }


        const req = {
            channel: formData.get("channel") as string,
            username: formData.get("username") as string,
            year: parseInt(formData.get("year") as string),
            month: parseInt(formData.get("month") as string),
            day: parseInt(formData.get("day") as string),
            model: formData.get("model") as string,
            query: formData.get("query") as string,
            limit: parseInt(formData.get("limit") as string),
        };
        setResp(".");
        loadRef.current = setInterval(() => {
            setResp(prev => prev + ".");
        }, 1000);
        doReq(req, (text: string) => {
            clearInterval(loadRef.current);
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
                                <input type="text" placeholder="Channel" name="channel" className="w-full bg-gray-800 p-2 rounded" />
                                <input type="text" placeholder="Username" name="username" className="w-full bg-gray-800 p-2 rounded" />
                                <select name="model" className="bg-gray-800 rounded appearance-none">
                                    <option value="mistral">mistral</option>
                                    <option value="llama">llama2</option>
                                    <option value="llama2:70b">llama2:70b</option>
                                    <option value="gemma:7b">gemma:7b</option>
                                    <option value="llama2:13b">llama2:13b</option>
                                    <option value="llama2-uncensored">llama2-uncensored (donk)</option>
                                </select>
                            </div>

                            <div className="flex gap-2 justify-center align-middle">
                                <input type="number" placeholder="Year" name="year" className="w-full bg-gray-800 p-2 rounded mt-2" />
                                <input type="number" placeholder="Month" name="month" className="w-full bg-gray-800 p-2 rounded mt-2" />
                                <input type="number" placeholder="Day" name="day" className="w-full bg-gray-800 p-2 rounded mt-2" />
                                <div className="flex justify-center align-middle pt-2">
                                    Max Tokens
                                </div>
                                <input type="number" placeholder="Max Tokens" name="limit" defaultValue={"4096"} className="w-full bg-gray-800 p-2 rounded mt-2" />
                            </div>
                        </div>
                        <input type="text" placeholder="Query" name="query" className="w-full bg-gray-800 p-2 rounded resize-none" />
                        <input type="submit" value="Send" className="bg-blue-500 py-2 px-5 rounded" />
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

