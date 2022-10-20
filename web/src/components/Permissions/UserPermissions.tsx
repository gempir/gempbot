import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { Permission, SetUserConfig, UserConfig } from "../../hooks/useUserConfig";
import { isNumeric } from "../../service/isNumeric";
import { XIcon } from "@heroicons/react/solid";

interface PermissionForm {
    permissions: Record<string | number, { User: string, Editor: boolean, Prediction: boolean }>
}

export function UserPermissions({ userConfig, setUserConfig, errorMessage, loading }: { userConfig: UserConfig, setUserConfig: SetUserConfig, errorMessage?: string, loading?: boolean }) {
    const [perms, setPerms] = useState(userConfig.Permissions);
    const { register, handleSubmit, setValue, reset, unregister, setFocus } = useForm();

    const [addCounter, setAddCounter] = useState(0);

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

            perms[perm.User.toLowerCase()] = { Editor: perm.Editor, Prediction: perm.Prediction };
        }

        setUserConfig({ ...userConfig, Permissions: perms })
    }

    const addRow = () => {
        const newPerms = { ...perms };
        newPerms["user" + addCounter] = { Editor: false, Prediction: true };

        setPerms(newPerms);
        setAddCounter(addCounter + 1);
    }

    const removeRow = (user: string, index: number) => {
        const newPerms = { ...perms };
        delete newPerms[user];

        unregister(`permissions.${index}.User`);
        unregister(`permissions.${index}.Editor`);
        unregister(`permissions.${index}.Prediction`);

        setPerms(newPerms);
    };


    useEffect(() => {
        if (addCounter > 0) {
            try {
                setFocus(`permissions.${Object.keys(perms).length - 1}.User`);
            } catch (e) {
            }
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [addCounter]);

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
                    <th className="hover:text-red-600 cursor-pointer" onClick={() => removeRow(user, index)}><XIcon className="h-6" /></th>
                    <th className="p-1"><input {...register(`permissions.${index}.User`)} className="p-1 bg-transparent leading-6" type="text" defaultValue={user} autoComplete={"off"} spellCheck={false} /> </th>
                    <th className="p-1"><input {...register(`permissions.${index}.Editor`)} className="p-1 bg-transparent leading-6" type="checkbox" defaultChecked={perms[user].Editor} /></th>
                    <th className="p-1"><input {...register(`permissions.${index}.Prediction`)} className="p-1 bg-transparent leading-6" type="checkbox" defaultChecked={perms[user].Prediction} /></th>
                </tr>)}
            </tbody>
        </table>
        <div className="hover:bg-gray-600 bg-gray-700 p-3 shadow rounded cursor-pointer font-bold mt-4 w-20 flex justify-center" onClick={addRow}>
            Add
        </div>
        <div className="flex justify-between items-center mt-4">
            <div className="text-red-700 max-w-xs">{errorMessage}</div>
            {!loading && <input type="submit" className="bg-green-700 hover:bg-green-600 p-2 rounded shadow block cursor-pointer" value="save" />}
            {loading && <span className="bg-blue-700 p-2 rounded shadow block">
                    <svg className="animate-spin mx-2 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                        <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                        <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
            </span>}
        </div>
    </form>
}