// ================================
const { SlashCommandBuilder } = require('@discordjs/builders')
const { createId, filterText } = require('../utils')
const { runQuery } = require('../structures/database')
const { MessageEmbed } = require('discord.js')
const config = require('../config')
const fetch = require('node-fetch')

// ================================
const data = new SlashCommandBuilder()
    .setName('report')
    .setDescription('Create a report.')
    .addStringOption(opt =>
        opt.setName('description')
            .setDescription('A brief description of your report.')
            .setRequired(true))

const execute = async function(client, interaction) {
    // Fetch the input/args
    const repDesc = await interaction.options.getString('description')

    const result = await runQuery('SELECT report_channel FROM servers WHERE id = $1::text', [interaction.guild.id])
    if (result == null || !result.rows.length) {
        await interaction.reply('ERROR: Please make sure an administrator has configured the report channel.')
        return
    }

    const repChannel = await interaction.guild.channels.fetch(result.rows[0].report_channel)
    if (repChannel == null) {
        await interaction.reply('ERROR: The configured report channel was not found.')
        return
    }

    const repId = createId('r_')

    const embed = new MessageEmbed()
        .setAuthor(interaction.author.tag, interaction.author.avatarURL())
        .setColor('#7583ff')
        .addField('Description', filterText(repDesc))
        .addField('Information', `**Status:** Open\n**ID":** ${repId}`)

    const msg = await repChannel.send({ embeds: [embed] })

    await interaction.reply('Your report has been submitted.')

    fetch(`${config.backend.url}/submit`, {
        method: 'POST',
        body: JSON.stringify({
            id: repId,
            context: repDesc,
            author: interaction.user.id,
            avatar: interaction.user.avatarURL(),
            guild: interaction.guildId,
            channel: repChannel.id,
            message: msg.id,
            status: 'Open'
        }),
        headers: {
            'Content-Type': 'application/json',
            'Api-Key': config.backend.apiKey
        }
    })
}

// ================================
module.exports.command = {
    isPremium: false,
    permLevel: 0,

    data: data,
    execute: execute
}