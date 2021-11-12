// ================================
const { promises } = require('fs')
const { Client } = require('discord.js')
const config = require('./config')

// ================================
const client = new Client({
    intents: ['GUILDS', 'GUILD_MESSAGES']
})

// ================================
const bindListeners = async function() {
    console.log('IS CLIENT READY? ', client.isReady());
    (await promises.readdir('./listeners')).forEach(file => {
        const obj = require(`./listeners/${file}`)
        if (obj.once) {
            client.once(file.split('.')[0], obj.bind(null, client))
        } else {
            client.on(file.split('.')[0], obj.bind(null, client))
        }
    })
}

bindListeners()

// ================================
client.login(config.botToken)