import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { Permission, SetUserConfig, UserConfig } from "../../hooks/useUserConfig";
import { Cross } from "../../icons/Cross";
import { isNumeric } from "../../service/isNumeric";

interface PermissionForm {
    permissions: Record<string | number, { User: string, Editor: boolean, Prediction: boolean }>
}

export function UserPermissions({ userConfig, setUserConfig }: { userConfig: UserConfig, setUserConfig: SetUserConfig }) {
    const [perms, setPerms] = useState(userConfig.Permissions);
    const { register, handleSubmit, setValue, reset, unregister } = useForm();

    useEffect(() => {
        reset({ permissions: userConfig.Permissions });
        setPerms(userConfig.Permissions);

        for (const [user, perm] of Object.entries(userConfig.Permissions)) {
            setValue(`permissions.${user}.User`, user)
            setValue(`permissions.${user}.Editor`, perm.Editor)
            setValue(`permissions.${user}.Prediction`, perm.Prediction)
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [JSON.stringify(userConfig.Permissions)]);

    const onSubmit = (data: PermissionForm) => {
        const perms: Record<string, Permission> = {};

        for (const [key, perm] of Object.entries(data.permissions)) {
            if (!isNumeric(key)) {
                continue
            }

            perms[perm.User] = { Editor: perm.Editor, Prediction: perm.Prediction };
        }

        setUserConfig({ ...userConfig, Permissions: perms })
    }

    const addRow = () => {
        const newPerms = { ...perms };
        newPerms["Username"] = { Editor: false, Prediction: true };

        setPerms(newPerms);
    }

    const removeRow = (user: string, index: number) => {
        const newPerms = { ...perms };
        delete newPerms[user];

        unregister(`permissions.${index}.User`);
        unregister(`permissions.${index}.Editor`);
        unregister(`permissions.${index}.Prediction`);

        setPerms(newPerms);
    };


    return <form onSubmit={handleSubmit(onSubmit)} className="p-4 bg-gray-800 rounded shadow relative">
        <h2 className="mb-4 text-xl">Permissions</h2>
        <table className="w-full">
            <thead>
                <tr className="border-b-8 border-transparent">
                    <th />
                    <th className="text-left pl-5">User</th>
                    <th className="px-5">Editor</th>
                    <th className="px-5">Prediction</th>
                </tr>
            </thead>
            <tbody>
                {Object.keys(perms).map((user, index) => <tr className={index % 2 ? "bg-gray-900" : ""} key={index}>
                    <th className="hover:text-red-600 cursor-pointer" onClick={() => removeRow(user, index)}><Cross /></th>
                    <th className="p-1"><input {...register(`permissions.${index}.User`)} className="p-1 bg-transparent leading-6" type="text" defaultValue={user} autoComplete={"off"} spellCheck={false} /> </th>
                    <th className="p-1"><input {...register(`permissions.${index}.Editor`)} className="p-1 bg-transparent leading-6" type="checkbox" defaultChecked={perms[user].Editor} /></th>
                    <th className="p-1"><input {...register(`permissions.${index}.Prediction`)} className="p-1 bg-transparent leading-6" type="checkbox" defaultChecked={perms[user].Prediction} /></th>
                </tr>)}
            </tbody>
        </table>
        <div className="hover:bg-gray-600 bg-gray-700 p-3 shadow rounded cursor-pointer font-bold mt-4 w-20 flex justify-center" onClick={addRow}>
            Add
        </div>
        <input type="submit" className="bg-green-700 hover:bg-green-600 p-2 rounded shadow block mt-3 cursor-pointer mr-0 ml-auto" value="save" />
    </form>
}