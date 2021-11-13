// =================================
const config = require('./config')
const Filter = require('bad-words')

// =================================
module.exports.createId = function (prefix) {
    const chars = config.ids.chars

    let result = ''
    for (let i = 0; i < config.ids.charLength; i++) {
        result += chars.charAt(Math.floor(Math.random() * chars.length))
    }

    return prefix + result
}

// =================================
const wordFilter = new Filter()

module.exports.filterText = function (text) {
    return wordFilter.clean(text)
}

// =================================
module.exports.infoLog = function (text) {
    const date = new Date()

    console.log(`[INFO] ${date.getHours()}:${date.getSeconds()} >> ${text}`)
}

module.exports.debugLog = function (text) {
    const date = new Date()

    console.debug(`[DEBUG] ${date.getHours()}:${date.getSeconds()} >> ${text}`)
}

module.exports.warnLog = function (text) {
    const date = new Date()

    console.warn(`[WARN] ${date.getHours()}:${date.getSeconds()} >> ${text}`)
}

module.exports.errorLog = function (text) {
    const date = new Date()

    console.error(`[ERROR] ${date.getHours()}:${date.getSeconds()} >> ${text}`)
}