// ================================
const { SlashCommandBuilder } = require('@discordjs/builders')
const { MessageEmbed } = require('discord.js')
const { getFromRedis, setInRedis } = require('../structures/cache')
const { runQuery } = require('../structures/database')
const config = require('../config')

// ================================
const data = new SlashCommandBuilder()
    .setName('setup')
    .setDescription('Setup the required settings for Suggestions to work properly.')

const execute = async function(_client, interaction) {
    // ================================
    const embed = new MessageEmbed()
        .setColor(config.embedColor.b)
        .setDescription('In which channel should suggestions show up? (Type: #channel)')

    await interaction.reply({ embeds: [embed] })

    // ================================
    // await filter
    const filter = msg => msg.author.id === interaction.user.id

    // ================================
    const chnAwait = await interaction.channel.awaitMessages({ filter, max: 1, time: 20000, errors: ['time'] })
    await chnAwait.first().delete()

    const channelId = chnAwait.first().content.replace('<#', '').replace('>', '')
    if (interaction.guild.channels.cache.get(channelId) == null) {
        embed.setDescription('That\'s not a valid channel, please run the command again.')
        embed.setColor(config.embedColor.r)
        await interaction.editReply({ embeds: [embed] })
        return
    }

    // ================================
    embed.setDescription('Which role should be able to review suggestions, so which role should have access to the `/setstatus` command? (Type: @role)')
    await interaction.editReply({ embeds: [embed] })

    const roleAwait = await interaction.channel.awaitMessages({ filter, max: 1, time: 20000, errors: ['time'] })
    await roleAwait.first().delete()

    const roleId = roleAwait.first().content.replace('<@&', '').replace('>', '')
    if (interaction.guild.roles.cache.get(roleId) == null) {
        embed.setDescription('That\'s not a valid role, please run the command again.')
        embed.setColor(config.embedColor.r)
        await interaction.editReply({ embeds: [embed] })
        return
    }

    // ================================
    embed.setDescription('That\'s all! You\'re done configurating suggestions! Please keep in mind that you will be able to change these and other settings using the `/config` command.')
    embed.setColor(config.embedColor.g)
    await interaction.editReply({ embeds: [embed] })

    // ================================
    // Saving the data
    var obj = await getFromRedis(interaction.guildId)
    obj.staffRole = roleId
    setInRedis(interaction.guildId, obj)

    runQuery('UPDATE servers SET suggestion_channel = $1::text, staff_role = $2::text WHERE id = $3::text', [channelId, roleId, interaction.guildId])
}

// ================================
module.exports.command = {
    isPremium: false,
    permLevel: 2,

    data: data,
    execute: execute
}