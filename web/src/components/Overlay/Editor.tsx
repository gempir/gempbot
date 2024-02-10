import { Editor, Tldraw, TldrawProps } from '@tldraw/tldraw';
import '@tldraw/tldraw/tldraw.css';
import { useYjsStore } from '../../hooks/useYjsStore';
import { useStore } from '../../store';


type Props = {
    roomId: string;
    readonly?: boolean;
}

export function CustomEditor(props: Partial<TldrawProps> & Props) {
    const yjsWsUrl = useStore(state => state.yjsWsUrl);
    const store = useYjsStore({
        roomId: props.roomId,
        hostUrl: yjsWsUrl,
    });
    
    const handleMount = (editor: Editor) => {
        editor.setCamera({ x: 0, y: 0, z: 1 });
        editor.updateInstanceState({ isReadonly: true, canMoveCamera: false })
		editor.selectNone();
	}

    return <Tldraw onMount={handleMount} inferDarkMode store={store} {...props} />
}