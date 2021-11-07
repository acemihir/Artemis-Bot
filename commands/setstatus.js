// ================================
const { SlashCommandBuilder } = require('@discordjs/builders')
const config = require('../config')

const fetch = (...args) => import('node-fetch').then(({default: fetch}) => fetch(...args))

// ================================
const data = new SlashCommandBuilder()
    .setName('setstatus')
    .setDescription('Change the status of a suggestion/report/poll.')
    .setDefaultPermission(false)

    // Suggestions
    .addSubcommand(scmd =>
        scmd.setName('suggestions').setDescription('Change the status for a suggestion.')
            .addStringOption(opt => opt.setName('id').setDescription('The ID of the suggestion.'))
            .addStringOption(opt => opt.setName('status').setDescription('The new status of the suggestion.')
                .addChoices([
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

const execute = async function(client, interaction) {
    const id = interaction.options.getString('id')
    var status = interaction.options.getString('status')

    const res = await fetch(`${config.backend.url}/setstatus`, {
        method: 'POST',
        body: JSON.stringify({
            id: id,
            guild: interaction.guildId,
            status: status.charAt(0).toUpperCase() + status.slice(1)
        }),
        headers: {
            'Content-Type': 'application/json',
            'Api-Key': config.backend.apiKey
        }
    })

    const body = await res.json()
    if (body["success"]) {
        interaction.reply(`Message status was successfully changed to ${status}.`)
        console.log(body)
    } else {
        interaction.reply(body.error)
    }
}

module.exports.statusUpdate = async function(msg, status) {
    const colour = config.embedColor.b

    if (status === 'Approved' || status === 'Resolved') {
        colour = config.embedColor.g
    } else if (status === 'Considering' || status === 'Progressing') {
        colour = config.embedColor.y  
    } else if (status === 'Rejected') {
        colour = config.embedColor.r
    }

    // ...
}

// ================================
module.exports.command = {
    isPremium: false,
    privileged: true,

    data: data,
    execute: execute
}