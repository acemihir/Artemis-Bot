// ================================
const { SlashCommandBuilder } = require('@discordjs/builders')
const { createId, filterText } = require('../utils')
const { runQuery } = require('../structures/database')
const { MessageEmbed, MessageActionRow, MessageButton } = require('discord.js')
const config = require('../config')
const fetch = require('node-fetch')

// ================================
const data = new SlashCommandBuilder()
    .setName('suggest')
    .setDescription('Create a suggestion.')
    .addStringOption(opt =>
        opt.setName('description')
            .setDescription('A brief description of your suggestion.')
            .setRequired(true))

const execute = async function(client, interaction) {
    // Fetch the input/args
    const sugDesc = await interaction.options.getString('description')

    const result = await runQuery('SELECT suggestion_channel FROM servers WHERE id = $1::text', [interaction.guild.id])
    if (result == null || !result.rows.length) {
        await interaction.reply('ERROR: Please make sure an administrator has configured the suggestion channel.')
        return
    }

    const sugChannel = await interaction.guild.channels.fetch(result.rows[0].suggestion_channel)
    if (sugChannel == null) {
        await interaction.reply('ERROR: The configured suggestion channel was not found.')
        return
    }

    const sugId = createId('s_')

    const embed = new MessageEmbed()
        .setAuthor(interaction.author.tag, interaction.author.avatarURL())
        .setColor('#7583ff')
        .addField('Description', filterText(sugDesc))
        .addField('Information', `**Status:** Open\n**ID":** ${sugId}\n\n*0 - upvotes | 0 - downvotes*`)

    const row = new MessageActionRow().addComponents(
        new MessageButton()
            .setCustomId('sug_upvote')
            .setLabel('Upvote')
            .setStyle('SUCCESS')
            .setEmoji('⬆️'),
        new MessageButton()
            .setCustomId('sug_downvote')
            .setLabel('Downvote')
            .setStyle('ERROR')
            .setEmoji('⬇️')
    )

    const msg = await sugChannel.send({ embeds: [embed], components: [row] })

    await interaction.reply('Your suggestion has been submitted.')

    fetch(`${config.backend.url}/submit`, {
        method: 'POST',
        body: JSON.stringify({
            id: sugId,
            context: sugDesc,
            author: interaction.user.id,
            avatar: interaction.user.avatarURL(),
            guild: interaction.guildId,
            channel: sugChannel.id,
            message: msg.id,
            status: 'Open'
        }),
        headers: {
            'Content-Type': 'application/json',
            'Api-Key': config.backend.apiKey
        }
    })
}

const buttonActions = [
    {
        id = 'sug_upvote',
        onClick = async function(client, interaction) {
            // Update the message and update the count
            console.log('DEBUG: ' + interaction)
        }
    },
    {
        id = 'sug_downvote',
        onClick = async function(client, interaction) {
            // Update the message and update the count
            console.log('DEBUG: ' + interaction)
        }
    }
]

// ================================
module.exports.command = {
    isPremium: false,
    permLevel: 0,

    data: data,
    execute: execute,
    buttons: buttonActions
}