#!/usr/bin/env node

import http from 'http';
import josnwebtoken from 'jsonwebtoken';
import WebSocket from 'ws';
import { setupWSConnection } from './util.cjs';

const wss = new WebSocket.Server({ noServer: true });

const host = process.env.HOST ?? '127.0.0.1'
const port = process.env.PORT ?? 1234
const jwtKey = parseEnv(process.env.SECRET ?? '')

const server = http.createServer((request, response) => {
    response.writeHead(200, { 'Content-Type': 'text/plain' })
    response.end('okay')
})

wss.on('connection', setupWSConnection)

server.on('upgrade', (request, socket, head) => {
    // You may check auth of request here..
    // See https://github.com/websockets/ws#client-authentication
    /**
     * @param {any} ws
     */
    const handleAuth = ws => {
        try {
            // // parse cookie 'scToken' 
            // const cookies = parseCookie(request.headers.cookie);
            // if (!cookies.scToken) {
            //     throw new Error('No cookie')
            // }

            // try {
            //     const payload = josnwebtoken.verify(cookies.scToken, jwtKey);

            //     // if user == managing etc..

            //     console.log('JWT Payload', payload);
            // } catch (e) {
            //     throw new Error('invalid scToken')
            // }

            wss.emit('connection', ws, request)
        } catch (e) {
            console.error(e.message);
            ws.close(1008, 'Not authorized')
        }
    }
    wss.handleUpgrade(request, socket, head, handleAuth)
})

server.listen(port, host, () => {
    console.log(`running at '${host}' on port ${port}`)
})

export function parseCookie(str) {
    if (!str || str.trim() === '') {
        return {};
    }

    return str
        .split(';')
        .map(v => v.split('='))
        .reduce((acc, v) => {
            acc[decodeURIComponent(v[0].trim())] = decodeURIComponent(v[1].trim());
            return acc;
        }, {});
}

function parseEnv(str) {
    if (typeof str !== 'string') {
        return '';
    }
    if (str.charAt(0) === '"' && str.charAt(str.length - 1) === '"') {
        return str.substr(1, str.length - 2);
    }
    return str;
}