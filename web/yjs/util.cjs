// @ts-nocheck
/// <reference lib="dom" />

const Y = require('yjs')
const syncProtocol = require('y-protocols/dist/sync.cjs')
const awarenessProtocol = require('y-protocols/dist/awareness.cjs')

const encoding = require('lib0/dist/encoding.cjs')
const decoding = require('lib0/dist/decoding.cjs')
const mutex = require('lib0/dist/mutex.cjs')
const map = require('lib0/dist/map.cjs')

const debounce = require('lodash.debounce')

const callbackHandler = require('./callback.cjs').callbackHandler
const isCallbackSet = require('./callback.cjs').isCallbackSet
const readSyncMessage = require('./readSyncMessageFork.cjs').readSyncMessage

const CALLBACK_DEBOUNCE_WAIT = parseInt(process.env.CALLBACK_DEBOUNCE_WAIT) || 2000
const CALLBACK_DEBOUNCE_MAXWAIT = parseInt(process.env.CALLBACK_DEBOUNCE_MAXWAIT) || 10000
let WRITE_STATE_DEBOUNCE_WAIT = 0

/**
 * @param {number} wait
 */
exports.setWriteStateDebounceWait = wait => {
    WRITE_STATE_DEBOUNCE_WAIT = wait
}

const wsReadyStateConnecting = 0
const wsReadyStateOpen = 1
const wsReadyStateClosing = 2 // eslint-disable-line
const wsReadyStateClosed = 3 // eslint-disable-line

// disable gc when using snapshots!
const gcEnabled = process.env.GC !== 'false' && process.env.GC !== '0'
const persistenceDir = process.env.YPERSISTENCE
/**
 * @type {{bindState: function(string,WSSharedDoc,WebSocket):void|Promise<void>, writeState:function(string,WSSharedDoc):Promise<any>, provider: any}|null}
 */
let persistence = null
if (typeof persistenceDir === 'string') {
    console.info('Persisting documents to "' + persistenceDir + '"')
    // @ts-ignore
    const LeveldbPersistence = require('y-leveldb').LeveldbPersistence
    const ldb = new LeveldbPersistence(persistenceDir)
    persistence = {
        provider: ldb,
        bindState: async (docName, ydoc) => {
            const persistedYdoc = await ldb.getYDoc(docName)
            const newUpdates = Y.encodeStateAsUpdate(ydoc)
            ldb.storeUpdate(docName, newUpdates)
            Y.applyUpdate(ydoc, Y.encodeStateAsUpdate(persistedYdoc))
            ydoc.on('update', update => {
                ldb.storeUpdate(docName, update)
            })
        },
        writeState: async (docName, ydoc) => { }
    }
}

/**
 * @param {{bindState: function(string,WSSharedDoc,WebSocket):Promise<void>,
 * writeState:function(string,WSSharedDoc):Promise<any>,provider:any}|null} persistence_
 */
exports.setPersistence = persistence_ => {
    persistence = persistence_
}

/**
 * @return {null|{bindState: function(string,WSSharedDoc,WebSocket):Promise<void>,
 * writeState:function(string,WSSharedDoc):Promise<any>,provider:any}|null} used persistence layer
 */
exports.getPersistence = () => persistence

/**
 * @type {Map<string,WSSharedDoc>}
 */
const docs = new Map()
// exporting docs so that others can use it
exports.docs = docs

const messageSync = 0
const messageAwareness = 1
// const messageAuth = 2

/**
 * @type {Map<string, number>}
 */
const closingTimeouts = new Map()
exports.closingTimeouts = closingTimeouts

/**
 * @param {Uint8Array} update
 * @param {any} origin
 * @param {WSSharedDoc} doc
 */
const updateHandler = (update, origin, doc) => {
    const encoder = encoding.createEncoder()
    encoding.writeVarUint(encoder, messageSync)
    syncProtocol.writeUpdate(encoder, update)
    const message = encoding.toUint8Array(encoder)
    if (origin !== 'bindState') {
        doc.conns.forEach((_, conn) => send(doc, conn, message))
    }
}

