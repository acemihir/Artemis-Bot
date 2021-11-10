// ================================
const { SlashCommandBuilder } = require('@discordjs/builders')
const { createId, filterText } = require('../utils')
const { MessageEmbed, MessageActionRow, MessageButton } = require('discord.js')
const config = require('../config')
const { getFromRedis } = require('../structures/cache')

const fetch = (...args) => import('node-fetch').then(({default: fetch}) => fetch(...args))

// ================================
const data = new SlashCommandBuilder()
    .setName('suggest')
    .setDescription('Create a suggestion.')
    .addStringOption(opt =>
        opt.setName('description')
            .setDescription('A brief description of your suggestion.')
            .setRequired(true))

const execute = async function(_client, interaction) {
    // Fetch the input/args
    const sugDesc = await interaction.options.getString('description')

    const cache = await getFromRedis(interaction.guildId)
    if (cache.sug_channel == null) {
        await interaction.reply('🟥 Please make sure an administrator has configured the suggestion channel.')
        return
    }

    const sugChannel = await interaction.guild.channels.fetch(cache.sug_channel)
    if (sugChannel == null) {
        await interaction.reply('🟥 The configured suggestion channel was not found.')
        return
    }

    const sugId = createId('s_')

    const embed = new MessageEmbed()
        .setAuthor(interaction.user.tag, interaction.user.avatarURL())
        .setColor(config.embedColor.b)
        .setDescription(`**Description:** ${filterText(sugDesc)}\n\n**Status:** Open\n**Id:** ${sugId}\n\n0 - upvotes | 0 - downvotes`)

    const approveEmoji = cache.approve_emoji == null ? '⬆️' : cache.approve_emoji
    const rejectEmoji = cache.reject_emoji == null ? '⬇️' : cache.reject_emoji

    const row = new MessageActionRow().addComponents(
        new MessageButton()
            .setCustomId('sug_upvote')
            .setLabel('Upvote')
            .setStyle('SUCCESS')
            .setEmoji(approveEmoji),
        new MessageButton()
            .setCustomId('sug_downvote')
            .setLabel('Downvote')
            .setStyle('DANGER')
            .setEmoji(rejectEmoji)
    )

    let msg
    try {
        msg = await sugChannel.send({ embeds: [embed], components: [row] })
    } catch (ex) {
        console.error(ex)
    }

    await interaction.reply('🟩 Your suggestion has been submitted.')

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

module.exports.buttons = [
    {
        id: 'sug_upvote',
        onClick: async function(_client, interaction) {
            const response = await fetch(`${config.backend.url}/suggestions/upvote`, {
                method: 'POST',
                body: JSON.stringify({
                    message: interaction.message.id,
                    guild: interaction.guildId,
                    user_id: interaction.user.id
                }),
                headers: {
                    'Content-Type': 'application/json',
                    'Api-Key': config.backend.apiKey
                }
            })

            const data = await response.json()
            if (data['success']) {
                // Get the old embed
                var embed = interaction.message.embeds[0]

                // Split it for each line
                const msgArray = embed.description.split('\n')
                // Manipulate the votes
                msgArray[msgArray.length - 1] = `${data["new_upvotes"]} - upvotes | ${data["new_downvotes"]} - downvotes`

                // Set the manipulated embed description
                embed.description = msgArray.join('\n')

                // Update the first message
                interaction.message.edit({ embeds: [embed] })

                interaction.deferUpdate()
            } else {
                interaction.reply({ content: data["error"], ephemeral: true })
            }
        }
    },
    {
        id: 'sug_downvote',
        onClick: async function(_client, interaction) {
            const response = await fetch(`${config.backend.url}/suggestions/downvote`, {
                method: 'POST',
                body: JSON.stringify({
                    message: interaction.message.id,
                    guild: interaction.guildId,
                    user_id: interaction.user.id
                }),
                headers: {
                    'Content-Type': 'application/json',
                    'Api-Key': config.backend.apiKey
                }
            })

            const data = await response.json()
            if (data['success']) {
                // Get the old embed
                var embed = interaction.message.embeds[0]

                // Split it for each line
                const msgArray = embed.description.split('\n')
                // Manipulate the votes
                msgArray[msgArray.length - 1] = `${data["new_upvotes"]} - upvotes | ${data["new_downvotes"]} - downvotes`

                // Set the manipulated embed description
                embed.description = msgArray.join('\n')

                // Update the first message
                interaction.message.edit({ embeds: [embed] })

                interaction.deferUpdate()
            } else {
                interaction.reply({ content: data["error"], ephemeral: true })
            }
        }
    }
]

// ================================
module.exports.command = {
    isPremium: false,
    privileged: true,

    data: data,
    execute: execute
}