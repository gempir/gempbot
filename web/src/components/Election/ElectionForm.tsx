import dayjs from "dayjs";
import { useEffect } from "react";
import { SubmitHandler, useForm } from "react-hook-form";
import { Election } from "../../types/Election";

type FormValues = {
    Hours: number;
    NominationCost: number;
    EmoteAmount: number;
    MaxNominationPerUser: number;
    VoteAmount: number;
    SpecificTime: string | undefined;
};

const defaultElection = {
    Hours: 24,
    NominationCost: 1000,
    EmoteAmount: 1,
    MaxNominationPerUser: 3,
    VoteAmount: 3,
    SpecificTime: undefined,
}

type Props = {
    election: Election | undefined,
    setElection: (election: Election) => void,
    deleteElection: () => void,
    electionErrorMessage: string | null,
    electionLoading: boolean
};

export function ElectionForm({ election, setElection, deleteElection, electionErrorMessage, electionLoading }: Props) {
    const { register, handleSubmit, setValue, formState: { errors } } = useForm<FormValues>();
    const onSubmit: SubmitHandler<FormValues> = data => {

        let specTime;
        if (data.SpecificTime) {
            const [hours, minutes] = data.SpecificTime?.split(":") ?? [0, 0];
            const parsed = dayjs().set("hour", Number(hours)).set("minute", Number(minutes));
            if (parsed) {
                specTime = parsed;
            }
        }
        setElection({
            ...election as Election,
            Hours: Number(data.Hours),
            NominationCost: Number(data.NominationCost),
            EmoteAmount: Number(data.EmoteAmount),
            MaxNominationPerUser: Number(data.MaxNominationPerUser),
            VoteAmount: Number(data.VoteAmount),
            SpecificTime: specTime,
        });
    };

    useEffect(() => {
        if (election) {
            setValue("Hours", election.Hours);
            setValue("NominationCost", election.NominationCost);
            setValue("EmoteAmount", election.EmoteAmount);
            setValue("MaxNominationPerUser", election.MaxNominationPerUser);
            setValue("VoteAmount", election.VoteAmount)
            setValue("SpecificTime", election.SpecificTime?.format("HH:mm"));
        } else {
            setValue("Hours", defaultElection.Hours);
            setValue("NominationCost", defaultElection.NominationCost);
            setValue("EmoteAmount", defaultElection.EmoteAmount);
            setValue("MaxNominationPerUser", defaultElection.MaxNominationPerUser);
            setValue("VoteAmount", defaultElection.VoteAmount);
            setValue("SpecificTime", defaultElection.SpecificTime);
        }
    }, [election]);


    return <form onSubmit={handleSubmit(onSubmit)} className="p-4 bg-gray-800 rounded shadow relative flex flex-col min-w-[28rem]">
        <div className="mb-5 flex items-start justify-between gap-4">
            <div>
                <h2 className="text-xl font-bold">Create a new 7TV emote election</h2>
                {election && election.CreatedAt && <div className="text-gray-400 text-sm">Created at {election.CreatedAt.format("L LT")}</div>}
            </div>
            <div className="min-w-[5rem] min-h-[42px] align-top text-right">
                {election?.ChannelTwitchID &&
                    <div className="bg-red-700 hover:bg-red-600 p-2 rounded shadow text-gray-100 inline-block cursor-pointer" onClick={deleteElection}>
                        Delete
                    </div>
                }
            </div>
        </div>
        {electionErrorMessage && <div className="bg-red-500 text-white p-2 rounded mb-5">{electionErrorMessage}</div>}
        <label>
            Cooldown
            <input type="number" defaultValue={election?.Hours ?? defaultElection.Hours} {...register("Hours", { required: true })} className="form-input border-none bg-gray-700 mx-2 py-2 rounded shadow" />
            Hours
        </label>
        <br />
        <label>
            Specific Time
            <input type="time" defaultValue={election?.SpecificTime?.format("HH:mm")} {...register("SpecificTime", { required: false })} className="form-input border-none bg-gray-700 mx-2 py-2 rounded shadow" />
            <span className="text-gray-400">Optional</span>
        </label>
        <br />
        <label>
            Top Emotes to add
            <input type="number" defaultValue={election?.EmoteAmount ?? defaultElection.EmoteAmount} {...register("EmoteAmount", { required: true })} className="form-input border-none bg-gray-700 mx-2 p-2 rounded shadow" />
        </label>
        <br />
        <label>
            Max Nominations per User
            <input type="number" defaultValue={election?.MaxNominationPerUser ?? defaultElection.MaxNominationPerUser} {...register("MaxNominationPerUser", { required: true })} className="form-input border-none bg-gray-700 mx-2 p-2 rounded shadow" />
        </label>
        <br />
        <label>
            Votes per User
            <input type="number" defaultValue={election?.VoteAmount ?? defaultElection.VoteAmount} {...register("VoteAmount", { required: true })} className="form-input border-none bg-gray-700 mx-2 p-2 rounded shadow" />
        </label>
        <br />
        <label>
            Nomination Cost
            <input type="number" defaultValue={election?.NominationCost ?? defaultElection.NominationCost} {...register("NominationCost", { required: true })} className="form-input border-none bg-gray-700 mx-2 p-2 rounded shadow" />
        </label>

        <input type="submit" className={`hover:bg-gray-600 bg-gray-700 p-3 shadow rounded cursor-pointer font-bold mt-4 w-20 flex justify-center ${electionLoading ? "animate-pulse" : ""}`} />
    </form>;
}
