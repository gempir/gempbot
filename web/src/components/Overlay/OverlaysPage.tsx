"use client";
import { Button, Table, TextInput, Tooltip, useMantineTheme, CopyButton } from "@mantine/core";
import Link from "next/link";
import { useOverlays } from "../../hooks/useOverlays";
import { useUserConfig } from "../../hooks/useUserConfig";

export function OverlaysPage() {
    const [overlays, addOverlay, deleteOverlay, errorMessage, loading] = useOverlays();
    const [user] = useUserConfig();

    const theme = useMantineTheme();

    return <div className="relative w-full h-[100vh] p-4">
        <div className="p-4 bg-gray-800 rounded shadow max-w-[800px]">
            <div className="group">
                <Tooltip label="only gempir can">
                    <Button variant="outline" onClick={addOverlay}>Add Overlay</Button>
                </Tooltip>
            </div>
            <Table verticalSpacing={"lg"}>
                <Table.Thead>
                    <Table.Tr>
                        <Table.Th>Id</Table.Th>
                        <Table.Th>Overlay Link</Table.Th>
                        <Table.Th>Symbol</Table.Th>
                    </Table.Tr>
                </Table.Thead>
                <Table.Tbody>
                    {overlays.map(overlay => <Table.Tr key={overlay.ID}>
                        <Table.Td>{overlay.ID}</Table.Td>
                        <Table.Td className="flex gap-3 min-w-[350px]">
                            <TextInput style={{ maxWidth: 300 }} value={`${window?.location?.href}/${overlay.RoomID}`} readOnly />
                            <CopyButton value={`${window?.location?.href}/${overlay.RoomID}`}>
                                {({ copied, copy }) => (
                                    <Button color={copied ? 'teal' : 'blue'} onClick={copy}>
                                        {copied ? 'Copied url' : 'Copy url'}
                                    </Button>
                                )}
                            </CopyButton>
                        </Table.Td>
                        <Table.Td>
                            <div className="flex gap-3">
                                <Tooltip label="only gempir can">
                                    <Button variant="contained" style={{ backgroundColor: theme.colors.red[9], opacity: 0.1 }}
                                        onClick={() => {
                                            confirm("Are you sure you want to delete this overlay?") && deleteOverlay(overlay.ID)
                                        }}
                                    >Delete</Button>
                                </Tooltip>
                                <Button component={Link} href={`/overlay/edit/${overlay.ID}`}>Edit</Button>
                            </div>
                        </Table.Td>
                    </Table.Tr>)}
                </Table.Tbody>
            </Table>
        </div>
    </div>;
}