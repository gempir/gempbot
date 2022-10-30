import { useRouter } from "next/router";
import { MediaPage } from "../../components/Media/MediaPage";
import { initializeStore } from "../../service/initializeStore";

export default function MediaChannel() {
    const router = useRouter()
    const { channel } = router.query

    return <MediaPage channel={String(channel)} />
}

export const getServerSideProps = initializeStore;