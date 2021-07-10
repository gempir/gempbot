import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { Permission, SetUserConfig, UserConfig } from "../../hooks/useUserConfig";
import { isNumeric } from "../../service/isNumeric";

interface PermissionForm {
    permissions: Record<string | number, { User: string, Prediction: boolean }>
}

export function Permissions({ userConfig, setUserConfig }: { userConfig: UserConfig, setUserConfig: SetUserConfig }) {
    const [perms, setPerms] = useState(userConfig.Permissions);

    const { register, handleSubmit, setValue } = useForm();

    useEffect(() => {
        console.log("changing");
        setPerms(userConfig.Permissions);

        for (const [user, perm] of Object.entries(userConfig.Permissions)) {
            setValue(`permissions.${user}.User`, user)
            setValue(`permissions.${user}.Prediction`, perm.Prediction)
        }
    }, [JSON.stringify(userConfig.Permissions)]);

    const onSubmit = (data: PermissionForm) => {
        const perms: Record<string, Permission> = {};

        for (const [key, perm] of Object.entries(data.permissions)) {
            if (!isNumeric(key)) {
                continue
            }

            perms[perm.User] = { Prediction: perm.Prediction };
        }

        setUserConfig({ ...userConfig, Permissions: perms })
    }

    const addRow = () => {
        const newPerms = { ...perms };
        newPerms["Username"] = { Prediction: true };

        setPerms(newPerms);
    }


    return <form onSubmit={handleSubmit(onSubmit)} className="m-4 p-4 bg-gray-800 rounded shadow w-96 relative">
        <h2 className="mb-4 text-xl">Permissions</h2>
        <table className="w-full">
            <thead>
                <tr className="border-b-8 border-transparent">
                    <th>User</th>
                    <th>Prediction</th>
                </tr>
            </thead>
            <tbody>
                {Object.keys(perms).map((user, index) => <tr className={index % 2 ? "bg-gray-900" : ""} key={index}>
                    <th className="p-1"><input {...register(`permissions.${index}.User`)} className="p-1 bg-transparent leading-6" type="text" defaultValue={user} spellCheck={false} /> </th>
                    <th className="p-1"><input {...register(`permissions.${index}.Prediction`)} className="p-1 bg-transparent leading-6" type="checkbox" defaultChecked={perms[user].Prediction} /></th>
                </tr>)}
            </tbody>
        </table>
        <div className="hover:bg-gray-700 p-3 rounded cursor-pointer font-bold mt-4 w-full flex justify-center" onClick={addRow}>
            Add
        </div>
        <input type="submit" className="bg-green-700 hover:bg-green-600 p-2 rounded shadow block mt-3 cursor-pointer absolute bottom-4 right-4" value="save" />
    </form>
}