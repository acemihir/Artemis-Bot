// ================================
const { Client, Options, MessageEmbed, MessageActionRow, MessageButton } = require('discord.js-light')
const config = require('./config')
const fs = require('fs')
const { botCache, getFromRedis } = require('./structures/cache')
const { REST } = require('@discordjs/rest')
const { Routes } = require('./node_modules/discord-api-types/v9')
const { printLog } = require('./utils')

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
client.on('ready', async (client) => {
    // Set an interval for the activity so all guilds are loaded/cached before counting.
    setInterval(async function () {
        const guilds_result = await client.shard.fetchClientValues('guilds.cache.size')
        const guildCount = guilds_result.reduce((prev, count) => prev + count, 0)

        client.user.setActivity(`${guildCount} servers | ${client.shard.count} shards`, {
            type: 'WATCHING'
        })
    }, 15 * 60 * 1000)

    printLog('Fully started.', 'INFO', client.shard.ids)
})

client.on('interactionCreate', async (interaction) => {
    if (interaction.isCommand()) {
        // Check if the used command is actually stored in the botCache object
        if (botCache.commands.has(interaction.commandName)) {
            // Retrieve the command data from the botCache object
            const obj = botCache.commands.get(interaction.commandName)

            // Fetch the guild data from the cache
            const cachedData = await getFromRedis(interaction.guildId)

            // Check if the command is a premium command
            if (obj.isPremium) {
                // Check if the guild does not have premium
                if (!cachedData.premium) {
                    // Construct the row
                    const row = new MessageActionRow().addComponents(new MessageButton()
                        .setURL('https://github.com/jerskisnow/Suggestions/wiki/Donating')
                        .setLabel('Donating')
                        .setEmoji('💰')
                        .setStyle('LINK'))

                    // Construct the embed
                    const embed = new MessageEmbed()
                    embed.setColor(config.embedColor.b)
                    embed.setTitle('Premium Command')
                    embed.setDescription('The command you tried to use is only for premium servers. See the button below for more information.')

                    // Send the message and return
                    return interaction.reply({ embeds: [embed], components: [row] })
                }
            }

            if (obj.execute.constructor.name === 'AsyncFunction') {
                await obj.execute(interaction.client, interaction)
            } else {
                obj.execute(interaction.client, interaction)
            }
        }
    } else if (interaction.isMessageComponent() && interaction.componentType === 'BUTTON') {

        // Check if the used button is actually stored in the botCache object
        if (botCache.buttons.has(interaction.customId)) {

            // Retrieve the interaction data from the botCache object and run the binded function
            botCache.buttons.get(interaction.customId)(interaction.client, interaction)
        }
    }
})

// ================================
const commandFiles = fs.readdirSync('./commands').filter(file => file.endsWith('.js'))
const commands = []

for (const file of commandFiles) {
    const cmdFile = require(`./commands/${file}`)
    const cmdName = file.split('.')[0]

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
        printLog('Started refreshing application (/) commands.', 'INFO', client.shard.ids)

        if (config.devMode) {
            await rest.put(Routes.applicationGuildCommands(config.botId, config.devGuild), { body: commands })
        } else {
            await rest.put(Routes.applicationCommands(config.botId), { body: commands })
        }

        printLog('Started refreshing application (/) commands.', 'INFO', client.shard.ids)
    } catch (error) {
        printLog(error, 'ERROR', client.shard.ids)
    }
})()

process.on('warning', w => printLog(w, 'WARN', client.shard.ids))

// ================================
client.login(config.botToken)