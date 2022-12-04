import { useEffect } from "react";
import { SubmitHandler, useForm } from "react-hook-form";
import { useElection } from "../../hooks/useElection";

type FormValues = {
    Hours: number;
    NominationCost: number;
};

export function ElectionForm() {
    const [election, setElection, deleteElection, errorMessage, loading] = useElection();
    const { register, handleSubmit, setValue, formState: { errors } } = useForm<FormValues>();
    const onSubmit: SubmitHandler<FormValues> = data => {
        setElection({ ...election, Hours: Number(data.Hours), NominationCost: Number(data.NominationCost) });
    };

    useEffect(() => {
        setValue("Hours", election.Hours);
        setValue("NominationCost", election.NominationCost);
    }, [election]);



    return <form onSubmit={handleSubmit(onSubmit)} className="p-4 bg-gray-800 rounded shadow relative flex flex-col">
        <div className="mb-5 flex items-start justify-between gap-3">
            <h2 className="text-xl font-bold">Create a new 7TV emote election</h2>
            <div className="min-w-[5rem] min-h-[42px] align-top text-right">
                {election?.ID &&
                    <div className="bg-red-700 hover:bg-red-600 p-2 rounded shadow text-gray-100 inline-block cursor-pointer" onClick={deleteElection}>
                        Delete
                    </div>
                }
            </div>
        </div>
        {errorMessage && <div className="bg-red-500 text-white p-2 rounded mb-5">{errorMessage}</div>}
        <label>
            Every
            <input defaultValue={election.Hours} {...register("Hours", { required: true })} className="form-input border-none bg-gray-700 mx-2 py-2 rounded shadow" />
            Hours
        </label>
        <br />
        <label>
            Nomination Cost
            <input defaultValue={election.NominationCost} {...register("NominationCost", { required: true })} className="form-input border-none bg-gray-700 mx-2 p-2 rounded shadow" />
        </label>

        <input type="submit" className="hover:bg-gray-600 bg-gray-700 p-3 shadow rounded cursor-pointer font-bold mt-4 w-20 flex justify-center" />
    </form>;
}
