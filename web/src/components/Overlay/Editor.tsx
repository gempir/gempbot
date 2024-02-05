import { Tldraw, TldrawProps } from '@tldraw/tldraw';
import '@tldraw/tldraw/tldraw.css';
import { useYjsStore } from '../../hooks/useYjsStore';
import { useStore } from '../../store';


export function Editor(props: Partial<TldrawProps>) {
    const yjsWsUrl = useStore(state => state.yjsWsUrl);
    const store = useYjsStore({
        roomId: 'example17',
        hostUrl: yjsWsUrl,
    });

    return <Tldraw inferDarkMode store={store} {...props} />
}