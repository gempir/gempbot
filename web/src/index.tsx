import { StrictMode, useContext } from 'react';
import ReactDOM from 'react-dom';
import { StateProvider, store } from './store';

function App() {
    const { state } = useContext(store);

    return <div>
        {state.apiBaseUrl}
    </div>
}

ReactDOM.render(
    <StrictMode>
        <StateProvider>
            <App />
        </StateProvider>
    </StrictMode>,
    document.getElementById('root')
);