import { useForm } from "react-hook-form";
import { useChannelPointReward } from "../../../hooks/useChannelPointReward";
import { SetUserConfig, UserConfig } from "../../../hooks/useUserConfig";
import { doFetch, Method } from "../../../service/doFetch";
import { ChannelPointReward, RewardTypes } from "../../../types/Rewards";

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

const defaultReward = {
    OwnerTwitchID: "",
    Type: RewardTypes.Bttv,
    Title: "BetterTTV Emote",
    Cost: 10000,
    Prompt: "",
    BackgroundColor: "",
    IsMaxPerStreamEnabled: false,
    MaxPerStream: 0,
    IsUserInputRequired: true,
    MaxPerUserPerStream: 0,
    IsMaxPerUserPerStreamEnabled: false,
    IsGlobalCooldownEnabled: false,
    GlobalCooldownSeconds: 0,
    ShouldRedemptionsSkipRequestQueue: false,
    Enabled: false
}

export function BttvForm({ userConfig }: { userConfig: UserConfig }) {
    const { register, handleSubmit, formState: { errors } } = useForm();
    const [reward, setReward] = useChannelPointReward(userConfig?.Protected.CurrentUserID, RewardTypes.Bttv, defaultReward);
    const onSubmit = (data: BttvRewardForm) => () => {
        const rewardData: ChannelPointReward = {
            OwnerTwitchID: userConfig?.Protected.CurrentUserID,
            Type: RewardTypes.Bttv,
            Title: data.title,
            Cost: Number(data.cost),
            BackgroundColor: data.backgroundColor,
            IsMaxPerStreamEnabled: Boolean(data.maxPerStream),
            MaxPerStream: Number(data.maxPerStream),
            IsUserInputRequired: true,
            MaxPerUserPerStream: Number(data.maxPerUserPerStream),
            IsMaxPerUserPerStreamEnabled: Boolean(data.maxPerUserPerStream),
            IsGlobalCooldownEnabled: Boolean(data.globalCooldownMinutes),
            GlobalCooldownSeconds: Number(data.globalCooldownMinutes) * 60,
            ShouldRedemptionsSkipRequestQueue: false,
            Enabled: data.enabled
        };

        setReward(rewardData);
    }


    return (
        <form onSubmit={handleSubmit(onSubmit)} className="max-w-5xl m-4 p-4 bg-gray-800 rounded shadow">
            <div className="flex items-center justify-between">
                <div className="flex items-center gap-5">
                    <img src="/images/bttv.png" alt="BetterTTV Logo" className="w-16" />
                    <h3 className="text-xl font-bold">BetterTTV Emote</h3>
                </div>
                <div className="text-gray-600">
                    {reward?.RewardID &&
                        <>{reward.RewardID}</>
                        // <div className="bg-red-700 hover:bg-red-600 p-2 rounded shadow mt-3 text-gray-100 inline-block ml-3 cursor-pointer"
                        //     onClick={() => doFetch(Method.DELETE, `/api/reward/${userConfig.Protected.CurrentUserID}/${reward.RewardID}`).then(fetchConfig)}>
                        //     Delete
                        // </div>
                    }
                </div>
            </div>
            <p className="my-2 mb-4 text-gray-400">
                <strong>Make sure <span className="text-green-600">gempbot</span> is BetterTTV editor</strong><br />
                This will swap out 1 emote constantly. If the previous emote is not found it will use a free slot or remove a random emote.
            </p>

            <label className="block">
                Title
                <input defaultValue={reward.Title} spellCheck={false} {...register("title", { required: true })} className="form-input border-none bg-gray-700 mx-2 p-2 rounded shadow" />
                {errors.title && <span className="text-red-700">required</span>}
            </label>

            <label className="block mt-3">
                Cost
                <input defaultValue={reward.Cost} type="number" spellCheck={false} {...register("cost", { required: true })} className="form-input border-none bg-gray-700 mx-2 p-2 rounded shadow" />
                {errors.cost && <span className="text-red-700">required</span>}
            </label>

            <label className="block mt-3">
                Prompt
                <input defaultValue={reward.Prompt} spellCheck={false} disabled {...register("prompt")} className="block cursor-not-allowed truncate form-input w-full opacity-25 border-none bg-gray-700 mt-2 p-2 rounded shadow" />
            </label>

            <label className="block mt-3">
                Background Color
                {/* <input defaultValue={reward.backgroundColor} placeholder="#FFFFFF" spellCheck={false} {...register("backgroundColor")} className="form-input border-none bg-gray-700 mx-2 p-2 rounded shadow" /> */}
            </label>

            <div className="mt-8 font-bold">Limits <span className="text-gray-500">(0 = unlimited)</span></div>

            <label className="flex items-center mt-3">
                Max per Stream
                {/* <input defaultValue={reward.maxPerStream} placeholder="0" type="number" spellCheck={false} {...register("maxPerStream")} className="form-input border-none bg-gray-700 mx-2 p-2 rounded shadow" /> */}
            </label>

            <label className="flex items-center mt-3">
                Max per User per Stream
                {/* <input defaultValue={reward.maxPerUserPerStream} placeholder="0" type="number" spellCheck={false} {...register("maxPerUserPerStream")} className="form-input border-none bg-gray-700 mx-2 p-2 rounded shadow" /> */}
            </label>

            <label className="flex items-center mt-3">
                Global Cooldown in Minutes
                {/* <input defaultValue={(reward.globalCooldownSeconds ?? 0) / 60} placeholder="0" type="number" spellCheck={false} {...register("globalCooldownMinutes")} className="form-input border-none bg-gray-700 mx-2 p-2 rounded shadow" /> */}
            </label>

            <div className="flex flex-row justify-between items-center select-none">
                <label className="flex items-center">
                    {/* <input defaultChecked={reward.enabled} type="checkbox" {...register("enabled")} className="form-checkbox rounded border-gray-300 text-indigo-600 shadow-sm focus:border-indigo-300 focus:ring focus:ring-offset-0 focus:ring-indigo-200 focus:ring-opacity-50" /> */}
                    <span className="ml-2">Enabled</span>
                </label>
                <input type="submit" className="bg-green-700 hover:bg-green-600 p-2 rounded shadow block mt-3" value="save" />
            </div>
        </form>
    );
}