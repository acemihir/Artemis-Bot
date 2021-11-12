// ================================
const { SlashCommandBuilder } = require('@discordjs/builders')
const { MessageEmbed } = require('discord.js')
const { getFromRedis, setInRedis, botCache } = require('../structures/cache')
const { runQuery } = require('../structures/database')
const config = require('../config')

// ================================
const data = new SlashCommandBuilder()
    .setName('setup')
    .setDescription('Setup the required settings for Suggestions to work properly.')

const execute = async function(client, interaction) {
    if (!interaction.member.permissions.has('ADMINISTRATOR')) {
        return interaction.reply({
            embeds: [new MessageEmbed()
                .setColor(config.embedColor.r)
                .setDescription('You need to have the `ADMINISTRATOR` permission to do that.')
            ]
        })
    }

    // ================================
    const filter = msg => msg.author.id === interaction.user.id

    // ================================
    // Suggestion channel
    const embed = new MessageEmbed()
        .setColor(config.embedColor.b)
        .setDescription('In which channel should suggestions show up? (Type: #channel)')

    await interaction.reply({ embeds: [embed] })

    const sugAwait = await interaction.channel.awaitMessages({ filter, max: 1, time: 20000, errors: ['time'] })
    await sugAwait.first().delete()

    const sugChannelId = sugAwait.first().content.replace('<#', '').replace('>', '')
    if (interaction.guild.channels.cache.get(sugChannelId) == null) {
        embed.setColor(config.embedColor.r)
        embed.setDescription('That\'s not a valid channel, please run the command again.')
        return interaction.editReply({ embeds: [embed] })
    }

    // ================================
    // Report channel
    embed.setDescription('In which channel should reports show up? (Type: #channel)')
    await interaction.editReply({ embeds: [embed] })

    const repAwait = await interaction.channel.awaitMessages({ filter, max: 1, time: 20000, errors: ['time'] })
    await repAwait.first().delete()

    const repChannelId = sugAwait.first().content.replace('<#', '').replace('>', '')
    if (interaction.guild.channels.cache.get(sugChannelId) == null) {
        embed.setColor(config.embedColor.r)
        embed.setDescription('That\'s not a valid channel, please run the command again.')
        return interaction.editReply({ embeds: [embed] })
    }   

    // ================================
    // Staff role
    embed.setDescription('Which role should be able to review suggestions, so which role should have access to the `/setstatus` command? (Type: @role)')
    await interaction.editReply({ embeds: [embed] })

    const roleAwait = await interaction.channel.awaitMessages({ filter, max: 1, time: 20000, errors: ['time'] })
    await roleAwait.first().delete()

    const roleId = roleAwait.first().content.replace('<@&', '').replace('>', '')
    if (interaction.guild.roles.cache.get(roleId) == null) {
        embed.setColor(config.embedColor.r)
        embed.setDescription('That\'s not a valid role, please run the command again.')
        return interaction.editReply({ embeds: [embed] })
    }

    // ================================
    embed.setColor(config.embedColor.g)
    embed.setDescription('That\'s all! You\'re done configurating suggestions! Please keep in mind that you will be able to change these and other settings using the `/config` command.')
    await interaction.editReply({ embeds: [embed] })

    // ================================
    // Update the cache
    var obj = await getFromRedis(interaction.guildId)
    obj['sug_channel'] = sugChannelId
    obj['rep_channel'] = repChannelId
    setInRedis(interaction.guildId, obj)

    // Update the discord command permission
    const commands = await interaction.guild.commands.fetch()
    const permissions  = []
    // Loop through all the commands
    for (const [key, value] of commands.entries()) {
        if (value.applicationId == interaction.applicationId) {
            if (botCache.privCommands.includes(value.name)) {
                permissions.push({
                    id: key,
                    permissions: [{
                        id: roleId,
                        type: 'ROLE',
                        permission: true,
                    }]
                })
            }
        }
    }
    // Set the actual permission
    await interaction.guild.commands.permissions.set({ fullPermissions: permissions })
    
    // Update the database values
    runQuery('UPDATE servers SET sug_channel = $1::text, rep_channel = $2::text WHERE id = $3::text', [sugChannelId, repChannelId, interaction.guildId])
}

// ================================
module.exports.command = {
    isPremium: false,
    privileged: false,

    data: data,
    execute: execute
}