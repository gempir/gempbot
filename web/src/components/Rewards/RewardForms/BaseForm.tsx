import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { useChannelPointReward } from "../../../hooks/useChannelPointReward";
import { UserConfig } from "../../../hooks/useUserConfig";
import { ChannelPointReward, RewardTypes } from "../../../types/Rewards";

interface BttvRewardForm {
    title: string;
    cost: string;
    prompt: string;
    backgroundColor: string;
    maxPerStream: string;
    maxPerUserPerStream: string;
    globalCooldownMinutes: string;
    enabled: boolean;
    isDefault: boolean;
    slots: number;
}

const defaultReward = {
    OwnerTwitchID: "",
    Type: RewardTypes.Bttv,
    Title: "BetterTTV Emote",
    Cost: 10000,
    Prompt: "Add a BetterTTV emote! In the text field, send a link to the BetterTTV emote. powered by bitraft.gempir.com",
    BackgroundColor: "",
    IsMaxPerStreamEnabled: false,
    MaxPerStream: 0,
    IsUserInputRequired: true,
    MaxPerUserPerStream: 0,
    IsMaxPerUserPerStreamEnabled: false,
    IsGlobalCooldownEnabled: false,
    GlobalCooldownSeconds: 0,
    ShouldRedemptionsSkipRequestQueue: false,
    Enabled: false,
    AdditionalOptionsParsed: { Slots: 1 }
}

export function BaseForm({ userConfig, type, header, description, children }: { userConfig: UserConfig, type: RewardTypes, header?: JSX.Element, description?: JSX.Element, children?: JSX.Element }) {
    const { register, handleSubmit, formState: { errors }, setValue } = useForm();
    const [loading, setLoading] = useState(false);

    const onUpdate = () => {
        setLoading(false);
    }

    const [reward, setReward, deleteReward] = useChannelPointReward(userConfig?.Protected.CurrentUserID, type, defaultReward, onUpdate);
    const onSubmit = (data: BttvRewardForm) => {
        const rewardData: ChannelPointReward = {
            OwnerTwitchID: userConfig?.Protected.CurrentUserID,
            Type: type,
            Title: data.title,
            Prompt: data.prompt,
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
            Enabled: data.enabled,
            AdditionalOptionsParsed: {
                Slots: Number(data.slots)
            }
        };

        setReward(rewardData);
        setLoading(true);
    }

    useEffect(() => {
        setValue("title", reward.Title);
        setValue("prompt", reward.Prompt);
        setValue("cost", reward.Cost);
        setValue("slots", reward.AdditionalOptionsParsed.Slots);
        setValue("backgroundColor", reward.BackgroundColor);
        setValue("maxPerStream", reward.MaxPerStream);
        setValue("maxPerUserPerStream", reward.MaxPerUserPerStream);
        setValue("globalCooldownMinutes", reward.GlobalCooldownSeconds / 60);
        setValue("enabled", reward.Enabled);
    }, [reward, setValue]);

    return (
        <form onSubmit={handleSubmit(onSubmit)} className="max-w-5xl p-4 bg-gray-800 rounded shadow">
            <div className="flex items-center justify-between">
                <div className="flex items-center gap-5">
                    {header}
                </div>
                <div className="text-gray-600">
                    {reward?.RewardID &&
                        <div className="bg-red-700 hover:bg-red-600 p-2 rounded shadow mt-3 text-gray-100 inline-block ml-3 cursor-pointer" onClick={deleteReward}>
                            Delete
                        </div>
                    }
                </div>
            </div>
            {description}
            <label className="block my-3">
                Slots
                <input defaultValue={reward.AdditionalOptionsParsed.Slots} type="number" spellCheck={false} {...register("slots", { required: true })} className="form-input border-none bg-gray-700 mx-2 p-2 rounded shadow" />
                {errors.cost && <span className="text-red-700">required</span>}
            </label>

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
                <input defaultValue={reward.Prompt} spellCheck={false} {...register("prompt")} className="block truncate form-input w-full border-none bg-gray-700 mt-2 p-2 rounded shadow" />
            </label>

            <label className="block mt-3">
                Background Color
                <input defaultValue={reward.BackgroundColor} placeholder="#FFFFFF" pattern="^#+([a-fA-F0-9]{6}|[a-fA-F0-9]{3})$" spellCheck={false} {...register("backgroundColor")} className="form-input border-none bg-gray-700 mx-2 p-2 rounded shadow" />
            </label>

            <div className="mt-8 font-bold">Limits <span className="text-gray-500">(0 = unlimited)</span></div>

            <label className="flex items-center mt-3">
                Max per Stream
                <input defaultValue={reward.MaxPerStream} placeholder="0" type="number" spellCheck={false} {...register("maxPerStream")} className="form-input border-none bg-gray-700 mx-2 p-2 rounded shadow" />
            </label>

            <label className="flex items-center mt-3">
                Max per User per Stream
                <input defaultValue={reward.MaxPerUserPerStream} placeholder="0" type="number" spellCheck={false} {...register("maxPerUserPerStream")} className="form-input border-none bg-gray-700 mx-2 p-2 rounded shadow" />
            </label>

            <label className="flex items-center mt-3">
                Global Cooldown in Minutes
                <input defaultValue={(reward.GlobalCooldownSeconds ?? 0) / 60} placeholder="0" type="number" spellCheck={false} {...register("globalCooldownMinutes")} className="form-input border-none bg-gray-700 mx-2 p-2 rounded shadow" />
            </label>

            {children}

            <div className="flex flex-row justify-between items-center select-none">
                <label className="flex items-center">
                    <input defaultChecked={reward.Enabled} type="checkbox" {...register("enabled")} className="form-checkbox rounded border-gray-300 text-indigo-600 shadow-sm focus:border-indigo-300 focus:ring focus:ring-offset-0 focus:ring-indigo-200 focus:ring-opacity-50" />
                    <span className="ml-2">Enabled</span>
                </label>
                {loading && <span className="bg-blue-700 p-2 rounded shadow block mt-3">
                    <svg className="animate-spin mx-2 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                        <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                        <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                </span>}
                {!loading && <input type="submit" className="bg-green-700 hover:bg-green-600 p-2 rounded shadow block mt-3 cursor-pointer" value="save" />}
            </div>
        </form>
    );
}