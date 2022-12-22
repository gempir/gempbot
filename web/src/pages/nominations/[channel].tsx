import { useRouter } from "next/router";
import { NominationsPage } from "../../components/Nominations/NominationsPage";
import { initializeStore } from "../../service/initializeStore";

export default function NominationsChannel() {
    const router = useRouter()
    const { channel } = router.query

    return <NominationsPage channel={String(channel)} />
}

export const getServerSideProps = initializeStore;