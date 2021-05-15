import { useForm } from "react-hook-form";

export function BttvForm() {
    const { register, handleSubmit, formState: { errors } } = useForm();
    // @ts-ignore
    const onSubmit = data => console.log(data);

    return (
        <form onSubmit={handleSubmit(onSubmit)} className="m-4 p-4 bg-gray-800 rounded shadow">
            <label className="block">
                Title
                <input defaultValue="Bttv Emote" spellCheck={false} {...register("title", { required: true })} className="form-input border-none bg-gray-700 mx-2 p-2 rounded shadow" />
                {errors.title && <span className="text-red-700">required</span>}
            </label>

            <label className="block mt-3">
                Cost 
                <input defaultValue="Bttv Emote" type="number" spellCheck={false} {...register("cost", { required: true })} className="form-input border-none bg-gray-700 mx-2 p-2 rounded shadow" />
                {errors.cost && <span className="text-red-700">required</span>}
            </label>

            <label className="block mt-3">
                Prompt
                <input spellCheck={false} disabled {...register("prompt")} className="block cursor-not-allowed truncate form-input w-full opacity-25 border-none bg-gray-700 mt-2 p-2 rounded shadow" defaultValue="Add a BetterTTV emote! In the text field, send a link to the BetterTTV emote. powered by spamchamp.gempir.com" />
            </label>

            <div className="flex flex-row justify-between items-center select-none">
                <label className="inline-flex items-center">
                    <input type="checkbox" className="form-checkbox rounded border-gray-300 text-indigo-600 shadow-sm focus:border-indigo-300 focus:ring focus:ring-offset-0 focus:ring-indigo-200 focus:ring-opacity-50" />
                    <span className="ml-2">Enabled</span>
                </label>
                <input type="submit" className="bg-green-700 hover:bg-green-600 p-2 rounded shadow block mt-3" value="save" />
            </div>
        </form>
    );
}