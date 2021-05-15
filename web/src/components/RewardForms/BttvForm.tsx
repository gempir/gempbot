import { useForm } from "react-hook-form";
import { SetUserConfig, UserConfig } from "../../hooks/useUserConfig";

interface BttvRewardForm {
    title: string;
    cost: string;
    backgroundColor: string;
    maxPerStream: string;
    maxPerUserPerStream: string;
    globalCooldownMinutes: string;
}

export function BttvForm({ userConfig, setUserConfig }: { userConfig: UserConfig, setUserConfig: SetUserConfig }) {
    const { register, handleSubmit, formState: { errors } } = useForm();
    const onSubmit = (data: BttvRewardForm) => setUserConfig(
        {
            ...userConfig,
            Rewards: {
                ...userConfig?.Rewards,
                Bttv: {
                    ...data,
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

            <label className="flex items-center mt-3">
                Max per Stream
                <input defaultValue={userConfig.Rewards.Bttv?.maxPerStream} placeholder="∞" type="number" spellCheck={false} {...register("maxPerStream")} className="form-input border-none bg-gray-700 mx-2 p-2 rounded shadow" />
            </label>

            <label className="flex items-center mt-3">
                Max per User per Stream
                <input defaultValue={userConfig.Rewards.Bttv?.maxPerUserPerStream} placeholder="∞" type="number" spellCheck={false} {...register("maxPerUserPerStream")} className="form-input border-none bg-gray-700 mx-2 p-2 rounded shadow" />
            </label>

            <label className="flex items-center mt-3">
                Global Cooldown in Minutes
                <input defaultValue={userConfig.Rewards.Bttv?.globalCooldownSeconds ?  userConfig.Rewards.Bttv?.globalCooldownSeconds / 60 : undefined} placeholder="0" type="number" spellCheck={false} {...register("globalCooldownMinutes")} className="form-input border-none bg-gray-700 mx-2 p-2 rounded shadow" />
            </label>

            <div className="flex flex-row justify-between items-center select-none">
                <label className="flex items-center">
                    <input defaultChecked={userConfig.Rewards.Bttv?.Enabled} type="checkbox" className="form-checkbox rounded border-gray-300 text-indigo-600 shadow-sm focus:border-indigo-300 focus:ring focus:ring-offset-0 focus:ring-indigo-200 focus:ring-opacity-50" />
                    <span className="ml-2">Enabled</span>
                </label>
                <input type="submit" className="bg-green-700 hover:bg-green-600 p-2 rounded shadow block mt-3" value="save" />
            </div>
        </form>
    );
}