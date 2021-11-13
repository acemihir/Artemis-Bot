// ================================
const { Client, Options } = require('discord.js-light')
const config = require('./config')
const fs = require('fs')
const { botCache } = require('./structures/cache')
const { REST } = require('@discordjs/rest')
const { Routes } = require('./node_modules/discord-api-types/v9')

const client = new Client({
    makeCache: Options.cacheWithLimits({
        ApplicationCommandManager: 0, // guild.commands
        BaseGuildEmojiManager: 0, // guild.emojis
        ChannelManager: 0, // client.channels
        GuildChannelManager: 0, // guild.channels
        GuildBanManager: 0, // guild.bans
        GuildInviteManager: 0, // guild.invites
        GuildManager: Infinity, // client.guilds
        GuildMemberManager: 0, // guild.members
        GuildStickerManager: 0, // guild.stickers
        MessageManager: 0, // channel.messages
        PermissionOverwriteManager: 0, // channel.permissionOverwrites
        PresenceManager: 0, // guild.presences
        ReactionManager: 0, // message.reactions
        ReactionUserManager: 0, // reaction.users
        RoleManager: 0, // guild.roles
        StageInstanceManager: 0, // guild.stageInstances
        ThreadManager: 0, // channel.threads
        ThreadMemberManager: 0, // threadchannel.members
        UserManager: 0, // client.users
        VoiceStateManager: 0 // guild.voiceStates
    }),
    intents: ['GUILDS', 'GUILD_MESSAGES'] 
})

// ================================
const eventFiles = fs.readdirSync('./listeners').filter(file => file.endsWith('.js'))
for (const file of eventFiles) {
    const event = require(`./listeners/${file}`)
    if (event.once) {
        client.once(file.split('.')[0], (...args) => event.execute(...args))
    } else {
        client.on(file.split('.')[0], (...args) => event.execute(...args))
    }
}

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

process.on('warning', console.warn)

// ================================
client.login(config.botToken)