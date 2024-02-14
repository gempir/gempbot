import Link from "next/link";
import { useOverlays } from "../../hooks/useOverlays";
import { useUserConfig } from "../../hooks/useUserConfig";

const gempirId = "77829817";

export function OverlaysPage() {
    const [overlays, addOverlay, deleteOverlay, errorMessage, loading] = useOverlays();
    const [user] = useUserConfig();


    return <div className="relative w-full h-[100vh] p-4">
        <div className="p-4 bg-gray-800 rounded shadow max-w-[800px]">
            <div className="group">
                <button disabled={user?.Protected.CurrentUserID !== gempirId} onClick={addOverlay} className="bg-green-700 hover:bg-green-600 p-2 rounded shadow block cursor-pointer">Add Overlay</button>
                <span className="absolute z-50 hidden mt-2 p-2 w-48 text-center bg-black/75 text-white rounded tooltip-text group-hover:block pointer-events-none">
                    Only gempir can
                </span>
            </div>
            <div className="mt-5">
                {overlays.map(overlay => <div key={overlay.ID} className="flex items-center justify-between p-4 bg-gray-900">
                    <div>
                        <div className="group">
                            <button className="bg-red-700 hover:bg-red-600 p-2 rounded shadow block cursor-pointer" disabled={user?.Protected.CurrentUserID !== gempirId} onClick={() => {
                                confirm("Are you sure you want to delete this overlay?") && deleteOverlay(overlay.ID)
                            }}>Delete</button>
                            <span className="absolute z-50 hidden mt-2 p-2 w-48 text-center bg-black/75 text-white rounded tooltip-text group-hover:block pointer-events-none">
                                Only gempir can
                            </span>
                        </div>
                    </div>
                    <div>{overlay.ID}</div>
                    <div>
                        <Link href={`/overlay/edit/${overlay.ID}`} className="bg-blue-700 hover:bg-blue-600 p-2 rounded shadow block cursor-pointer">Edit</Link>
                    </div>
                    <div>
                        <input type="text" value={`${window?.location?.href}/${overlay.RoomID}`} readOnly className="bg-gray-900" />
                    </div>
                </div>)}
            </div>
        </div>
    </div>;
}