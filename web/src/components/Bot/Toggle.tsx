

export function Toggle({ checked, onChange }: { checked: boolean, onChange: (checked: boolean) => void }) {
    return <div className="flex flex-col justify-center items-center">
        {!checked && <div className="flex justify-center items-center cursor-pointer" onClick={() => onChange(!checked)}>
            <div className="w-14 h-7 flex items-center bg-gray-300 rounded-full px-1">
                <div className="bg-white w-5 h-5 rounded-full shadow-md transform" />
            </div>
        </div>}
        {checked &&
            <div className="flex justify-center items-center cursor-pointer" onClick={() => onChange(!checked)}>
                <div className="w-14 h-7 flex items-center rounded-full px-1 bg-blue-700">
                    <div className="bg-white w-5 h-5 rounded-full shadow-md transform translate-x-7" />
                </div>
            </div>}
    </div>
}