class WSSharedDoc extends Y.Doc {
    /**
     * @param {string} name
     * @param {WebSocket} [conn]
     */
    constructor(name, conn) {
        super({ gc: gcEnabled })
        this.name = name
        this.mux = mutex.createMutex()
        /**
         * Maps from conn to set of controlled user ids. Delete all user ids from awareness when this conn is closed
         * @type {Map<WebSocket, Set<number>>}
         */
        this.conns = new Map()
        /**
         * @type {import('y-protocols/awareness.js').Awareness}
         */
        this.awareness = new awarenessProtocol.Awareness(this)
        this.awareness.setLocalState(null)
        /**
         * @type {Promise<void>|void}
         */
        this.whenSynced = void 0
        /**
         * @param {{ added: Array<number>, updated: Array<number>, removed: Array<number> }} changes
         * @param {Object | null} conn Origin is the connection that made the change
         */
        const awarenessChangeHandler = ({ added, updated, removed }, conn) => {
            const changedClients = added.concat(updated, removed)
            if (conn !== null) {
                const connControlledIDs = /** @type {Set<number>} */ (this.conns.get(conn))
                if (connControlledIDs !== undefined) {
                    added.forEach(clientID => { connControlledIDs.add(clientID) })
                    removed.forEach(clientID => { connControlledIDs.delete(clientID) })
                }
            }
            // broadcast awareness update
            const encoder = encoding.createEncoder()
            encoding.writeVarUint(encoder, messageAwareness)
            encoding.writeVarUint8Array(encoder, awarenessProtocol.encodeAwarenessUpdate(this.awareness, changedClients))
            const buff = encoding.toUint8Array(encoder)
            this.conns.forEach((_, c) => {
                send(this, c, buff)
            })
        }
        this.awareness.on('update', awarenessChangeHandler)
        this.on('update', updateHandler)
        if (isCallbackSet) {
            this.on('update', debounce(
                callbackHandler,
                CALLBACK_DEBOUNCE_WAIT,
                { maxWait: CALLBACK_DEBOUNCE_MAXWAIT }
            ))
        }

        if (persistence !== null) {
            this.whenSynced = persistence.bindState(name, this, conn)
        }
    }
}

/**
 * @type {WSSharedDoc}
 */
exports.WSSharedDoc = WSSharedDoc

/**
 * Gets a Y.Doc by name, whether in memory or on disk
 *
 * @param {string} docname - the name of the Y.Doc to find or create
 * @param {boolean} gc - whether to allow gc on the doc (applies only when created)
 * @param {WebSocket} conn - (applies only when created)
 * @return {WSSharedDoc}
 */
const getYDoc = (docname, gc = true, conn) => map.setIfUndefined(docs, docname, () => {
    const doc = new WSSharedDoc(docname, conn)
    doc.gc = gc
    docs.set(docname, doc)
    return doc
})

exports.getYDoc = getYDoc

/**
 * @param {any} conn
 * @param {WSSharedDoc} doc
 * @param {Uint8Array} message
 */
const messageListener = async (conn, doc, message) => {
    try {
        const encoder = encoding.createEncoder()
        const decoder = decoding.createDecoder(message)
        const messageType = decoding.readVarUint(decoder)
        switch (messageType) {
            case messageSync:
                conn.readOnly = true
                // await the doc state being updated from persistence, if available, otherwise
                // we may send sync step 2 too early
                if (doc.whenSynced) {
                    await doc.whenSynced
                }
                encoding.writeVarUint(encoder, messageSync)
                const messageType = readSyncMessage(decoder, encoder, doc, conn.readOnly, null)
                if (encoding.length(encoder) > 1) {
                    send(doc, conn, encoding.toUint8Array(encoder))
                }
                if (typeof conn.reportStats === 'function' && messageType !== void 0) {
                    conn.reportStats({
                        docName: doc.name,
                        messageType: messageType,
                        bytes: message.length
                    })
                }
                break
            case messageAwareness: {
                awarenessProtocol.applyAwarenessUpdate(doc.awareness, decoding.readVarUint8Array(decoder), conn)
                break
            }
        }
    } catch (err) {
        doc.emit('error', [err])
    }
}

