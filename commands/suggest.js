// ================================
const { SlashCommandBuilder } = require('@discordjs/builders')
const { createId, filterText } = require('../utils')
const { MessageEmbed, MessageActionRow, MessageButton } = require('discord.js-light')
const config = require('../config')
const { getFromRedis } = require('../structures/cache')

const fetch = (...args) => import('node-fetch').then(({ default: fetch }) => fetch(...args))

// ================================
const data = new SlashCommandBuilder()
    .setName('suggest')
    .setDescription('Create a suggestion.')
    .addStringOption(opt =>
        opt.setName('description')
            .setDescription('A brief description of your suggestion.')
            .setRequired(true))

const execute = async function (interaction) {
    await interaction.deferReply()

    const cache = await getFromRedis(interaction.guildId)
    if (cache.sug_channel == null) {
        return interaction.editReply({
            embeds: [new MessageEmbed()
                .setColor(config.embedColor.r)
                .setDescription('Please make sure an administrator has configured the suggestion channel.')
            ]
        })
    }

    const sugChannel = await interaction.guild.channels.fetch(cache.sug_channel)
    if (sugChannel == null) {
        return interaction.editReply({
            embeds: [new MessageEmbed()
                .setColor(config.embedColor.r)
                .setDescription('The configured suggestion channel was not found.')
            ]
        })
    }

    const sugId = createId('s_')
    const sugDesc = interaction.options.getString('description')

    let msg
    try {
        msg = await sugChannel.send({
            embeds: [new MessageEmbed()
                .setAuthor(interaction.user.tag, interaction.user.avatarURL())
                .setColor(config.embedColor.b)
                .setDescription(`**Description:** ${filterText(sugDesc)}\n\n**Status:** Open\n**Id:** ${sugId}\n\n0 - upvotes | 0 - downvotes`)],
            components: [new MessageActionRow().addComponents(
                new MessageButton()
                    .setCustomId('sug_upvote')
                    .setLabel('Upvote')
                    .setStyle('SUCCESS')
                    .setEmoji(cache.approve_emoji),
                new MessageButton()
                    .setCustomId('sug_downvote')
                    .setLabel('Downvote')
                    .setStyle('DANGER')
                    .setEmoji(cache.reject_emoji)
            )]
        })
    } catch (ex) {
        return interaction.reply({
            embeds: [new MessageEmbed()
                .setColor(config.embedColor.r)
                .setDescription('I could not send the suggestion message. (This could be permission related)')
            ]
        })
    }

    await fetch(`${config.backend.url}/submit`, {
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

    await interaction.editReply({
        embeds: [new MessageEmbed()
            .setColor(config.embedColor.g)
            .setDescription('Your suggestion has been submitted.')
        ]
    })
}

module.exports.buttons = [
    {
        id: 'sug_upvote',
        onClick: async function (interaction) {
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
                msgArray[msgArray.length - 1] = `${data['new_upvotes']} - upvotes | ${data['new_downvotes']} - downvotes`

                // Set the manipulated embed description
                embed.description = msgArray.join('\n')

                // Update the first message
                interaction.message.edit({ embeds: [embed] })

                interaction.deferUpdate()
            } else {
                interaction.reply({ content: data['error'], ephemeral: true })
            }
        }
    },
    {
        id: 'sug_downvote',
        onClick: async function (interaction) {
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
                msgArray[msgArray.length - 1] = `${data['new_upvotes']} - upvotes | ${data['new_downvotes']} - downvotes`

                // Set the manipulated embed description
                embed.description = msgArray.join('\n')

                // Update the first message
                interaction.message.edit({ embeds: [embed] })

                interaction.deferUpdate()
            } else {
                interaction.reply({ content: data['error'], ephemeral: true })
            }
        }
    }
]

// ================================
module.exports.command = {
    data: data,
    execute: execute
}