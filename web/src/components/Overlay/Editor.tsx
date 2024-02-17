import { Editor, Tldraw, TldrawProps } from '@tldraw/tldraw';
import '@tldraw/tldraw/tldraw.css';
import { useYjsStore } from '../../hooks/useYjsStore';


type Props = {
    readonly?: boolean;
}

export function CustomEditor(props: Partial<TldrawProps> & Props) {
    const store = useYjsStore();

    const handleMount = (editor: Editor) => {
        console.log('editor mounted', props.readonly, editor);
        if (props.readonly) {
            editor.setCamera({ x: 0, y: 0, z: 1 });
            editor.updateInstanceState({ isReadonly: true, canMoveCamera: false })
            editor.selectNone();
        }
    }

    return <Tldraw onMount={handleMount} inferDarkMode store={store} {...props} />
}