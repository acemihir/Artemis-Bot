// ================================
const { SlashCommandBuilder } = require('@discordjs/builders')

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
            .addStringOption(opt => opt.setName('id').setDescription('The ID of the report.'))
            .addStringOption(opt => opt.setName('status').setDescription('The new status of the report.')
                .addChoices([
                    ['Open', 'open'],
                    ['Resolved', 'resolved'],
                    ['Progressing', 'progressing']
                ]))
    )

const execute = async function(client, interaction) {
    await interaction.reply('to be done')
}

// ================================
module.exports.command = {
    isPremium: false,
    privileged: true,

    data: data,
    execute: execute
}