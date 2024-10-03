import { useSync } from '@tldraw/sync';
import 'tldraw/tldraw.css';
import { AssetRecordType, Editor, TLAssetStore, TLBookmarkAsset, Tldraw, TldrawProps, getHashForString, uniqueId } from 'tldraw';

type Props = {
    readonly?: boolean;
    roomID: string;
}

export function CustomEditor(props: Partial<TldrawProps> & Props) {
    // Create a store connected to multiplayer.
    const store = useSync({
        // We need to know the websocket's URI...
        uri: `https://bot-tldraw.gempir.com/connect/${props.roomID}`,
        // ...and how to handle static assets like images & videos
        assets: multiplayerAssets,
    })

    // const upload = useAssetUploader();

    const handleMount = (editor: Editor) => {
        // @ts-expect-error
        window.editor = editor;
        console.log('editor mounted', props.readonly, editor);
        if (props.readonly) {
            editor.setCamera({ x: 0, y: 0, z: 1 });
            editor.setCameraOptions({ isLocked: true })
            editor.updateInstanceState({ isReadonly: true })
            editor.selectNone();
        } else {
            editor.registerExternalAssetHandler('url', unfurlBookmarkUrl)
        }
    }

    return <Tldraw onMount={handleMount} inferDarkMode store={store} {...props} Background={<></>} />
}

// How does our server handle assets like images and videos?
const multiplayerAssets: TLAssetStore = {
    // to upload an asset, we prefix it with a unique id, POST it to our worker, and return the URL
    async upload(_asset, file) {
        const id = uniqueId()

        const objectName = `${id}-${file.name}`
        const url = `https://bot-tldraw.gempir.com/uploads/${encodeURIComponent(objectName)}`

        const response = await fetch(url, {
            method: 'PUT',
            body: file,
        })

        if (!response.ok) {
            throw new Error(`Failed to upload asset: ${response.statusText}`)
        }

        return url
    },
    // to retrieve an asset, we can just use the same URL. you could customize this to add extra
    // auth, or to serve optimized versions / sizes of the asset.
    resolve(asset) {
        return asset.props.src
    },
}

// How does our server handle bookmark unfurling?
async function unfurlBookmarkUrl({ url }: { url: string }): Promise<TLBookmarkAsset> {
    const asset: TLBookmarkAsset = {
        id: AssetRecordType.createId(getHashForString(url)),
        typeName: 'asset',
        type: 'bookmark',
        meta: {},
        props: {
            src: url,
            description: '',
            image: '',
            favicon: '',
            title: '',
        },
    }

    try {
        const response = await fetch(`https://bot-tldraw.gempir.com/unfurl?url=${encodeURIComponent(url)}`)
        const data = await response.json()

        asset.props.description = data?.description ?? ''
        asset.props.image = data?.image ?? ''
        asset.props.favicon = data?.favicon ?? ''
        asset.props.title = data?.title ?? ''
    } catch (e) {
        console.error(e)
    }

    return asset
}
