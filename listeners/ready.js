// ================================
const { promises } = require('fs')
const { REST } = require('@discordjs/rest');
const { Routes } = require('discord-api-types/v9');
const config = require('../config')
const { botCache } = require('../structures/cache')

const rest = new REST({ version: '9' }).setToken(config.botToken);

// ================================
module.exports = async function (client) {
    // Set an interval for the activity so all guilds are loaded/cached before counting.
    setInterval(async function () {
        const guilds_result = await client.shard.fetchClientValues('guilds.cache.size')
        const guildCount = guilds_result.reduce((prev, count) => prev + count, 0)

        await client.user.setActivity(`${guildCount} servers | ${client.shard.count} shards`, {
            type: 'WATCHING'
        })
    }, config.activityUpdateInterval * 60 * 1000)

    // Register the commands
    await registerCommands(client)

    console.log('Fully started.')
}

// ================================
async function registerCommands(client) {
    const commands = [];

    (await promises.readdir('./commands')).forEach(file => {
        const cmdFile = require('../commands/' + file)

        botCache.commands.set(file.split('.')[0], cmdFile.command)
        if (cmdFile.buttons != null) {
            for (let i = 0; i < cmdFile.buttons.length; i++) {
                botCache.buttonActions.set(cmdFile.buttons[i].id, cmdFile.buttons[i].onClick)
            }
        }

        commands.push(cmdFile.command.data.toJSON())
    })

    console.log(`Started refreshing application (/) commands. (DevMode: ${config.devMode ? 'Enabled' : 'Disabled'})`);

    try {
        if (config.devMode) {
            await rest.put(Routes.applicationGuildCommands(client.user.id, config.devGuild), { body: commands })
        } else {
            await rest.put(Routes.applicationCommands(client.user.id), { body: commands })
        }
    } catch (ex) {
        console.error(ex)
    }

    console.log(`Successfully reloaded application (/) commands. (DevMode: ${config.devMode ? 'Enabled' : 'Disabled'})`);
}

module.exports.once = true