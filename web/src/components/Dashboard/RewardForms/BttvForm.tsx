import { useForm } from "react-hook-form";
import { SetUserConfig, UserConfig } from "../../../hooks/useUserConfig";
import { doFetch, Method } from "../../../service/doFetch";

interface BttvRewardForm {
    title: string;
    cost: string;
    backgroundColor: string;
    maxPerStream: string;
    maxPerUserPerStream: string;
    globalCooldownMinutes: string;
    enabled: boolean;
    isDefault: boolean
}

export function BttvForm({ userConfig, setUserConfig, fetchConfig }: { userConfig: UserConfig, setUserConfig: SetUserConfig, fetchConfig: () => void }) {
    const { register, handleSubmit, formState: { errors } } = useForm();
    const onSubmit = (data: BttvRewardForm) => setUserConfig(
        {
            ...userConfig,
            Rewards: {
                ...userConfig?.Rewards,
                Bttv: {
                    ...data,
                    isDefault: false,
                    cost: Number(data.cost),
                    maxPerStream: Number(data.maxPerStream),
                    maxPerUserPerStream: Number(data.maxPerUserPerStream),
                    globalCooldownSeconds: Number(data.globalCooldownMinutes) * 60
                }
            }
        }
    );

    return (
        <form onSubmit={handleSubmit(onSubmit)} className="max-w-5xl m-4 p-4 bg-gray-800 rounded shadow">
            <div className="flex items-center justify-between">
                <div className="flex items-center gap-5">
                    <img src="/images/bttv.png" alt="BetterTTV Logo" className="w-16" />
                    <h3 className="text-xl font-bold">BetterTTV Emote</h3>
                </div>
                <div className="text-gray-600">
                    {userConfig.Rewards.Bttv?.ID &&
                        <div className="bg-red-700 hover:bg-red-600 p-2 rounded shadow mt-3 text-gray-100 inline-block ml-3 cursor-pointer"
                            onClick={() => doFetch(Method.DELETE, `/api/reward/${userConfig.CurrentUserID}/${userConfig.Rewards.Bttv?.ID}`).then(fetchConfig)}>
                            Delete
                        </div>
                    }
                </div>
            </div>
            <p className="my-2 mb-4 text-gray-400">
                <strong>Make sure <span className="text-green-600">gempbot</span> is BetterTTV editor</strong><br />
                This will swap out 1 emote constantly. If the previous emote is not found it will use a free slot or remove a random emote.
            </p>

            <label className="block">
                Title
                <input defaultValue={userConfig.Rewards.Bttv?.title} spellCheck={false} {...register("title", { required: true })} className="form-input border-none bg-gray-700 mx-2 p-2 rounded shadow" />
                {errors.title && <span className="text-red-700">required</span>}
            </label>

            <label className="block mt-3">
                Cost
                <input defaultValue={userConfig.Rewards.Bttv?.cost} type="number" spellCheck={false} {...register("cost", { required: true })} className="form-input border-none bg-gray-700 mx-2 p-2 rounded shadow" />
                {errors.cost && <span className="text-red-700">required</span>}
            </label>

            <label className="block mt-3">
                Prompt
                <input defaultValue={userConfig.Rewards.Bttv?.prompt} spellCheck={false} disabled {...register("prompt")} className="block cursor-not-allowed truncate form-input w-full opacity-25 border-none bg-gray-700 mt-2 p-2 rounded shadow" />
            </label>

            <label className="block mt-3">
                Background Color
                <input defaultValue={userConfig.Rewards.Bttv?.backgroundColor} placeholder="#FFFFFF" spellCheck={false} {...register("backgroundColor")} className="form-input border-none bg-gray-700 mx-2 p-2 rounded shadow" />
            </label>

            <div className="mt-8 font-bold">Limits <span className="text-gray-500">(0 = unlimited)</span></div>

            <label className="flex items-center mt-3">
                Max per Stream
                <input defaultValue={userConfig.Rewards.Bttv?.maxPerStream} placeholder="0" type="number" spellCheck={false} {...register("maxPerStream")} className="form-input border-none bg-gray-700 mx-2 p-2 rounded shadow" />
            </label>

            <label className="flex items-center mt-3">
                Max per User per Stream
                <input defaultValue={userConfig.Rewards.Bttv?.maxPerUserPerStream} placeholder="0" type="number" spellCheck={false} {...register("maxPerUserPerStream")} className="form-input border-none bg-gray-700 mx-2 p-2 rounded shadow" />
            </label>

            <label className="flex items-center mt-3">
                Global Cooldown in Minutes
                <input defaultValue={(userConfig.Rewards.Bttv?.globalCooldownSeconds ?? 0) / 60} placeholder="0" type="number" spellCheck={false} {...register("globalCooldownMinutes")} className="form-input border-none bg-gray-700 mx-2 p-2 rounded shadow" />
            </label>

            <div className="flex flex-row justify-between items-center select-none">
                <label className="flex items-center">
                    <input defaultChecked={userConfig.Rewards.Bttv?.enabled} type="checkbox" {...register("enabled")} className="form-checkbox rounded border-gray-300 text-indigo-600 shadow-sm focus:border-indigo-300 focus:ring focus:ring-offset-0 focus:ring-indigo-200 focus:ring-opacity-50" />
                    <span className="ml-2">Enabled</span>
                </label>
                <input type="submit" className="bg-green-700 hover:bg-green-600 p-2 rounded shadow block mt-3 cursor-pointer" value="save" />
            </div>
        </form>
    );
}