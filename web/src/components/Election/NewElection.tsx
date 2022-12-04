import { SubmitHandler, useForm } from "react-hook-form";

type FormValues = {
    hours: number;
    nominationCost: number;
};

export function NewElection() {
    const { register, handleSubmit, watch, formState: { errors } } = useForm<FormValues>();
    const onSubmit: SubmitHandler<FormValues> = data => console.log(data);

    return <form onSubmit={handleSubmit(onSubmit)} className="p-4 bg-gray-800 rounded shadow relative flex flex-col">
        <h2 className="mb-5 text-xl font-bold">Create a new 7tv emote election</h2>
        <label>
            Every X Hours
            <input defaultValue="24" {...register("hours", { required: true })} className="form-input border-none bg-gray-700 mx-2 p-2 rounded shadow" />
        </label>
        <br />
        <label>
            Nomination Cost
            <input defaultValue="5000" {...register("nominationCost", { required: true })} className="form-input border-none bg-gray-700 mx-2 p-2 rounded shadow" />
        </label>

        <input type="submit" className="hover:bg-gray-600 bg-gray-700 p-3 shadow rounded cursor-pointer font-bold mt-4 w-20 flex justify-center" />
    </form>;
}
