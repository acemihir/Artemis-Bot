// ================================
const { promises } = require('fs')
const { Client } = require('discord.js')
const config = require('./config')
const fs = require('fs')
const { botCache } = require('./structures/cache')
const { REST } = require('@discordjs/rest')
const { Routes } = require('./node_modules/discord-api-types/v9')

// ================================
const client = new Client({
    intents: ['GUILDS', 'GUILD_MESSAGES']
})

// ================================
async function bindListeners() {
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
const commands = []
const commandFiles = fs.readdirSync('./commands').filter(file => file.endsWith('.js'))

for (const file of commandFiles) {
    const cmdFile = require(`./commands/${file}`)
    const cmdName = file.split('.')[0]

    // Check if the command is privileged
    if (cmdFile.command.privileged) {
        // Add the commandname to the privCommands array in the botCache
        botCache.privCommands.push(cmdName)
    }

    // delete cmdFile.command.privileged

    // Set the command
    botCache.commands.set(cmdName, cmdFile.command)
    // Check if there are any buttons
    if (cmdFile.buttons != null) {
        // Loop over the buttons
        for (let i = 0; i < cmdFile.buttons.length; i++) {
            // Set the (button) interaction
            botCache.buttons.set(cmdFile.buttons[i].id, cmdFile.buttons[i].onClick)
        }
    }

    commands.push(cmdFile.command.data.toJSON())
}

const rest = new REST({ version: '9' }).setToken(config.botToken);

(async () => {
    try {
        console.log('Started refreshing application (/) commands.')

        if (config.devMode) {
            await rest.put(Routes.applicationGuildCommands(config.botId, config.devGuild), { body: commands })
        } else {
            await rest.put(Routes.applicationCommands(config.botId), { body: commands })
        }

        console.log('Successfully reloaded application (/) commands.')
    } catch (error) {
        console.error(error)
    }
})()

// ================================
client.login(config.botToken)