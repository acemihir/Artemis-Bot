// ================================
const { SlashCommandBuilder } = require('@discordjs/builders')
const { Constants, MessageEmbed } = require('discord.js-light')
const config = require('../config')

const fetch = (...args) => import('node-fetch').then(({default: fetch}) => fetch(...args))

// ================================
const data = new SlashCommandBuilder()
    .setName('move')
    .setDescription('Move a suggestion/report.')
    .setDefaultPermission(false)  
    .addStringOption(opt => opt.setName('id').setDescription('The ID of the suggestion/report.').setRequired(true))
    .addChannelOption(opt => opt.setName('channel').setDescription('The channel where the message should be moved to.').setRequired(true))

const execute = async function(_client, interaction) {
    const id = interaction.options.getString('id')

    // We must get the message id first
    const res = await fetch(`${config.backend.url}/fetch/${interaction.guildId}/${id}`, {
        headers: {
            'Content-Type': 'application/json',
            'Api-Key': config.backend.apiKey
        }
    })

    const body = await res.json()
    if (!body['success']) {
        return interaction.reply(body['error'])
    }

    // Check old channel
    const channel = await interaction.guild.channels.fetch(body.data[0].channel)
    if (channel == null) {
        return interaction.reply({
            embeds: [new MessageEmbed()
                .setColor(config.embedColor.r)
                .setDescription('Couldn\'t find the channel the corresponding message was placed in.')
            ]
        })
    }

    // Try to get old message
    let msg
    try {
        msg = await channel.messages.fetch(body.data[0].message)
    } catch (ex) {
        if (ex.code !== Constants.APIErrors.UNKNOWN_MESSAGE) {
            console.log(ex)
        }
        return interaction.reply({
            embeds: [new MessageEmbed()
                .setColor(config.embedColor.r)
                .setDescription('Couldn\'t find the corresponding message.')
            ]
        })
    }

    // New channel checks
    const newChannel = interaction.options.getChannel('channel')
    if (newChannel == null || newChannel.deleted) {
        return interaction.reply({
            embeds: [new MessageEmbed()
                .setColor(config.embedColor.r)
                .setDescription('Couldn\'t find that channel.')
            ]
        })
    }

    // Get the old message
    const embed = msg.embeds[0]
    const row = msg.components[0]

    // Send new message
    let newMsg
    try {
        newMsg = await newChannel.send({ embeds: [embed], components: [row] })
    } catch (ex) {
        console.error(ex)
        return interaction.reply({
            embeds: [new MessageEmbed()
                .setColor(config.embedColor.r)
                .setDescription('Something went wrong while creating the new message.')
            ]
        })
    }

    await fetch(`${config.backend.url}/move`, {
        method: 'POST',
        body: JSON.stringify({
            id: id,
            guild: interaction.guildId,
            channel: newChannel.id,
            message: newMsg.id
        }),
        headers: {
            'Content-Type': 'application/json',
            'Api-Key': config.backend.apiKey
        }
    })

    // Delete old message
    try {
        await msg.delete()
    } catch (ex) {
        console.error(ex)
        return interaction.reply({
            embeds: [new MessageEmbed()
                .setColor(config.embedColor.r)
                .setDescription('Could not delete the message, delete it manually.')
            ]
        })
    }

    interaction.reply({
        embeds: [new MessageEmbed()
            .setColor(config.embedColor.g)
            .setDescription('Successfully moved the message to another channel.')
        ]
    })
}

// ================================
module.exports.command = {
    isPremium: false,
    privileged: true,

    data: data,
    execute: execute
}