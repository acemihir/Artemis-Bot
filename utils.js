const config = require('./config')
const Filter = require('bad-words')

const wordFilter = new Filter()

module.exports.createId = function (prefix) {
    const chars = config.ids.chars

    let result = ''
    for (let i = 0; i < config.ids.charLength; i++) {
        result += chars.charAt(Math.floor(Math.random() * chars.length))
    }

    return prefix + result
}

module.exports.filterText = function (text) {
    return wordFilter.clean(text)
}