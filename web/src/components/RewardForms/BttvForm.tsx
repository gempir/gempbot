import { useForm } from "react-hook-form";

export function BttvForm() {
    const { register, handleSubmit, watch, formState: { errors } } = useForm();
    // @ts-ignore
    const onSubmit = data => console.log(data);

    console.log(watch("example")); // watch input value by passing the name of it

    return (
        /* "handleSubmit" will validate your inputs before invoking "onSubmit" */
        <form onSubmit={handleSubmit(onSubmit)} className="m-4 p-4 bg-gray-800 rounded shadow">
            {/* register your input into the hook by invoking the "register" function */}
            <input defaultValue="Bttv Emote" {...register("name")} className="bg-gray-700 p-2 rounded shadow block" />

            {/* include validation with required or other standard HTML validation rules */}
            <input {...register("exampleRequired", { required: true })} className="bg-gray-700 p-2 rounded shadow block mt-3" />
            {/* errors will return when field validation fails  */}
            {errors.exampleRequired && <span>This field is required</span>}

            <div className="flex justify-end">
                <input type="submit" className="bg-green-700 hover:bg-green-600 p-2 rounded shadow block mt-3" value="save" />
            </div>
        </form>
    );
}