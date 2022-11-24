import { SubmitHandler, useForm } from "react-hook-form";
import React from "react";

type FormValues = {
    title: string;
};

export function NewElection() {
    const { register, handleSubmit, watch, formState: { errors } } = useForm<FormValues>();
    const onSubmit: SubmitHandler<FormValues> = data => console.log(data);

    return <form onSubmit={handleSubmit(onSubmit)}>
        <input defaultValue="test" {...register("title", { required: true })} className="form-input border-none bg-gray-700 mx-2 p-2 rounded shadow" />
        {errors.title && <span>This field is required</span>}

        <input type="submit" />
    </form>;
}
