// ================================
const { SlashCommandBuilder } = require('@discordjs/builders')
const { Constants, MessageEmbed } = require('discord.js-light')
const config = require('../config')

const fetch = (...args) => import('node-fetch').then(({default: fetch}) => fetch(...args))

// ================================
const data = new SlashCommandBuilder()
    .setName('setstatus')
    .setDescription('Change the status of a suggestion/report/poll.')

    // Suggestions
    .addSubcommand(scmd =>
        scmd.setName('suggestions').setDescription('Change the status for a suggestion.')
            .addStringOption(opt => opt.setName('id').setDescription('The ID of the suggestion.'))
            .addStringOption(opt => opt.setName('status').setDescription('The new status of the suggestion.').addChoices([
                ['Open', 'open'],
                ['Approved', 'approved'],
                ['Rejected', 'rejected'],
                ['Considering', 'considering']
            ]))
    )

    // Reports
    .addSubcommand(scmd =>
        scmd.setName('reports').setDescription('Change the status for a report.')
            .addStringOption(opt => opt.setName('id').setDescription('The ID of the report.').setRequired(true))
            .addStringOption(opt => opt.setName('status').setDescription('The new status of the report.').addChoices([
                ['Open', 'open'],
                ['Resolved', 'resolved'],
                ['Progressing', 'progressing']
            ]).setRequired(true))
    )

const execute = async function(_client, interaction) {
    const id = interaction.options.getString('id')

    const status = interaction.options.getString('status')
    const desStatus = status.charAt(0).toUpperCase() + status.slice(1)

    const res = await fetch(`${config.backend.url}/setstatus`, {
        method: 'POST',
        body: JSON.stringify({
            id: id,
            guild: interaction.guildId,
            status: desStatus
        }),
        headers: {
            'Content-Type': 'application/json',
            'Api-Key': config.backend.apiKey
        }
    })

    const body = await res.json()
    if (body['success']) {
        const msgId = body['messageId']
        const chnId = body['channelId']

        const channel = await interaction.guild.channels.fetch(chnId)
        if (channel == null) {
            return interaction.reply({
                embeds: [new MessageEmbed()
                    .setColor(config.embedColor.r)
                    .setDescription('Couldn\'t find the channel the corresponding message was placed in.')
                ]
            })
        }

        let msg
        try {
            msg = await channel.messages.fetch(msgId)
        } catch (ex) {
            if (ex.code !== Constants.APIErrors.UNKNOWN_MESSAGE) {
                console.log(ex)
            }
        }

        if (msg == null || msg.deleted) {
            return interaction.reply({
                embeds: [new MessageEmbed()
                    .setColor(config.embedColor.r)
                    .setDescription('Couldn\'t find the corresponding message.')
                ]
            })
        }

        await statusUpdate(msg, desStatus)
        interaction.reply({
            embeds: [new MessageEmbed()
                .setColor(config.embedColor.g)
                .setDescription(`Message status was successfully changed to ${status}.`)
            ]
        })
    } else {
        interaction.reply(body['error'])
    }
}

// ================================
module.exports.command = {
    privileged: true,

    data: data,
    execute: execute
}

// ================================
const statusUpdate = async (msg, status) => {
    let colour = config.embedColor.b

    if (status === 'Approved' || status === 'Resolved') {
        colour = config.embedColor.g
    } else if (status === 'Considering' || status === 'Progressing') {
        colour = config.embedColor.y  
    } else if (status === 'Rejected') {
        colour = config.embedColor.r
    }

    const embed = msg.embeds[0]
    
    const descArray = embed.description.split('\n')

    const statusline = descArray[2]
    const slArray = statusline.split(' ')

    descArray[2] = slArray[0] + ' ' + status

    embed.description = descArray.join('\n')
    embed.color = parseInt(colour.slice(1), 16)

    await msg.edit({ embeds: [embed] })
}

module.exports.statusUpdate = statusUpdate