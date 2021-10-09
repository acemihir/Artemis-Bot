// ================================
const { SlashCommandBuilder } = require('@discordjs/builders')
const { createId, filterText } = require('../utils')
const { MessageEmbed } = require('discord.js')
const config = require('../config')
const { getFromRedis } = require('../structures/cache')

const fetch = (...args) => import('node-fetch').then(({default: fetch}) => fetch(...args));

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

    const cache = await getFromRedis(interaction.guildId)
    if (cache.rep_channel == null) {
        await interaction.reply('ERROR: Please make sure an administrator has configured the report channel.')
        return
    }

    const repChannel = await interaction.guild.channels.fetch(cache.rep_channel)
    if (repChannel == null) {
        await interaction.reply('ERROR: The configured report channel was not found.')
        return
    }

    const repId = createId('r_')

    const embed = new MessageEmbed()
        .setAuthor(interaction.user.tag, interaction.user.avatarURL())
        .setColor(config.embedColor.b)
        .setDescription(`**Description:** ${filterText(repDesc)}\n\n**Status:** Open\n**Id:** ${repId}`)

    let msg
    try {
        msg = await repChannel.send({ embeds: [embed] })
    } catch (ex) {
        console.error(ex)
    }

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
    privileged: false,

    data: data,
    execute: execute
}