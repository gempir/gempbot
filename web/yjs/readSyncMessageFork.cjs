// @ts-nocheck
const syncProtocol = require('y-protocols/dist/sync.cjs')

const decoding = require('lib0/dist/decoding.cjs')

/**
 * @param {decoding.Decoder} decoder A message received from another client
 * @param {encoding.Encoder} encoder The reply message. Will not be sent if empty.
 * @param {Y.Doc} doc
 * @param {boolean} readOnly If true, updates will be silently ignored instead of applied.
 * @param {any} transactionOrigin
 * @returns {number|undefined}
 */
const readSyncMessage = (decoder, encoder, doc, readOnly = false, transactionOrigin) => {
  const messageType = decoding.readVarUint(decoder)
  switch (messageType) {
    case syncProtocol.messageYjsSyncStep1:
      syncProtocol.readSyncStep1(decoder, encoder, doc)
      break
    case syncProtocol.messageYjsSyncStep2:
      if (readOnly) return
      syncProtocol.readSyncStep2(decoder, doc, transactionOrigin)
      break
    case syncProtocol.messageYjsUpdate:
      if (readOnly) return
      syncProtocol.readUpdate(decoder, doc, transactionOrigin)
      break
    default:
      throw new Error('Unknown message type')
  }
  return messageType
}

module.exports = { readSyncMessage }