/**
 * @param {WSSharedDoc} doc
 * @param {any} conn
 */
const closeConn = (doc, conn) => {
    if (doc.conns.has(conn)) {
        /**
         * @type {Set<number>}
         */
        // @ts-ignore
        const controlledIds = doc.conns.get(conn)
        doc.conns.delete(conn)
        awarenessProtocol.removeAwarenessStates(doc.awareness, Array.from(controlledIds), null)
        if (doc.conns.size === 0 && persistence !== null) {
            cancelClosing(doc.name)
            closingTimeouts.set(
                doc.name,
                setTimeout(
                    () => {
                        closingTimeouts.delete(doc.name)
                        if (doc.conns.size === 0) {
                            // if persisted, we store state and destroy ydocument
                            persistence.writeState(doc.name, doc)
                                .then(
                                    async () => {
                                        if (doc.conns.size === 0) {
                                            docs.delete(doc.name)
                                            doc.destroy()
                                        }
                                    }
                                )
                        }
                    },
                    WRITE_STATE_DEBOUNCE_WAIT
                )
            )
        }
    }
    conn.close()
}

/**
 * @param {string} docName
 */
const cancelClosing = (docName) => {
    const timeout = closingTimeouts.get(docName)
    if (timeout) {
        clearTimeout(timeout)
        closingTimeouts.delete(docName)
    }
}

/**
 * @param {WSSharedDoc} doc
 * @param {any} conn
 * @param {Uint8Array} m
 */
const send = (doc, conn, m) => {
    if (conn.readyState !== wsReadyStateConnecting && conn.readyState !== wsReadyStateOpen) {
        closeConn(doc, conn)
    }
    try {
        conn.send(m, /** @param {any} err */ err => { err != null && closeConn(doc, conn) })
    } catch (e) {
        closeConn(doc, conn)
    }
}

const pingTimeout = 30000

/**
 * @param {WebSocket} conn
 * @param {import('http').IncomingMessage} req
 * @param {object} [opts]
 * @param {string} [opts.docName]
 * @param {boolean} [opts.gc]
 */
exports.setupWSConnection = async (conn, req, { docName = req.url.slice(1).split('?')[0], gc = true } = {}) => {
    conn.binaryType = 'arraybuffer'
    cancelClosing(docName)
    // get doc, initialize if it does not exist yet
    const doc = getYDoc(docName, gc, conn)
    doc.conns.set(conn, new Set())

    // listen and reply to events
    conn.on('message', /** @param {ArrayBuffer} message */ message => messageListener(conn, doc, new Uint8Array(message)))

    // Check if connection is still alive
    let pongReceived = true
    const pingInterval = setInterval(() => {
        if (!pongReceived) {
            if (doc.conns.has(conn)) {
                closeConn(doc, conn)
            }
            clearInterval(pingInterval)
        } else if (doc.conns.has(conn)) {
            pongReceived = false
            try {
                conn.ping()
            } catch (e) {
                closeConn(doc, conn)
                clearInterval(pingInterval)
            }
        }
    }, pingTimeout)
    conn.on('close', () => {
        closeConn(doc, conn)
        clearInterval(pingInterval)
    })
    conn.on('pong', () => {
        pongReceived = true
    })
    // put the following in a variables in a block so the interval handlers don't keep in in
    // scope
    {
        // await the doc state being updated from persistence, if available, otherwise
        // we may send sync step 1 too early
        if (doc.whenSynced) {
            await doc.whenSynced
        }

        // send sync step 1
        const encoder = encoding.createEncoder()
        encoding.writeVarUint(encoder, messageSync)
        syncProtocol.writeSyncStep1(encoder, doc)
        send(doc, conn, encoding.toUint8Array(encoder))
        const awarenessStates = doc.awareness.getStates()
        if (awarenessStates.size > 0) {
            const encoder = encoding.createEncoder()
            encoding.writeVarUint(encoder, messageAwareness)
            encoding.writeVarUint8Array(encoder, awarenessProtocol.encodeAwarenessUpdate(doc.awareness, Array.from(awarenessStates.keys())))
            send(doc, conn, encoding.toUint8Array(encoder))
        }
    }
}