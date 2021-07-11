import { UserConfig } from "../../../hooks/useUserConfig";
import { RewardTypes } from "../../../types/Rewards";
import { BaseForm } from "./BaseForm";

export function BttvForm({ userConfig }: { userConfig: UserConfig }) {
    return <BaseForm type={RewardTypes.Bttv} userConfig={userConfig} header={<><img height="56" src="/images/bttv.png" alt="BetterTTV Logo" className="w-16" />
        <h3 className="text-xl font-bold">BetterTTV Emote</h3></>}

        description={<p className="my-2 mb-4 text-gray-400">
            <strong>Make sure <span className="text-green-600">gempbot</span> is BetterTTV editor</strong><br />
                This will swap out 1 emote constantly. If the previous emote is not found it will use a free slot or remove a random emote.
            </p>}
    />
}