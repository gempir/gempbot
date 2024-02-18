import { AssetRecordType, Editor, MediaHelpers, TLAsset, TLAssetId, Tldraw, TldrawProps, getHashForString } from '@tldraw/tldraw';
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
                const assetId: TLAssetId = AssetRecordType.createId(getHashForString(uploadedAsset.id))

                const size = uploadedAsset.isVideo ? await MediaHelpers.getVideoSizeFromSrc(uploadedAsset.url) : await MediaHelpers.getImageSizeFromSrc(uploadedAsset.url);

                const asset: TLAsset = AssetRecordType.create({
                    id: assetId,
                    type: uploadedAsset.isVideo ? 'video' : 'image',
                    typeName: 'asset',
                    props: {
                        name: uploadedAsset.id,
                        src: uploadedAsset.url,
                        w: size.w,
                        h: size.h,
                        mimeType: uploadedAsset.mimeType,
                        isAnimated: uploadedAsset.isAnimated,
                    },
                })

                return asset
            })
        }
    }

    return <Tldraw onMount={handleMount} inferDarkMode store={store} {...props} />
}