// =================================
const config = require('./config')
const Filter = require('bad-words')
const { promises } = require('fs')

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
module.exports.printLog = function (txt, type, shard = null) {
    const d = new Date()

    let func, pref

    if (type === 'DEBUG') {
        func = console.debug
        pref = shard === null ? '[DEBUG]' : `[DEBUG-${shard}]`
    } else if (type === 'WARN') {
        func = console.warn
        pref = shard === null ? '[WARN ]' : `[WARN-${shard} ]`
    } else if (type === 'ERROR') {
        func = console.error
        pref = shard === null ? '[ERROR]' : `[ERROR-${shard}]`
    } else {
        func = console.info
        pref = shard === null ? '[INFO ]' : `[INFO-${shard} ]`
    }

    func(`${pref} ${d.getHours()}:${d.getSeconds()}: ${txt}`)
}

// =================================
module.exports.setPrivPermissions = async function (commands, appId, roleId) {
    const privCommands = []

    await promises.readdir(file => {
        if (require(`./commands/${file}`).command.privileged) {
            privCommands.push(file.split('.')[0])
        }
    })

    const permissions = []
    for (const [k, v] of commands.entries()) {
        if (v.applicationId === appId && privCommands.includes(v.name)) {
            permissions.push({
                id: k,
                permissions: [{
                    id: roleId,
                    type: 'ROLE',
                    permission: true,
                }]
            })
        }
    }

    await commands.permissions.set({ fullPermissions: permissions })
}