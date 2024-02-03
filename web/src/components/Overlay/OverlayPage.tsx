'use client';

import dynamic from 'next/dynamic';

const Editor = dynamic(() => import('./Editor'), { ssr: false })

export function OverlayPage() {
    return (
        <div style={{ position: 'fixed', inset: 0 }}>
			<Editor />
		</div>
    );
}