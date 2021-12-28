// ================================
const { Client, Options, MessageEmbed, MessageActionRow, MessageButton } = require('discord.js-light');
const config = require('./config');
const fs = require('fs');
const { botCache, getFromRedis } = require('./structures/cache');
const { printLog } = require('./utils');

const fetch = (...args) => import('node-fetch').then(({ default: fetch }) => fetch(...args));

// ================================
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
});

// ================================
client.on('ready', async (client) => {
    // Set an interval for the activity so all guilds are loaded/cached before counting.
    setInterval(async function () {
        const guilds_result = await client.shard.fetchClientValues('guilds.cache.size');
        const guildCount = guilds_result.reduce((prev, count) => prev + count, 0);

        client.user.setActivity(`${guildCount} servers | ${client.shard.count} shards`, {
            type: 'WATCHING'
        });
    }, 15 * 60 * 1000);

    printLog('\u001b[32mFully initialized and ready.\u001b[0m', 'INFO', client.shard.ids);
});

client.on('interactionCreate', async (interaction) => {
    if (interaction.isCommand() && botCache.commands.has(interaction.commandName)) {
        const obj = botCache.commands.get(interaction.commandName);

        // This looks odd but saves us from calling getFromRedis twice
        if (obj.isPremium || obj.privileged) {
            const cachedData = await getFromRedis(interaction.guildId);

            if (obj.isPremium && !cachedData['premium']) {
                return interaction.reply({
                    embeds: [new MessageEmbed()
                        .setColor(config.embedColor.b)
                        .setTitle('Premium Command')
                        .setDescription('The command you tried to use is only for premium servers. See the button below for more information.')],
                    components: [new MessageActionRow().addComponents(new MessageButton()
                        .setURL('https://github.com/jerskisnow/Suggestions/wiki/Donating')
                        .setLabel('Donating')
                        .setEmoji('💰')
                        .setStyle('LINK'))]
                });
            }

            if (obj.privileged) {
                const embed = new MessageEmbed();
                embed.setColor(config.embedColor.r);
                embed.setTitle('Privileged Command');

                if (!cachedData['staff_role']) {
                    embed.setDescription('This command and privileged and therefore only usable by members with the staff role. (No staff role set)');
                    return interaction.reply({ embeds: [embed] });
                }

                const staffRole = await interaction.guild.roles.fetch(cachedData['staff_role']);
                if (!staffRole) {
                    embed.setDescription('This command and privileged and therefore only usable by members with the staff role. (Invalid staff role set)');
                    return interaction.reply({ embeds: [embed] });
                }

                if (!interaction.member.roles.cache.has(staffRole.id)) {
                    embed.setDescription('This command and privileged and therefore only usable by members with the staff role.');
                    return interaction.reply({ embeds: [embed] });
                }
            }
        }

        await obj.execute(interaction);

    } else if (interaction.isButton() && botCache.buttons.has(interaction.customId)) {
        // Retrieve the interaction data from the botCache object and run the binded function
        await botCache.buttons.get(interaction.customId)(interaction);
    }
});

// ================================
const commandFiles = fs.readdirSync('./commands').filter(file => file.endsWith('.js'));
const commands = [];

for (const file of commandFiles) {
    const cmdFile = require(`./commands/${file}`);
    const cmdName = file.split('.')[0];

    // Set the command
    botCache.commands.set(cmdName, cmdFile.command);
    // Check if there are any buttons
    if (cmdFile.buttons != null) {
        // Loop over the buttons
        for (let i = 0; i < cmdFile.buttons.length; i++) {
            // Set the (button) interaction
            botCache.buttons.set(cmdFile.buttons[i].id, cmdFile.buttons[i].onClick);
        }
    }

    commands.push(cmdFile.command.data.toJSON());
}

(async () => {
    try {
        printLog('Started refreshing application (/) commands.', 'INFO', client.shard.ids);

        const url = config.devMode ?
            `https://discord.com/api/v8/applications/${config.botId}/guilds/${config.devGuild}/commands` :
            `https://discord.com/api/v8/applications/${config.botId}/commands`;

        const res = await fetch(url, {
            method: 'PUT',
            body: JSON.stringify(commands),
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bot ' + config.botToken
            }
        });

        printLog(`Application (/) commands PUT response: ${res.statusText} (${res.status})`, 'INFO', client.shard.ids);
    } catch (error) {
        printLog(error, 'ERROR', client.shard.ids);
    }
})();

// ================================
// Some logging
process.on('warning', console.warn);
process.on('unhandledRejection', console.error);
process.on('uncaughtException', console.error);

// Client login
client.login(config.botToken);