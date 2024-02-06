import { Tldraw, TldrawProps, useEditor } from '@tldraw/tldraw';
import '@tldraw/tldraw/tldraw.css';
import { useYjsStore } from '../../hooks/useYjsStore';
import { useStore } from '../../store';


type Props = {
    roomId: string;
    readonly?: boolean;
}

export function Editor(props: Partial<TldrawProps> & Props) {
    const yjsWsUrl = useStore(state => state.yjsWsUrl);
    const store = useYjsStore({
        roomId: props.roomId,
        hostUrl: yjsWsUrl,
    });
    const editor = useEditor();


    if (props.readonly && editor) {
        console.log(editor);
        // editor.setCamera({ x: 0, y: 0, z: 1 });
        // editor.updateInstanceState({ isReadonly: true, canMoveCamera: false })
    }

    return <Tldraw inferDarkMode store={store} {...props} />
}