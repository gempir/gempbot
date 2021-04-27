import React, { StrictMode } from 'react';
import ReactDOM from 'react-dom';
import { createGlobalStyle } from 'styled-components';
import { App } from './App';
import { StateProvider } from './store';

const GlobalStyle = createGlobalStyle`
  body {
    --bg: #0e0e10;
    --bg-bright: #18181b;
    --bg-brighter: #3d4146;
    --bg-dark: #121416;
    --theme: #00CC66;
    --theme-bright: #00FF80;
    --theme2: #2980b9;
    --theme2-bright: #3498db;
    --text: #F5F5F5;
    --text-dark: #616161;
    --twitch: #6441a5;

    background: var(--bg);
    margin: 0;
    padding: 0;
    color: var(--text);
    margin: 0;
    font-family: Helvetica, Arial, sans-serif;
    height: 100%;
    width: 100%;
  }
`

ReactDOM.render(
    <StrictMode>
        <StateProvider>
            <React.Fragment>
                <GlobalStyle />
                <App />
            </React.Fragment>
        </StateProvider>
    </StrictMode>,
    document.getElementById('root')
);