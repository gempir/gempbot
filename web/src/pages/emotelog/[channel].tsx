import { useRouter } from "next/router";
import { EmoteLogPage } from "../../components/EmoteLog/EmoteLogPage";
import { initializeStore } from "../../service/initializeStore";

export default function MediaChannel() {
    const router = useRouter()
    const { channel } = router.query

    return <EmoteLogPage channel={String(channel)} />
}

export const getServerSideProps = initializeStore;