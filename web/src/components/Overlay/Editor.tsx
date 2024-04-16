import { AssetRecordType, Editor, TLAsset, TLAssetId, Tldraw, TldrawProps, getHashForString } from '@tldraw/tldraw';
import { MediaHelpers, isGifAnimated } from 'tldraw';
import '@tldraw/tldraw/tldraw.css';
import { useAssetUploader } from '../../hooks/useAssetUploader';
import { useYjsStore } from '../../hooks/useYjsStore';


type Props = {
    readonly?: boolean;
}

export function CustomEditor(props: Partial<TldrawProps> & Props) {
    const store = useYjsStore();
    const upload = useAssetUploader();

    const handleMount = (editor: Editor) => {
        console.log('editor mounted', props.readonly, editor);
        if (props.readonly) {
            editor.setCamera({ x: 0, y: 0, z: 1 });
            editor.updateInstanceState({ isReadonly: true, canMoveCamera: false })
            editor.selectNone();
        } else {
            editor.registerExternalAssetHandler('file', async ({ file }: { type: 'file'; file: File }) => {
                const uploadedAsset = await upload(file);
                //[b]
                const assetId: TLAssetId = AssetRecordType.createId(getHashForString(uploadedAsset.url))
    
                let size: {
                    w: number
                    h: number
                }
                let isAnimated: boolean
                let shapeType: 'image' | 'video'
    
                //[c]
                if (['image/jpeg', 'image/png', 'image/gif', 'image/svg+xml'].includes(file.type)) {
                    shapeType = 'image'
                    size = await MediaHelpers.getImageSize(file)
                    isAnimated = file.type === 'image/gif' && (await isGifAnimated(file))
                } else {
                    shapeType = 'video'
                    isAnimated = true
                    size = await MediaHelpers.getVideoSize(file)
                }
                //[d]
                const asset: TLAsset = AssetRecordType.create({
                    id: assetId,
                    type: shapeType,
                    typeName: 'asset',
                    props: {
                        name: file.name,
                        src: uploadedAsset.url,
                        w: size.w,
                        h: size.h,
                        mimeType: file.type,
                        isAnimated,
                    },
                })
    
                return asset
            })
        }
    }

    return <Tldraw onMount={handleMount} inferDarkMode store={store} {...props} />
}